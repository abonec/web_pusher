package main

import (
	"github.com/abonec/web_pusher"
	"github.com/abonec/web_pusher/frontend/ws_frontend"
	"github.com/abonec/web_pusher/backend/http_backend"
	"github.com/abonec/web_pusher/logger/stdout_logger"
	"net/http"
	"log"
)

func main() {
	server := web_pusher.NewServer(&App{})
	server.Start()

	back := http_backend.NewBackend(server, stdout_logger.StandardLogger{})
	front := ws_frontend.NewFrontend(server, stdout_logger.StandardLogger{})

	http.HandleFunc("/send-to-actor/", back.SendToUser)
	http.HandleFunc("/ws", front.Handle)

	log.Println("Server started at localhost:8083")
	log.Fatal(http.ListenAndServe(":8083", nil))
}
