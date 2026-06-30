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

func ToASCII(domain string) (string, error) {
	domain = strings.TrimSpace(domain)
	domain = strings.TrimSuffix(domain, ".")

	return profile.ToASCII(domain)
}

func ToUnicode(domain string) (string, error) {
	domain = strings.TrimSpace(domain)
	domain = strings.TrimSuffix(domain, ".")

	return profile.ToUnicode(domain)
}
