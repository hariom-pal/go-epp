package unit

import (
	"strings"
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/extensions/launch"
	"github.com/hariom-pal/go-epp/extensions/secdns"
	"github.com/hariom-pal/go-epp/types"
	"github.com/hariom-pal/go-epp/validation"
)

func TestSessionValidation(t *testing.T) {
	if err := validation.ValidateLogin("registrar", "secret"); err != nil {
		t.Fatalf("valid login rejected: %v", err)
	}
	if err := validation.ValidateLogin("", "secret"); err == nil {
		t.Fatal("expected empty login username to be rejected")
	}
	if err := validation.ValidatePoll(types.PollRequest{Operation: constants.PollAcknowledge}); err == nil {
		t.Fatal("expected poll ack without message ID to be rejected")
	}
	if err := validation.ValidateGreeting(types.Greeting{
		ServerID: "registry",
		Versions: []string{"1.0"},
		Languages: []string{
			"en",
		},
	}); err != nil {
		t.Fatalf("valid greeting rejected: %v", err)
	}
}

func TestDomainValidation(t *testing.T) {
	validCreate := types.DomainCreateRequest{
		Domain:     "example.com",
		Period:     1,
		Unit:       "y",
		Registrant: "CONTACT-1",
		AuthInfo:   "secret",
		NameServers: []string{
			"ns1.example.com",
		},
	}
	if err := validation.ValidateDomainCreate(validCreate); err != nil {
		t.Fatalf("valid domain create rejected: %v", err)
	}

	invalidRenew := types.DomainRenewRequest{
		Domain:            "example.com",
		CurrentExpiryDate: time.Now(),
		Period:            100,
		Unit:              "y",
	}
	if err := validation.ValidateDomainRenew(invalidRenew); err == nil {
		t.Fatal("expected oversized renew period to be rejected")
	}

	if err := validation.ValidateDomainTransferRequest(types.DomainTransferRequest{
		Domain: "example.com",
	}); err == nil || !strings.Contains(err.Error(), "authInfo") {
		t.Fatalf("expected transfer request authInfo error, got %v", err)
	}
}

func TestContactValidation(t *testing.T) {
	req := types.ContactCreateRequest{
		ContactID: "CONTACT-1",
		InternationalPostalInfo: &types.PostalInfo{
			Name:        "Example Registrant",
			City:        "Mumbai",
			CountryCode: "IN",
		},
		Voice:    types.Phone{Number: "+91.9876543210"},
		Email:    "admin@example.com",
		AuthInfo: "secret",
	}
	if err := validation.ValidateContactCreate(req); err != nil {
		t.Fatalf("valid contact create rejected: %v", err)
	}
	if err := validation.ValidateContactUpdate(types.ContactUpdateRequest{
		ContactID: "CONTACT-1",
	}); err == nil {
		t.Fatal("expected empty contact update to be rejected")
	}
}

func TestHostValidation(t *testing.T) {
	if err := validation.ValidateHostCreate(types.HostCreateRequest{
		HostName: "ns1.example.com",
		Addresses: []types.HostAddress{
			{IPVersion: "v4", Address: "192.0.2.10"},
		},
	}); err != nil {
		t.Fatalf("valid host create rejected: %v", err)
	}
	if err := validation.ValidateHostUpdate(types.HostUpdateRequest{
		HostName:    "ns1.example.com",
		AddStatuses: []string{"not-a-status"},
	}); err == nil {
		t.Fatal("expected invalid host status to be rejected")
	}
}

func TestExtensionValidation(t *testing.T) {
	if err := validation.ValidateDNSSECCreate(&secdns.CreateRequest{
		Data: secdns.Data{
			DSData: []secdns.DSData{{KeyTag: 1}},
			KeyData: []secdns.KeyData{{
				Flags: 257,
			}},
		},
	}); err == nil {
		t.Fatal("expected mixed secDNS create interfaces to be rejected")
	}

	if err := validation.ValidateLaunchInfo(&launch.InfoRequest{
		Phase: launch.Phase{Value: launch.PhaseSunrise},
	}); err != nil {
		t.Fatalf("valid launch info rejected: %v", err)
	}
}
