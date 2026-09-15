package epp

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

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
	ErrorKindValidation         ErrorKind = "validation"
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
// the registry but no authoritative response was received. The object may or
// may not have changed, so the caller must reconcile against registry state
// before retrying.
type AmbiguousTransformError struct {
	Command string

	// ClientTRID is the transaction identifier the command was sent with.
	// Reconciliation depends on it: it is what identifies this command in the
	// registry's transaction record and in any poll message it produced, and
	// it is the only way to tell this attempt apart from a retry.
	ClientTRID string

	Err error
}

func (e *AmbiguousTransformError) Error() string {
	trid := ""
	if e.ClientTRID != "" {
		trid = fmt.Sprintf(" (ClientTRID=%s)", e.ClientTRID)
	}
	if e.Err == nil {
		return fmt.Sprintf("%s error: transform command %s has ambiguous outcome%s", ErrorKindAmbiguousTransform, e.Command, trid)
	}
	return fmt.Sprintf("%s error: transform command %s has ambiguous outcome%s: %v", ErrorKindAmbiguousTransform, e.Command, trid, e.Err)
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
	Kind    ErrorKind
	Code    int
	Message string

	// Values carries the server's <value>/<extValue> diagnostics for the
	// failing result. These name the exact command element the server
	// rejected and are included in Error().
	Values []types.ResultValue

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
		"%s error: EPP result [%d]: %s%s (ClientTRID=%s, ServerTRID=%s)",
		kind,
		e.Code,
		e.Message,
		formatResultValues(e.Values),
		e.ClientTRID,
		e.ServerTRID,
	)
}

// formatResultValues renders server diagnostics for inclusion in an error
// string. It returns an empty string when the server supplied none.
func formatResultValues(values []types.ResultValue) string {
	details := make([]string, 0, len(values))
	for _, value := range values {
		switch {
		case value.Reason != "" && value.Value != "":
			details = append(details, fmt.Sprintf("%s: %s", value.Value, value.Reason))
		case value.Reason != "":
			details = append(details, value.Reason)
		case value.Value != "":
			details = append(details, value.Value)
		}
	}
	if len(details) == 0 {
		return ""
	}
	return " [" + strings.Join(details, "; ") + "]"
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

// newValidationError builds an error for a request the SDK rejected locally,
// before anything was sent. It carries the EPP result code the registry would
// be expected to use so callers can switch on Code, but its Kind marks it as a
// client-side failure: a caller must not treat it as a registry response, and
// in particular must not reconcile against registry state because of it.
func newValidationError(code int, message string) *Error {
	return &Error{
		Kind:    ErrorKindValidation,
		Code:    code,
		Message: message,
	}
}

// IsValidationError reports whether the SDK rejected the request locally
// rather than the registry rejecting it.
func (e *Error) IsValidationError() bool {
	return e.Kind == ErrorKindValidation
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
