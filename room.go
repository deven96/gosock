package main

import(
//	"os"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/deven96/gosock/pkg/custlog"
)

func init(){
	// append to the main log
	def_writers := custlog.DefaultWriters("main.log", true)
	//TRACE will be Discarded, while the rest will be routed accordingly
	custlog.LogInit(def_writers)	
	
}

type room struct {
	//forward is then channel holding the incoming messages
	// that should be forwarded to other clients
	forward chan []byte
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
				custlog.Info.Printf("Client %v is joining...", client)
			case client:= <- r.leave:
				// leaving
				// delete client
				custlog.Warning.Printf("Deleting Client %v from room clients and closing send channel...", client)
				delete(r.clients, client)
				// close send channel of client
				close(client.send)
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
		custlog.Error.Println(err)
		return
	} 
	client := &client {
		socket : socket,
		send: make(chan []byte, messageBufferSize),
		room: r,
	}
	r.join <- client
	defer func() {r.leave <- client} ()
	go client.write()
	client.read()
}

// creates a new room
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
	}
}
