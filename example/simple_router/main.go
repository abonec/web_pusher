package main

import (
	"github.com/abonec/web_pusher/backend/http_backend"
	"github.com/abonec/web_pusher/backend/redis_pubsub"
	"github.com/abonec/web_pusher/frontend/ws_frontend"
	"github.com/abonec/web_pusher/logger/stdout_logger"
	"github.com/abonec/web_pusher/server"
	"log"
	"net/http"
)

func main() {
	server := server.NewServer(&App{})
	server.Start()

	logger := stdout_logger.StandardLogger{}

	back := http_backend.NewBackend(server, logger)
	redisBack := redis_pubsub.NewBackend(server, logger)
	err := redisBack.Start("sendToActor:*")
	if err != nil {
		logger.Printf("error while starting redis pubsub: %s", err)
	}
	front := ws_frontend.NewFrontend(server, stdout_logger.StandardLogger{})

	http.HandleFunc("/send-to-actor/", back.SendToUser)
	http.HandleFunc("/ws", front.Handle)

	log.Println("Server started at localhost:8083")
	log.Fatal(http.ListenAndServe(":8083", nil))
}
