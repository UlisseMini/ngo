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
	configs := []config{
		config{
			proto:   "tcp",
			addr:    "127.0.0.1:31893",
			timeout: 1 * time.Second,
		},
		config{
			proto:   "tcp",
			addr:    "127.0.0.1:40913",
			timeout: 1 * time.Second,
		},
	}

	for _, conf := range configs {
		errChan := make(chan error, 1)
		go func() {
			l, err := net.Listen(conf.proto, conf.addr)
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
			// give the listener time to get ready (should really find better way ;p
			time.Sleep(10 * time.Millisecond)

			conn, err := connect(conf)
			if err != nil {
				t.Fatalf("client: %v", err)
			}
			defer conn.Close()
		}
	}
}
