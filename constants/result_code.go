package constants

const (
	// ResultSuccess indicates that the EPP command completed successfully.
	ResultSuccess = 1000

	// ResultSuccessPending indicates that the command completed but action is pending.
	ResultSuccessPending = 1001

	// ResultNoMessages indicates that the command completed and no messages are queued.
	ResultNoMessages = 1300

	// ResultAckToDequeue indicates that the command completed and a message is available.
	ResultAckToDequeue = 1301

	// ResultEndingSession indicates that the command completed successfully and the
	// server is ending the session. RFC 5730 requires this code in response to logout.
	ResultEndingSession = 1500

	// ResultUnknownCommand indicates that the server does not recognize the command.
	ResultUnknownCommand = 2000

	// ResultSyntaxError indicates command XML syntax is invalid.
	ResultSyntaxError = 2001

	// ResultUseError indicates the command was used incorrectly.
	ResultUseError = 2002

	// ResultParameterError indicates one or more command parameters are invalid.
	ResultParameterError = 2005

	// ResultAuthenticationError indicates authentication failed.
	ResultAuthenticationError = 2501

	// ResultAuthorizationError indicates the client is not authorized for the command.
	ResultAuthorizationError = 2502

	// ResultObjectExists indicates the requested object already exists.
	ResultObjectExists = 2302

	// ResultObjectDoesNotExist indicates the requested object does not exist.
	ResultObjectDoesNotExist = 2303

	// ResultObjectStatusProhibits indicates the object's status prohibits the command.
	ResultObjectStatusProhibits = 2304
)

// IsSuccessResultCode reports whether code is an EPP success result code.
//
// RFC 5730 section 2.6 classifies result codes by their first digit: every code
// in the 1xxx range is a positive completion reply. Enumerating only the codes
// the SDK happens to know would reject valid successes such as 1500 ("command
// completed successfully; ending session"), which RFC 5730 mandates as the
// response to logout, so the whole range is accepted.
func IsSuccessResultCode(code int) bool {
	return code >= 1000 && code < 2000
}
