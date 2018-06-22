# Simple router

This simple example just show how to use standard http backend and websocket client.

In this example implemented App (app.go) with no auth logic and User (user.go). All other parts from the given package.

## Usage
>$ go run *.go

## Frontend
Connect to ws://localhost:8083/ws via websockets and send any message.

## Backend
>$ curl -XPOST   --data '{"hello":"world"}' http://localhost:8083/send-to-actor/1137

