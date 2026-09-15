// Package validation offers optional pre-flight checks for EPP requests.
//
// The client in package epp already validates every request before sending it,
// so calling this package is never required: a request the client accepts is
// one this package accepts, and vice versa. That equivalence is enforced by
// TestValidationParity in test/integration, so the two cannot drift apart.
//
// Use it when you want to reject a request earlier than the call itself, for
// example while validating an API payload or a queued job before a session is
// available. Errors here are plain errors; the same failure raised by the
// client arrives as an *epp.Error whose Kind is epp.ErrorKindValidation.
package validation
