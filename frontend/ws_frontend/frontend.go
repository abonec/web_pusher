package ws_frontend

import (
	"net/http"
	"github.com/gorilla/websocket"
	"time"
	"github.com/abonec/web_pusher"
	"github.com/abonec/web_pusher/logger"
	"github.com/abonec/web_pusher/logger/nope_logger"
	"github.com/abonec/web_pusher/application"
)

const (
	pongWait = 30 * time.Second
	pingWait = (pongWait * 8) / 10 // should be less than pongWait
)

type Frontend struct {
	server *web_pusher.Server
	logger logger.Logger
}

type WebSocketConnection struct {
	conn     *websocket.Conn
	frontend *Frontend
	user     application.User
}

func NewFrontend(server *web_pusher.Server, logger logger.Logger) *Frontend {
	if logger == nil {
		logger = nope_logger.NopeLogger{}
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
	front.logf(logger.VERBOSE_LOGGING, "Incoming connection %s\n", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		front.logf(logger.VERBOSE_LOGGING, "Upgrade error for %s: %s\n", r.RemoteAddr, err.Error())
		return
	}
	defer conn.Close()
	front.logf(logger.VERBOSE_LOGGING, "Upgrade successful")

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, message, err := conn.ReadMessage()
	if err != nil {
		front.logf(logger.VERBOSE_LOGGING, "Error while reading auth message: %s", err.Error())
		return
	}
	conn.SetReadDeadline(time.Now().Add(pongWait))

	ws := NewWebSocketConnection(conn, front)
	u, err := front.server.Auth(ws, message)
	if err != nil {
		conn.WriteJSON(NewFrontendErrorMessage(AuthError, err))
		front.logf(logger.VERBOSE_LOGGING, "Error while authentication: %s", err.Error())
		return
	}
	ws.user = u
	front.logf(logger.VERBOSE_LOGGING, "Authentication successful for user id #%s, online %d/%d", u.Id(), front.server.OnlineUsers(), front.server.OnlineConnections())
	conn.WriteJSON(NewFrontendSuccessMessage(u.Id()))
	conn.SetCloseHandler(func(code int, text string) error {
		closeConnection(conn, front, u)
		return nil
	})

	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	shutdownPinger := make(chan interface{})
	go func() {
		t := time.Tick(pingWait)

		for {
			select {
			case <-t:
				err := conn.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					closeConnection(conn, front, u)
					return
				}
			case <-shutdownPinger:
				return
			}
		}
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			close(shutdownPinger)
			closeConnection(conn, front, u)
			break
		}
	}
}

func closeConnection(conn *websocket.Conn, front *Frontend, u application.User) {
	front.server.Close(u)
	front.logf(logger.VERBOSE_LOGGING, "Connection for #%s(%s) was closed, online %d/%d", u.Id(), conn.RemoteAddr(), front.server.OnlineUsers(), front.server.OnlineConnections())
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
	conn.frontend.logf(logger.VERBOSE_LOGGING, "Send %s to %s(%s)", msg, conn.user.Id(), conn.conn.RemoteAddr())
	return true
}

func (conn *WebSocketConnection) Close(msg []byte) {
}
