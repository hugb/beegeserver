package controllers

import (
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/hugb/beegeserver/models"
)

type WSController struct {
	baseController
}

func (this *WSController) Get() {
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if e, ok := err.(websocket.HandshakeError); ok {
		log.Error("Not a websocket handshake:%s.", e)
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Error("Websocket upgrade error,%s.", err)
		http.Error(this.Ctx.ResponseWriter, "Websocket upgrade fails.", 400)
		return
	}

	log.Trace("New one websocket connection.")
	conn := models.NewWSConn(ws)
	models.WSSwitcher.RegisterConn(conn)

	log.Trace("Start to read and write data.")
	go conn.WritePump()
	conn.ReadPump()
}
