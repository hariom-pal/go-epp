# Security Policy

## Supported Versions

The project is preparing for v1.0.0. Until v1.0.0 is released, security fixes are applied to the main development branch.

## Reporting a Vulnerability

Please report suspected vulnerabilities privately to the project maintainer.

Include:

- A clear description of the issue.
- Steps to reproduce, if available.
- Affected commands or packages.
- Redacted XML payloads, logs, or stack traces when helpful.

Do not publish real EPP credentials, client certificates, private keys, authInfo values, or production registry payloads in public issues.

## Scope

Security-sensitive areas include:

- TLS transport and certificate handling.
- EPP framing and response parsing.
- Authentication and authInfo handling.
- XML generation and parsing.
- Error handling that might expose sensitive data.
