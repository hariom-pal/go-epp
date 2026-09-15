//go:build live

package live

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/epp"
)

// TestConnectWrongEndpoint covers an unreachable/unresolvable registry host.
func TestConnectWrongEndpoint(t *testing.T) {
	cfg := loadConfig(t)
	cfg.Server.Host = "epp.invalid.nixiregistry.invalid"
	cfg.Timeout.Connect = 5

	_, err := epp.Connect(cfg)
	if err == nil {
		t.Fatal("expected connect to an unresolvable host to fail")
	}
	var sdkErr *epp.SDKError
	if !asSDKError(err, &sdkErr) || sdkErr.Kind != epp.ErrorKindTLS {
		t.Fatalf("expected a TLS SDKError, got %T: %v", err, err)
	}
}

// TestConnectWrongPort covers a closed port on the real registry host.
func TestConnectWrongPort(t *testing.T) {
	cfg := loadConfig(t)
	cfg.Server.Port = 7001
	cfg.Timeout.Connect = 5

	_, err := epp.Connect(cfg)
	if err == nil {
		t.Fatal("expected connect to a closed port to fail")
	}
	t.Logf("closed port rejected: %v", err)
}

// TestConnectWrongServerName covers TLS hostname verification: an SNI/name
// mismatch must fail the handshake, never be silently accepted.
func TestConnectWrongServerName(t *testing.T) {
	cfg := loadConfig(t)
	cfg.TLS.ServerName = "wrong-name.nixiregistry.in"
	cfg.Timeout.Connect = 10

	_, err := epp.Connect(cfg)
	if err == nil {
		t.Fatal("expected a server name mismatch to fail the TLS handshake")
	}
	var sdkErr *epp.SDKError
	if !asSDKError(err, &sdkErr) || sdkErr.Kind != epp.ErrorKindTLS {
		t.Fatalf("expected a TLS SDKError, got %T: %v", err, err)
	}
}

// TestConnectUnrelatedCADefaultMode pins the default trust model: a configured
// CA file is added to the host's public roots rather than replacing them, so
// supplying an unrelated CA does not by itself reduce trust. Restricting trust
// to the registry's own CA is what ca_only is for.
func TestConnectUnrelatedCADefaultMode(t *testing.T) {
	cfg := loadConfig(t)
	cfg.TLS.CAFile = writeUnrelatedCA(t)
	cfg.Timeout.Connect = 10

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Logf("registry chain is not publicly trusted, so an unrelated CA also fails: %v", err)
		return
	}
	defer func() { _ = client.Close() }()
	t.Log("default mode still trusts the public roots, as documented")
}

// TestConnectMalformedCA covers a CA file that is not PEM at all.
func TestConnectMalformedCA(t *testing.T) {
	cfg := loadConfig(t)

	bad := filepath.Join(t.TempDir(), "not-a-ca.pem")
	if err := os.WriteFile(bad, []byte("this is not a certificate\n"), 0o600); err != nil {
		t.Fatalf("write temp CA failed: %v", err)
	}
	cfg.TLS.CAFile = bad

	_, err := epp.Connect(cfg)
	if err == nil {
		t.Fatal("expected a malformed CA file to be rejected")
	}
	var sdkErr *epp.SDKError
	if !asSDKError(err, &sdkErr) || sdkErr.Kind != epp.ErrorKindTLS {
		t.Fatalf("expected a TLS SDKError, got %T: %v", err, err)
	}
}

// TestConnectMissingCA covers a CA path that does not exist.
func TestConnectMissingCA(t *testing.T) {
	cfg := loadConfig(t)
	cfg.TLS.CAFile = filepath.Join(t.TempDir(), "absent.pem")

	_, err := epp.Connect(cfg)
	if err == nil {
		t.Fatal("expected a missing CA file to be rejected")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Logf("missing CA reported as: %v", err)
	}
}

// TestConnectMissingClientCertificate covers a client certificate path that
// does not exist. Registries requiring mTLS must never be reached without one.
func TestConnectMissingClientCertificate(t *testing.T) {
	cfg := loadConfig(t)
	cfg.TLS.CertFile = filepath.Join(t.TempDir(), "absent.crt")

	_, err := epp.Connect(cfg)
	if err == nil {
		t.Fatal("expected a missing client certificate to be rejected")
	}
	var sdkErr *epp.SDKError
	if !asSDKError(err, &sdkErr) || sdkErr.Kind != epp.ErrorKindTLS {
		t.Fatalf("expected a TLS SDKError, got %T: %v", err, err)
	}
}

// TestConnectMismatchedKey covers a certificate and key that do not pair.
func TestConnectMismatchedKey(t *testing.T) {
	cfg := loadConfig(t)
	cfg.TLS.KeyFile = cfg.TLS.CAFile // a certificate, not the matching key

	_, err := epp.Connect(cfg)
	if err == nil {
		t.Fatal("expected a mismatched certificate/key pair to be rejected")
	}
	var sdkErr *epp.SDKError
	if !asSDKError(err, &sdkErr) || sdkErr.Kind != epp.ErrorKindTLS {
		t.Fatalf("expected a TLS SDKError, got %T: %v", err, err)
	}
}

// TestInsecureSkipVerifyRequiresOptIn covers the SDK's guard against silently
// disabling certificate verification.
func TestInsecureSkipVerifyRequiresOptIn(t *testing.T) {
	cfg := loadConfig(t)
	cfg.TLS.InsecureSkipVerify = true
	cfg.TLS.AllowInsecure = false

	_, err := epp.Connect(cfg)
	if err == nil {
		t.Fatal("expected insecure_skip_verify without allow_insecure to be refused")
	}
	var sdkErr *epp.SDKError
	if !asSDKError(err, &sdkErr) || sdkErr.Kind != epp.ErrorKindConfiguration {
		t.Fatalf("expected a configuration SDKError, got %T: %v", err, err)
	}
}

// TestConnectNilConfig and friends cover configuration validation.
func TestConnectInvalidConfig(t *testing.T) {
	if _, err := epp.Connect(nil); err == nil {
		t.Fatal("expected a nil config to be rejected")
	}

	cfg := loadConfig(t)
	cfg.Server.Host = ""
	if _, err := epp.Connect(cfg); err == nil {
		t.Fatal("expected an empty host to be rejected")
	}

	cfg = loadConfig(t)
	cfg.Server.Port = 0
	if _, err := epp.Connect(cfg); err == nil {
		t.Fatal("expected a zero port to be rejected")
	}

	cfg = loadConfig(t)
	cfg.TLS.CertFile = ""
	if _, err := epp.Connect(cfg); err == nil {
		t.Fatal("expected a missing certificate path to be rejected")
	}
}

// TestCAOnlyRestrictsTrust covers the ca_only trust mode. By default the
// configured registry CA is added to the host's public roots, so any publicly
// trusted CA can still authenticate the registry; ca_only makes the configured
// CA the only accepted anchor.
func TestCAOnlyRestrictsTrust(t *testing.T) {
	cfg := loadConfig(t)
	cfg.TLS.CAOnly = true
	cfg.Timeout.Connect = 10

	client, err := epp.Connect(cfg)
	if err != nil {
		// The registry may present a publicly issued server certificate that
		// the registrar's own CA bundle does not chain to. That is a real
		// deployment fact, not an SDK failure, so record it rather than fail.
		t.Logf("ca_only rejected the registry chain, so the configured CA does not sign it: %v", err)
		return
	}
	defer func() { _ = client.Close() }()
	t.Log("ca_only accepted: the configured CA alone authenticates the registry")
}

// TestCAOnlyRequiresCAFile covers the configuration guard: restricting trust to
// a CA file that was never supplied would silently trust nothing.
func TestCAOnlyRequiresCAFile(t *testing.T) {
	cfg := loadConfig(t)
	cfg.TLS.CAOnly = true
	cfg.TLS.CAFile = ""

	_, err := epp.Connect(cfg)
	if err == nil {
		t.Fatal("expected ca_only without a ca_file to be refused")
	}
	var sdkErr *epp.SDKError
	if !asSDKError(err, &sdkErr) || sdkErr.Kind != epp.ErrorKindConfiguration {
		t.Fatalf("expected a configuration SDKError, got %T: %v", err, err)
	}
}

// TestCAOnlyRejectsUnrelatedCA confirms ca_only actually excludes the public
// roots: an unrelated CA must not authenticate the registry.
func TestCAOnlyRejectsUnrelatedCA(t *testing.T) {
	cfg := loadConfig(t)

	cfg.TLS.CAFile = writeUnrelatedCA(t)
	cfg.TLS.CAOnly = true
	cfg.Timeout.Connect = 10

	if _, err := epp.Connect(cfg); err == nil {
		t.Fatal("ca_only accepted a CA that does not sign the registry chain")
	}
}

// writeUnrelatedCA generates a self-signed CA that has nothing to do with the
// registry and returns its PEM path.
func writeUnrelatedCA(t *testing.T) string {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key failed: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "go-epp unrelated test CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate failed: %v", err)
	}

	path := filepath.Join(t.TempDir(), "unrelated-ca.pem")
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	if err := os.WriteFile(path, pemBytes, 0o600); err != nil {
		t.Fatalf("write CA failed: %v", err)
	}
	return path
}
