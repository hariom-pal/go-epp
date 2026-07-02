package test

import (
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/extensions/launch"
	"github.com/hariom-pal/go-epp/types"
)

func TestLaunchDomainCreateXMLAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, launchDomainCreateResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainCreate(types.DomainCreateRequest{
		Domain:     "example.com",
		Period:     1,
		Unit:       "y",
		Registrant: "REG001",
		AuthInfo:   "secret",
		Launch: &launch.CreateRequest{
			Type: launch.ObjectApplication,
			Phase: launch.Phase{
				Value: launch.PhaseSunrise,
			},
			CodeMarks: []launch.CodeMark{
				{
					Code: &launch.Code{
						Value:       "49FD46E6C4B45C55D4AC",
						ValidatorID: "tmch",
					},
				},
			},
			Notices: []launch.Notice{
				{
					ID:           "370d8f6a1b2c3",
					ValidatorID:  "tmch",
					NotAfter:     "2026-07-08T10:00:00Z",
					AcceptedDate: "2026-07-01T10:00:00Z",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("domain create failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `<launch:create xmlns:launch="urn:ietf:params:xml:ns:launch-1.0" type="application">`)
	assertContains(t, requestXML, `<launch:phase>sunrise</launch:phase>`)
	assertContains(t, requestXML, `<launch:codeMark>`)
	assertContains(t, requestXML, `<launch:code validatorID="tmch">49FD46E6C4B45C55D4AC</launch:code>`)
	assertContains(t, requestXML, `<launch:notice>`)
	assertContains(t, requestXML, `<launch:noticeID validatorID="tmch">370d8f6a1b2c3</launch:noticeID>`)
	assertContains(t, requestXML, `<launch:notAfter>2026-07-08T10:00:00Z</launch:notAfter>`)
	assertContains(t, requestXML, `<launch:acceptedDate>2026-07-01T10:00:00Z</launch:acceptedDate>`)

	if resp.Result.Launch.Phase.Value != launch.PhaseSunrise ||
		resp.Result.Launch.ApplicationID != "2393-9323-E08C-03B1" {

		t.Fatalf("unexpected launch create data: %+v", resp.Result.Launch)
	}
}

func TestLaunchDomainCreateSignedMarkXML(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, domainCreateResponse("example.com"))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	_, err = client.DomainCreate(types.DomainCreateRequest{
		Domain:     "example.com",
		Period:     1,
		Unit:       "y",
		Registrant: "REG001",
		AuthInfo:   "secret",
		Launch: &launch.CreateRequest{
			Phase: launch.Phase{
				Value: launch.PhaseSunrise,
			},
			SignedMarkXML: []string{
				`<smd:signedMark><smd:id>signed-mark-1</smd:id></smd:signedMark>`,
			},
		},
	})
	if err != nil {
		t.Fatalf("domain create failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `xmlns:smd="urn:ietf:params:xml:ns:signedMark-1.0"`)
	assertContains(t, requestXML, `<smd:signedMark><smd:id>signed-mark-1</smd:id></smd:signedMark>`)
}

func TestLaunchDomainInfoXMLAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, launchDomainInfoResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainInfo(types.DomainInfoRequest{
		Domain: "example.com",
		Launch: &launch.InfoRequest{
			IncludeMark:   true,
			ApplicationID: "abc123",
			Phase: launch.Phase{
				Value: launch.PhaseClaims,
				Name:  launch.PhaseLandrush,
			},
		},
	})
	if err != nil {
		t.Fatalf("domain info failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `<launch:info xmlns:launch="urn:ietf:params:xml:ns:launch-1.0" includeMark="true">`)
	assertContains(t, requestXML, `<launch:phase name="landrush">claims</launch:phase>`)
	assertContains(t, requestXML, `<launch:applicationID>abc123</launch:applicationID>`)

	if resp.Result.LaunchData.Phase.Value != launch.PhaseClaims ||
		resp.Result.LaunchData.Phase.Name != launch.PhaseLandrush ||
		resp.Result.LaunchData.ApplicationID != "abc123" ||
		resp.Result.LaunchData.Status.Status != launch.StatusPendingValidation {

		t.Fatalf("unexpected launch info data: %+v", resp.Result.LaunchData)
	}

	if resp.Result.Launch.Phase != launch.PhaseClaims ||
		resp.Result.Launch.Status != launch.StatusPendingValidation ||
		resp.Result.Launch.ApplicationID != "abc123" {

		t.Fatalf("legacy launch fields were not populated: %+v", resp.Result.Launch)
	}

	assertContains(t, resp.Result.LaunchData.RawXML, `<mark:mark xmlns:mark="urn:ietf:params:xml:ns:mark-1.0">`)
}

func TestLaunchDomainUpdateXML(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, domainUpdateResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	_, err = client.DomainUpdate(types.DomainUpdateRequest{
		Domain: "example.com",
		Launch: &launch.UpdateRequest{
			ApplicationID: "abc123",
			Phase: launch.Phase{
				Value: launch.PhaseSunrise,
			},
		},
	})
	if err != nil {
		t.Fatalf("domain update failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `<domain:chg></domain:chg>`)
	assertContains(t, requestXML, `<launch:update xmlns:launch="urn:ietf:params:xml:ns:launch-1.0">`)
	assertContains(t, requestXML, `<launch:phase>sunrise</launch:phase>`)
	assertContains(t, requestXML, `<launch:applicationID>abc123</launch:applicationID>`)
}

func TestLaunchDomainDeleteXML(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, domainDeleteResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	_, err = client.DomainDelete(types.DomainDeleteRequest{
		Domain: "example.com",
		Launch: &launch.DeleteRequest{
			ApplicationID: "abc123",
			Phase: launch.Phase{
				Value: launch.PhaseSunrise,
			},
		},
	})
	if err != nil {
		t.Fatalf("domain delete failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `<launch:delete xmlns:launch="urn:ietf:params:xml:ns:launch-1.0">`)
	assertContains(t, requestXML, `<launch:phase>sunrise</launch:phase>`)
	assertContains(t, requestXML, `<launch:applicationID>abc123</launch:applicationID>`)
}

func TestLaunchValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainCreate(types.DomainCreateRequest{
		Domain:     "example.com",
		Period:     1,
		Unit:       "y",
		Registrant: "REG001",
		AuthInfo:   "secret",
		Launch:     &launch.CreateRequest{},
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)

	_, err = client.DomainUpdate(types.DomainUpdateRequest{
		Domain: "example.com",
		Launch: &launch.UpdateRequest{
			Phase: launch.Phase{Value: launch.PhaseSunrise},
		},
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)

	_, err = client.DomainDelete(types.DomainDeleteRequest{
		Domain: "example.com",
		Launch: &launch.DeleteRequest{
			Phase: launch.Phase{Value: launch.PhaseSunrise},
		},
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func launchDomainCreateResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1001">
            <msg>Command completed successfully; action pending</msg>
        </result>
        <resData>
            <domain:creData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
                <domain:name>example.com</domain:name>
                <domain:crDate>2026-07-01T10:00:00Z</domain:crDate>
            </domain:creData>
        </resData>
        <extension>
            <launch:creData xmlns:launch="urn:ietf:params:xml:ns:launch-1.0">
                <launch:phase>sunrise</launch:phase>
                <launch:applicationID>2393-9323-E08C-03B1</launch:applicationID>
            </launch:creData>
        </extension>
        <trID>
            <clTRID>DOMAIN-CREATE-TEST</clTRID>
            <svTRID>SERVER-DOMAIN-CREATE</svTRID>
        </trID>
    </response>
</epp>`
}

func launchDomainInfoResponse() string {
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
                <domain:status s="pendingCreate"/>
                <domain:clID>Registrar-OTE</domain:clID>
                <domain:crID>Registrar-OTE</domain:crID>
                <domain:crDate>2026-07-01T10:00:00Z</domain:crDate>
                <domain:authInfo>
                    <domain:pw>secret</domain:pw>
                </domain:authInfo>
            </domain:infData>
        </resData>
        <extension>
            <launch:infData xmlns:launch="urn:ietf:params:xml:ns:launch-1.0">
                <launch:phase name="landrush">claims</launch:phase>
                <launch:applicationID>abc123</launch:applicationID>
                <launch:status s="pendingValidation" lang="en">Pending validation</launch:status>
                <mark:mark xmlns:mark="urn:ietf:params:xml:ns:mark-1.0">
                    <mark:id>mark-1</mark:id>
                </mark:mark>
            </launch:infData>
        </extension>
        <trID>
            <clTRID>DOMAIN-INFO-TEST</clTRID>
            <svTRID>SERVER-DOMAIN-INFO</svTRID>
        </trID>
    </response>
</epp>`
}
