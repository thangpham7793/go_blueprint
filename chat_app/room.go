package main

import (
	"log"
	"net/http"

	"chat_app/trace"

	"github.com/gorilla/websocket"
)

type room struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	tracer  trace.Tracer
}

//constructor
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		//only one block of case will run, so only one goroutine will modify the map
		select {
		//add the client pointer to the map and mark them as active
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("New Client Joined!")

		//remove a client pointer from the map when they leave (as in they join the leave channel)
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left!")
		//forward all messages received to the send channel of each client
		case msg := <-r.forward:
			r.tracer.Trace("Received: ", string(msg))
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace("Sent to client")
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: messageBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	client := &client{socket, make(chan []byte, messageBufferSize), r}
	r.join <- client
	defer func() {
		r.leave <- client
	}()
	go client.writeToClientSocket()
	client.readFromClientSocket()
}
