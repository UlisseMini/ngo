package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/UlisseMini/ngo/internal/exec"
	"github.com/UlisseMini/ngo/internal/tlsconfig"
	"github.com/UlisseMini/utils/cmd"
	log "github.com/sirupsen/logrus"
)

// hostname for our tls certificate
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

	rw := io.ReadWriter(conn)

	// if there is a command to execute over the connection
	if conf.cmdStr != "" {
		log.Infof("executing: %q over the connection", conf.cmdStr)
		cmd := cmd.Parse(conf.cmdStr)
		exec.Spawn(rw, cmd)
		return
	}

	// Don't force handleConn to close the connection, i want to keep it
	// as flexible as possible
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

	if !conf.listen {
		if conf.ssl {
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
		if conf.udp {
			log.Infof("Sending to %s", conf.addr)
		} else {
			log.Infof("Connected to %s", conf.addr)
		}
	} else {
		// listening
		var l net.Listener
		if conf.ssl {
			config, err := tlsconfig.Get(hostname)
			if err != nil {
				return nil, fmt.Errorf("getting tlsconfig: %v", err)
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
// the i/o file descriptors are in conf
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
	err = <-done
	log.Debugf("handleConn exiting: %#v", err)
	return err
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
