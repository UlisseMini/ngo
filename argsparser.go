// Command line argument parser for ncat
package main

import (
	// TODO find a better library for args parsing
	// flag is utter trash :l
	"errors"
	"github.com/UlisseMini/ngo/internal/aes"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"time"
)

var (
	// on/off options
	listen *bool = flag.BoolP("listen", "l", false, "listen for connections")
	udp    *bool = flag.BoolP("udp", "u", false, "use udp instead of tcp")
	ssl    *bool = flag.Bool("ssl", false, "enable ssl [WARNING NO HOST VERIFICATION]")

	// int options
	timeoutFlag *int  = flag.IntP("timeout", "t", 10, "connection timeout in seconds")
	debugLevel  *uint = flag.UintP("debug", "d", 4, "logging level 0-6")

	// handled in main.go
	cmdStr *string = flag.StringP("exec", "e", "",
		"run command and redirect file descriptors to the connection")
	aesKey *string = flag.StringP("aes", "a", "",
		"encrypt the connection using a key and AES")

	// parsed options
	addr    string
	proto   string
	timeout time.Duration
)

func parseArgs() error {
	// manage aes flag
	defaultKey := `Not entering an AES key is very bad, luckly ngo is smarter then you`
	flag.Lookup("aes").NoOptDefVal = defaultKey

	flag.Parse()
	// TODO allow them to supply in another format other then ip:port
	if len(flag.Args()) != 1 {
		return errors.New("You must specify a host to connect to")
	}

	// set addr to the first argument
	addr = flag.Arg(0)
	proto = "tcp"
	if *udp == true {
		proto = "udp"
	}

	// set the dial timeout
	timeout = time.Duration(*timeoutFlag) * time.Second

	// set the logging levels
	log.SetLevel(log.Level(*debugLevel))
	aes.SetLoggingLevel(*debugLevel)

	return nil
}
