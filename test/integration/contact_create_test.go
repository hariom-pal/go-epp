package test

import (
	"strings"
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func TestContactCreateXMLGenerationAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, contactCreateResponse("CNT001"))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.ContactCreate(types.ContactCreateRequest{
		ContactID: "CNT001",
		InternationalPostalInfo: &types.PostalInfo{
			Type:          "int",
			Name:          "John Doe",
			Organization:  "Example Inc.",
			Street:        []string{"123 Example Street", "Suite 100"},
			City:          "Dulles",
			StateProvince: "VA",
			PostalCode:    "20166",
			CountryCode:   "US",
		},
		LocalizedPostalInfo: &types.PostalInfo{
			Type:          "loc",
			Name:          "Local Name",
			Organization:  "Local Org",
			Street:        []string{"Local Street"},
			City:          "Dulles",
			StateProvince: "VA",
			PostalCode:    "20166",
			CountryCode:   "US",
		},
		Voice: types.Phone{
			Number:    "+1.7035555555",
			Extension: "1234",
		},
		Fax: types.Phone{
			Number:    "+1.7035555556",
			Extension: "5678",
		},
		Email:    "john@example.test",
		AuthInfo: "contactSecret",
		Disclosure: &types.ContactDisclosure{
			Flag: true,
			Name: types.ContactDisclosurePostal{
				International: true,
				Localized:     true,
			},
			Organization: types.ContactDisclosurePostal{
				International: true,
			},
			Address: types.ContactDisclosurePostal{
				Localized: true,
			},
			Voice: true,
			Fax:   true,
			Email: true,
		},
	})
	if err != nil {
		t.Fatalf("contact create failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `xmlns:contact="urn:ietf:params:xml:ns:contact-1.0"`)
	assertContains(t, requestXML, `<contact:create>`)
	assertContains(t, requestXML, `<contact:id>CNT001</contact:id>`)
	assertContains(t, requestXML, `<contact:postalInfo type="int">`)
	assertContains(t, requestXML, `<contact:postalInfo type="loc">`)
	assertContains(t, requestXML, `<contact:name>John Doe</contact:name>`)
	assertContains(t, requestXML, `<contact:org>Example Inc.</contact:org>`)
	assertContains(t, requestXML, `<contact:street>123 Example Street</contact:street>`)
	assertContains(t, requestXML, `<contact:street>Suite 100</contact:street>`)
	assertContains(t, requestXML, `<contact:street>Local Street</contact:street>`)
	assertContains(t, requestXML, `<contact:city>Dulles</contact:city>`)
	assertContains(t, requestXML, `<contact:sp>VA</contact:sp>`)
	assertContains(t, requestXML, `<contact:pc>20166</contact:pc>`)
	assertContains(t, requestXML, `<contact:cc>US</contact:cc>`)
	assertContains(t, requestXML, `<contact:voice x="1234">+1.7035555555</contact:voice>`)
	assertContains(t, requestXML, `<contact:fax x="5678">+1.7035555556</contact:fax>`)
	assertContains(t, requestXML, `<contact:email>john@example.test</contact:email>`)
	assertContains(t, requestXML, `<contact:authInfo>`)
	assertContains(t, requestXML, `<contact:pw>contactSecret</contact:pw>`)
	assertContains(t, requestXML, `<contact:disclose flag="1">`)
	assertContains(t, requestXML, `<contact:name type="int"></contact:name>`)
	assertContains(t, requestXML, `<contact:name type="loc"></contact:name>`)
	assertContains(t, requestXML, `<contact:org type="int"></contact:org>`)
	assertContains(t, requestXML, `<contact:addr type="loc"></contact:addr>`)
	assertContains(t, requestXML, `<contact:voice></contact:voice>`)
	assertContains(t, requestXML, `<contact:fax></contact:fax>`)
	assertContains(t, requestXML, `<contact:email></contact:email>`)
	assertContains(t, requestXML, `<clTRID>CREATE-`)

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ClientTRID != "CONTACT-CREATE-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-CONTACT-CREATE" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}

	if resp.Result.ContactID != "CNT001" {
		t.Fatalf("unexpected contact ID: %s", resp.Result.ContactID)
	}

	expectedCreated := time.Date(2026, 7, 1, 10, 30, 0, 0, time.UTC)
	if !resp.Result.CreatedDate.Equal(expectedCreated) {
		t.Fatalf("unexpected created date: %s", resp.Result.CreatedDate.Format(time.RFC3339))
	}
}

func TestContactCreateMissingOptionalFields(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, contactCreateResponse("CNT002"))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	_, err = client.ContactCreate(types.ContactCreateRequest{
		ContactID: "CNT002",
		InternationalPostalInfo: &types.PostalInfo{
			Type:        "int",
			Name:        "Jane Doe",
			Street:      []string{"1 Main Street"},
			City:        "Dulles",
			CountryCode: "US",
		},
		Voice: types.Phone{
			Number: "+1.7035550000",
		},
		Email:    "jane@example.test",
		AuthInfo: "contactSecret",
	})
	if err != nil {
		t.Fatalf("contact create failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `<contact:postalInfo type="int">`)
	assertContains(t, requestXML, `<contact:name>Jane Doe</contact:name>`)
	assertContains(t, requestXML, `<contact:street>1 Main Street</contact:street>`)
	assertContains(t, requestXML, `<contact:city>Dulles</contact:city>`)
	assertContains(t, requestXML, `<contact:cc>US</contact:cc>`)
	assertContains(t, requestXML, `<contact:voice>+1.7035550000</contact:voice>`)
	assertNotContains(t, requestXML, `<contact:postalInfo type="loc">`)
	assertNotContains(t, requestXML, `<contact:org>`)
	assertNotContains(t, requestXML, `<contact:fax`)
	assertNotContains(t, requestXML, `<contact:disclose`)
}

func TestContactCreateEmptyRequestValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.ContactCreate(types.ContactCreateRequest{})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func assertNotContains(t *testing.T, value string, unexpected string) {
	t.Helper()

	if strings.Contains(value, unexpected) {
		t.Fatalf("expected XML not to contain %q\nXML:\n%s", unexpected, value)
	}
}

func contactCreateResponse(contactID string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <contact:creData xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
                <contact:id>` + contactID + `</contact:id>
                <contact:crDate>2026-07-01T10:30:00Z</contact:crDate>
            </contact:creData>
        </resData>
        <trID>
            <clTRID>CONTACT-CREATE-TEST</clTRID>
            <svTRID>SERVER-CONTACT-CREATE</svTRID>
        </trID>
    </response>
</epp>`
}
