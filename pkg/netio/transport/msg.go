package transport

// Msg is the interface used by Channel to exchange message frames with the transport layer
// e.g. socket, websocket...
type Msg interface {
	ReadMsg() (payload []byte, err error)
	WriteMsg(payload []byte) (err error)
	Close() error
}
