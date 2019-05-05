/*
 Simple Chat application using sockets in GoLang
*/


package main

import (
	"flag"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
	"strings"
	"github.com/deven96/gosock/pkg/custlog"
)

// declare global variables for use throughout the main package
var TEMPLATE_DIR = filepath.Join("templates")
var ASSETS_DIR = filepath.Join("assets")
// command line arguments and defaults
var LOG_FILE = flag.String("log", "gosock.log", "Name of the log file to save to")
var SERVER_LOCATION = flag.String("addr", ":8008", "The addr of the application.")

// templateHandler represents a single template
type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}


// ServeHTTP handles the HTTPRequest
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func(){
		filearr := []string{TEMPLATE_DIR, t.filename}
		filepath := strings.Join(filearr, "/")
		t.templ = template.Must(template.ParseFiles(filepath))
	})
	t.templ.Execute(w, r)
}

func main() {
	flag.Parse() // parse the flags
	def_writers := custlog.DefaultWriters(*LOG_FILE, false)
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
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir(ASSETS_DIR))))
	


	
	//start the room
	custlog.Info.Println("Initializing Room...")
	go r.run()
	//start the webserver
	custlog.Info.Printf("Running server started on %s", *SERVER_LOCATION)
	
	if err := http.ListenAndServe(*SERVER_LOCATION, nil); err != nil {
		custlog.Error.Println(err)
	}
}
