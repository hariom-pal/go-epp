# Supported RFCs

This document tracks protocol and extension coverage.

## Core Protocol

| Specification | Coverage |
| --- | --- |
| RFC5730 | Core EPP envelope, login, logout, hello, poll, result handling, TRID handling |
| RFC5731 | Domain check, info, create, update, delete, renew, transfer |
| RFC5732 | Host check, info, create, update, delete |
| RFC5733 | Contact check, info, create, update, delete |
| RFC5734 | TCP transport framing over TLS |

## Extensions

| Specification | Package | Coverage |
| --- | --- | --- |
| fee-0.7 | `extensions/fee` | Fee check and transform request/response models |
| RFC3915 | `extensions/rgp` | RGP info parsing and restore request/report update extension |
| RFC5910 | `extensions/secdns` | secDNS create, update, and info DS/key data models |
| RFC8334 | `extensions/launch` | Launch create, info, update, delete, application ID, notice, code mark, status models |

## Non-Goals

- Registry-specific pricing policy.
- Registry-specific launch policy.
- Registry-specific DNSSEC algorithm restrictions.
- Registry adapters that hardcode individual registry quirks.

Those concerns should live in application code or optional adapter packages.
