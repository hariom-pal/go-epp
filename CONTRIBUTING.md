# Contributing

Thank you for helping improve `go-epp`.

## Development Setup

Use the Go version declared in `go.mod`.

Run the full local quality gate before opening a pull request:

```bash
gofmt -w .
go vet ./...
go test ./...
go build ./...
```

## Design Guidelines

- Keep RFC5730-RFC5734 protocol code in `epp/`.
- Keep optional extensions under `extensions/`.
- Avoid registry-specific behavior in reusable packages.
- Reuse existing request, response, XML, constants, error, and TRID patterns.
- Preserve backward compatibility for public request and response types when possible.
- Add focused tests for XML generation, parsing, and validation.

## Pull Requests

Pull requests should include:

- A short description of the protocol behavior being added or changed.
- Links to the relevant RFC or registry specification.
- Tests that prove generated XML and parsed response data.
- Confirmation that `gofmt`, `go vet`, `go test`, and `go build` pass.

## Security

Do not include real registry credentials, client certificates, private keys, production domains, or unredacted EPP payloads in issues or pull requests.
