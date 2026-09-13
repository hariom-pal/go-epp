package epp

import (
	"context"
	"encoding/xml"
	"errors"
	"time"
)

// Execute sends raw EPP XML and returns the raw EPP response XML.
func (c *Client) Execute(xml []byte) ([]byte, error) {
	return c.ExecuteContext(context.Background(), xml)
}

// ExecuteContext sends raw EPP XML and returns the raw EPP response XML.
// A single Client serializes the complete write/read transaction to prevent
// request/response interleaving on one EPP session.
func (c *Client) ExecuteContext(ctx context.Context, xml []byte) ([]byte, error) {
	return c.executeCommandContext(ctx, xml, "raw", false)
}

func (c *Client) executeCommandContext(ctx context.Context, xml []byte, command string, transform bool) ([]byte, error) {
	if c == nil || c.conn == nil {
		return nil, newSDKError(ErrorKindSession, "EPP session is not connected", nil)
	}
	if c.isClosed() {
		return nil, newSDKError(ErrorKindSession, "EPP session is closed", ErrSessionClosed)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return nil, contextError(err)
	}

	start := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()

	// Send Request
	if err := c.setWriteDeadline(ctx); err != nil {
		return nil, err
	}
	if err := WriteFrame(c.conn, xml); err != nil {
		classified := classifyTransportError(err)
		c.emit(Event{Type: EventCommand, Command: command, Duration: time.Since(start), Err: classified})
		return nil, classified
	}

	// Receive Response
	if err := c.setReadDeadline(ctx); err != nil {
		return nil, err
	}
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = c.conn.SetReadDeadline(time.Now())
		case <-done:
		}
	}()
	response, err := ReadFrameWithMax(c.conn, maxFrameSize(c.config))
	close(done)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil && (errors.Is(ctxErr, context.Canceled) || errors.Is(ctxErr, context.DeadlineExceeded)) {
			err = contextError(ctxErr)
		}
		if transform {
			err = &AmbiguousTransformError{Command: command, Err: err}
		}
		c.emit(Event{Type: EventCommand, Command: command, Duration: time.Since(start), Err: err})
		return nil, err
	}
	metadata := responseMetadata(response)
	c.emit(Event{
		Type:       EventCommand,
		Command:    command,
		Duration:   time.Since(start),
		ResultCode: metadata.ResultCode,
		ClientTRID: metadata.ClientTRID,
		ServerTRID: metadata.ServerTRID,
	})

	return response, nil
}

type commandMetadata struct {
	ResultCode int
	ClientTRID string
	ServerTRID string
}

func responseMetadata(responseXML []byte) commandMetadata {
	var response struct {
		Response struct {
			Result struct {
				Code int `xml:"code,attr"`
			} `xml:"result"`
			TRID struct {
				ClientTRID string `xml:"clTRID"`
				ServerTRID string `xml:"svTRID"`
			} `xml:"trID"`
		} `xml:"response"`
	}
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return commandMetadata{}
	}
	return commandMetadata{
		ResultCode: response.Response.Result.Code,
		ClientTRID: response.Response.TRID.ClientTRID,
		ServerTRID: response.Response.TRID.ServerTRID,
	}
}

func (c *Client) setReadDeadline(ctx context.Context) error {
	if c.config == nil {
		return nil
	}
	return setContextDeadline(ctx, c.conn.SetReadDeadline, c.config.Timeout.Read)
}

func (c *Client) setWriteDeadline(ctx context.Context) error {
	if c.config == nil {
		return nil
	}
	return setContextDeadline(ctx, c.conn.SetWriteDeadline, c.config.Timeout.Write)
}
