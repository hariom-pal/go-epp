# NIXI UAT API Test Plan

Status: executed 2026-09-13 against NIXI OT&E (`epp.ote.nixiregistry.in:700`, svID `nixi EPP Server - ote`, client `csc_egov_in_a`)
Environment: NIXI UAT/OT&E only

Server greeting advertises: domain-1.0, contact-1.0, host-1.0; secDNS-1.1, rgp-1.0, idn-1.0, fee-0.7, launch-1.0.
All five extensions negotiate successfully at login.

Nine SDK defects were found and fixed; see "SDK Defects Found" below. A
reusable live suite now exercises the SDK against the registry:

```bash
go test -race -tags live ./test/live/ -v
```

Latest run: 115 live assertions, 0 failures, 0 skips, 0 data races.

This checklist is the evidence log for testing `go-epp` against NIXI UAT/OT&E. Do not run transform operations unless the config target is clearly UAT/OT&E/test and the command includes `--confirm-transform`.

Use configured UAT-only object identifiers. Do not use production domains, contacts, hosts, passwords, authInfo values, private keys, or production registry endpoints in this document.

## Safe Execution Rules

- Use `configs/nixi-uat.yaml` or another ignored local config copied from `configs/nixi-uat.yaml.example`.
- Keep real credentials, certificates, private keys, authInfo, and UAT logs out of git.
- Run one API operation per command in `--uat` mode.
- Capture evidence with `--capture-dir .uat-captures` only after confirming that directory remains uncommitted.
- Record response result codes/messages manually in the table after each test.
- If a transform command returns an ambiguous transport error, reconcile with info/check/transfer query/poll before retrying.

## Phase 1 - Connectivity/Session

| # | API | Test Object | Result Code | Status | Notes |
| - | --- | ----------- | ----------- | ------ | ----- |
| 1 | connect/greeting | NIXI OT&E endpoint | - | PASS | TLS + greeting OK. svID `nixi EPP Server - ote`. |
| 2 | hello | NIXI OT&E endpoint | - | PASS | 3 objects, 5 extensions, lang en, version 1.0. |
| 3 | login | csc_egov_in_a | 1000 | PASS | All 3 objects + all 5 extension URIs negotiated. |
| 4 | logout | OT&E session | 1500 | PASS | **Defect 1**: 1500 was treated as an error until fixed. |

## Phase 2 - Query Operations

| # | API | Test Object | Result Code | Status | Notes |
| - | --- | ----------- | ----------- | ------ | ----- |
| 5 | domain-check | cscuatd01x.in + others | 1000 | PASS | avail/reason parsed. Most dictionary-like `.in` names are registry-reserved in OT&E. |
| 6 | domain-info | cscuatd01x.in | 1000 | PASS | Registrant, 3 contact types, 2 NS, statuses, dates. authInfo redacted in output. |
| 7 | contact-check | cscuatc01x/02x/03x | 1000 | PASS | All available before create. |
| 8 | contact-info | cscuatc01x | 1000 | PASS | ROID `CO_...-NIXI`, int postalInfo, voice, email. |
| 9 | host-check | ns1/ns2.cscuatd01x.in | 1000 | PASS | Both available before create. |
| 10 | host-info | ns1.cscuatd01x.in | 1000 | PASS | ROID, status ok, v4 + v6 addresses. |
| 11 | poll request | OT&E poll queue | 1300 | PASS | Queue empty (no messages). |

## Phase 3 - Contact Lifecycle

| # | API | Test Object | Result Code | Status | Notes |
| - | --- | ----------- | ----------- | ------ | ----- |
| 12 | contact-create | cscuatc01x, cscuatc02x | 1000 | PASS | int postalInfo, IN country, E.164 voice. |
| 13 | contact-info-after-create | cscuatc01x | 1000 | PASS | All submitted fields round-tripped. |
| 14 | contact-update | cscuatc01x | 1000 | PASS | Email + voice changed; verified by re-info. |
| 15 | contact-transfer-query | cscuatc01x | 2301 | PASS | Correct "Object not pending transfer" for a non-pending object. |
| 16 | contact-transfer-request | run-scoped contact | 2304 | PARTIAL | Encoding verified live; a completed transfer needs a second registrar. |
| 17 | contact-transfer-approve | run-scoped contact | 2301 | PARTIAL | Registry parsed the command; no pending transfer to approve. |
| 18 | contact-transfer-reject | run-scoped contact | 2301 | PARTIAL | Registry parsed the command; no pending transfer to reject. |
| 19 | contact-transfer-cancel | run-scoped contact | 2301 | PARTIAL | Registry parsed the command; no pending transfer to cancel. |
| 20 | contact-delete | cscuatc03x | 1000 | PASS | Disposable unlinked contact created then deleted. |

## Phase 4 - Host Lifecycle

| # | API | Test Object | Result Code | Status | Notes |
| - | --- | ----------- | ----------- | ------ | ----- |
| 21 | host-create | ns1/ns2.cscuatd01x.in | 1000 | PASS | Registry rejects RFC5737 documentation IPs (2306); routable addresses accepted. |
| 22 | host-info-after-create | ns1.cscuatd01x.in | 1000 | PASS | Address and ROID verified. |
| 23 | host-update | ns1.cscuatd01x.in | 1000 | PASS | IPv6 added; verified by re-info. |
| 24 | host-delete | ns2.cscuatd01x.in | 1000 | PASS | Detached from domain first, then deleted. |

## Phase 5 - Domain Lifecycle

| # | API | Test Object | Result Code | Status | Notes |
| - | --- | ----------- | ----------- | ------ | ----- |
| 25 | domain-create | cscuatd01x.in | 1000 | PASS | 1 year; registrant + admin/tech/billing; expiry 2027-09-13. |
| 26 | domain-info-after-create | cscuatd01x.in | 1000 | PASS | Full object verified. |
| 27 | domain-update | cscuatd01x.in | 1000 | PASS | NS added, then one removed. |
| 28 | domain-renew | cscuatd01x.in | 1000 | PASS | +2 years, 2027-09-13 -> 2029-09-13. |
| 29 | domain-delete | cscuatd02x.in | 1000 | PASS | Deleted inside add grace period, so removed immediately rather than entering redemption. |

## Phase 6 - Domain Transfer

| # | API | Test Object | Result Code | Status | Notes |
| - | --- | ----------- | ----------- | ------ | ----- |
| 30 | domain-transfer-query | cscuatd01x.in | 2301 | PASS | **Defect 3**: failed client-side before the fix. |
| 31 | domain-transfer-request | run-scoped domain | 2304 | PARTIAL | Encoding verified live; a completed transfer needs a second registrar. |
| 32 | domain-transfer-approve | run-scoped domain | 2301 | PARTIAL | Registry parsed the command; no pending transfer to approve. |
| 33 | domain-transfer-reject | run-scoped domain | 2301 | PARTIAL | Registry parsed the command; no pending transfer to reject. |
| 34 | domain-transfer-cancel | run-scoped domain | 2301 | PARTIAL | Registry parsed the command; no pending transfer to cancel. |

## Phase 7 - Poll

| # | API | Test Object | Result Code | Status | Notes |
| - | --- | ----------- | ----------- | ------ | ----- |
| 35 | poll request | OT&E poll queue | 1300 | PASS | Still empty after the full lifecycle; own actions raise no messages. |
| 36 | poll ack | OT&E poll msgID | 1000 | PASS | Queue filled during lifecycle testing; ack verified, count 4 -> 3. |

## Phase 8 - Extensions

| # | API | Test Object | Result Code | Status | Notes |
| - | --- | ----------- | ----------- | ------ | ----- |
| 37 | secDNS domain-update/info | cscuatd01x.in | 1000 | PASS | DS record added, read back by domain-info, then removed. |
| 38 | RGP restore request/report | cscuatd02x.in | 2303 | BLOCKED | Add-grace-period deletion purges immediately, so redemption is never entered. Needs a domain older than the add grace period. |
| 39 | fee check/create | cscuatfee01x.in | 1000 | PASS | fee-0.7: INR 500/yr, class standard, refundable. creData fee returned on IDN create. |
| 40 | IDN check/create/info | jaanchcsc.bharat (Devanagari) | 1000 | PASS | **Defect 4**: idn-1.0 extension was missing entirely. Registered with table `hi`. |
| 41 | launch-1.0 | run-scoped domain | 2306 | PASS | **Defect 5**: emitted a bogus `<ChoiceXML>` element. Now parsed by the registry; no launch phase is open to complete a sunrise. |

## Coverage Beyond the Original Plan

The live suite also covers: host rename (`chg`), domain registrant change,
authInfo rotation, typed contact add/remove, `hosts` attribute variants
(all/del/sub/none), localized postalInfo, contact `disclose`, contact status
add/remove, hostAttr delegation (refused 2102, encoding correct), every contact
and domain transfer operation, concurrency on one session and across sessions,
goroutine-leak checks, reconnect and recovery, transform idempotency and
reconciliation, and negative paths for TLS, authentication, framing and
malformed server responses.

## SDK Defects Found and Fixed

All nine were found by running against a real registry, and all are generic
(they would affect any registry, not just NIXI).

1. **Logout always reported as an error.** `constants.IsSuccessResultCode`
   enumerated only 1000/1001/1300/1301, so result code 1500 -- which RFC 5730
   *requires* in response to `<logout>` -- was classified as a failure. Every
   clean session close returned an error. Fixed by classifying the whole 1xxx
   range as positive completion, per RFC 5730 section 2.6.

2. **Server diagnostics were discarded.** `<result>` was parsed for `code` and
   `msg` only, dropping the `<value>` and `<extValue><reason>` elements that
   RFC 5730 section 2.6 defines to say *which* part of a command was rejected.
   A real failure surfaced as the useless "Parameter value policy error"
   instead of "address 203.0.113.11 forbidden by rule for 203.0.113.0/24,
   Reserved for TEST-NET-3 RFC5737". Fixed by parsing both into
   `types.Result.Values` / `epp.Error.Values` and including them in `Error()`.
   This fix is what made defects 4 and the RGP finding diagnosable at all.

3. **Transfer commands rejected a defaulted period unit.** An optional period
   was only treated as absent when *both* value and unit were empty, so any
   caller that defaults the unit to "y" (as the bundled CLI does) got a
   client-side "period must be between 1 and 99" and could never issue a
   transfer query/approve/reject/cancel. A unit without a value carries no
   meaning; fixed to treat an optional period as absent whenever no value is set.

4. **No idn-1.0 extension.** The SDK converted names to punycode but had no way
   to send the IDN table tag, so IDN registration failed with
   "idn_table required for idn" at any registry that requires one -- while the
   SDK happily negotiated the idn-1.0 URI at login. Added `extensions/idn`,
   wired it into `DomainCreateRequest.IDN`, and exposed `-idn-table` /
   `-idn-uname` on the CLI.

5. **Launch extension emitted a bogus element.** `launch:create` carried
   `<ChoiceXML>` in the EPP namespace instead of the caller's signed mark:
   `encoding/xml` honours the `,innerxml` flag on marshal only for `string` and
   `[]byte` fields, and silently emits any other type as an ordinary element
   named after the Go field. Signed marks were both mis-wrapped and dropped, so
   sunrise registration was impossible at any registry. The pre-existing test
   passed because it only asserted the mark substring was *present*, which it
   was, inside the wrapper. See RFC 8334 section 3.1.

6. **IDN table and U-label were silently dropped.** Domain info parsed only
   `<idn:infData>`, while this registry (and others) send `<idn:data>`, so
   `DomainInfoResult.IDN` was always empty. A registrar could not tell which
   IDN table a domain was registered under. Both spellings are now accepted and
   the U-label is exposed.

7. **Client transaction identifiers collided across sessions.** The identifier
   was a command prefix, a one-second timestamp and a per-client counter that
   always starts at zero, so two sessions opened in the same second, the normal
   case for a connection pool, emitted identical values. RFC 5730 section 2.5
   requires uniqueness because the identifier is how a lost transform is
   reconciled. A per-client random nonce now makes them unique across sessions,
   processes and hosts, within the 64-character clTRIDType limit. Relatedly,
   `AmbiguousTransformError` now carries the `ClientTRID` it was sent with,
   without which the reconciliation the SDK documents is impossible.

8. **Client-side rejections impersonated registry rejections.** All 72 local
   validation failures were built without a `Kind`, so `Error()` defaulted them
   to `epp_result` with an EPP result code and empty TRIDs. A registrar keying
   retry, alerting or reconciliation on "the registry answered" would
   misattribute a request that never left the process. They now carry
   `ErrorKindValidation`.

9. **Incomplete DNSSEC records were silently dropped.** A `dsData` record with
   an empty digest was removed from the request rather than rejected, so a
   caller with one malformed record got a command carrying no DNSSEC data and
   no indication of it, while believing the domain was signed. Zero `alg` and
   `digestType` were not checked at all and reached the registry as
   schema-invalid XML. RFC 5910 section 4 mandatory fields are now validated,
   and nothing is dropped.

Two further hardening changes accompany the fixes:

- **`<msgQ>` is now read from every response.** RFC 5730 section 2.9.2.3 lets a
  server report queued messages on any response; the SDK read it only on poll,
  so the signal was lost. `types.Response.MessageQueue` now carries it.
- **`tls.ca_only`** restricts trust to the configured CA instead of adding it to
  the host's public roots. The NIXI OT&E chain is signed by a private registry
  CA and is not publicly trusted, so `ca_only: true` is both safer and correct
  here. The default is unchanged.

Regression tests live in `test/integration/uat_regression_test.go`,
`test/integration/protocol_robustness_test.go`,
`test/integration/validation_parity_test.go`, `epp/trid_test.go` and
`extensions/secdns/request_test.go`. Each was verified to fail when its fix is
reverted.

## Not Covered, and Why

- **Completing a transfer (#16-19, #31-34)** needs a second registrar. Every
  transfer operation is encoded correctly and parsed by the registry, but a
  transfer cannot be completed against oneself. The `csc_egov_in_b` OT&E
  certificate and private key are present locally and its mTLS handshake
  succeeds; its EPP password is not held here (login returns 2200). One attempt
  was made and not repeated, to avoid any risk of locking the account. Request
  the `csc_egov_in_b` OT&E password from NIXI to finish these eight.
- **RGP restore (#38)** needs a domain past the add grace period. A domain
  deleted within that window is purged immediately rather than entering
  redemption, so the restore path cannot be reached inside a single session.
- **launch-1.0 (#41)** is advertised but no launch phase is open in OT&E.

## First Command

After creating a local ignored config with NIXI UAT details, the first safe external command should be:

```bash
go run ./cmd/epp-cli \
  -config configs/nixi-uat.yaml \
  -uat \
  -env ote \
  -connect-only \
  -capture-dir .uat-captures
```
