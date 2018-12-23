// Command line argument parser for ncat
package main

import (
	// TODO find a better library for args parsing
	// flag is utter trash :l
	"errors"
	"flag"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	// on/off options
	listen *bool = flag.Bool("l", false, "listen for connections")
	udp    *bool = flag.Bool("u", false, "use udp instead of tcp")
	ssl    *bool = flag.Bool("ssl", false, "enable ssh [WARNING NO HOST VALIDATION]")

	// int options
	timeoutFlag *int  = flag.Int("t", 10, "timeout in seconds")
	debugLevel  *uint = flag.Uint("d", 6, "debug level 0-6")

	// handled in main.go
	cmdStr *string = flag.String("e", "",
		"run command and redirect file descriptors to the connection")
	// parsed options
	addr    string
	proto   string
	timeout time.Duration
)

func parseArgs() error {
	flag.Parse()
	// TODO allow them to supply in another format other then ip:port
	if len(flag.Args()) != 1 {
		return errors.New("You must specify a host to connect to. QUITTING.")
	}

	// set addr to the first argument
	addr = flag.Arg(0)
	proto = "tcp"
	if *udp == true {
		proto = "udp"
	}

	// set the dial timeout
	timeout = time.Duration(*timeoutFlag) * time.Second

	// set the logging level
	level := log.Level(*debugLevel)
	log.SetLevel(level)

	return nil
}
