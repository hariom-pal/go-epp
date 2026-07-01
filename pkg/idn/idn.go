package idn

import (
	"strings"

	"golang.org/x/net/idna"
)

var profile = idna.New(
	idna.ValidateForRegistration(),
	idna.StrictDomainName(true),
	idna.Transitional(false),
)

// ToASCII converts a Unicode domain name to its IDNA ASCII form.
func ToASCII(domain string) (string, error) {
	domain = strings.TrimSpace(domain)
	domain = strings.TrimSuffix(domain, ".")

	return profile.ToASCII(domain)
}

// ToUnicode converts an IDNA ASCII domain name to Unicode.
func ToUnicode(domain string) (string, error) {
	domain = strings.TrimSpace(domain)
	domain = strings.TrimSuffix(domain, ".")

	return profile.ToUnicode(domain)
}
