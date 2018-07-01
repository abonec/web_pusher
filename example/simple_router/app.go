package main

import (
	"github.com/abonec/web_pusher/application"
)

type App struct {
}

func (app *App) Auth(msg []byte) (application.AuthUser, []string, error) {
	return NewUser(), []string{"main"}, nil
}
