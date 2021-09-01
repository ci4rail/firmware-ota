package socket

import (
	"net"
)

type Socket struct {
	Connection net.Conn
	Buffer     []byte
}

func WaitForConnect(port string) (*Socket, error) {

	l, err := net.Listen("tcp4", port)
	if err != nil {
		return &Socket{}, err
	}

	conn, err := l.Accept()
	if err != nil {
		return &Socket{}, err
	}

	return &Socket{
		Connection: conn,
		Buffer:     make([]byte, 2000),
	}, nil

}

func (s *Socket) Read() ([]byte, error) {
	// make sure we have the magic bytes
	err := s.readMagicBytes()
	if err != nil {
		return nil, err
	}

	length, err := s.readLength()
	if err != nil {
		return nil, err
	}
	// log.Println("Length: ", length)
	payload, err := s.readPayload(length)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// readMagicBytes blocks until it receives the magic bytes 0xFE, 0xED from s.Connection.
func (s *Socket) readMagicBytes() error {
	// block until we get the magic bytes
	for {
		magicBytes := make([]byte, 2)
		for i := 0; i < 2; i++ {
			b := make([]byte, 1)

			_, err := s.Connection.Read(b)
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
func (s *Socket) readLength() (uint, error) {
	lengthBytes := make([]byte, 4)
	_, err := s.Connection.Read(lengthBytes)
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
func (s *Socket) readPayload(length uint) ([]byte, error) {
	payload := make([]byte, length)
	_, err := s.Connection.Read(payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
