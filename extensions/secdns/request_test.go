package secdns

import (
	"strings"
	"testing"
)

const validDigest = "49FD46E6C4B45C55D4AC69CBD3CD34AC1AFE51DE0EE1E5A2B3C4D5E6F708192A"

// TestIncompleteDSDataIsRejectedNotDropped covers the NIXI OT&E finding that an
// incomplete DS record was silently removed from the request instead of
// rejected. RFC 5910 section 4.1 makes keyTag, alg, digestType and digest
// mandatory, and dropping the record left a caller believing a domain had been
// signed while the command carried no DNSSEC data at all.
func TestIncompleteDSDataIsRejectedNotDropped(t *testing.T) {
	cases := []struct {
		name string
		ds   DSData
	}{
		{"empty digest", DSData{KeyTag: 12345, Algorithm: 8, DigestType: 2}},
		{"odd length digest", DSData{KeyTag: 12345, Algorithm: 8, DigestType: 2, Digest: "ABC"}},
		{"non hex digest", DSData{KeyTag: 12345, Algorithm: 8, DigestType: 2, Digest: "ZZZZ"}},
		{"zero algorithm", DSData{KeyTag: 12345, DigestType: 2, Digest: validDigest}},
		{"zero digest type", DSData{KeyTag: 12345, Algorithm: 8, Digest: validDigest}},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			create := &CreateRequest{Data: Data{DSData: []DSData{testCase.ds}}}
			if ValidCreate(create) {
				t.Fatal("an incomplete DS record must be rejected")
			}

			update := &UpdateRequest{Add: &UpdateAdd{Data: Data{DSData: []DSData{testCase.ds}}}}
			if ValidUpdate(update) {
				t.Fatal("an incomplete DS record must be rejected on update")
			}
		})
	}
}

// TestCompleteDSDataIsSerialised confirms a valid record still reaches the wire
// intact, so the stricter validation did not start dropping good data.
func TestCompleteDSDataIsSerialised(t *testing.T) {
	request := &CreateRequest{Data: Data{DSData: []DSData{{
		KeyTag: 12345, Algorithm: 8, DigestType: 2, Digest: validDigest,
	}}}}

	if !ValidCreate(request) {
		t.Fatal("a complete DS record must be accepted")
	}

	created := NewCreate(request)
	if created == nil {
		t.Fatal("a complete DS record produced no extension XML")
	}
	if len(created.DSData) != 1 {
		t.Fatalf("expected 1 serialised DS record, got %d", len(created.DSData))
	}
	if !strings.EqualFold(created.DSData[0].Digest, validDigest) {
		t.Fatalf("digest did not survive serialisation: %q", created.DSData[0].Digest)
	}
}

// TestIncompleteKeyDataIsRejected covers the key data interface, where an
// empty public key was dropped the same way.
func TestIncompleteKeyDataIsRejected(t *testing.T) {
	cases := []struct {
		name string
		key  KeyData
	}{
		{"empty public key", KeyData{Flags: 257, Protocol: 3, Algorithm: 8}},
		{"zero protocol", KeyData{Flags: 257, Algorithm: 8, PublicKey: "AwEAAb=="}},
		{"zero algorithm", KeyData{Flags: 257, Protocol: 3, PublicKey: "AwEAAb=="}},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if ValidCreate(&CreateRequest{Data: Data{KeyData: []KeyData{testCase.key}}}) {
				t.Fatal("an incomplete key data record must be rejected")
			}
		})
	}
}

// TestMixedInterfacesStillRejected confirms the RFC 5910 section 4 rule that a
// request uses either the DS data interface or the key data interface, never
// both, survives the stricter validation.
func TestMixedInterfacesStillRejected(t *testing.T) {
	request := &CreateRequest{Data: Data{
		DSData:  []DSData{{KeyTag: 1, Algorithm: 8, DigestType: 2, Digest: validDigest}},
		KeyData: []KeyData{{Flags: 257, Protocol: 3, Algorithm: 8, PublicKey: "AwEAAb=="}},
	}}

	if ValidCreate(request) {
		t.Fatal("mixing the DS and key data interfaces must be rejected")
	}
}
