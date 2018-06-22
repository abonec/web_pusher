package main

import "github.com/abonec/web_pusher"

type App struct {
}

func (app *App) Auth(msg []byte) (web_pusher.AuthUser, []string, error) {
	return NewUser(), []string{"main"}, nil
}
