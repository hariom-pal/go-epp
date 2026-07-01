package test

import (
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func TestContactInfoXMLGenerationAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, contactInfoResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.ContactInfo(types.ContactInfoRequest{
		ContactID: "CNT001",
	})
	if err != nil {
		t.Fatalf("contact info failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `xmlns:contact="urn:ietf:params:xml:ns:contact-1.0"`)
	assertContains(t, requestXML, `<contact:info>`)
	assertContains(t, requestXML, `<contact:id>CNT001</contact:id>`)
	assertContains(t, requestXML, `<clTRID>INFO-`)

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ClientTRID != "CONTACT-INFO-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-CONTACT-INFO" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}

	contact := resp.Contact
	if contact.ContactID != "CNT001" {
		t.Fatalf("unexpected contact ID: %s", contact.ContactID)
	}

	if contact.ROID != "CNT001-ROID" {
		t.Fatalf("unexpected ROID: %s", contact.ROID)
	}

	if len(contact.Statuses) != 2 ||
		contact.Statuses[0] != "ok" ||
		contact.Statuses[1] != "linked" {

		t.Fatalf("unexpected statuses: %#v", contact.Statuses)
	}

	if contact.InternationalPostalInfo == nil {
		t.Fatal("expected international postal info")
	}

	assertPostalInfo(t, contact.InternationalPostalInfo, "int", "John Doe", "Example Inc.", []string{
		"123 Example Street",
		"Suite 100",
	})

	if contact.LocalizedPostalInfo == nil {
		t.Fatal("expected localized postal info")
	}

	assertPostalInfo(t, contact.LocalizedPostalInfo, "loc", "Local Name", "Local Org", []string{
		"Local Street",
	})

	if contact.Voice.Number != "+1.7035555555" ||
		contact.Voice.Extension != "1234" {

		t.Fatalf("unexpected voice: %#v", contact.Voice)
	}

	if contact.Fax.Number != "+1.7035555556" ||
		contact.Fax.Extension != "5678" {

		t.Fatalf("unexpected fax: %#v", contact.Fax)
	}

	if contact.Email != "john@example.test" {
		t.Fatalf("unexpected email: %s", contact.Email)
	}

	if contact.ClientID != "ClientY" ||
		contact.CreatedBy != "ClientX" ||
		contact.UpdatedBy != "ClientZ" {

		t.Fatalf("unexpected client fields: clID=%s crID=%s upID=%s",
			contact.ClientID,
			contact.CreatedBy,
			contact.UpdatedBy,
		)
	}

	if !contact.CreatedDate.Equal(time.Date(2026, 6, 30, 10, 30, 0, 0, time.UTC)) {
		t.Fatalf("unexpected created date: %s", contact.CreatedDate.Format(time.RFC3339))
	}

	if !contact.UpdatedDate.Equal(time.Date(2026, 7, 1, 10, 30, 0, 0, time.UTC)) {
		t.Fatalf("unexpected updated date: %s", contact.UpdatedDate.Format(time.RFC3339))
	}

	if !contact.TransferDate.Equal(time.Date(2026, 7, 2, 10, 30, 0, 0, time.UTC)) {
		t.Fatalf("unexpected transfer date: %s", contact.TransferDate.Format(time.RFC3339))
	}

	if contact.AuthInfo != "contactSecret" {
		t.Fatalf("unexpected authInfo: %s", contact.AuthInfo)
	}
}

func TestContactInfoMissingOptionalFields(t *testing.T) {
	cfg, _, cleanup := startDomainCreateServer(t, contactInfoMinimalResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.ContactInfo(types.ContactInfoRequest{
		ContactID: "CNT002",
	})
	if err != nil {
		t.Fatalf("contact info failed: %v", err)
	}

	contact := resp.Contact

	if contact.LocalizedPostalInfo != nil {
		t.Fatalf("expected missing localized postal info, got %#v", contact.LocalizedPostalInfo)
	}

	if contact.InternationalPostalInfo == nil {
		t.Fatal("expected international postal info")
	}

	if contact.Fax.Number != "" || contact.Fax.Extension != "" {
		t.Fatalf("expected missing fax, got %#v", contact.Fax)
	}

	if !contact.UpdatedDate.IsZero() {
		t.Fatalf("expected zero updated date, got %s", contact.UpdatedDate.Format(time.RFC3339))
	}

	if !contact.TransferDate.IsZero() {
		t.Fatalf("expected zero transfer date, got %s", contact.TransferDate.Format(time.RFC3339))
	}

	if contact.AuthInfo != "" {
		t.Fatalf("expected missing authInfo, got %s", contact.AuthInfo)
	}
}

func TestContactInfoEmptyRequestValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.ContactInfo(types.ContactInfoRequest{})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func assertPostalInfo(
	t *testing.T,
	info *types.PostalInfo,
	postalType string,
	name string,
	organization string,
	streets []string,
) {
	t.Helper()

	if info.Type != postalType {
		t.Fatalf("expected postal type %s, got %s", postalType, info.Type)
	}

	if info.Name != name {
		t.Fatalf("expected postal name %s, got %s", name, info.Name)
	}

	if info.Organization != organization {
		t.Fatalf("expected postal org %s, got %s", organization, info.Organization)
	}

	if len(info.Street) != len(streets) {
		t.Fatalf("expected %d street lines, got %#v", len(streets), info.Street)
	}

	for i, street := range streets {
		if info.Street[i] != street {
			t.Fatalf("expected street line %d to be %s, got %s", i, street, info.Street[i])
		}
	}

	if info.City != "Dulles" {
		t.Fatalf("unexpected city: %s", info.City)
	}

	if info.StateProvince != "VA" {
		t.Fatalf("unexpected state: %s", info.StateProvince)
	}

	if info.PostalCode != "20166" {
		t.Fatalf("unexpected postal code: %s", info.PostalCode)
	}

	if info.CountryCode != "US" {
		t.Fatalf("unexpected country code: %s", info.CountryCode)
	}
}

func contactInfoResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <contact:infData xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
                <contact:id>CNT001</contact:id>
                <contact:roid>CNT001-ROID</contact:roid>
                <contact:status s="ok"/>
                <contact:status s="linked"/>
                <contact:postalInfo type="int">
                    <contact:name>John Doe</contact:name>
                    <contact:org>Example Inc.</contact:org>
                    <contact:addr>
                        <contact:street>123 Example Street</contact:street>
                        <contact:street>Suite 100</contact:street>
                        <contact:city>Dulles</contact:city>
                        <contact:sp>VA</contact:sp>
                        <contact:pc>20166</contact:pc>
                        <contact:cc>US</contact:cc>
                    </contact:addr>
                </contact:postalInfo>
                <contact:postalInfo type="loc">
                    <contact:name>Local Name</contact:name>
                    <contact:org>Local Org</contact:org>
                    <contact:addr>
                        <contact:street>Local Street</contact:street>
                        <contact:city>Dulles</contact:city>
                        <contact:sp>VA</contact:sp>
                        <contact:pc>20166</contact:pc>
                        <contact:cc>US</contact:cc>
                    </contact:addr>
                </contact:postalInfo>
                <contact:voice x="1234">+1.7035555555</contact:voice>
                <contact:fax x="5678">+1.7035555556</contact:fax>
                <contact:email>john@example.test</contact:email>
                <contact:clID>ClientY</contact:clID>
                <contact:crID>ClientX</contact:crID>
                <contact:crDate>2026-06-30T10:30:00Z</contact:crDate>
                <contact:upID>ClientZ</contact:upID>
                <contact:upDate>2026-07-01T10:30:00Z</contact:upDate>
                <contact:trDate>2026-07-02T10:30:00Z</contact:trDate>
                <contact:authInfo>
                    <contact:pw>contactSecret</contact:pw>
                </contact:authInfo>
            </contact:infData>
        </resData>
        <trID>
            <clTRID>CONTACT-INFO-TEST</clTRID>
            <svTRID>SERVER-CONTACT-INFO</svTRID>
        </trID>
    </response>
</epp>`
}

func contactInfoMinimalResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <contact:infData xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
                <contact:id>CNT002</contact:id>
                <contact:roid>CNT002-ROID</contact:roid>
                <contact:postalInfo type="int">
                    <contact:name>Jane Doe</contact:name>
                    <contact:addr>
                        <contact:street>1 Main Street</contact:street>
                        <contact:city>Dulles</contact:city>
                        <contact:sp>VA</contact:sp>
                        <contact:pc>20166</contact:pc>
                        <contact:cc>US</contact:cc>
                    </contact:addr>
                </contact:postalInfo>
                <contact:voice>+1.7035550000</contact:voice>
                <contact:email>jane@example.test</contact:email>
                <contact:clID>ClientY</contact:clID>
                <contact:crID>ClientX</contact:crID>
                <contact:crDate>2026-06-30T10:30:00Z</contact:crDate>
            </contact:infData>
        </resData>
        <trID>
            <clTRID>CONTACT-INFO-MINIMAL-TEST</clTRID>
            <svTRID>SERVER-CONTACT-INFO-MINIMAL</svTRID>
        </trID>
    </response>
</epp>`
}
