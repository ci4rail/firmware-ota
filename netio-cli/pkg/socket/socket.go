package socket

import (
	"net"
)

type Socket struct {
	Connection *net.TCPConn
	Buffer     []byte
}

func NewConnection(address string) (*Socket, error) {
	addr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return &Socket{}, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return &Socket{}, err
	}

	return &Socket{
		Connection: conn,
		Buffer:     make([]byte, 2000),
	}, nil

}

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
