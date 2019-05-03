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
	def_writers := custlog.DefaultWriters("gosock.log", true)
	//TRACE will be Discarded, while the rest will be routed accordingly
	custlog.LogInit(def_writers)	
	
}

// generate random color string
func randHex() (string, error) {
	bytes := make([]byte, 3)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	ret := "#" + hex.EncodeToString(bytes)
	custlog.Info.Printf("Assigning %s color to new client", ret)
	return ret, nil
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
				custlog.Info.Printf("Client %s joined the room...", client.color)
				txt := fmt.Sprintf("ADMIN: Client %s joined the room", client.color)
				msg := []byte(txt)
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
				msg := []byte(fmt.Sprintf("ADMIN: Client %s left the room", client.color))
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
		custlog.Error.Println(err)
		return
	} 
	gen_color, _ := randHex()
	// pass client address to room
	client := &client {
		socket : socket,
		send: make(chan []byte, messageBufferSize),
		room: r,
		color: gen_color,
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
