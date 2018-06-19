package web_pusher

import (
	"net/http"
	"github.com/gorilla/websocket"
)

type Frontend struct {
	server *Server
	logger Logger
}

type WebSocketConnection struct {
	conn     *websocket.Conn
	frontend *Frontend
	user     User
}

func NewFrontend(server *Server, logger Logger) *Frontend {
	if logger == nil {
		logger = NopeLogger{}
	}
	return &Frontend{server, logger}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (front *Frontend) Handle(w http.ResponseWriter, r *http.Request) {
	front.logf(VERBOSE_LOGGING, "Incoming connection %s\n", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		front.logf(VERBOSE_LOGGING, "Upgrade error for %s: %s\n", r.RemoteAddr, err.Error())
		return
	}
	defer conn.Close()
	front.logf(VERBOSE_LOGGING, "Upgrade successful")

	_, message, err := conn.ReadMessage()
	if err != nil {
		front.logf(VERBOSE_LOGGING, "Error while reading auth message: %s", err.Error())
		return
	}

	ws := NewWebSocketConnection(conn, front)
	u, err := front.server.Auth(ws, message)
	if err != nil {
		conn.WriteJSON(NewFrontendErrorMessage(AuthError, err))
		front.logf(VERBOSE_LOGGING, "Error while authentication: %s", err.Error())
		return
	}
	ws.user = u
	front.logf(VERBOSE_LOGGING, "Authentication successful for user id #%s", u.Id())
	conn.WriteJSON(NewFrontendSuccessMessage(u.Id()))
	conn.SetCloseHandler(func(code int, text string) error {
		front.server.Close(u)
		front.logf(VERBOSE_LOGGING, "Connection for #%s(%s) was closed", u.Id(), conn.RemoteAddr())
		return nil
	})

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		err = conn.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func (front *Frontend) logf(logLevel int, format string, v ...interface{}) {
	front.logger.Printf("[FRON] "+format, v...)
}

func NewWebSocketConnection(conn *websocket.Conn, frontend *Frontend) *WebSocketConnection {
	return &WebSocketConnection{conn: conn, frontend: frontend}
}

func (conn *WebSocketConnection) Send(msg []byte) bool {
	conn.conn.WriteMessage(websocket.TextMessage, msg)
	conn.frontend.logf(VERBOSE_LOGGING, "Send %s to %s(%s)", msg, conn.user.Id(), conn.conn.RemoteAddr())
	return true
}
func (conn *WebSocketConnection) Close(msg []byte) {

}
