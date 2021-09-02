package socket

import (
	"net"
)

type Listener struct {
	Listener *net.TCPListener
}

type Socket struct {
	Connection *net.TCPConn
}

// Server
func NewListener(port string) (*Listener, error) {
	addr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		return &Listener{}, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return &Listener{}, err
	}
	return &Listener{
		Listener: l,
	}, nil
}

func WaitForConnect(l *Listener) (*Socket, error) {

	conn, err := l.Listener.AcceptTCP()
	if err != nil {
		return &Socket{}, err
	}

	return &Socket{
		Connection: conn,
	}, nil

}

// Client
func NewConnection(address string) (*Socket, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return &Socket{}, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return &Socket{}, err
	}

	return &Socket{
		Connection: conn,
	}, nil

}

// Common
func (s *Socket) Write(payload []byte) error {
	// make sure we have the magic bytes
	err := s.writeMagicBytes()
	if err != nil {
		return err
	}

	length := uint(len(payload))
	err = s.writeLength(length)
	if err != nil {
		return err
	}

	err = s.writePayload(payload)
	if err != nil {
		return err
	}
	return nil
}

// writeMagicBytes write the magic bytes 0xFE, 0xED to s.Connection.
func (s *Socket) writeMagicBytes() error {
	magicBytes := []byte{0xFE, 0xED}

	err := s.writeBytesSafe(magicBytes)
	return err
}

// writeLength writes 4 bytes to s.Connection with the length
func (s *Socket) writeLength(length uint) error {
	lengthBytes := make([]byte, 4)

	lengthBytes[0] = byte(length & 0xFF)
	lengthBytes[1] = byte((length >> 8) & 0xFF)
	lengthBytes[2] = byte((length >> 16) & 0xFF)
	lengthBytes[3] = byte((length >> 24) & 0xFF)

	err := s.writeBytesSafe(lengthBytes)
	return err
}

// writePayload write the payload to s.Connection.
func (s *Socket) writePayload(payload []byte) error {
	err := s.writeBytesSafe(payload)
	return err
}

func (s *Socket) writeBytesSafe(payload []byte) error {
	for {
		written, err := s.Connection.Write(payload)
		if err != nil {
			return err
		}
		if written == len(payload) {
			return nil
		}
		payload = payload[written:]
	}
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

func (s *Socket) Close() {
	s.Connection.Close()
}
