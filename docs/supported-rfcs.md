# Supported RFCs

This document tracks protocol coverage for the SDK.

## Core EPP

- RFC5730: EPP core, login, logout, hello, poll, transaction IDs, results, and errors.
- RFC5731: Domain check, info, create, update, delete, renew, and transfer.
- RFC5732: Host check, info, create, update, and delete.
- RFC5733: Contact check, info, create, update, and delete.
- RFC5734: TCP transport framing over TLS.

## Extensions

- RFC3915: Redemption Grace Period extension under `extensions/rgp`.
- RFC5910: DNSSEC secDNS-1.1 extension under `extensions/secdns`.
- RFC8334: Launch Phase Mapping extension under `extensions/launch`.
- fee-0.7: Fee extension under `extensions/fee`.

## Notes

The SDK aims to implement registry-independent RFC models. Registry-specific policy, pricing, validation, and operational rules should live in application code or registry adapters, not in the core SDK.
