// Command line argument parser for ncat
package main

import (
	// TODO find a better library for args parsing
	// flag is utter trash :l
	"flag"
	"time"
)

var (
	// on/off options
	listen  *bool = flag.Bool("l", false, "listen for connections")
	udp     *bool = flag.Bool("u", false, "use udp instead of tcp")
	verbose *bool = flag.Bool("v", false, "be verbose")
	debug   *bool = flag.Bool("d", false, "debug mode, more logging")

	// int options
	timeoutFlag *int = flag.Int("t", 10, "timeout in seconds")

	// handled in main.go
	cmdStr *string = flag.String("e", "",
		"run command and redirect file descriptors to the connection")
	// parsed options
	addr    string
	proto   string
	timeout time.Duration
)

func ParseArgs() error {
	flag.Parse()
	// TODO allow them to supply in another format other then ip:port
	if len(flag.Args()) != 1 {
		return NoHostErr
	}

	// set addr to the first argument
	addr = flag.Arg(0)
	proto = "tcp"
	if *udp == true {
		proto = "udp"
	}

	// set the dial timeout
	timeout = time.Duration(*timeoutFlag) * time.Second

	return nil
}
