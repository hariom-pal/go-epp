package test

import (
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func TestContactUpdateXMLGenerationAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, contactUpdateResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.ContactUpdate(types.ContactUpdateRequest{
		ContactID:      "CNT001",
		AddStatuses:    []string{"clientDeleteProhibited"},
		RemoveStatuses: []string{"clientUpdateProhibited"},
		InternationalPostalInfo: &types.PostalInfo{
			Type:          "int",
			Name:          "John Updated",
			Organization:  "Updated Inc.",
			Street:        []string{"456 Updated Street", "Floor 2"},
			City:          "Dulles",
			StateProvince: "VA",
			PostalCode:    "20166",
			CountryCode:   "US",
		},
		Voice: &types.Phone{
			Number:    "+1.7035551111",
			Extension: "321",
		},
		Fax: &types.Phone{
			Number: "+1.7035552222",
		},
		Email:    "updated@example.test",
		AuthInfo: "newSecret",
		Disclosure: &types.ContactDisclosure{
			Flag: false,
			Name: types.ContactDisclosurePostal{
				International: true,
			},
			Voice: true,
			Email: true,
		},
	})
	if err != nil {
		t.Fatalf("contact update failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `xmlns:contact="urn:ietf:params:xml:ns:contact-1.0"`)
	assertContains(t, requestXML, `<contact:update>`)
	assertContains(t, requestXML, `<contact:id>CNT001</contact:id>`)
	assertContains(t, requestXML, `<contact:add>`)
	assertContains(t, requestXML, `<contact:status s="clientDeleteProhibited"></contact:status>`)
	assertContains(t, requestXML, `<contact:rem>`)
	assertContains(t, requestXML, `<contact:status s="clientUpdateProhibited"></contact:status>`)
	assertContains(t, requestXML, `<contact:chg>`)
	assertContains(t, requestXML, `<contact:postalInfo type="int">`)
	assertContains(t, requestXML, `<contact:name>John Updated</contact:name>`)
	assertContains(t, requestXML, `<contact:org>Updated Inc.</contact:org>`)
	assertContains(t, requestXML, `<contact:street>456 Updated Street</contact:street>`)
	assertContains(t, requestXML, `<contact:street>Floor 2</contact:street>`)
	assertContains(t, requestXML, `<contact:city>Dulles</contact:city>`)
	assertContains(t, requestXML, `<contact:sp>VA</contact:sp>`)
	assertContains(t, requestXML, `<contact:pc>20166</contact:pc>`)
	assertContains(t, requestXML, `<contact:cc>US</contact:cc>`)
	assertContains(t, requestXML, `<contact:voice x="321">+1.7035551111</contact:voice>`)
	assertContains(t, requestXML, `<contact:fax>+1.7035552222</contact:fax>`)
	assertContains(t, requestXML, `<contact:email>updated@example.test</contact:email>`)
	assertContains(t, requestXML, `<contact:authInfo>`)
	assertContains(t, requestXML, `<contact:pw>newSecret</contact:pw>`)
	assertContains(t, requestXML, `<contact:disclose flag="0">`)
	assertContains(t, requestXML, `<contact:name type="int"></contact:name>`)
	assertContains(t, requestXML, `<contact:voice></contact:voice>`)
	assertContains(t, requestXML, `<contact:email></contact:email>`)
	assertContains(t, requestXML, `<clTRID>UPDATE-`)

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ClientTRID != "CONTACT-UPDATE-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-CONTACT-UPDATE" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}
}

func TestContactUpdateStatusAddRemoveOnly(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, contactUpdateResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	_, err = client.ContactUpdate(types.ContactUpdateRequest{
		ContactID:      "CNT002",
		AddStatuses:    []string{"clientTransferProhibited"},
		RemoveStatuses: []string{"clientDeleteProhibited"},
	})
	if err != nil {
		t.Fatalf("contact update failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `<contact:add>`)
	assertContains(t, requestXML, `<contact:status s="clientTransferProhibited"></contact:status>`)
	assertContains(t, requestXML, `<contact:rem>`)
	assertContains(t, requestXML, `<contact:status s="clientDeleteProhibited"></contact:status>`)
	assertNotContains(t, requestXML, `<contact:chg>`)
}

func TestContactUpdateEmptyRequestValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.ContactUpdate(types.ContactUpdateRequest{})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func contactUpdateResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <trID>
            <clTRID>CONTACT-UPDATE-TEST</clTRID>
            <svTRID>SERVER-CONTACT-UPDATE</svTRID>
        </trID>
    </response>
</epp>`
}
