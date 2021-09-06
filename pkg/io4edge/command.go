package io4edge

import (
	"time"

	"google.golang.org/protobuf/proto"
)

// Command issues a command cmd to a channel, waits for the devices response and returns it in res
func (c *Channel) Command(cmd proto.Message, res proto.Message, timeout time.Duration) error {
	err := c.WriteMessage(cmd)
	if err != nil {
		return err
	}
	return c.ReadMessage(res, timeout)
}
