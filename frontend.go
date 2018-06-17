package web_pusher

import (
	"net/http"
	"github.com/gorilla/websocket"
	"log"
)

type Frontend struct {
	server *Server
}

type WebSocketConnection struct {
	conn *websocket.Conn
}

func NewFrontend(server *Server) *Frontend {
	return &Frontend{server}
}

var upgrader = websocket.Upgrader{}

func (front *Frontend) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer conn.Close()

	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("error while authenticating", err)
		return
	}

	_, err = front.server.Auth(NewWebSocketConnection(conn), message)
	if err != nil {
		log.Println("error while auth:", err)
		return
	}

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("error while read message:", err)
			break
		}

		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("error while send message:", err)
		}
	}
}

func NewWebSocketConnection(conn *websocket.Conn) *WebSocketConnection {
	return &WebSocketConnection{conn}
}

func (conn *WebSocketConnection) Send(msg []byte) bool {
	conn.conn.WriteMessage(websocket.TextMessage, msg)
	return true
}
func (conn *WebSocketConnection) Close(msg []byte) {

}
