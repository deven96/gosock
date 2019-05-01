package main

import (
//	"os"
	"github.com/gorilla/websocket"
	"github.com/deven96/gosock/pkg/custlog"
)

func init(){
	def_writers := custlog.DefaultWriters("main.log")
	//TRACE will be Discarded, while the rest will be routed accordingly
	custlog.LogInit(def_writers)	
	
}

type client struct {
	// websocket per client
	socket *websocket.Conn
	// send is a channel on which messages are sent
	send chan []byte
	// room is the room this client is chatting on
	room *room
}

// read from the websocket
func (c *client) read(){
	defer c.socket.Close()

	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			custlog.Error.Println(err)
			return
		}
		// send msg to the forward channel of the client
		c.room.forward <- msg
		custlog.Info.Printf("Reading message %v", msg)
	}
}

// write to the websocket
func (c *client) write(){
	defer c.socket.Close()

	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {			custlog.Error.Println(err)
			return
		}
	}
	custlog.Info.Printf("Writing message %v", c.send)
}
