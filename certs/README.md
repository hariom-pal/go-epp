# NIXI EPP Certificates

Place your registry-issued TLS certificates in this directory.

Example:

client.crt
client.key
root-ca.crt

These files are intentionally excluded from Git.

Do NOT commit:
- Private keys
- Client certificates
- Root CA certificates issued by the registry

Each developer or deployment environment must provide its own certificates.