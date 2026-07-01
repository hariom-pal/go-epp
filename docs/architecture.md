# Architecture

`go-epp` keeps the EPP core protocol and optional registry extensions separate.

## Package Layout

- `epp/` contains core RFC5730-RFC5734 command implementations, XML request/response models, transport execution, and parsing glue.
- `types/` contains public SDK request and response types.
- `constants/` contains protocol constants, result codes, statuses, namespaces, and reusable validation helpers.
- `extensions/` contains optional EPP extension packages. Extension packages expose reusable request, response, and XML builders without owning transport execution.
- `cmd/epp-cli/` contains the operational CLI used for registry testing.
- `test/` contains integration-style XML generation, parsing, and validation tests using a local TLS EPP server.

## Extension Boundary

Extension packages are intentionally outside `epp/`.

The core `epp` package remains focused on standard EPP operations. Optional registry features such as fee, secDNS, RGP, and launch live under `extensions/` and are attached to core commands only when requested by public request types.

This keeps the stable core small while allowing registries to support different extension sets.

## Request Flow

1. Public request types are passed to a client method such as `DomainCreate`, `DomainUpdate`, or `DomainInfo`.
2. The `epp` package validates core command inputs.
3. Optional extension builders validate and render their own XML models.
4. The command is serialized with `encoding/xml`.
5. `Execute` sends the framed EPP command and returns the response XML.
6. Response XML is parsed into internal XML models and converted into public response types.

## Compatibility

Public APIs should remain backward compatible. When a reusable extension model supersedes an older convenience field, both should be populated when practical. For example, domain info exposes structured `RGP` data while retaining legacy `RGPStatuses`.
