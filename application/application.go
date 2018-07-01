package application

type Application interface {
	Auth(msg []byte) (AuthUser, []string, error)
}

