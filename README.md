[![Build Status](https://travis-ci.org/abonec/web_pusher.svg?branch=master)](https://travis-ci.org/abonec/web_pusher)
[![Coverage Status](https://img.shields.io/coveralls/github/abonec/web_pusher/master.svg)](https://coveralls.io/github/abonec/web_pusher?branch=master)
# web_pusher

Framework to connect backends to frontend for notification

The library was broken into three main parts: frontend, backend, and server. Web pusher was developed keeping modularity in mind. For use web_pusher you need to implement the only App with simple Auth method for authorization incoming connections and User with simplest interface.
All parts can be replaced with own implementation. You can use different methods for message dispatching.
Frontend used to connect web client to the app. Standard implementation use websockets to it. A client should connect to the frontend, send auth message as the first message and wait for a messages from the server.
Backend used to send messages to the connected clients. Standard implementation use simple http post handler.
The server handles and store all incoming frontend connections and send messages from backend to client.
