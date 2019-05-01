package custlog

import (
	"log"
	"io"
	"io/ioutil"
	"os"
)

// logging variables
var (
	Trace, Info, Warning, Error *log.Logger
	Logfile string
)

/** initialize a handler for every trace writer e.g os.Stderr
	Args:
		tracehandle
		infohandle
		warninghandle
		errorhandle
		log_name: name of the output log file
**/

// declare a new writer struct whenever you want to change default writers
type Writers struct {
	Tracehandle, Infohandle, Warninghandle, Errorhandle io.Writer
	Append bool
	Logfile string //name of logfile
}

//default writer values
func DefaultWriters (file_name string, append bool) Writers {
	return Writers{
		Tracehandle: ioutil.Discard,
		Infohandle: os.Stdout,
		Warninghandle: os.Stdout,
		Errorhandle: os.Stderr,
		Logfile: file_name,
		Append: append,		
	}
}

func LogInit(w Writers){

		// gets current working directory
		dir, err := os.Getwd()
		// creates a log file and sets it into append mode
		var file *os.File
		file, err = os.OpenFile(w.Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if w.Append == false{
			file, err = os.Create(w.Logfile)
		}
		if err != nil {
		    log.Fatalln("Failed to open log file ", w.Logfile, ":", err)
		}
		//empty log file
		//file.Write([]byte(``,))
		//export logfile name
		Logfile = dir + "/" + w.Logfile
		
		//create all log objects writing to the multiWriter i.e `handle` and a file
		Trace = log.New(
			io.MultiWriter(file, w.Tracehandle),
			"TRACE: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		)

		Info = log.New(
					io.MultiWriter(file, w.Infohandle),
					"INFO: ",
					log.Ldate|log.Ltime|log.Lshortfile,
		)

		Warning = log.New(
					io.MultiWriter(file, w.Warninghandle),
					"WARNING: ",
					log.Ldate|log.Ltime|log.Lshortfile,
		)

		Error = log.New(
					io.MultiWriter(file, w.Errorhandle),
					"ERROR: ",
					log.Ldate|log.Ltime|log.Lshortfile,
		)
}
