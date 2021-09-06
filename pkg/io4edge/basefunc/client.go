package basefunc

import (
	"github.com/ci4rail/firmware-ota/pkg/io4edge"
)

// Client represents a client for the io4edge base function
type Client struct {
	ch *io4edge.Channel
}

// NewClient creates a new client for the base function
func NewClient(c *io4edge.Channel) (*Client, error) {
	return &Client{ch: c}, nil
}
