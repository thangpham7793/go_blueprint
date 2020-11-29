package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

func (c *client) readFromClientSocket() {
	defer c.socket.Close()
	//endless loop
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		//send the client's message to the room's forward channel, which will then broadcast the msg to other clients
		c.room.forward <- msg
	}
}

func (c *client) writeToClientSocket() {
	defer c.socket.Close()
	//for every msg received by the channel, write the message to the socket so that it sends to the client
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
