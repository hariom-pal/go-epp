package constants

const (
	// ResultSuccess indicates that the EPP command completed successfully.
	ResultSuccess = 1000

	// ResultSuccessPending indicates that the command completed but action is pending.
	ResultSuccessPending = 1001

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
