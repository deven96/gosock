package main

import (
	"time"
	"github.com/deven96/gosock/pkg/custlog"
)

func init(){
	//initialize log file
	def_writers := custlog.DefaultWriters(*LOG_FILE, true)
	custlog.LogInit(def_writers)
}

type message struct {
	Code string
	Message string
	When time.Time
}
