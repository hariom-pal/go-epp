package epp

import (
	"fmt"

	"github.com/hariom-pal/go-epp/constants"
)

// Error represents a structured EPP response error.
type Error struct {
	Code       int
	Message    string
	ClientTRID string
	ServerTRID string
}

// Error returns a human-readable EPP error string.
func (e *Error) Error() string {
	return fmt.Sprintf(
		"EPP Error [%d]: %s (ClientTRID=%s, ServerTRID=%s)",
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
