package epp

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

// ErrorKind classifies SDK errors for registrar applications.
type ErrorKind string

const (
	ErrorKindConfiguration      ErrorKind = "configuration"
	ErrorKindTLS                ErrorKind = "tls"
	ErrorKindTransport          ErrorKind = "transport"
	ErrorKindTimeout            ErrorKind = "timeout"
	ErrorKindCancellation       ErrorKind = "cancellation"
	ErrorKindFraming            ErrorKind = "framing"
	ErrorKindProtocol           ErrorKind = "protocol"
	ErrorKindEPPResult          ErrorKind = "epp_result"
	ErrorKindSession            ErrorKind = "session"
	ErrorKindAmbiguousTransform ErrorKind = "ambiguous_transform"
)

var (
	ErrInvalidFrameLength = errors.New("invalid EPP frame length")
	ErrFrameTooLarge      = errors.New("EPP frame length exceeds maximum")
	ErrSessionClosed      = errors.New("EPP session is closed")
)

// SDKError wraps non-EPP-result failures with a stable classification.
type SDKError struct {
	Kind    ErrorKind
	Message string
	Err     error
}

// AmbiguousTransformError reports that a transform command may have reached
// the registry but no authoritative response was received.
type AmbiguousTransformError struct {
	Command string
	Err     error
}

func (e *AmbiguousTransformError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s error: transform command %s has ambiguous outcome", ErrorKindAmbiguousTransform, e.Command)
	}
	return fmt.Sprintf("%s error: transform command %s has ambiguous outcome: %v", ErrorKindAmbiguousTransform, e.Command, e.Err)
}

func (e *AmbiguousTransformError) Unwrap() error {
	return e.Err
}

func (e *SDKError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s error: %s", e.Kind, e.Message)
	}
	return fmt.Sprintf("%s error: %s: %v", e.Kind, e.Message, e.Err)
}

func (e *SDKError) Unwrap() error {
	return e.Err
}

// Error represents a structured EPP response error.
type Error struct {
	Kind       ErrorKind
	Code       int
	Message    string
	Results    []types.Result
	ClientTRID string
	ServerTRID string
}

// Error returns a human-readable EPP error string.
func (e *Error) Error() string {
	kind := e.Kind
	if kind == "" {
		kind = ErrorKindEPPResult
	}
	return fmt.Sprintf(
		"%s error: EPP result [%d]: %s (ClientTRID=%s, ServerTRID=%s)",
		kind,
		e.Code,
		e.Message,
		e.ClientTRID,
		e.ServerTRID,
	)
}

// IsSuccess reports whether the EPP result code is a success code.
func (e *Error) IsSuccess() bool {
	return constants.IsSuccessResultCode(e.Code)
}

// IsObjectExists reports whether the error indicates an existing object.
func (e *Error) IsObjectExists() bool {
	return e.Code == constants.ResultObjectExists
}

// IsObjectNotFound reports whether the error indicates a missing object.
func (e *Error) IsObjectNotFound() bool {
	return e.Code == constants.ResultObjectDoesNotExist
}

// IsAuthenticationError reports whether authentication failed.
func (e *Error) IsAuthenticationError() bool {
	return e.Code == constants.ResultAuthenticationError
}

// IsAuthorizationError reports whether authorization failed.
func (e *Error) IsAuthorizationError() bool {
	return e.Code == constants.ResultAuthorizationError
}

// IsObjectStatusProhibited reports whether object status prohibits the command.
func (e *Error) IsObjectStatusProhibited() bool {
	return e.Code == constants.ResultObjectStatusProhibits
}

func newSDKError(kind ErrorKind, message string, err error) error {
	if err == nil {
		return &SDKError{Kind: kind, Message: message}
	}
	var sdkErr *SDKError
	if errors.As(err, &sdkErr) {
		return err
	}
	return &SDKError{Kind: kind, Message: message, Err: err}
}

func classifyTransportError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, os.ErrDeadlineExceeded) {
		return newSDKError(ErrorKindTimeout, "network deadline exceeded", err)
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return newSDKError(ErrorKindTimeout, "network timeout", err)
	}
	return newSDKError(ErrorKindTransport, "network operation failed", err)
}
