package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than readWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from client.
	maxMessageSize = 512
)

// connection is an middleman between the websocket connection and the hub.
type Connection struct {
	// The websocket connection.
	ws *websocket.Conn

	auth bool

	// Buffered channel of outbound messages.
	send chan []byte

	data map[string]string
}

type WSMessage struct {
	Cmd  string      `json:"cmd"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type WSHandler func(conn *Connection, message *WSMessage)

type switcher struct {
	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *Connection

	// Unregister requests from connections.
	unregister chan *Connection

	handlers map[string]WSHandler

	// Registered connections.
	connections map[string]*Connection
}

var WSSwitcher = &switcher{
	broadcast:   make(chan []byte, maxMessageSize),
	register:    make(chan *Connection, 1),
	unregister:  make(chan *Connection, 1),
	handlers:    make(map[string]WSHandler),
	connections: make(map[string]*Connection),
}

func init() {
	log.Trace("Run switcher.")
	go WSSwitcher.run()
}

func (this *switcher) run() {
	for {
		select {
		case c := <-this.register:
			this.connections[c.ws.RemoteAddr().String()] = c
		case c := <-this.unregister:
			delete(this.connections, c.ws.RemoteAddr().String())
			close(c.send)
		case m := <-this.broadcast:
			for _, c := range this.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(this.connections, c.ws.RemoteAddr().String())
				}
			}
		}
	}
}

func (this *switcher) RegisterConn(conn *Connection) {
	WSSwitcher.register <- conn
}

func (this *switcher) Broadcast(message []byte) {
	WSSwitcher.broadcast <- message
}

func (this *switcher) Unicast(address string, message []byte) {
	if conn, exist := this.connections[address]; exist {
		conn.send <- message
	}
}

func (this *switcher) RegisterHandler(name string, handler WSHandler) error {
	if _, exists := this.handlers[name]; exists {
		return fmt.Errorf("can't overwrite handler for command %s", name)
	} else {
		this.handlers[name] = handler
	}
	return nil
}

// new websocket connection
func NewWSConn(ws *websocket.Conn) *Connection {
	return &Connection{
		ws:   ws,
		send: make(chan []byte, 256),
		data: make(map[string]string),
	}
}

func (this *Connection) Send(message []byte) {
	this.send <- message
}

func (this *Connection) SendMessage(message *WSMessage) {
	if b, err := json.Marshal(message); err == nil {
		this.send <- b
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (this *Connection) ReadPump() {
	defer func() {
		WSSwitcher.unregister <- this
		this.ws.Close()
	}()
	this.ws.SetReadLimit(maxMessageSize)
	this.ws.SetReadDeadline(time.Now().Add(pongWait))
	this.ws.SetPongHandler(func(string) error {
		this.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		msgType, message, err := this.ws.ReadMessage()
		if err != nil {
			log.Error("Read message fails:%s.", err)
			break
		}
		log.Trace("Recive message,type:%d,content:%s.", msgType, string(message))

		msg := &WSMessage{}
		if err = json.Unmarshal(message, msg); err != nil {
			log.Error("Receive message is decoded fails:%s.", err)
			continue
		}
		log.Trace("cmd:%s,code:%d,data:%v.", msg.Cmd, msg.Code, msg.Data)

		if msg.Cmd != "login" && !this.auth {
			msg.Cmd = "login"
			msg.Code = 1
			if b, err := json.Marshal(msg); err == nil {
				this.Send(b)
			}
			continue
		}

		if fct, exist := WSSwitcher.handlers[msg.Cmd]; exist {
			fct(this, msg)
		}
	}
	log.Trace("Websocket connection is closed.")
}

// write writes a message with the given opCode and payload.
func (this *Connection) write(opCode int, payload []byte) error {
	log.Trace("Write a message,operate code:%d,content:%s.", opCode, string(payload))
	this.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return this.ws.WriteMessage(opCode, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (this *Connection) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		this.ws.Close()
	}()
	for {
		select {
		case message, ok := <-this.send:
			if !ok {
				this.write(websocket.CloseMessage, []byte{})
				log.Debug("Send close message.")
				return
			}
			if err := this.write(websocket.TextMessage, message); err != nil {
				log.Error("Send text message fails:%s.", err)
				return
			}
		case <-ticker.C:
			if err := this.write(websocket.PingMessage, []byte{}); err != nil {
				log.Error("Send ping message fails.", err)
				return
			}
		}
	}
	log.Trace("Websocker connection is closed.")
}
