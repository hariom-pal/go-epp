# SDK Validation Report

## SDK Version

- Module: `github.com/hariom-pal/go-epp`
- Version: `Unreleased` development series, current changelog baseline `0.0.0`
- Validation package: Package-30, Production Validation & Registrar Integration

## RFC Coverage

| Area | Coverage |
| --- | --- |
| RFC5730 | Greeting, hello, login, logout, poll request, poll ack, result handling, TRID handling |
| RFC5731 | Domain check, info, create, update, delete, renew, transfer request/query/approve/reject/cancel |
| RFC5732 | Host check, info, create, update, delete |
| RFC5733 | Contact check, info, create, update, delete |
| RFC5734 | TCP/TLS frame reader and frame writer |

No new RFC was implemented in this package.

## Extension Coverage

| Extension | Coverage |
| --- | --- |
| fee-0.7 | Fee check and transform validation hooks |
| RFC5910 secDNS-1.1 | Create/update validation using DS/key interface rules |
| RFC3915 RGP | Restore request/report update validation |
| RFC8334 launch-1.0 | Create, info, update, delete validation |

No new extension was implemented in this package.

## Architecture Score

- Score: 8.5/10
- Rationale: Core protocol code, public request/response types, constants, and extension packages remain separated. The new `validation` package is additive and does not alter XML generation or public SDK behavior.

## Code Quality

- Validation is reusable by registrar platform services before sending commands to a registry.
- Existing command builders continue to enforce their current runtime checks.
- Integration tests are separated under `test/integration/`.
- New focused validation unit tests are under `test/unit/`.
- Benchmarks cover frame reader, frame writer, XML marshal, XML unmarshal, domain check, and domain info.

## Coverage

- `go test ./... -cover`: passes.
- `go test ./... -coverprofile=coverage.out`: passes and generates a profile.
- `go tool cover -html=coverage.out`: supported by the generated profile.
- Default profile total: 0.0%, because tests are intentionally isolated under `test/` packages and Go's default coverage only instruments each package's own test binary.
- Aggregate check used for readiness visibility: `go test ./... -coverpkg=./... -coverprofile=coverage_all.out`
- Aggregate total from this run: 53.3%.
- Target: 90%+.
- Status: Not yet met.

## Benchmarks

Latest local run:

| Benchmark | Result |
| --- | --- |
| `BenchmarkFrameReader` | 2216 ns/op, 300 B/op, 4 allocs/op |
| `BenchmarkFrameWriter` | 1707 ns/op, 240 B/op, 1 alloc/op |
| `BenchmarkXMLMarshal` | 79745 ns/op, 5632 B/op, 16 allocs/op |
| `BenchmarkXMLUnmarshal` | 329207 ns/op, 6928 B/op, 176 allocs/op |
| `BenchmarkDomainCheck` | 435309 ns/op, 14935 B/op, 209 allocs/op |
| `BenchmarkDomainInfo` | 778457 ns/op, 20344 B/op, 315 allocs/op |

## Missing Features

- 90%+ aggregate coverage has not been reached.
- Registry-specific adapters, policy engines, pricing rules, and launch policy decisions remain application-level responsibilities.
- Live registry OT&E certification matrix is not included in this SDK package.

## Known Limitations

- CLI info commands redact authInfo values that may be returned by registries.
- Validation intentionally avoids registry-specific policy checks such as premium pricing eligibility, launch phase calendars, DNSSEC algorithm policy, and reserved-name policy.
- Coverage visibility requires `-coverpkg=./...` when tests remain in external `test/` packages.
