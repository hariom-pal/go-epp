# Architecture

`go-epp` separates the stable EPP protocol core from optional registry extensions. The goal is to keep the SDK predictable for v1.0 while still supporting registries with different extension sets.

## Package Layout

- `epp/`: RFC5730-RFC5734 commands, XML request/response models, TRID generation, `Execute`, session handling, and RFC5734 framing.
- `types/`: public request and response structs used by applications.
- `constants/`: result codes, statuses, transfer operations, object constants, and XML namespaces.
- `extensions/`: optional reusable extension packages.
- `cmd/epp-cli/`: operational CLI for OT&E and smoke testing.
- `examples/`: compile-checked library usage examples.
- `test/`: integration-style tests using a local TLS EPP server.

## Extension Boundary

Extension code lives under `extensions/`, not under `epp/`.

```text
epp/
  core commands and transport

extensions/
  common/
  fee/
  launch/
  rgp/
  secdns/
```

The `epp` package attaches extension XML only when a public request contains extension data. This preserves existing behavior for SDK users who only use core RFC5730-RFC5734 commands.

## Request Flow

```text
types.DomainCreateRequest
        |
        v
epp.DomainCreate
        |
        +-- core validation
        +-- optional extension validation
        +-- XML marshal
        |
        v
epp.Execute
        |
        +-- WriteFrame
        +-- ReadFrame
        |
        v
parse XML response into public response types
```

## Response Flow

1. XML response is unmarshaled into internal XML structs.
2. EPP result codes are checked.
3. Registry errors are returned as `*epp.Error`.
4. Successful responses are converted into public `types` models.
5. Extension data is parsed by its owning extension package.

## Compatibility Policy

The v1.0 API should avoid breaking changes. Additive fields are preferred over replacements. Compatibility fields remain populated when structured extension models are introduced; for example, `DomainInfoResult.RGPStatuses` remains populated alongside `DomainInfoResult.RGP`.

## Operational Guidance

- Reuse a client for one EPP session.
- Login before issuing commands that require authentication.
- Logout and close the connection when finished.
- Send extension XML only when the registry supports the namespace and policy.
- Treat registry policy as application-level configuration, not SDK behavior.
