// Command line argument parser for ncat
package main

import (
	// TODO find a better library for args parsing
	// flag is utter trash :l
	"errors"
	"flag"
	"github.com/UlisseMini/ngo/internal/aes"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	// on/off options
	listen *bool = flag.Bool("l", false, "listen for connections")
	udp    *bool = flag.Bool("u", false, "use udp instead of tcp")
	ssl    *bool = flag.Bool("ssl", false, "enable ssl [WARNING NO HOST VALIDATION]")

	// int options
	timeoutFlag *int  = flag.Int("t", 10, "timeout in seconds")
	debugLevel  *uint = flag.Uint("d", 4, "debug level 0-6")

	// handled in main.go
	cmdStr *string = flag.String("e", "",
		"run command and redirect file descriptors to the connection")
	aesKey *string = flag.String("a", "", "encrypt the connection using a key and AES")

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

	// set the aes logging level
	aes.SetLoggingLevel(*debugLevel)

	return nil
}
