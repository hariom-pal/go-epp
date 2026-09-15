package test

import (
	"testing"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/extensions/secdns"
	"github.com/hariom-pal/go-epp/types"
	"github.com/hariom-pal/go-epp/validation"
)

// TestValidationParity compares the standalone validation package against what
// the client itself enforces. The package mirrors every command, so a request
// one accepts and the other rejects means a caller that pre-validates gets a
// different answer from the call itself.
func TestValidationParity(t *testing.T) {
	// A disconnected client performs its request validation before any I/O,
	// so a validation failure is distinguishable from a session failure.
	client := &epp.Client{}

	cases := []struct {
		name      string
		validate  func() error
		send      func() error
		expectBad bool
	}{
		{
			name:      "domain create: no domain",
			validate:  func() error { return validation.ValidateDomainCreate(types.DomainCreateRequest{}) },
			send:      func() error { _, err := client.DomainCreate(types.DomainCreateRequest{}); return err },
			expectBad: true,
		},
		{
			name: "domain create: period 0",
			validate: func() error {
				return validation.ValidateDomainCreate(types.DomainCreateRequest{
					Domain: "example.in", Registrant: "reg1", AuthInfo: "secret123",
				})
			},
			send: func() error {
				_, err := client.DomainCreate(types.DomainCreateRequest{
					Domain: "example.in", Registrant: "reg1", AuthInfo: "secret123",
				})
				return err
			},
			expectBad: true,
		},
		{
			name: "domain create: period 100",
			validate: func() error {
				return validation.ValidateDomainCreate(types.DomainCreateRequest{
					Domain: "example.in", Period: 100, Unit: "y", Registrant: "reg1", AuthInfo: "secret123",
				})
			},
			send: func() error {
				_, err := client.DomainCreate(types.DomainCreateRequest{
					Domain: "example.in", Period: 100, Unit: "y", Registrant: "reg1", AuthInfo: "secret123",
				})
				return err
			},
			expectBad: true,
		},
		{
			name: "domain create: bad period unit",
			validate: func() error {
				return validation.ValidateDomainCreate(types.DomainCreateRequest{
					Domain: "example.in", Period: 1, Unit: "decades", Registrant: "reg1", AuthInfo: "secret123",
				})
			},
			send: func() error {
				_, err := client.DomainCreate(types.DomainCreateRequest{
					Domain: "example.in", Period: 1, Unit: "decades", Registrant: "reg1", AuthInfo: "secret123",
				})
				return err
			},
			expectBad: true,
		},
		{
			name: "domain create: no authInfo",
			validate: func() error {
				return validation.ValidateDomainCreate(types.DomainCreateRequest{
					Domain: "example.in", Period: 1, Unit: "y", Registrant: "reg1",
				})
			},
			send: func() error {
				_, err := client.DomainCreate(types.DomainCreateRequest{
					Domain: "example.in", Period: 1, Unit: "y", Registrant: "reg1",
				})
				return err
			},
			expectBad: true,
		},
		{
			name: "domain create: invalid secDNS digest",
			validate: func() error {
				return validation.ValidateDomainCreate(types.DomainCreateRequest{
					Domain: "example.in", Period: 1, Unit: "y", Registrant: "reg1", AuthInfo: "secret123",
					SecDNS: &secdns.CreateRequest{Data: secdns.Data{DSData: []secdns.DSData{{
						KeyTag: 1, Algorithm: 8, DigestType: 2, Digest: "",
					}}}},
				})
			},
			send: func() error {
				_, err := client.DomainCreate(types.DomainCreateRequest{
					Domain: "example.in", Period: 1, Unit: "y", Registrant: "reg1", AuthInfo: "secret123",
					SecDNS: &secdns.CreateRequest{Data: secdns.Data{DSData: []secdns.DSData{{
						KeyTag: 1, Algorithm: 8, DigestType: 2, Digest: "",
					}}}},
				})
				return err
			},
			expectBad: true,
		},
		{
			name:      "host create: no name",
			validate:  func() error { return validation.ValidateHostCreate(types.HostCreateRequest{}) },
			send:      func() error { _, err := client.HostCreate(types.HostCreateRequest{}); return err },
			expectBad: true,
		},
		{
			name: "host create: bad IP",
			validate: func() error {
				return validation.ValidateHostCreate(types.HostCreateRequest{
					HostName:  "ns1.example.in",
					Addresses: []types.HostAddress{{Address: "999.999.999.999", IPVersion: "v4"}},
				})
			},
			send: func() error {
				_, err := client.HostCreate(types.HostCreateRequest{
					HostName:  "ns1.example.in",
					Addresses: []types.HostAddress{{Address: "999.999.999.999", IPVersion: "v4"}},
				})
				return err
			},
			expectBad: true,
		},
		{
			name:      "contact create: no ID",
			validate:  func() error { return validation.ValidateContactCreate(types.ContactCreateRequest{}) },
			send:      func() error { _, err := client.ContactCreate(types.ContactCreateRequest{}); return err },
			expectBad: true,
		},
		{
			name:      "domain check: no names",
			validate:  func() error { return validation.ValidateDomainCheck(types.DomainCheckRequest{}) },
			send:      func() error { _, err := client.DomainCheck(types.DomainCheckRequest{}); return err },
			expectBad: true,
		},
		{
			name: "domain transfer: no operation",
			validate: func() error {
				return validation.ValidateDomainTransfer(types.DomainTransferRequest{DomainName: "example.in"})
			},
			send: func() error {
				_, err := client.DomainTransfer(types.DomainTransferRequest{DomainName: "example.in"})
				return err
			},
			expectBad: true,
		},
		{
			name: "poll ack: no message ID",
			validate: func() error {
				return validation.ValidatePoll(types.PollRequest{Operation: "ack"})
			},
			send: func() error {
				_, err := client.Poll(types.PollRequest{Operation: "ack"})
				return err
			},
			expectBad: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			validationErr := testCase.validate()
			sendErr := testCase.send()

			validationRejected := validationErr != nil
			// A disconnected client reports a session error only after its
			// request validation passed, so that is the "accepted" signal.
			sendRejected := sendErr != nil && !isSessionError(sendErr)

			if validationRejected != sendRejected {
				t.Fatalf("validation package and client disagree: validation rejected=%v (%v), client rejected=%v (%v)",
					validationRejected, validationErr, sendRejected, sendErr)
			}
			if testCase.expectBad && !validationRejected {
				t.Fatalf("expected this request to be rejected, but both accepted it")
			}
		})
	}
}

func isSessionError(err error) bool {
	var sdkErr *epp.SDKError
	if !asSDK(err, &sdkErr) {
		return false
	}
	return sdkErr.Kind == epp.ErrorKindSession
}
