# Changelog

All notable changes to this project will be documented here.

The format follows Keep a Changelog, and this project aims to use semantic versioning from v1.0.0 onward.

## Unreleased

### Added

- Release-readiness documentation, issue templates, support files, lint configuration, examples, and benchmarks.
- RFC5730 poll request and acknowledge support.
- RFC5730 hello support using the shared greeting parser.
- fee-0.7 reusable extension package.
- RFC5910 secDNS-1.1 reusable extension package.
- RFC3915 RGP reusable extension package.
- RFC8334 launch reusable extension package.
- Optional extension integration for supported domain operations.
- Repository documentation for architecture, supported RFCs, and registry compatibility.

### Changed

- Extension code is organized under top-level `extensions/` packages to keep the core `epp` package focused on standard EPP operations.

## 0.0.0

- Initial development series.
