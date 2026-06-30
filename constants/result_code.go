package constants

const (
	// Success
	ResultSuccess        = 1000
	ResultSuccessPending = 1001

	// Command Errors
	ResultUnknownCommand = 2000
	ResultSyntaxError    = 2001
	ResultUseError       = 2002
	ResultParameterError = 2005

	// Authorization
	ResultAuthenticationError = 2501
	ResultAuthorizationError  = 2502

	// Object Errors
	ResultObjectExists          = 2302
	ResultObjectDoesNotExist    = 2303
	ResultObjectStatusProhibits = 2304
)
