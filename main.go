package main

import (
	"crypto/tls"
	"fmt"
	"github.com/UlisseMini/ngo/internal/aes"
	"github.com/UlisseMini/ngo/internal/exec"
	"github.com/UlisseMini/utils/cmd"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
)

func init() {
	// Parse commandline arguments (argsparser.go)
	err := parseArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	conn, err := connect()
	mustNot(err)

	var rw io.ReadWriter
	// encrypt the connection with AES (if needs be)
	if *aesKey != "" {
		rw, err = aes.NewReadWriter(conn, *aesKey)
		mustNot(err)
	}

	// if there is a command to execute over the connection
	if *cmdStr != "" {
		log.Infof("executing: %q over the connection", *cmdStr)
		cmd := cmd.Parse(*cmdStr)
		exec.Spawn(rw, cmd)
	}

	func() {
		defer conn.Close()
		handleConn(rw)
	}()
}

// connect will connect to the correct host using the correct settings
func connect() (net.Conn, error) {
	var (
		err  error
		conn net.Conn
	)

	if !*listen {
		if *ssl {
			// connect with ssl / tls
			config := &tls.Config{InsecureSkipVerify: true}
			conn, err := tls.Dial(proto, addr, config)
			if err != nil {
				return nil, err
			}
			return conn, nil
		}
		conn, err = net.DialTimeout(proto, addr, timeout)
		if err != nil {
			return nil, err
		}

		// print the connected message (diferent depending on proto)
		if *udp {
			log.Info("Sending to", addr)
		} else {
			log.Info("Connected to", addr)
		}
	} else {
		// listening
		l, err := net.Listen(proto, addr)
		if err != nil {
			return nil, err
		}

		log.Infof("Listening on %s", addr)
		conn, err = l.Accept()
		log.Infof("Connection from %s", conn.RemoteAddr().String())
		if err != nil {
			return nil, err
		}

	}

	return conn, nil
}

// handleConn connects the two connections file descriptors.
func handleConn(conn io.ReadWriter) {
	done := make(chan struct{})

	// connect conn to stdout
	go func() {
		n, err := io.Copy(os.Stdout, conn)
		errPrint(err)

		log.Debugf("Read %d bytes\n", n)
		done <- struct{}{}
	}()

	// connect stdin to conn
	go func() {
		n, err := io.Copy(conn, os.Stdin)
		errPrint(err)

		log.Debugf("Wrote %d bytes\n", n)
		done <- struct{}{}
	}()

	// wait for one of the goroutines to finish
	<-done
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
