package epp

import (
	"net"
	"sync/atomic"

	"github.com/hariom-pal/go-epp/internal/config"
)

// Client is an authenticated-capable EPP TCP/TLS client.
type Client struct {
	conn     net.Conn
	config   *config.Config
	greeting []byte
	sequence atomic.Uint64
}

// Close closes the underlying EPP connection.
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// Greeting returns the raw server greeting received after connection.
func (c *Client) Greeting() []byte {
	return c.greeting
}

// Config returns the configuration used to create the client.
func (c *Client) Config() *config.Config {
	return c.config
}
