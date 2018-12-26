package main

import (
	"bytes"
	"io"
	"testing"
)

// Test_handleConn test's the handleConn function by simulating a connection
func Test_handleConn(t *testing.T) {
	conn := bytes.NewBuffer(nil)
	monitored := bytes.NewBuffer(nil)

	// create a config for handleConn's local file descriptors
	conf := config{
		in:  monitored,
		out: monitored,
	}

	go func() {
		err := handleConn(conf, conn)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// now pretend to be the connected person
	msg := "ngo is pretty sweet!"
	_, err := conn.Write([]byte(msg))
	if err != nil {
		t.Fatal("Failed writing: ", err)
	}

	// see if it showed up on the other end.
	buf := make([]byte, 1024)
	n, err := monitored.Read(buf)
	if err != io.EOF && err != nil {
		t.Fatal(err)
	}

	received := string(buf[:n])
	// make sure they received the correct string.
	if received != msg {
		t.Fatalf("Expected: %s\ngot: %s\nGot %d bytes", msg, received, n)
	}
}
