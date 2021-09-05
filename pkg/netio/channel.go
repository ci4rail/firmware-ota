package netio

import (
	"errors"

	"github.com/ci4rail/firmware-ota/pkg/netio/transport"
	"google.golang.org/protobuf/proto"
)

// Channel holds the channels variables
type Channel struct {
	trans transport.Transport
}

// NewChannel creates a new channel using the transport mechanism in t
func NewChannel(t transport.Transport) (*Channel, error) {
	return &Channel{trans: t}, nil
}

// Close closes the transport stream
func (c *Channel) Close() {
	c.trans.Close()
}

// WriteMessage encodes m using protobuf and sends the encoded value through the transport stream
func (c *Channel) WriteMessage(m proto.Message) error {
	payload, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return c.Write(payload)
}

// ReadMessage waits for a new message in transport stream and decodes it via protobuf
func (c *Channel) ReadMessage(m proto.Message) error {
	payload, err := c.Read()
	if err != nil {
		return err
	}

	return proto.Unmarshal(payload, m)
}

// Write writes Netio standard message to the socket s
func (c *Channel) Write(payload []byte) error {
	// make sure we have the magic bytes
	err := c.writeMagicBytes()
	if err != nil {
		return err
	}

	length := uint(len(payload))
	err = c.writeLength(length)
	if err != nil {
		return err
	}

	err = c.writePayload(payload)
	if err != nil {
		return err
	}
	return nil
}

// writeMagicBytes write the magic bytes 0xFE, 0xED to s.Connection.
func (c *Channel) writeMagicBytes() error {
	magicBytes := []byte{0xFE, 0xED}

	err := c.writeBytesSafe(magicBytes)
	return err
}

// writeLength writes 4 bytes to s.Connection with the length
func (c *Channel) writeLength(length uint) error {
	lengthBytes := make([]byte, 4)

	lengthBytes[0] = byte(length & 0xFF)
	lengthBytes[1] = byte((length >> 8) & 0xFF)
	lengthBytes[2] = byte((length >> 16) & 0xFF)
	lengthBytes[3] = byte((length >> 24) & 0xFF)

	err := c.writeBytesSafe(lengthBytes)
	return err
}

// writePayload write the payload to s.Connection.
func (c *Channel) writePayload(payload []byte) error {
	err := c.writeBytesSafe(payload)
	return err
}

// writeBytesSafe retries writing to socket until all bytes are written
func (c *Channel) writeBytesSafe(payload []byte) error {
	for {
		written, err := c.trans.Write(payload)
		if err != nil {
			return err
		}
		if written == len(payload) {
			return nil
		}
		payload = payload[written:]
	}
}

// Read reads a Netio standard message from socket s
func (c *Channel) Read() ([]byte, error) {
	// make sure we have the magic bytes
	err := c.readMagicBytes()
	if err != nil {
		return nil, err
	}

	length, err := c.readLength()
	if err != nil {
		return nil, err
	}
	// log.Println("Length: ", length)
	payload, err := c.readPayload(length)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// readMagicBytes blocks until it receives the magic bytes 0xFE, 0xED from s.Connection.
func (c *Channel) readMagicBytes() error {
	// block until we get the magic bytes
	for {
		magicBytes := make([]byte, 2)
		for i := 0; i < 2; i++ {
			b := make([]byte, 1)

			_, err := c.trans.Read(b)
			if err != nil {
				return err
			}
			magicBytes[i] = b[0]
		}
		if magicBytes[0] == 0xFE && magicBytes[1] == 0xED {
			// log.Println(magicBytes[0], magicBytes[1], magicBytes[2], magicBytes[3])
			return nil
		}
	}
}

// readLength reads 4 bytes from s.Connection and returns the length as uint of the message.
func (c *Channel) readLength() (uint, error) {
	lengthBytes := make([]byte, 4)
	_, err := c.trans.Read(lengthBytes)
	if err != nil {
		return 0, err
	}
	length := uint(lengthBytes[0])
	length |= uint(lengthBytes[1]) << 8
	length |= uint(lengthBytes[2]) << 16
	length |= uint(lengthBytes[3]) << 24
	return length, nil
}

// readPayload reads the payload from s.Connection and returns it as []byte.
func (c *Channel) readPayload(length uint) ([]byte, error) {
	payload := make([]byte, length)
	n, err := c.trans.Read(payload)
	if err != nil {
		return nil, err
	}
	if n != int(length) {
		return nil, errors.New("read too few bytes")
	}
	return payload, nil
}
