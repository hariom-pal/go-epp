package epp

import (
	"net"
	"sync/atomic"

	"github.com/hariom-pal/go-epp/internal/config"
)

type Client struct {
	conn     net.Conn
	config   *config.Config
	greeting []byte
	sequence atomic.Uint64
}

func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Client) Greeting() []byte {
	return c.greeting
}

func (c *Client) Config() *config.Config {
	return c.config
}
