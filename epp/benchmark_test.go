package epp

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"io"
	"net"
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

var benchmarkPayload = []byte(`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><check><domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>example.com</domain:name></domain:check></check><clTRID>CHECK-BENCH</clTRID></command></epp>`)

func BenchmarkFrameReader(b *testing.B) {
	frame := benchmarkFrame(benchmarkPayload)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		conn := &readOnlyConn{reader: bytes.NewReader(frame)}
		if _, err := ReadFrame(conn); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFrameWriter(b *testing.B) {
	conn := discardConn{}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := WriteFrame(conn, benchmarkPayload); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkXMLMarshal(b *testing.B) {
	req := domainCheckRequestXML{
		XMLNS:       constants.EPPNamespace,
		DomainXMLNS: constants.DomainNamespace,
		Command: domainCheckCommandXML{
			ClientTRID: "CHECK-BENCH",
			Check: domainCheckXML{
				Domain: domainCheckNamesXML{
					Names: []string{"example.com", "example.net"},
				},
			},
		},
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := xml.Marshal(req); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkXMLUnmarshal(b *testing.B) {
	responseXML := []byte(benchmarkDomainCheckResponseXML)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var response domainCheckResponseXML
		if err := xml.Unmarshal(responseXML, &response); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDomainCheck(b *testing.B) {
	responseFrame := benchmarkFrame([]byte(benchmarkDomainCheckResponseXML))
	client := &Client{
		conn: &cyclicResponseConn{
			response: responseFrame,
		},
	}
	req := types.DomainCheckRequest{
		Domains: []string{"example.com", "example.net"},
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := client.DomainCheck(req); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkFrame(payload []byte) []byte {
	frame := make([]byte, len(payload)+eppFrameHeaderLength)
	binary.BigEndian.PutUint32(frame[:eppFrameHeaderLength], uint32(len(frame)))
	copy(frame[eppFrameHeaderLength:], payload)
	return frame
}

type readOnlyConn struct {
	reader *bytes.Reader
}

func (c *readOnlyConn) Read(p []byte) (int, error)       { return c.reader.Read(p) }
func (c *readOnlyConn) Write(_ []byte) (int, error)      { return 0, io.ErrClosedPipe }
func (c *readOnlyConn) Close() error                     { return nil }
func (c *readOnlyConn) LocalAddr() net.Addr              { return benchmarkAddr("local") }
func (c *readOnlyConn) RemoteAddr() net.Addr             { return benchmarkAddr("remote") }
func (c *readOnlyConn) SetDeadline(time.Time) error      { return nil }
func (c *readOnlyConn) SetReadDeadline(time.Time) error  { return nil }
func (c *readOnlyConn) SetWriteDeadline(time.Time) error { return nil }

type discardConn struct{}

func (discardConn) Read(_ []byte) (int, error)       { return 0, io.EOF }
func (discardConn) Write(p []byte) (int, error)      { return len(p), nil }
func (discardConn) Close() error                     { return nil }
func (discardConn) LocalAddr() net.Addr              { return benchmarkAddr("local") }
func (discardConn) RemoteAddr() net.Addr             { return benchmarkAddr("remote") }
func (discardConn) SetDeadline(time.Time) error      { return nil }
func (discardConn) SetReadDeadline(time.Time) error  { return nil }
func (discardConn) SetWriteDeadline(time.Time) error { return nil }

type cyclicResponseConn struct {
	response []byte
	offset   int
}

func (c *cyclicResponseConn) Read(p []byte) (int, error) {
	if c.offset >= len(c.response) {
		c.offset = 0
	}
	n := copy(p, c.response[c.offset:])
	c.offset += n
	return n, nil
}

func (c *cyclicResponseConn) Write(p []byte) (int, error)      { return len(p), nil }
func (c *cyclicResponseConn) Close() error                     { return nil }
func (c *cyclicResponseConn) LocalAddr() net.Addr              { return benchmarkAddr("local") }
func (c *cyclicResponseConn) RemoteAddr() net.Addr             { return benchmarkAddr("remote") }
func (c *cyclicResponseConn) SetDeadline(time.Time) error      { return nil }
func (c *cyclicResponseConn) SetReadDeadline(time.Time) error  { return nil }
func (c *cyclicResponseConn) SetWriteDeadline(time.Time) error { return nil }

type benchmarkAddr string

func (a benchmarkAddr) Network() string { return string(a) }
func (a benchmarkAddr) String() string  { return string(a) }

const benchmarkDomainCheckResponseXML = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <domain:chkData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
                <domain:cd>
                    <domain:name avail="1">example.com</domain:name>
                </domain:cd>
                <domain:cd>
                    <domain:name avail="0">example.net</domain:name>
                    <domain:reason>In use</domain:reason>
                </domain:cd>
            </domain:chkData>
        </resData>
        <trID>
            <clTRID>CHECK-BENCH</clTRID>
            <svTRID>SERVER-BENCH</svTRID>
        </trID>
    </response>
</epp>`
