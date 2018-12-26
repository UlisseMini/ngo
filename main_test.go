package main

import (
	"github.com/UlisseMini/ngo/internal/aes"
	log "github.com/sirupsen/logrus"
	"net"
	"testing"
	"time"
)

// initalize the tests by disabling logging from the main program.
func init() {
	// set the logging levels
	log.SetLevel(0)
	aes.SetLoggingLevel(0)
}

// Test the connect function. creates a listener in a new goroutine
// then connects using the `connect` function from `main.go`
// TODO
// Add more test cases
// Add verification of a correct connection by transmiting a packet back and forth
func Test_connect(t *testing.T) {
	testAddr := "127.0.0.1:31893"
	conf := config{
		proto:   "tcp",
		addr:    testAddr,
		timeout: 1 * time.Second,
	}

	errChan := make(chan error, 1)
	go func() {
		l, err := net.Listen("tcp", testAddr)
		if err != nil {
			errChan <- err
			return
		}
		defer l.Close()

		conn, err := l.Accept()
		if err != nil {
			errChan <- err
			return
		}
		err = conn.Close()
		errChan <- err
	}()

	select {
	case err := <-errChan:
		t.Fatalf("server: %v", err)
	default:
		// give it time to get ready
		time.Sleep(100 * time.Millisecond)

		conn, err := connect(conf)
		if err != nil {
			t.Fatalf("client: %v", err)
		}
		defer conn.Close()
	}
}
