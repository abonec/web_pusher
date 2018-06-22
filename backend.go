package web_pusher

import (
	"net/http"
	"strings"
	"bytes"
	"time"
)

type backend struct {
	server *Server
	logger Logger
}

func NewBackend(server *Server, logger Logger) *backend {
	if logger == nil {
		logger = NopeLogger{}
	}
	return &backend{server, logger}
}

// Handler accepts request with POST body that would be sent to the given client as is
func (b *backend) SendToUser(rw http.ResponseWriter, r *http.Request) {
	var statusCode int
	defer b.completeLog(time.Now(), &statusCode)

	rw.Header().Set("Content-Type", "application/json")
	userId := strings.Split(r.URL.Path, "/")[2]
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	msg := buf.Bytes()
	b.startLog(r, string(msg))

	if len(msg) == 0 {
		statusCode = http.StatusBadRequest
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	b.server.sendToUser(userId, buf.Bytes())
	statusCode = http.StatusOK
	rw.WriteHeader(http.StatusOK)
}

func (b *backend) SendToChannel(rw http.ResponseWriter, r *http.Request) {
}

func (b *backend) completeLog(startTime time.Time, statusCode *int) {
	b.logger.Printf("[BACK] Completed %d %s in %s", *statusCode, http.StatusText(*statusCode), time.Since(startTime))
}

func (b *backend) startLog(r *http.Request, msg string) {
	b.logger.Printf("[BACK] %s %s %s %s\n", r.RemoteAddr, r.Method, r.URL, msg)
}
