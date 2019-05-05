package main

import (
	"time"
	"github.com/gorilla/websocket"
	"github.com/deven96/gosock/pkg/custlog"
)

func init(){
	// log file flag
	defwriters := custlog.DefaultWriters(*LogFile, true)
	//TRACE will be Discarded, while the rest will be routed accordingly
	custlog.LogInit(defwriters)	
	
}


type client struct {
	// websocket per client
	socket *websocket.Conn
	// send is a channel on which messages are sent
	send chan *message
	// room is the room this client is chatting on
	room *room
	// color assigned to the client
	color string
}

// read from the websocket
func (c *client) read(){
	defer c.socket.Close()

	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			custlog.Trace.Println(err)
			return
		}
		// send msg to the forward channel of the client
		custlog.Info.Printf("Message received from %s: %s", c.color, msg.Message)
		msg.When = time.Now()
		msg.Code = c.color
		c.room.forward <- msg
	}
}

// write to the websocket
func (c *client) write(){
	defer c.socket.Close()

	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			custlog.Error.Println(err)
			return
		}
	}
}
