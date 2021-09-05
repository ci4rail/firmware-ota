package basefunc

import (
	"github.com/ci4rail/firmware-ota/pkg/netio"
)

// Client represents a client for the netio base function
type Client struct {
	ch *netio.Channel
}

// NewClient creates a new client for the base function
func NewClient(c *netio.Channel) (*Client, error) {
	return &Client{ch: c}, nil
}
