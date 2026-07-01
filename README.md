# go-epp

`go-epp` is a production-oriented Go SDK for the Extensible Provisioning Protocol (EPP). It provides typed request and response models, RFC5734 TLS framing, structured EPP errors, IDN conversion, a small operational CLI, and reusable extension packages for common registry launch and policy workflows.

The project is preparing for a public v1.0 release. The core API is designed to stay stable while optional registry behavior is kept in dedicated extension packages.

## Project Overview

EPP is the protocol used by domain registries and registrars to provision domains, contacts, hosts, and registry-specific extensions. `go-epp` focuses on SDK use in registrar systems, registry integrations, OT&E automation, and operational tooling.

The SDK separates the stable protocol core from optional extensions:

- `epp/` implements RFC5730-RFC5734 commands and transport.
- `types/` contains public request and response models.
- `extensions/` contains reusable extension packages.
- `cmd/epp-cli/` provides an operational command-line client.
- `examples/` contains compile-checked library usage examples.

## Features

- RFC5734 EPP frame reader and writer over TLS.
- Login, logout, hello, and poll.
- Domain check, info, create, update, renew, transfer, and delete.
- Contact check, info, create, update, and delete.
- Host check, info, create, update, and delete.
- Structured response types with result code, result message, client TRID, and server TRID.
- Structured `epp.Error` helpers for common EPP failure classes.
- IDNA/Punycode conversion for internationalized domain and host names.
- Optional fee, secDNS, RGP, and launch extension support.
- CLI for registry OT&E and production smoke testing.
- Tests for XML generation, response parsing, validation, framing, and extension integration.

## Supported RFCs

- RFC5730: Extensible Provisioning Protocol core.
- RFC5731: Domain name mapping.
- RFC5732: Host mapping.
- RFC5733: Contact mapping.
- RFC5734: TCP transport and framing over TLS.
- RFC3915: Redemption Grace Period extension.
- RFC5910: DNSSEC secDNS extension.
- RFC8334: Launch Phase Mapping extension.

## Supported Extensions

- `extensions/fee`: fee-0.7 check and transform extension models.
- `extensions/secdns`: RFC5910 secDNS-1.1 DS/key data models for domain create, update, and info.
- `extensions/rgp`: RFC3915 RGP info and restore request/report models.
- `extensions/launch`: RFC8334 launch phase models for create, info, update, and delete.
- `extensions/common`: shared extension helpers.

Registry-specific policy is intentionally not hardcoded. Applications should decide which extensions to send based on registry documentation and greeting-advertised namespaces.

## Installation

```bash
go get github.com/hariom-pal/go-epp
```

## Quick Start

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

resp, err := client.DomainCheck(types.DomainCheckRequest{
    Domains: []string{"example.com"},
})
if err != nil {
    return err
}

fmt.Println(resp.Results[0].Domain, resp.Results[0].Available)
```

## Configuration

The CLI and examples use YAML configuration:

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

## Connect

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
```

## Login

```go
if err := client.Login(); err != nil {
    return err
}
defer client.Logout()
```

## Domain Examples

```go
check, err := client.DomainCheck(types.DomainCheckRequest{
    Domains: []string{"example.com", "example.net"},
})
```

```go
created, err := client.DomainCreate(types.DomainCreateRequest{
    Domain:        "example.com",
    Period:        1,
    Unit:          "y",
    Registrant:    "REG123",
    AdminContacts: []string{"CNT-ADMIN"},
    TechContacts:  []string{"CNT-TECH"},
    NameServers:   []string{"ns1.example.com", "ns2.example.com"},
    AuthInfo:      "change-me",
})
```

```go
updated, err := client.DomainUpdate(types.DomainUpdateRequest{
    Domain:         "example.com",
    AddNameServers: []string{"ns3.example.com"},
})
```

## Contact Examples

```go
contact, err := client.ContactInfo(types.ContactInfoRequest{
    ContactID: "CNT123",
})
```

```go
created, err := client.ContactCreate(types.ContactCreateRequest{
    ContactID: "CNT123",
    InternationalPostalInfo: &types.PostalInfo{
        Type:        "int",
        Name:        "Example User",
        Street:      []string{"1 Example Street"},
        City:        "Dulles",
        CountryCode: "US",
    },
    Voice:    types.Phone{Number: "+1.7035555555"},
    Email:    "user@example.test",
    AuthInfo: "change-me",
})
```

## Host Examples

```go
host, err := client.HostCreate(types.HostCreateRequest{
    HostName: "ns1.example.com",
    Addresses: []types.HostAddress{
        {IPVersion: "v4", Address: "192.0.2.1"},
        {IPVersion: "v6", Address: "2001:db8::1"},
    },
})
```

## Poll

```go
msg, err := client.Poll(types.PollRequest{
    Operation: constants.PollRequest,
})
```

```go
ack, err := client.Poll(types.PollRequest{
    Operation: constants.PollAcknowledge,
    MessageID: "12345",
})
```

## Hello

```go
greeting, err := client.Hello()
```

`Hello` sends `<hello/>` and parses the server greeting using the same greeting model used when the connection opens.

## Fee

```go
resp, err := client.DomainCreate(types.DomainCreateRequest{
    Domain:     "example.com",
    Period:     1,
    Unit:       "y",
    Registrant: "REG123",
    AuthInfo:   "change-me",
    Fee: &fee.TransformRequest{
        Currency: "USD",
        Fees: []fee.Fee{
            {Amount: "5.00"},
        },
    },
})
```

## DNSSEC

```go
resp, err := client.DomainCreate(types.DomainCreateRequest{
    Domain:     "example.com",
    Period:     1,
    Unit:       "y",
    Registrant: "REG123",
    AuthInfo:   "change-me",
    SecDNS: &secdns.CreateRequest{
        Data: secdns.Data{
            DSData: []secdns.DSData{
                {KeyTag: 12345, Algorithm: 8, DigestType: 2, Digest: "49FD46E6C4B45C55D4AC"},
            },
        },
    },
})
```

## RGP

```go
resp, err := client.DomainUpdate(types.DomainUpdateRequest{
    Domain: "example.com",
    RGP: &rgp.UpdateRequest{
        Restore: &rgp.Restore{Operation: rgp.OperationRequest},
    },
})
```

## Launch

```go
resp, err := client.DomainCreate(types.DomainCreateRequest{
    Domain:     "example.com",
    Period:     1,
    Unit:       "y",
    Registrant: "REG123",
    AuthInfo:   "change-me",
    Launch: &launch.CreateRequest{
        Type:  launch.ObjectApplication,
        Phase: launch.Phase{Value: launch.PhaseSunrise},
        CodeMarks: []launch.CodeMark{
            {Code: &launch.Code{Value: "49FD46E6C4B45C55D4AC", ValidatorID: "tmch"}},
        },
    },
})
```

## CLI Usage

```bash
go run ./cmd/epp-cli -config configs/config.yaml -hello
go run ./cmd/epp-cli -config configs/config.yaml -poll
go run ./cmd/epp-cli -config configs/config.yaml -poll-ack 12345
go run ./cmd/epp-cli -config configs/config.yaml -check example.com
go run ./cmd/epp-cli -config configs/config.yaml -info example.com
```

Create a domain:

```bash
go run ./cmd/epp-cli \
  -config configs/config.yaml \
  -create example.com \
  -period 1 \
  -unit y \
  -registrant REG123 \
  -admin CNT-ADMIN \
  -tech CNT-TECH \
  -ns ns1.example.com \
  -ns ns2.example.com \
  -authInfo change-me
```

## Library Usage

The SDK is designed for direct use from registrar services and integration tooling:

- Reuse one `epp.Client` per EPP session.
- Call `Login` before transform/query commands when required by the registry.
- Always call `Logout` and `Close` when a session is complete.
- Inspect `*epp.Error` with `errors.As` for EPP result-code-specific behavior.
- Send optional extensions only when supported by the target registry.

## Architecture Diagram

```text
Application / Registrar System
            |
            v
      github.com/hariom-pal/go-epp/types
            |
            v
      github.com/hariom-pal/go-epp/epp
      - command builders
      - XML parsing
      - TRID generation
      - Execute()
      - RFC5734 framing
            |
            +-----------------------------+
            | optional extension packages |
            | fee | secdns | rgp | launch |
            +-----------------------------+
            |
            v
        TLS EPP Registry
```

## Roadmap

- v1.0 API freeze and compatibility policy.
- Additional registry-specific examples and compatibility notes.
- Broader benchmark coverage for high-volume registrar workloads.
- Optional registry adapter packages outside the protocol core.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

Before submitting a change:

```bash
gofmt -w .
go vet ./...
go test ./...
go build ./...
```

## License

Apache License 2.0. See [LICENSE](LICENSE).
