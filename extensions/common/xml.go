package common

// BoolString returns the XML Schema lexical value for a boolean.
func BoolString(value bool) string {
	if value {
		return "true"
	}

	return "false"
}
