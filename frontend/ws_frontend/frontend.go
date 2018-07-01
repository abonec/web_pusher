package ws_frontend

import (
	"net/http"
	"github.com/gorilla/websocket"
	"time"
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

// Handler accepts websocket connections.
// First message will be dispatched to the app Auth method.
// In case of failure auth connection will be closed.
// In case of auth message will not be sent in 5 seconds after establish connection it will be closed.
// All further incoming messages will be ignored.
func (front *Frontend) Handle(w http.ResponseWriter, r *http.Request) {
	front.logf(VERBOSE_LOGGING, "Incoming connection %s\n", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		front.logf(VERBOSE_LOGGING, "Upgrade error for %s: %s\n", r.RemoteAddr, err.Error())
		return
	}
	defer conn.Close()
	front.logf(VERBOSE_LOGGING, "Upgrade successful")

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, message, err := conn.ReadMessage()
	if err != nil {
		front.logf(VERBOSE_LOGGING, "Error while reading auth message: %s", err.Error())
		return
	}
	conn.SetReadDeadline(time.Time{})

	ws := NewWebSocketConnection(conn, front)
	u, err := front.server.Auth(ws, message)
	if err != nil {
		conn.WriteJSON(NewFrontendErrorMessage(AuthError, err))
		front.logf(VERBOSE_LOGGING, "Error while authentication: %s", err.Error())
		return
	}
	ws.user = u
	front.logf(VERBOSE_LOGGING, "Authentication successful for user id #%s, online %d/%d", u.Id(), front.server.OnlineUsers(), front.server.OnlineConnections())
	conn.WriteJSON(NewFrontendSuccessMessage(u.Id()))
	conn.SetCloseHandler(func(code int, text string) error {
		front.server.Close(u)
		front.logf(VERBOSE_LOGGING, "Connection for #%s(%s) was closed, online %d/%d", u.Id(), conn.RemoteAddr(), front.server.OnlineUsers(), front.server.OnlineConnections())
		return nil
	})

	for {
		_, _, err := conn.ReadMessage()
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

// Send message to the client as is
func (conn *WebSocketConnection) Send(msg []byte) bool {
	conn.conn.WriteMessage(websocket.TextMessage, msg)
	conn.frontend.logf(VERBOSE_LOGGING, "Send %s to %s(%s)", msg, conn.user.Id(), conn.conn.RemoteAddr())
	return true
}

func (conn *WebSocketConnection) Close(msg []byte) {
}
