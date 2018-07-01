package application

type User interface {
	Id() string
	Send(msg []byte) bool
}

type AuthUser interface {
	Id() string
}
