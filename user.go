package web_pusher

type User interface {
	Id() string
	Send(msg []byte) bool
}

type user struct {
	id   string
	conn Connection
}

func (u *user) Id() string {
	return u.id
}

func (u *user) Send(msg []byte) bool {
	u.conn.Send(msg)
	return true
}

func NewUser(id string, conn Connection) *user {
	return &user{id, conn}
}

type AuthUser interface {
	Id() string
}
