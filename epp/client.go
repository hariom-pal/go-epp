package epp

import (
	"net"
	"sync"
	"sync/atomic"

	"github.com/hariom-pal/go-epp/internal/config"
	"github.com/hariom-pal/go-epp/types"
)

// Client is an authenticated-capable EPP TCP/TLS client.
type Client struct {
	conn         net.Conn
	config       *config.Config
	greeting     []byte
	greetingInfo *types.Greeting
	sequence     atomic.Uint64
	mu           sync.Mutex
	stateMu      sync.RWMutex
	closed       bool
	loggedIn     bool
	logger       Logger

	// nonce is this client's random component for client transaction
	// identifiers, generated lazily under stateMu.
	nonce string
}

// NewClient creates a new EPP client connection.
// It establishes a connection to the EPP server based on the provided configuration
// and reads the initial server greeting.
func NewClient(cfg *config.Config) (*Client, error) {
	return Connect(cfg)
}

// Close closes the underlying EPP connection.
func (c *Client) Close() error {
	c.stateMu.Lock()
	if c.closed {
		c.stateMu.Unlock()
		return nil
	}
	c.closed = true
	c.loggedIn = false
	c.stateMu.Unlock()

	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.emit(Event{Type: EventClose})
	return classifyTransportError(err)
}

// Greeting returns the raw server greeting received after connection.
func (c *Client) Greeting() []byte {
	return c.greeting
}

// Config returns the configuration used to create the client.
func (c *Client) Config() *config.Config {
	return c.config
}

// SetLogger installs a lightweight event sink. Events never contain raw XML or secrets.
func (c *Client) SetLogger(logger Logger) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()

	c.logger = logger
}

func (c *Client) isClosed() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()

	return c.closed
}

func (c *Client) setLoggedIn(value bool) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()

	c.loggedIn = value
}

func (c *Client) emit(event Event) {
	c.stateMu.RLock()
	logger := c.logger
	c.stateMu.RUnlock()

	if logger != nil {
		logger.EPPEvent(event)
	}
}
