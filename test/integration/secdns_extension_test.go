package test

import (
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/extensions/secdns"
	"github.com/hariom-pal/go-epp/types"
)

func TestSecDNSDomainCreateDSDataXML(t *testing.T) {
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
		SecDNS: &secdns.CreateRequest{
			Data: secdns.Data{
				MaxSigLife: 604800,
				DSData: []secdns.DSData{
					{
						KeyTag:     12345,
						Algorithm:  8,
						DigestType: 2,
						Digest:     "49FD46E6C4B45C55D4AC",
						KeyData: &secdns.KeyData{
							Flags:     257,
							Protocol:  3,
							Algorithm: 8,
							PublicKey: "AwEAAaExampleKey",
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("domain create failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `<extension>`)
	assertContains(t, requestXML, `<secDNS:create xmlns:secDNS="urn:ietf:params:xml:ns:secDNS-1.1">`)
	assertContains(t, requestXML, `<secDNS:maxSigLife>604800</secDNS:maxSigLife>`)
	assertContains(t, requestXML, `<secDNS:dsData>`)
	assertContains(t, requestXML, `<secDNS:keyTag>12345</secDNS:keyTag>`)
	assertContains(t, requestXML, `<secDNS:alg>8</secDNS:alg>`)
	assertContains(t, requestXML, `<secDNS:digestType>2</secDNS:digestType>`)
	assertContains(t, requestXML, `<secDNS:digest>49FD46E6C4B45C55D4AC</secDNS:digest>`)
	assertContains(t, requestXML, `<secDNS:keyData>`)
	assertContains(t, requestXML, `<secDNS:flags>257</secDNS:flags>`)
	assertContains(t, requestXML, `<secDNS:protocol>3</secDNS:protocol>`)
	assertContains(t, requestXML, `<secDNS:pubKey>AwEAAaExampleKey</secDNS:pubKey>`)
}

func TestSecDNSDomainCreateKeyDataXML(t *testing.T) {
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
		SecDNS: &secdns.CreateRequest{
			Data: secdns.Data{
				KeyData: []secdns.KeyData{
					{
						Flags:     256,
						Protocol:  3,
						Algorithm: 13,
						PublicKey: "AwEAAcExampleKey",
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("domain create failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `<secDNS:create xmlns:secDNS="urn:ietf:params:xml:ns:secDNS-1.1">`)
	assertContains(t, requestXML, `<secDNS:keyData>`)
	assertContains(t, requestXML, `<secDNS:flags>256</secDNS:flags>`)
	assertContains(t, requestXML, `<secDNS:protocol>3</secDNS:protocol>`)
	assertContains(t, requestXML, `<secDNS:alg>13</secDNS:alg>`)
	assertContains(t, requestXML, `<secDNS:pubKey>AwEAAcExampleKey</secDNS:pubKey>`)
}

func TestSecDNSDomainUpdateXML(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, domainUpdateResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	all := true
	_, err = client.DomainUpdate(types.DomainUpdateRequest{
		Domain: "example.com",
		SecDNS: &secdns.UpdateRequest{
			Urgent: true,
			Remove: &secdns.UpdateRemove{
				All: &all,
			},
			Add: &secdns.UpdateAdd{
				Data: secdns.Data{
					DSData: []secdns.DSData{
						{
							KeyTag:     12345,
							Algorithm:  8,
							DigestType: 2,
							Digest:     "49FD46E6C4B45C55D4AC",
						},
					},
				},
			},
			Change: &secdns.UpdateChange{
				MaxSigLife: 605900,
			},
		},
	})
	if err != nil {
		t.Fatalf("domain update failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `<extension>`)
	assertContains(t, requestXML, `<secDNS:update xmlns:secDNS="urn:ietf:params:xml:ns:secDNS-1.1" urgent="true">`)
	assertContains(t, requestXML, `<secDNS:rem>`)
	assertContains(t, requestXML, `<secDNS:all>true</secDNS:all>`)
	assertContains(t, requestXML, `<secDNS:add>`)
	assertContains(t, requestXML, `<secDNS:keyTag>12345</secDNS:keyTag>`)
	assertContains(t, requestXML, `<secDNS:chg>`)
	assertContains(t, requestXML, `<secDNS:maxSigLife>605900</secDNS:maxSigLife>`)
}

func TestSecDNSDomainInfoParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, secDNSDomainInfoResponse())
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

	if resp.Result.SecDNS.MaxSigLife != 604800 {
		t.Fatalf("unexpected max sig life: %d", resp.Result.SecDNS.MaxSigLife)
	}

	if len(resp.Result.SecDNS.DSData) != 1 {
		t.Fatalf("unexpected DS data: %+v", resp.Result.SecDNS.DSData)
	}
	ds := resp.Result.SecDNS.DSData[0]
	if ds.KeyTag != 12345 ||
		ds.Algorithm != 8 ||
		ds.DigestType != 2 ||
		ds.Digest != "49FD46E6C4B45C55D4AC" {

		t.Fatalf("unexpected DS data: %+v", ds)
	}
	if ds.KeyData == nil ||
		ds.KeyData.Flags != 257 ||
		ds.KeyData.Protocol != 3 ||
		ds.KeyData.Algorithm != 8 ||
		ds.KeyData.PublicKey != "AwEAAaExampleKey" {

		t.Fatalf("unexpected DS key data: %+v", ds.KeyData)
	}

	if len(resp.Result.SecDNS.KeyData) != 1 ||
		resp.Result.SecDNS.KeyData[0].Flags != 256 ||
		resp.Result.SecDNS.KeyData[0].Algorithm != 13 ||
		resp.Result.SecDNS.KeyData[0].PublicKey != "AwEAAcExampleKey" {

		t.Fatalf("unexpected key data: %+v", resp.Result.SecDNS.KeyData)
	}

	if resp.Result.DNSSEC.MaxSigLife != resp.Result.SecDNS.MaxSigLife ||
		len(resp.Result.DNSSEC.DSData) != 1 ||
		len(resp.Result.DNSSEC.KeyData) != 1 {

		t.Fatalf("legacy DNSSEC fields were not populated: %+v", resp.Result.DNSSEC)
	}
}

func TestSecDNSValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainCreate(types.DomainCreateRequest{
		Domain:     "example.com",
		Period:     1,
		Unit:       "y",
		Registrant: "REG001",
		AuthInfo:   "secret",
		SecDNS: &secdns.CreateRequest{
			Data: secdns.Data{
				DSData: []secdns.DSData{
					{KeyTag: 12345, Algorithm: 8, DigestType: 2, Digest: "49FD46E6C4B45C55D4AC"},
				},
				KeyData: []secdns.KeyData{
					{Flags: 256, Protocol: 3, Algorithm: 13, PublicKey: "AwEAAcExampleKey"},
				},
			},
		},
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)

	all := true
	_, err = client.DomainUpdate(types.DomainUpdateRequest{
		Domain: "example.com",
		SecDNS: &secdns.UpdateRequest{
			Remove: &secdns.UpdateRemove{
				All: &all,
				DSData: []secdns.DSData{
					{KeyTag: 12345, Algorithm: 8, DigestType: 2, Digest: "49FD46E6C4B45C55D4AC"},
				},
			},
		},
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func secDNSDomainInfoResponse() string {
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
                <domain:status s="ok"/>
                <domain:clID>Registrar-OTE</domain:clID>
                <domain:crID>Registrar-OTE</domain:crID>
                <domain:crDate>2026-06-30T09:30:00Z</domain:crDate>
                <domain:authInfo>
                    <domain:pw>secret</domain:pw>
                </domain:authInfo>
            </domain:infData>
        </resData>
        <extension>
            <secDNS:infData xmlns:secDNS="urn:ietf:params:xml:ns:secDNS-1.1">
                <secDNS:maxSigLife>604800</secDNS:maxSigLife>
                <secDNS:dsData>
                    <secDNS:keyTag>12345</secDNS:keyTag>
                    <secDNS:alg>8</secDNS:alg>
                    <secDNS:digestType>2</secDNS:digestType>
                    <secDNS:digest>49FD46E6C4B45C55D4AC</secDNS:digest>
                    <secDNS:keyData>
                        <secDNS:flags>257</secDNS:flags>
                        <secDNS:protocol>3</secDNS:protocol>
                        <secDNS:alg>8</secDNS:alg>
                        <secDNS:pubKey>AwEAAaExampleKey</secDNS:pubKey>
                    </secDNS:keyData>
                </secDNS:dsData>
                <secDNS:keyData>
                    <secDNS:flags>256</secDNS:flags>
                    <secDNS:protocol>3</secDNS:protocol>
                    <secDNS:alg>13</secDNS:alg>
                    <secDNS:pubKey>AwEAAcExampleKey</secDNS:pubKey>
                </secDNS:keyData>
            </secDNS:infData>
        </extension>
        <trID>
            <clTRID>DOMAIN-INFO-TEST</clTRID>
            <svTRID>SERVER-DOMAIN-INFO</svTRID>
        </trID>
    </response>
</epp>`
}
