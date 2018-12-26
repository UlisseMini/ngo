package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
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
	const packetSize = 16 // size of the packet in bytes
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
		// generate the packet
		packet := make([]byte, packetSize)
		n, err := rand.Read(packet)
		if err != nil {
			// this should never error, but just in case
			t.Error(err)
			continue
		} else if n < packetSize {
			t.Errorf("n (%d) is less then packetSize (%d)", n, packetSize)
			continue
		}

		errChan := make(chan error)
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
			defer conn.Close()

			// read the packet
			buf := make([]byte, packetSize)
			n, err := conn.Read(buf)
			if err != nil {
				errChan <- err
				return
			}
			if n < packetSize {
				errChan <- fmt.Errorf("server: n(%d) is less then packetSize(%d)",
					n, packetSize)
				return
			}
			// now make sure its valid
			if bytes.Compare(buf, packet) != 0 {
				errChan <- fmt.Errorf("server: buf (%X) != packet (%X)", buf, packet)
			}
		}()

		select {
		case err := <-errChan:
			t.Fatalf("server returned: %v", err)
		default:
			// give the listener time to get ready (should really find better way ;p
			time.Sleep(10 * time.Millisecond)

			conn, err := connect(conf)
			if err != nil {
				t.Fatalf("connect returned: %v", err)
			}
			defer conn.Close()

			// now we have the connection send the packet
			n, err := conn.Write(packet)
			if err != nil {
				t.Fatalf("failed to send packet: %v", err)
			}
			if n < packetSize {
				t.Fatalf("n (%d) < packetSize (%d)", n, packetSize)
			}
			// TODO listen for the packet back
		}
	}
}
