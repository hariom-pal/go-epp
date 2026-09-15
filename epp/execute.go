package epp

import (
	"context"
	"encoding/xml"
	"errors"
	"strings"
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
	// Abort a blocked read when the caller's context ends. The watcher must
	// not outlive this command: once the read returns, close(done) and
	// ctx.Done() can both be ready, and a select picks between ready cases at
	// random. A watcher that took the ctx branch after this command returned
	// would push the shared connection's read deadline into the past and
	// abort whichever command ran next on this session, which on a transform
	// would surface as a spurious ambiguous outcome. Joining it here keeps any
	// deadline change confined to the command that owns the context.
	done := make(chan struct{})
	watcherDone := make(chan struct{})
	go func() {
		defer close(watcherDone)
		select {
		case <-ctx.Done():
			_ = c.conn.SetReadDeadline(time.Now())
		case <-done:
		}
	}()
	response, err := ReadFrameWithMax(c.conn, maxFrameSize(c.config))
	close(done)
	<-watcherDone
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil && (errors.Is(ctxErr, context.Canceled) || errors.Is(ctxErr, context.DeadlineExceeded)) {
			err = contextError(ctxErr)
		}
		if transform {
			err = &AmbiguousTransformError{
				Command:    command,
				ClientTRID: requestClientTRID(xml),
				Err:        err,
			}
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

// requestClientTRID extracts the clTRID from an outgoing command so an
// ambiguous transform can report the identifier the caller needs in order to
// reconcile it against registry state.
func requestClientTRID(requestXML []byte) string {
	const open = "<clTRID>"
	const close = "</clTRID>"

	value := string(requestXML)
	start := strings.Index(value, open)
	if start < 0 {
		return ""
	}
	value = value[start+len(open):]
	end := strings.Index(value, close)
	if end < 0 {
		return ""
	}
	return strings.TrimSpace(value[:end])
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
