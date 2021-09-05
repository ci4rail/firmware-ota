package client

import (
	"errors"

	"github.com/ci4rail/firmware-ota/pkg/netio"
	"github.com/ci4rail/firmware-ota/pkg/netio/basefunc"
	"github.com/ci4rail/firmware-ota/pkg/netio/transport"
)

// NewClient creates a new base function client from address. Currently sockets is used as transport
func NewClient(address string) (*basefunc.Client, error) {
	t, err := transport.NewConnection(address)
	if err != nil {
		return nil, errors.New("can't create connection: " + err.Error())
	}
	ch, err := netio.NewChannel(t)
	if err != nil {
		return nil, errors.New("can't create channel: " + err.Error())
	}
	c, err := basefunc.NewClient(ch)
	if err != nil {
		return nil, errors.New("can't create client: " + err.Error())
	}
	return c, nil
}
