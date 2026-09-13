# Production Safety

This SDK is a generic EPP implementation. Registry-specific policy must live outside the core package.

## Transport

Use `epp.Connect` or `epp.ConnectContext`. Both open TLS connections and read the server greeting before returning a client.

`epp.NewClient` is retained for compatibility, but it now delegates to the same TLS-only connection path. It must not silently fall back to plain TCP.

Client certificate and key paths are required. A custom CA file can be configured when a registry uses a private or OT&E CA.

`insecure_skip_verify` is blocked unless `allow_insecure` is also set. That opt-in is only for local mock servers and development tests. Production registrar configuration should keep both values false unless the registry explicitly requires a private CA, in which case configure `ca_file`.

## Frame Size

EPP uses a four-byte RFC 5734 frame length. The reader validates the frame length before allocating memory.

The default maximum frame size is 16 MiB. Override `transport.max_frame_size` only when a registry documents a larger legitimate response size.

Malformed frames are classified as framing/protocol errors. Oversized frames are rejected before payload allocation.

## Context

The SDK keeps existing compatibility methods such as `DomainCheck` and `DomainCreate`.

Context-aware methods are available for registrar services, including:

- `ConnectContext`
- `ExecuteContext`
- `LoginContext`
- `LogoutContext`
- `HelloContext`
- `PollContext`
- `Domain*Context`
- `Contact*Context`
- `Host*Context`

Context cancellation and deadlines interrupt network reads/writes through connection deadlines.

## Command Serialization

One EPP TCP/TLS session is request/response ordered. A single `Client` serializes the complete write-request/read-response transaction. This prevents concurrent goroutines from interleaving commands on one connection.

Use multiple clients only when the registry allows multiple sessions and your registrar application can track each session separately.

## Session Lifecycle

The intended lifecycle is:

```text
connect -> greeting -> login -> commands -> logout -> close
```

Executing a command after `Close` returns a session-classified error.

## Login Services And Extensions

The SDK starts with domain, contact, and host object URIs by default, then compares them with the greeting-advertised object URIs. Extension URIs are not advertised blindly.

Configure desired objects in `login.object_uris` and desired extensions in `login.extension_uris`. The SDK compares requested extensions with the greeting-advertised extension URIs. With `require_supported_objects: true` or `require_supported_extensions: true`, unsupported requested services are rejected before login XML is sent.

Registry adapters should choose extension URIs based on registry documentation and the live greeting.

## Reconnect And Ambiguous Transforms

The SDK does not automatically replay EPP commands.

Read/query operations include check, info, hello, and poll. These can usually be retried by the registrar application with normal idempotency care.

Transform operations include create, update, delete, renew, and transfer request/approve/reject/cancel. If the connection fails after a transform request is written but before the response is read, the outcome is ambiguous. The SDK returns `AmbiguousTransformError`.

Registrar services should reconcile ambiguous outcomes using registry state, for example:

- object `info`
- object `check`
- poll messages
- local transaction records
- registry support guidance

## Error Classification

Non-EPP result errors are returned as `*epp.SDKError` with a stable kind:

- configuration
- TLS
- transport
- timeout
- cancellation
- framing
- protocol
- session

EPP result errors remain `*epp.Error` and preserve result code, client TRID, and server TRID.

RFC 5730 permits multiple result elements. Public response structs preserve all results in `Response.Results` while keeping `ResultCode` and `ResultMsg` as first-result compatibility fields. EPP result errors also expose all parsed results.

Ambiguous transform failures are returned as `*epp.AmbiguousTransformError`.

## Sensitive Data

SDK logging events never include raw XML, login passwords, authInfo values, private keys, or certificate material.

Registrar applications must also avoid logging request structs directly because those can contain authInfo or credentials.

## Logging Hooks

Install an optional logger with `client.SetLogger(logger)`. The logger receives redaction-safe events for connection close, login/logout, command type, duration, result metadata when available, and errors.

The SDK does not depend on any logging framework.

## Generic SDK Boundary

Keep this SDK registry-neutral.

The intended production architecture for NIXI or any other registry is:

```text
Registrar Service -> Registry Adapter -> go-epp
```

The adapter should handle:

- registry-specific extensions
- policy
- IDN tables
- premium and reserved behavior
- registry-specific contact requirements
- OT&E and certification differences

Do not hardcode NIXI, `.in`, `.भारत`, or any registry policy into the generic EPP package.
