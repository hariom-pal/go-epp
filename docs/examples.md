# Examples

The `examples/` package contains compile-checked snippets for common SDK workflows. They are not standalone binaries; they are intentionally small functions that demonstrate request construction while remaining part of `go build ./...`.

## Connection Lifecycle

- `connect.go`: load YAML configuration and open a TLS EPP connection.
- `login.go`: connect and login.
- `logout.go`: logout and close.
- `hello.go`: send RFC5730 hello and parse the greeting.

## Poll

- `poll.go`: request the next message and acknowledge a queued message.

## Domain

- `domain_check.go`
- `domain_info.go`
- `domain_create.go`
- `domain_update.go`
- `domain_renew.go`
- `domain_transfer.go`
- `domain_delete.go`

## Contact

- `contact_check.go`
- `contact_info.go`
- `contact_create.go`
- `contact_update.go`
- `contact_delete.go`

## Host

- `host_check.go`
- `host_info.go`
- `host_create.go`
- `host_update.go`
- `host_delete.go`

## Extensions

- `fee.go`: attach fee-0.7 data to a domain create.
- `dnssec.go`: attach secDNS DS data to a domain create.
- `rgp.go`: submit an RGP restore request.
- `launch.go`: submit a launch application create with a code mark.

## Running Against a Registry

Use `cmd/epp-cli` for live registry testing. The example package is for library usage patterns, while the CLI handles command-line flags, configuration loading, and output formatting.
