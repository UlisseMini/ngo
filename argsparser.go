// Command line argument parser for ncat
package main

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"os"
	"time"
)

type config struct {
	listen bool
	udp    bool
	ssl    bool

	addr    string
	proto   string
	timeout time.Duration

	cmdStr string

	// file descriptors to be connected to the connection.
	in  io.Reader
	out io.Writer
}

// return a config based on commandline arguments,
// if commandline arguments are invalid it will call os.Exit(1)
func parseArgs() config {
	var (
		// on/off options
		listen *bool = flag.BoolP("listen", "l", false, "listen for connections")
		udp    *bool = flag.BoolP("udp", "u", false, "use udp instead of tcp")
		ssl    *bool = flag.Bool("ssl", false, "enable ssl")

		// int options
		timeoutFlag *time.Duration = flag.DurationP("timeout", "t", 10*time.Second,
			"connection timeout in seconds")
		debugLevel *uint = flag.UintP("debug", "d", 4, "logging level 0-6")

		// handled in main.go
		cmdStr *string = flag.StringP("exec", "e", "",
			"run command and redirect file descriptors to the connection")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <ip:port> [flags...]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	// instance of config
	conf := config{}

	// TODO allow them to supply in another format other then ip:port
	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	// set addr to the first argument
	conf.addr = flag.Arg(0)
	conf.proto = "tcp"
	if *udp == true {
		conf.proto = "udp"
	}

	// set the dial timeout
	conf.timeout = *timeoutFlag

	// set the logging levels
	log.SetLevel(log.Level(*debugLevel))

	// set in and out to the default file descriptors
	conf.in = os.Stdin
	conf.out = os.Stdout

	conf.listen = *listen
	conf.udp = *udp
	conf.ssl = *ssl
	conf.cmdStr = *cmdStr

	return conf
}
