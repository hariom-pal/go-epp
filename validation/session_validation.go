package validation

import (
	"fmt"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

// ValidateGreeting validates parsed RFC5730 greeting data.
func ValidateGreeting(greeting types.Greeting) error {
	if strings.TrimSpace(greeting.ServerID) == "" {
		return required("greeting server ID")
	}
	if len(nonEmptyStrings(greeting.Versions)) == 0 {
		return required("greeting service version")
	}
	if len(nonEmptyStrings(greeting.Languages)) == 0 {
		return required("greeting language")
	}
	return nil
}

// ValidateHello validates a hello command.
func ValidateHello() error {
	return nil
}

// ValidateLogin validates login credentials before building an EPP login command.
func ValidateLogin(username string, password string) error {
	if strings.TrimSpace(username) == "" {
		return required("login username")
	}
	if strings.TrimSpace(password) == "" {
		return required("login password")
	}
	return nil
}

// ValidateLogout validates a logout command.
func ValidateLogout() error {
	return nil
}

// ValidatePollRequest validates a poll request command.
func ValidatePollRequest() error {
	return nil
}

// ValidatePollAck validates a poll acknowledgement command.
func ValidatePollAck(messageID string) error {
	if strings.TrimSpace(messageID) == "" {
		return required("poll message ID")
	}
	return nil
}

// ValidatePoll validates an RFC5730 poll request.
func ValidatePoll(req types.PollRequest) error {
	operation := strings.ToLower(strings.TrimSpace(req.Operation))
	switch operation {
	case constants.PollRequest, "req":
		return ValidatePollRequest()
	case constants.PollAcknowledge:
		return ValidatePollAck(req.MessageID)
	default:
		return fmt.Errorf("poll operation must be %q or %q", constants.PollRequest, constants.PollAcknowledge)
	}
}

func required(field string) error {
	return fmt.Errorf("%s is required", field)
}

func nonEmptyStrings(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}
