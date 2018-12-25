package main

import (
	"crypto/tls"
	"fmt"
	"github.com/UlisseMini/ngo/internal/aes"
	"github.com/UlisseMini/ngo/internal/exec"
	"github.com/UlisseMini/ngo/internal/tlsconfig"
	"github.com/UlisseMini/utils/cmd"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
)

const hostname = "ngo"

func main() {
	// Parse commandline arguments (argsparser.go)
	conf, err := parseArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// connect
	conn, err := connect(conf)
	mustNot(err)

	var rw io.ReadWriter = io.ReadWriter(conn)
	// encrypt the connection with AES (if needs be)
	if *aesKey != "" {
		log.Tracef("AES mode enabled key: %s", *aesKey)
		rw, err = aes.NewReadWriter(conn, *aesKey)
		mustNot(err)
	}

	// if there is a command to execute over the connection
	if *cmdStr != "" {
		log.Infof("executing: %q over the connection", *cmdStr)
		cmd := cmd.Parse(*cmdStr)
		exec.Spawn(rw, cmd)
		return
	}

	// Don't force handleConn to close the connection, since i can't be
	// bothered to implement `close` in internal/aes
	func() {
		defer conn.Close()
		handleConn(conf, rw)
	}()
}

func connect(conf config) (net.Conn, error) {
	var (
		err  error
		conn net.Conn
	)

	if !*listen {
		if *ssl {
			// connect with ssl / tls
			tlsconf := &tls.Config{InsecureSkipVerify: true}
			conn, err := tls.Dial(conf.proto, conf.addr, tlsconf)
			if err != nil {
				return nil, err
			}
			return conn, nil
		}
		conn, err = net.DialTimeout(conf.proto, conf.addr, conf.timeout)
		if err != nil {
			return nil, err
		}

		// print the connected message (diferent depending on proto)
		if *udp {
			log.Info("Sending to", conf.addr)
		} else {
			log.Info("Connected to", conf.addr)
		}
	} else {
		// listening
		var l net.Listener
		if *ssl {
			config, err := tlsconfig.Get(hostname)
			if err != nil {
				return nil, err
			}

			l, err = tls.Listen(conf.proto, conf.addr, config)
			if err != nil {
				return nil, err
			}

		} else {
			l, err = net.Listen(conf.proto, conf.addr)
			if err != nil {
				return nil, err
			}
		}

		log.Infof("Listening on %s", conf.addr)
		conn, err = l.Accept()
		log.Infof("Connection from %s", conn.RemoteAddr().String())
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}

// handleConn connects the two connections file descriptors.
// fd is not a file descriptor, but that is what i'll usually be passing.
// (except in tests)
func handleConn(conf config, conn io.ReadWriter) (err error) {
	done := make(chan error)

	// connect conn to stdout
	go func() {
		n, err := io.Copy(conf.out, conn)
		errPrint(err)

		log.Debugf("Read %d bytes\n", n)
		done <- err
	}()

	// connect stdin to conn
	go func() {
		n, err := io.Copy(conn, conf.in)
		errPrint(err)

		log.Debugf("Wrote %d bytes\n", n)
		done <- err
	}()

	// wait for one of the goroutines to finish
	return <-done
}

// errPrint prints an error using the Error logger
// if the error is not nil
func errPrint(err error) {
	if err != nil {
		log.Error(err)
	}
}

// mustNot prints and exits the program on error.
func mustNot(err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
