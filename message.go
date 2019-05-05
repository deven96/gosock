package main

import (
	"time"
	"github.com/deven96/gosock/pkg/custlog"
)

func init(){
	//initialize log file
	defwriters := custlog.DefaultWriters(*LogFile, true)
	custlog.LogInit(defwriters)
}

type message struct {
	Code string
	Message string
	When time.Time
}
