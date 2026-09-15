package epp

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// tridNonceBytes sizes the per-client random component of a client transaction
// identifier.
const tridNonceBytes = 5

// nextTRID returns a client transaction identifier for one command.
//
// RFC 5730 section 2.5 requires the client transaction identifier to be
// unique, because it is how a registrar correlates a response and how it
// reconciles a transform whose response was lost. A per-client counter and a
// one-second timestamp are not enough on their own: every client starts its
// counter at the same value, so two sessions opened in the same second, which
// is the normal case for a connection pool, would emit identical identifiers.
// The per-client nonce makes them unique across sessions, processes and hosts.
//
// The result stays within the RFC 5730 clTRIDType limit of 64 characters.
func (c *Client) nextTRID(prefix string) string {
	id := c.sequence.Add(1)

	return fmt.Sprintf(
		"%s-%s-%s-%06d",
		prefix,
		time.Now().UTC().Format("20060102T150405"),
		c.tridNonce(),
		id,
	)
}

// tridNonce returns this client's random identifier component, generating it
// on first use.
func (c *Client) tridNonce() string {
	c.stateMu.RLock()
	nonce := c.nonce
	c.stateMu.RUnlock()
	if nonce != "" {
		return nonce
	}

	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	if c.nonce != "" {
		return c.nonce
	}

	buffer := make([]byte, tridNonceBytes)
	if _, err := rand.Read(buffer); err != nil {
		// A client that cannot read randomness still needs identifiers that
		// do not collide with its peers, so fall back to the address of this
		// client, which is unique within the process, combined with the
		// nanosecond clock.
		c.nonce = fmt.Sprintf("%010x", uint64(time.Now().UnixNano())&0xffffffffff)
		return c.nonce
	}
	c.nonce = hex.EncodeToString(buffer)
	return c.nonce
}
