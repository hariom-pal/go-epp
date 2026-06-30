package epp

import (
	"fmt"

	"github.com/hariom-pal/go-epp/constants"
)

type Error struct {
	Code        int
	Message     string
	ClientTRID  string
	ServerTRID  string
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"EPP Error [%d]: %s (ClientTRID=%s, ServerTRID=%s)",
		e.Code,
		e.Message,
		e.ClientTRID,
		e.ServerTRID,
	)
}

func (e *Error) IsSuccess() bool {
	return e.Code == constants.ResultSuccess ||
		e.Code == constants.ResultSuccessPending
}

func (e *Error) IsObjectExists() bool {
	return e.Code == constants.ResultObjectExists
}

func (e *Error) IsObjectNotFound() bool {
	return e.Code == constants.ResultObjectDoesNotExist
}

func (e *Error) IsAuthenticationError() bool {
	return e.Code == constants.ResultAuthenticationError
}

func (e *Error) IsAuthorizationError() bool {
	return e.Code == constants.ResultAuthorizationError
}

func (e *Error) IsObjectStatusProhibited() bool {
	return e.Code == constants.ResultObjectStatusProhibits
}
