package web_pusher

import (
	"net/http"
	"strings"
	"bytes"
)

type Backend interface {
	Listen()
}

type backend struct {
	server *Server
}

func NewBackend(server *Server) *backend {
	return &backend{server}
}

func (b *backend) Listen() {
	http.HandleFunc("/send-to-user", b.sendToUser)
	http.HandleFunc("/send-to-channel", b.sendToChannel)
}

func (b *backend) sendToUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	userId := strings.Split(r.URL.Path, "/")[2]
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	msg := buf.Bytes()
	if len(msg) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	b.server.sendToUser(userId, buf.Bytes())
	rw.WriteHeader(http.StatusOK)
}

func (b *backend) sendToChannel(rw http.ResponseWriter, r *http.Request) {
}
