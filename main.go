package main

import (
	"fmt"
	"github.com/UlisseMini/ngo/exec"
	"github.com/UlisseMini/utils/cmd"
	"io"
	"net"
	"os"
)

func init() {
	// Parse commandline arguments (argsparser.go)
	err := ParseArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// initalize loggers (logging.go)
	InitLoggers()
}

func main() {
	var (
		err  error
		conn net.Conn
	)

	if !*listen {
		conn, err = net.DialTimeout(proto, addr, timeout)
		mustNot(err)

		// print the connected message (diferent depending on proto)
		if *udp {
			Info.Println("Sending to", addr)
		} else {
			Info.Println("Connected to", addr)
		}
	} else {
		// listening
		l, err := net.Listen(proto, addr)
		mustNot(err)

		Info.Printf("Listening on %s", addr)
		conn, err = l.Accept()
		Info.Printf("Connection from %s", conn.RemoteAddr().String())
		mustNot(err)
	}
	defer conn.Close()

	// if there is a command to execute over the conn
	if *cmdStr != "" {
		cmd := cmd.Parse(*cmdStr)
		err := exec.Spawn(conn, cmd)
		errPrint(err)
		return
	}

	// otherwise use default
	handleConn(conn)
}

// handleConn connects the two connections file descriptors.
func handleConn(conn net.Conn) {
	done := make(chan struct{})

	// connect conn to stdout
	go func() {
		n, err := io.Copy(os.Stdout, conn)
		errPrint(err)

		Debug.Printf("Read %d bytes\n", n)
		done <- struct{}{}
	}()

	// connect stdin to conn
	go func() {
		n, err := io.Copy(conn, os.Stdin)
		errPrint(err)

		Debug.Printf("Wrote %d bytes\n", n)
		done <- struct{}{}
	}()

	// wait for one of the goroutines to finish
	<-done
}

// errPrint prints an error using the Error logger
// if the error is not nil
func errPrint(err error) {
	if err != nil {
		Error.Println(err)
	}
}

// mustNot prints and exits the program on error.
func mustNot(err error) {
	if err != nil {
		Error.Println(err)
		os.Exit(1)
	}
}
