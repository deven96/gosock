/*
 Simple Chat application using sockets in GoLang
*/
package main

import (
//	"os"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
	"strings"
	"github.com/deven96/gosock/pkg/custlog"
)

var TEMPLATE_FOLDER string = filepath.Join("templates")

// templateHandler represents a single template
type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}

// ServeHTTP handles the HTTPRequest
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func(){
		filearr := []string{TEMPLATE_FOLDER, t.filename}
		filepath := strings.Join(filearr, "/")
		t.templ = template.Must(template.ParseFiles(filepath))
	})
	t.templ.Execute(w, nil)
}

func main() {
	// set log name along with default outputs
	def_writers := custlog.DefaultWriters("main.log", false)
	//TRACE will be Discarded, while the rest will be routed accordingly
	custlog.LogInit(def_writers)	
	custlog.Trace.Println("Imported Custom Logging")
	custlog.Info.Println("Log file can be found at ", custlog.Logfile)
	//os.Setenv("GOSOCK_LOG", custlog.Logfile)

	// create a room
	r := newRoom()
	
	/* Routes */
	// Handle function for route "/"
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	//start the room
	custlog.Info.Println("Initializing Room...")
	go r.run()
	// port variable
	port := ":8008"
	//start the webserver
	custlog.Info.Println("Running server started on Port", port)
	
	if err := http.ListenAndServe(port, nil); err != nil {
		custlog.Error.Println(err)
	}
}
