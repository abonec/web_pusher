package web_pusher

type Connection interface {
	Send(msg []byte) bool
	Close(msg []byte)
}

type connection struct {
}
