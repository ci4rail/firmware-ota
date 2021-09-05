package transport

import (
	"net"
)

// NewListener creates a Listener on a socket on a TCP socket
// port should be the port to listen to e.g. ":9999"
// pass the Listener to WaitForConnect
func NewListener(port string) (*net.TCPListener, error) {
	addr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		return nil, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// WaitForConnect waits for a client to connect to the TCP socket
// There is no timeout
func WaitForConnect(l *net.TCPListener) (*net.TCPConn, error) {

	conn, err := l.AcceptTCP()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// NewConnection connects to a TCP server at address
func NewConnection(address string) (*net.TCPConn, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
