package main

import (
	"github.com/abonec/web_pusher"
	"net/http"
	"log"
)

func main() {
	server := web_pusher.NewServer(&App{})
	server.Start()

	back := web_pusher.NewBackend(server, web_pusher.StandardLogger{})
	front := web_pusher.NewFrontend(server, web_pusher.StandardLogger{})

	http.HandleFunc("/send-to-actor/", back.SendToUser)
	http.HandleFunc("/ws", front.Handle)

	log.Println("Server started at localhost:8083")
	log.Fatal(http.ListenAndServe(":8083", nil))
}
