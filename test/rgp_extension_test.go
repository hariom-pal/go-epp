package test

import (
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/extensions/rgp"
	"github.com/hariom-pal/go-epp/types"
)

func TestRGPDomainUpdateRestoreRequestXMLAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, rgpRestoreRequestResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain: "example.com",
		RGP: &rgp.UpdateRequest{
			Restore: &rgp.Restore{
				Operation: rgp.OperationRequest,
			},
		},
	})
	if err != nil {
		t.Fatalf("domain update failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `<domain:update>`)
	assertContains(t, requestXML, `<domain:name>example.com</domain:name>`)
	assertContains(t, requestXML, `<domain:chg></domain:chg>`)
	assertContains(t, requestXML, `<extension>`)
	assertContains(t, requestXML, `<rgp:update xmlns:rgp="urn:ietf:params:xml:ns:rgp-1.0">`)
	assertContains(t, requestXML, `<rgp:restore op="request"></rgp:restore>`)

	if len(resp.Result.RGP.Statuses) != 1 ||
		resp.Result.RGP.Statuses[0].Status != rgp.StatusPendingRestore {

		t.Fatalf("unexpected RGP update data: %+v", resp.Result.RGP)
	}
}

func TestRGPDomainUpdateRestoreReportXML(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, domainUpdateResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	_, err = client.DomainUpdate(types.DomainUpdateRequest{
		Domain: "example.com",
		RGP: &rgp.UpdateRequest{
			Restore: &rgp.Restore{
				Operation: rgp.OperationReport,
				Report: &rgp.RestoreReport{
					PreDataXML:  `<domain:name>example.com</domain:name>`,
					PostData:    `Post & restore registration data`,
					DeleteTime:  "2026-06-01T10:00:00Z",
					RestoreTime: "2026-06-02T10:00:00Z",
					RestoreReason: rgp.Text{
						Value: "Registrant error.",
					},
					Statements: []rgp.Text{
						{
							Lang:  "en",
							Value: "The registrar has not restored the domain for itself or a third party.",
						},
						{
							Value: "The information in this report is true to the registrar's knowledge.",
						},
					},
					Other: "Supporting information.",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("domain update failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `<rgp:restore op="report">`)
	assertContains(t, requestXML, `<rgp:report>`)
	assertContains(t, requestXML, `<rgp:preData><domain:name>example.com</domain:name></rgp:preData>`)
	assertContains(t, requestXML, `<rgp:postData>Post &amp; restore registration data</rgp:postData>`)
	assertContains(t, requestXML, `<rgp:delTime>2026-06-01T10:00:00Z</rgp:delTime>`)
	assertContains(t, requestXML, `<rgp:resTime>2026-06-02T10:00:00Z</rgp:resTime>`)
	assertContains(t, requestXML, `<rgp:resReason>Registrant error.</rgp:resReason>`)
	assertContains(t, requestXML, `<rgp:statement lang="en">The registrar has not restored the domain for itself or a third party.</rgp:statement>`)
	assertContains(t, requestXML, `<rgp:statement>The information in this report is true to the registrar&#39;s knowledge.</rgp:statement>`)
	assertContains(t, requestXML, `<rgp:other>Supporting information.</rgp:other>`)
}

func TestRGPDomainInfoParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, rgpDomainInfoResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainInfo(types.DomainInfoRequest{
		Domain: "example.com",
	})
	if err != nil {
		t.Fatalf("domain info failed: %v", err)
	}

	requestXML := readRequest(t, requests)
	assertContains(t, requestXML, `<domain:info>`)

	if len(resp.Result.RGP.Statuses) != 3 {
		t.Fatalf("unexpected RGP statuses: %+v", resp.Result.RGP.Statuses)
	}

	expected := []string{
		rgp.StatusRedemption,
		rgp.StatusPendingRestore,
		rgp.StatusPendingDelete,
	}
	for i, status := range expected {
		if resp.Result.RGP.Statuses[i].Status != status ||
			resp.Result.RGPStatuses[i] != status {

			t.Fatalf("unexpected RGP status at %d: structured=%+v legacy=%+v", i, resp.Result.RGP.Statuses, resp.Result.RGPStatuses)
		}
	}

	if resp.Result.RGP.Statuses[1].Lang != "en" ||
		resp.Result.RGP.Statuses[1].Text != "Restore requested" {

		t.Fatalf("unexpected RGP status metadata: %+v", resp.Result.RGP.Statuses[1])
	}
}

func TestRGPValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain: "example.com",
		RGP: &rgp.UpdateRequest{
			Restore: &rgp.Restore{
				Operation: "invalid",
			},
		},
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)

	_, err = client.DomainUpdate(types.DomainUpdateRequest{
		Domain: "example.com",
		RGP: &rgp.UpdateRequest{
			Restore: &rgp.Restore{
				Operation: rgp.OperationReport,
				Report: &rgp.RestoreReport{
					PreData: "pre",
				},
			},
		},
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func rgpRestoreRequestResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <extension>
            <rgp:upData xmlns:rgp="urn:ietf:params:xml:ns:rgp-1.0">
                <rgp:rgpStatus s="pendingRestore"/>
            </rgp:upData>
        </extension>
        <trID>
            <clTRID>DOMAIN-UPDATE-TEST</clTRID>
            <svTRID>SERVER-DOMAIN-UPDATE</svTRID>
        </trID>
    </response>
</epp>`
}

func rgpDomainInfoResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <domain:infData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
                <domain:name>example.com</domain:name>
                <domain:roid>D123456-COM</domain:roid>
                <domain:registrant>REG001</domain:registrant>
                <domain:status s="pendingDelete"/>
                <domain:clID>Registrar-OTE</domain:clID>
                <domain:crID>Registrar-OTE</domain:crID>
                <domain:crDate>2026-06-30T09:30:00Z</domain:crDate>
                <domain:authInfo>
                    <domain:pw>secret</domain:pw>
                </domain:authInfo>
            </domain:infData>
        </resData>
        <extension>
            <rgp:infData xmlns:rgp="urn:ietf:params:xml:ns:rgp-1.0">
                <rgp:rgpStatus s="redemptionPeriod"/>
                <rgp:rgpStatus s="pendingRestore" lang="en">Restore requested</rgp:rgpStatus>
                <rgp:rgpStatus s="pendingDelete"/>
            </rgp:infData>
        </extension>
        <trID>
            <clTRID>DOMAIN-INFO-TEST</clTRID>
            <svTRID>SERVER-DOMAIN-INFO</svTRID>
        </trID>
    </response>
</epp>`
}
