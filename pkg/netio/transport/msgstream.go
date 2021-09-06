// Implements message framing on a streaming transport (e.g. sockets)

package transport

import (
	"errors"
)

// MsgStream represents a stream with message semantics
type MsgStream struct {
	trans Transport
}

// WriteMsg writes Netio standard message to the transport stream
func (ms *MsgStream) WriteMsg(payload []byte) error {
	// make sure we have the magic bytes
	err := ms.writeMagicBytes()
	if err != nil {
		return err
	}

	length := uint(len(payload))
	err = ms.writeLength(length)
	if err != nil {
		return err
	}

	err = ms.writePayload(payload)
	if err != nil {
		return err
	}
	return nil
}

// writeMagicBytes write the magic bytes 0xFE, 0xED to transport stream
func (ms *MsgStream) writeMagicBytes() error {
	magicBytes := []byte{0xFE, 0xED}

	err := ms.writeBytesSafe(magicBytes)
	return err
}

// writeLength writes 4 bytes to transport stream with the length
func (ms *MsgStream) writeLength(length uint) error {
	lengthBytes := make([]byte, 4)

	lengthBytes[0] = byte(length & 0xFF)
	lengthBytes[1] = byte((length >> 8) & 0xFF)
	lengthBytes[2] = byte((length >> 16) & 0xFF)
	lengthBytes[3] = byte((length >> 24) & 0xFF)

	err := ms.writeBytesSafe(lengthBytes)
	return err
}

// writePayload write the payload to transport stream.
func (ms *MsgStream) writePayload(payload []byte) error {
	err := ms.writeBytesSafe(payload)
	return err
}

// writeBytesSafe retries writing to transport stream until all bytes are written
func (ms *MsgStream) writeBytesSafe(payload []byte) error {
	for {
		written, err := ms.trans.Write(payload)
		if err != nil {
			return err
		}
		if written == len(payload) {
			return nil
		}
		payload = payload[written:]
	}
}

// ReadMsg reads a Netio standard message from transport stream
func (ms *MsgStream) ReadMsg() ([]byte, error) {
	// make sure we have the magic bytes
	err := ms.readMagicBytes()
	if err != nil {
		return nil, err
	}

	length, err := ms.readLength()
	if err != nil {
		return nil, err
	}
	payload, err := ms.readPayload(length)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// readMagicBytes blocks until it receives the magic bytes 0xFE, 0xED from transport stream.
func (ms *MsgStream) readMagicBytes() error {
	// block until we get the magic bytes
	for {
		magicBytes := make([]byte, 2)
		for i := 0; i < 2; i++ {
			b := make([]byte, 1)

			_, err := ms.trans.Read(b)
			if err != nil {
				return err
			}
			magicBytes[i] = b[0]
		}
		if magicBytes[0] == 0xFE && magicBytes[1] == 0xED {
			return nil
		}
	}
}

// readLength reads 4 bytes from transport stream and returns the length as uint of the message.
func (ms *MsgStream) readLength() (uint, error) {
	lengthBytes := make([]byte, 4)
	_, err := ms.trans.Read(lengthBytes)
	if err != nil {
		return 0, err
	}
	length := uint(lengthBytes[0])
	length |= uint(lengthBytes[1]) << 8
	length |= uint(lengthBytes[2]) << 16
	length |= uint(lengthBytes[3]) << 24
	return length, nil
}

// readPayload reads the payload from transport stream and returns it as []byte.
func (ms *MsgStream) readPayload(length uint) ([]byte, error) {
	payload := make([]byte, length)
	n, err := ms.trans.Read(payload)
	if err != nil {
		return nil, err
	}
	if n != int(length) {
		return nil, errors.New("read too few bytes")
	}
	return payload, nil
}

// Close closes the transport stream
func (ms *MsgStream) Close() error {
	return ms.trans.Close()
}
