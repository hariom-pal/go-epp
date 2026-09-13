package epp

import (
	"context"
)

// Reconnect opens a fresh TLS session using the existing configuration.
// It does not replay any previous EPP command.
func (c *Client) Reconnect() (*Client, error) {
	return c.ReconnectContext(context.Background())
}

// ReconnectContext opens a fresh TLS session using the existing configuration.
// Transform commands with lost responses must be reconciled by the caller using
// registry state such as info/check/poll and transaction records.
func (c *Client) ReconnectContext(ctx context.Context) (*Client, error) {
	if c == nil || c.config == nil {
		return nil, newSDKError(ErrorKindSession, "EPP session has no reusable configuration", nil)
	}

	c.stateMu.RLock()
	logger := c.logger
	c.stateMu.RUnlock()

	next, err := ConnectContext(ctx, c.config)
	if err != nil {
		return nil, err
	}
	next.SetLogger(logger)
	return next, nil
}
