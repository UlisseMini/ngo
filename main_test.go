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
	// disable the logging levels
	log.SetLevel(0)
	aes.SetLoggingLevel(0)
}

// getpkt returns a byte array of random bytes
func getpkt(size int) []byte {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return buf
}

type testcase struct {
	// the config to pass to connect
	conf config

	// should our side listen?
	listen bool

	packetSize int
	packet     []byte
}

// Test the connect function. creates a listener in a new goroutine
// then connects using the `connect` function from `main.go`
func Test_connect(t *testing.T) {
	const packetSize = 16 // size of the packet in bytes
	tt := []testcase{
		testcase{
			conf: config{
				proto:   "tcp",
				addr:    "127.0.0.1:31893",
				timeout: 1 * time.Second,
			},
			listen:     true,
			packetSize: packetSize,
			packet:     getpkt(packetSize),
		},
		testcase{
			conf: config{
				proto:   "tcp",
				addr:    "127.0.0.1:40913",
				timeout: 1 * time.Second,
				listen:  true,
			},
			listen:     false,
			packetSize: packetSize,
			packet:     getpkt(packetSize),
		},
	}

	for _, tc := range tt {
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

		errChan := make(chan error, 1)
		if tc.listen {
			// listen for the client
			go listen(tc, errChan)
			time.Sleep(10 * time.Millisecond)
		} else {
			// connect to the client (that is listening)
			go func() {
				time.Sleep(10 * time.Millisecond)
				connectClient(tc, errChan)
			}()
		}

		select {
		case err := <-errChan:
			t.Fatalf("errChan: %v", err)
		default:
			t.Logf("Testing connect with: %#v", tc.conf)
			conn, err := connect(tc.conf)
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
		}
	}
}

func listen(tc testcase, errChan chan error) {
	l, err := net.Listen(tc.conf.proto, tc.conf.addr)
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

	err = readPacket(conn, tc)
	if err != nil {
		errChan <- err
	}
}

func connectClient(tc testcase, errChan chan error) {
	conn, err := net.DialTimeout(tc.conf.proto, tc.conf.addr, 10*time.Millisecond)
	if err != nil {
		errChan <- err
		return
	}
	defer conn.Close()

	err = readPacket(conn, tc)
	if err != nil {
		errChan <- err
	}
}

func readPacket(conn net.Conn, tc testcase) error {
	buf := make([]byte, tc.packetSize)
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}
	if n < tc.packetSize {
		return fmt.Errorf("readPacket: n (%d) is less then packetSize (%d)",
			n, tc.packetSize)
	}
	// now make sure its valid
	if bytes.Compare(buf, tc.packet) != 0 {
		return fmt.Errorf("readPacket: buf (%X) != packet (%X)", buf, tc.packet)
	}

	return nil
}
