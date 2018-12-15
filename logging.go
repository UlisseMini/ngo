// logging.go contains loggers used by the rest of the program
package main

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

func InitLoggers() {
	// write logs to standard error by default.
	writer := os.Stderr

	// manage verbose and info flags
	infoWriter := ioutil.Discard
	debugWriter := ioutil.Discard
	if *verbose {
		infoWriter = os.Stderr
	}
	if *debug {
		debugWriter = os.Stderr
	}

	// 16 = log.Lshortfile
	Warn = log.New(writer, "[WARN] ", 16)
	Error = log.New(writer, "[ERROR] ", 16)
	Info = log.New(infoWriter, "[INFO] ", 16)
	Debug = log.New(debugWriter, "[DEBUG] ", 16)
}
