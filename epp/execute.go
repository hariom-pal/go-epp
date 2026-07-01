package epp

// Execute sends raw EPP XML and returns the raw EPP response XML.
func (c *Client) Execute(xml []byte) ([]byte, error) {

	// Send Request
	if err := WriteFrame(c.conn, xml); err != nil {
		return nil, err
	}

	// Receive Response
	response, err := ReadFrame(c.conn)
	if err != nil {
		return nil, err
	}

	return response, nil
}
