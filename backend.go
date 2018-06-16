package web_pusher

import "net/http"

type Backend interface {
	Listen()
}

type backend struct {
	server Server
}

func newBackend(server Server) *backend {
	return &backend{server}
}

func (b *backend) Listen() {
	http.HandleFunc("/send_to_channel", b.sendToChannel)
	http.HandleFunc("/send_to_user", b.sendToUser)
}

func (b *backend) sendToUser(rw http.ResponseWriter, r *http.Request) {
}

func (b *backend) sendToChannel(rw http.ResponseWriter, r *http.Request) {
}
