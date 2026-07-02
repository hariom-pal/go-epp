package epp

import "time"

// Execute sends raw EPP XML and returns the raw EPP response XML.
func (c *Client) Execute(xml []byte) ([]byte, error) {

	// Send Request
	if err := c.setWriteDeadline(); err != nil {
		return nil, err
	}
	if err := WriteFrame(c.conn, xml); err != nil {
		return nil, err
	}

	// Receive Response
	if err := c.setReadDeadline(); err != nil {
		return nil, err
	}
	response, err := ReadFrame(c.conn)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) setReadDeadline() error {
	if c.config == nil {
		return nil
	}
	return setDeadline(c.conn.SetReadDeadline, c.config.Timeout.Read)
}

func (c *Client) setWriteDeadline() error {
	if c.config == nil {
		return nil
	}
	return setDeadline(c.conn.SetWriteDeadline, c.config.Timeout.Write)
}

func setDeadline(setter func(time.Time) error, seconds int) error {
	if seconds <= 0 {
		return setter(time.Time{})
	}
	return setter(time.Now().Add(time.Duration(seconds) * time.Second))
}
