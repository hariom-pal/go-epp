# Registry Compatibility

`go-epp` is registry independent. It implements EPP commands and extension XML models but does not hardcode registry policy.

## Registry Differences

Registries commonly differ in:

- advertised extension namespaces;
- required login service extensions;
- contact requirements;
- host object vs host attribute policy;
- IDN table policy;
- fee currency and command policy;
- secDNS DS-data vs key-data interface policy;
- RGP restore report validation;
- launch phases, application behavior, and claims notice rules.

## Recommended Integration Process

1. Connect to the registry OT&E environment.
2. Inspect the greeting and supported extension namespaces.
3. Login with the service extensions required by the registry.
4. Use the CLI to run smoke tests for check/info commands.
5. Validate create/update/delete XML in OT&E before production.
6. Capture registry-specific rules in application configuration or adapter code.

## Extension Guidelines

- Do not send an extension unless the registry supports it.
- Do not assume a registry supports all fields in an RFC-defined extension.
- Prefer registry documentation over inferred behavior.
- Keep registry-specific workarounds outside the core SDK.

## Compatibility Issue Reports

When reporting a compatibility issue, include:

- registry name and environment, if shareable;
- SDK version or commit;
- redacted command XML;
- redacted response XML;
- expected behavior;
- actual behavior;
- relevant registry documentation.
