package main

import(
	"fmt"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/deven96/gosock/pkg/custlog"
)

func init(){
	// append to the main log
	defwriters := custlog.DefaultWriters(*LogFile, true)
	//TRACE will be Discarded, while the rest will be routed accordingly
	custlog.LogInit(defwriters)	
	
}

// generate random color string
func randHex() (string, error) {
	bytez := make([]byte, 3)
	if _, err := rand.Read(bytez); err != nil {
		return "", err
	}
	ret := "#" + hex.EncodeToString(bytez)
	custlog.Info.Printf("Assigning %s color to new client", ret)
	return ret, nil
}

type room struct {
	//forward is then channel holding the incoming messages
	// that should be forwarded to other clients
	forward chan *message
	// clients that want to join
	join chan *client
	// clients that want to leave
	leave chan *client
	// holds all current clients in the room
	clients map[*client]bool

}

func (r *room) run(){
	for {
		select {
			case client:= <- r.join:
				// joining
				r.clients[client] = true
				custlog.Info.Printf("Client %s joined the room...", client.color)
				msg := &message{
					Code: "#ffffff",
					Message: fmt.Sprintf("Client %s joined the room", client.color),
				}
				for client := range r.clients {
					client.send <- msg
				}
			case client:= <- r.leave:
				// leaving
				// delete client
				custlog.Warning.Printf("Removing Client %s from room clients and closing send channel...", client.color)
				delete(r.clients, client)
				// close send channel of client
				close(client.send)
				msg := &message{
					Code: "#ffffff",
					Message: fmt.Sprintf("Client %s left the room", client.color),
				}
				for client := range r.clients {
					client.send <- msg
				}
			case msg := <- r.forward:
				for client := range r.clients {
					client.send <- msg
				}
		}
	}
}


// message constants
const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

// set read and write size
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

// ServeHTTP for room struct
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request){
	socket, err := upgrader.Upgrade(w, req, nil)
	if err!= nil {
		custlog.Trace.Println(err)
		return
	} 
	gencolor, _ := randHex()
	// pass client address to room
	client := &client {
		socket : socket,
		send: make(chan *message, messageBufferSize),
		room: r,
		color: gencolor,
	}
	r.join <- client
	defer func() {r.leave <- client} ()
	go client.write()
	client.read()
}

// creates a new room
func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
	}
}
