// aes manages encryption for ngo.
package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	log "github.com/sirupsen/logrus"
	"io"
)

func SetLoggingLevel(i uint) {
	// set the logging level
	level := log.Level(i)
	log.SetLevel(level)
}

// Packet is a a packet to the other side, encoded using gob
type Packet struct {
	// Data holds the encrypted data
	Data []byte

	// nonce does not need to be kept secret,
	// it needs to be shared for the remote end to decrypt our data
	Nonce []byte
}

// readWriter
type readWriter struct {
	// the underlying writer (intended to be the connection)
	encoder *gob.Encoder
	decoder *gob.Decoder

	// the cipher block
	aesgcm cipher.AEAD
}

// The readWriter method
// TODO:
// * Correctly return number of bytes written, currently just returns zero.
func (self *readWriter) Write(data []byte) (n int, err error) {
	// First Get the nonce
	nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return n, err
	}
	log.Tracef("length of nonce: %d", len(nonce))
	log.Tracef("nonce: %v", nonce)

	// Encrypt it
	encrypted := self.aesgcm.Seal(nil, nonce, data, nil)
	log.Tracef("encrypted data: %v", encrypted)

	// Encode into the packet
	p := &Packet{
		Data:  encrypted,
		Nonce: nonce,
	}

	log.Tracef("encoding %#v", p)
	err = self.encoder.Encode(p)
	if err != nil {
		return n, err
	}

	return len(data), err
}

func (self *readWriter) Read(data []byte) (n int, err error) {
	p := &Packet{}
	log.Tracef("err = self.decoder.Decode(p)")
	err = self.decoder.Decode(p)
	if err != nil {
		return 0, err
	}

	log.Tracef("decrypting self.Data with '%d' nonce", p.Nonce)
	msg, err := self.aesgcm.Open(nil, p.Nonce, p.Data, nil)
	if err != nil {
		return 0, err
	}

	n = copy(data, msg)
	log.Tracef("Read %d bytes into data", n)
	return n, nil
}

// NewReadWriter is a high level wrapper that creates an readWriter
// TODO:
// Make into a readwriter and change the name.
func NewReadWriter(rw io.ReadWriter, k string) (io.ReadWriter, error) {
	// compute the hash of the key, this will make it the correct length
	keyBytes := []byte(k)
	key := sha256.Sum256(keyBytes[:])

	// Create the new cipher using the key
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &readWriter{
		// the encoder to use \w gob
		encoder: gob.NewEncoder(rw),
		decoder: gob.NewDecoder(rw),

		// the cipher block
		aesgcm: aesgcm,
	}, nil
}
