# go-epp

`go-epp` is a Go SDK for building registry and registrar integrations over the Extensible Provisioning Protocol (EPP). The SDK focuses on production-quality RFC5730-family support, clear request and response types, IDN handling, structured EPP errors, and a small CLI for operational testing.

The project is under active development and is preparing for a v1.0 architecture freeze.

## Features

- TLS EPP session support using RFC5734 framing
- Login, logout, hello, and poll commands
- Domain check, info, create, update, renew, transfer, and delete
- Contact check, info, create, update, and delete
- Host check, info, create, update, and delete
- Reusable fee, secDNS, and RGP extension packages
- Unicode domain and host conversion through IDNA/Punycode
- Structured response types with result codes and transaction IDs
- Structured `epp.Error` values for registry errors
- CLI for OT&E and registry integration testing
- Focused tests for XML generation, response parsing, and validation

## Supported RFCs

- RFC5730: Extensible Provisioning Protocol core
- RFC5731: Domain name mapping
- RFC5732: Host mapping
- RFC5733: Contact mapping
- RFC5734: TCP transport
- RFC3915: Redemption Grace Period extension
- RFC5910: DNSSEC secDNS extension
- fee-0.7: Fee extension

Extension namespaces are advertised at login where required by supported registries. Launch extension support is planned under `extensions/launch`.

## Installation

```bash
go get github.com/hariom-pal/go-epp
```

## Quick Start

```go
package main

import (
	"fmt"
	"log"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func main() {
	cfg, err := epp.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	client, err := epp.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if err := client.Login(); err != nil {
		log.Fatal(err)
	}
	defer client.Logout()

	resp, err := client.DomainCheck(types.DomainCheckRequest{
		Domains: []string{"example.in"},
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range resp.Results {
		fmt.Printf("%s available=%t\n", result.Domain, result.Available)
	}
}
```

## Configuration

The CLI and examples use a YAML configuration file:

```yaml
server:
  host: ote.example.registry
  port: 700

authentication:
  username: registrar-id
  password: registrar-password

tls:
  cert_file: certs/client.crt
  key_file: certs/client.key
  ca_file: certs/ca.crt
  insecure_skip_verify: false

timeout:
  connect: 30
  read: 30
  write: 30
```

## Login Example

```go
cfg, err := epp.LoadConfig("configs/config.yaml")
if err != nil {
	return err
}

client, err := epp.Connect(cfg)
if err != nil {
	return err
}
defer client.Close()

if err := client.Login(); err != nil {
	return err
}
defer client.Logout()
```

## Domain Example

```go
resp, err := client.DomainCreate(types.DomainCreateRequest{
	Domain:          "example.in",
	Period:          1,
	Unit:            "y",
	Registrant:      "REG123",
	AdminContacts:   []string{"CNT001"},
	TechContacts:    []string{"CNT002"},
	BillingContacts: []string{"CNT003"},
	NameServers:     []string{"ns1.example.in", "ns2.example.in"},
	AuthInfo:        "change-me",
})
if err != nil {
	return err
}

fmt.Println(resp.Result.Domain)
```

## Contact Example

```go
resp, err := client.ContactInfo(types.ContactInfoRequest{
	ContactID: "CNT001",
})
if err != nil {
	return err
}

fmt.Println(resp.Contact.Email)
```

## Host Example

```go
resp, err := client.HostCreate(types.HostCreateRequest{
	HostName: "ns1.example.in",
	Addresses: []types.HostAddress{
		{IPVersion: "v4", Address: "192.0.2.1"},
		{IPVersion: "v6", Address: "2001:db8::1"},
	},
})
if err != nil {
	return err
}

fmt.Println(resp.Result.HostName)
```

## CLI Examples

Check a domain:

```bash
go run ./cmd/epp-cli \
  -config configs/config.yaml \
  -check example.in
```

Create a domain:

```bash
go run ./cmd/epp-cli \
  -config configs/config.yaml \
  -create example.in \
  -period 1 \
  -unit y \
  -registrant REG123 \
  -admin CNT001 \
  -tech CNT002 \
  -billing CNT003 \
  -ns ns1.example.in \
  -ns ns2.example.in \
  -authInfo change-me
```

Request a domain transfer:

```bash
go run ./cmd/epp-cli \
  -config configs/config.yaml \
  -domain-transfer example.in \
  -op request \
  -authInfo change-me
```

Create a contact:

```bash
go run ./cmd/epp-cli \
  -config configs/config.yaml \
  -contact-create CNT001 \
  -contact-name "Example User" \
  -contact-street "1 Example Street" \
  -contact-city Dulles \
  -contact-cc US \
  -contact-voice +1.7035555555 \
  -contact-email user@example.test \
  -contact-authInfo change-me
```

Create a host:

```bash
go run ./cmd/epp-cli \
  -config configs/config.yaml \
  -host-create ns1.example.in \
  -ipv4 192.0.2.1 \
  -ipv6 2001:db8::1
```

Update a host:

```bash
go run ./cmd/epp-cli \
  -config configs/config.yaml \
  -host-update ns1.example.in \
  -host-add-ipv4 192.0.2.10 \
  -host-rem-ipv4 192.0.2.1 \
  -host-add-status clientUpdateProhibited \
  -host-new-name ns2.example.in
```

## Examples

The `examples/` package contains small compile-checked examples for login/logout, domain, contact, and host operations. They are intentionally written as functions rather than standalone binaries so they build as part of `go build ./...`.

## Roadmap

- Poll message support
- RGP extension commands
- DNSSEC extension commands
- Fee extension commands
- Launch extension commands
- Additional registry-specific extension packages
- v1.0 API freeze and compatibility policy

## Contributing

Contributions should preserve the existing architecture and request/response patterns. Before submitting a change, run:

```bash
gofmt -w .
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run
go vet ./...
go test ./...
go build ./...
```

New EPP operations should include:

- Request and response types in `types`
- XML request/response models in `epp`
- A `Client` method returning `(*Response, error)`
- CLI support when useful
- Tests for XML generation, parsing, and validation

## License

Apache License 2.0. See [LICENSE](LICENSE).
