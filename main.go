package main

import (
	"net/http"
	custlog "lib/custlog"
)

func main() {
	// set log name along with default outputs
	var def_writers custlog.Writers = custlog.DefaultWriters("test.log")
	//TRACE will be Discarded, while the rest will be routed accordingly
	custlog.LogInit(def_writers)	
	custlog.Trace.Println("Imported Custom Logging")
	custlog.Info.Println("Log file can be found at ", custlog.Logfile)
	// Handle function for root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte (`
			<html>
				<head>
					<title> GoSOCK </title>
				</head>
				<body>
					Let's Discuss Go!
				</body>
			</html>`,
			))
	})
	// port variable
	port := ":8008"
	//start the webserver
	custlog.Info.Println("Running server started on Port", port)
	
	if err := http.ListenAndServe(port, nil); err != nil {
		custlog.Error.Println(err)
	}
}
