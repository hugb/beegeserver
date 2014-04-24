package models

func init() {
	WSSwitcher.RegisterHandler("menu", getMenu)
	WSSwitcher.RegisterHandler("tree", getTree)
}

type menu struct {
	Type  int    `json:"type"`
	Title string `json:"title"`
}

func getMenu(conn *Connection, message *WSMessage) {
	message.Code = 0
	message.Data = []menu{
		menu{Type: 1, Title: "Image"},
		menu{Type: 1, Title: "Container"},
	}
	conn.SendMessage(message)
}

type tree struct {
	Type  int    `json:"type"`
	Leaf  bool   `json:"leaf"`
	Title string `json:"title"`
}

func getTree(conn *Connection, message *WSMessage) {
	message.Code = 0
	message.Data = []tree{
		tree{Type: 1, Leaf: true, Title: "Image"},
		tree{Type: 1, Leaf: true, Title: "Container"},
	}
	conn.SendMessage(message)
}
