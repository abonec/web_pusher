package web_pusher

type Connection interface {
	Send(msg []byte)
	Close(msg []byte)
}

type connection struct {
}

