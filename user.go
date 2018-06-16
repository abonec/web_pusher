package web_pusher

type User interface {
	Id() string
}

type user struct {
	id   string
	conn Connection
}

func (u *user) Id() string {
	return u.id
}

func NewUser(id string, conn Connection) *user {
	return &user{id, conn}
}
