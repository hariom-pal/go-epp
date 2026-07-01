# Registry Compatibility

`go-epp` is designed as a registry-independent SDK. It implements EPP protocol and extension models, while leaving registry policy choices to callers.

## Compatibility Model

Registries differ in supported extensions, command policies, grace-period windows, fee behavior, IDN tables, launch phases, and validation rules. The SDK therefore keeps request models generic and avoids hardcoding registry-specific behavior.

## Extension Use

Optional extension requests are attached only when the caller supplies extension data:

- Fee extension data is optional for supported domain operations.
- secDNS data is optional for domain create, update, and info parsing.
- RGP data is optional for domain info parsing and domain update restore operations.

If a registry does not advertise or support an extension namespace, callers should not send that extension.

## Testing Against Registries

Use the CLI against OT&E environments before production use. Confirm:

- Advertised extension namespaces in the greeting.
- Required login service extensions.
- Domain/contact/host object policy.
- Required contacts and name server formats.
- Extension-specific policy such as fee commands, DNSSEC interface choice, and RGP restore report requirements.

## Reporting Compatibility Issues

When reporting a registry compatibility issue, include:

- Registry name and environment, if shareable.
- Command XML with credentials removed.
- Response XML with sensitive values removed.
- Expected behavior and actual behavior.
- Any registry documentation that describes the required behavior.
