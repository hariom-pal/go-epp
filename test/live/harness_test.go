//go:build live

// Package live runs the go-epp SDK against a real EPP registry.
//
// These tests are excluded from the default build because they require a
// live registry, client certificates, and credentials. Run them with:
//
//	go test -tags live ./test/live/ -v
//
// The target is read from EPP_LIVE_CONFIG (default configs/nixi-uat.yaml) and
// must be an OT&E/UAT/test environment; the harness refuses to run transform
// commands against anything it cannot positively classify as non-production.
package live

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

const defaultConfigPath = "../../configs/nixi-uat.yaml"

var (
	runTagOnce sync.Once
	runTag     string
)

// tag returns a short identifier unique to this test binary run, so repeated
// runs never collide on registry object identifiers.
func tag() string {
	runTagOnce.Do(func() {
		runTag = strconv.FormatInt(time.Now().Unix()%1000000, 36)
	})
	return runTag
}

// contactID builds a run-scoped contact ID inside the RFC 5733 clIDType
// length limit of 16 characters.
func contactID(suffix string) string { return "c" + tag() + suffix }

// domainName builds a run-scoped domain name.
func domainName(suffix string) string { return "uat" + tag() + suffix + ".in" }

// hostName builds a run-scoped host name under the given domain.
func hostName(prefix, domain string) string { return prefix + "." + domain }

func configPath() string {
	if path := strings.TrimSpace(os.Getenv("EPP_LIVE_CONFIG")); path != "" {
		return path
	}
	return defaultConfigPath
}

// loadConfig reads the live configuration and enforces the safety gate: the
// environment must be positively identifiable as OT&E/UAT/test.
func loadConfig(t *testing.T) *epp.Config {
	t.Helper()

	path := configPath()
	cfg, err := epp.LoadConfig(path)
	if err != nil {
		t.Skipf("live config %s unavailable: %v", path, err)
	}

	environment := strings.ToLower(strings.TrimSpace(cfg.Environment))
	safe := false
	for _, marker := range []string{"ote", "ot&e", "uat", "test", "sandbox"} {
		if strings.Contains(environment, marker) {
			safe = true
			break
		}
	}
	if !safe {
		t.Fatalf("REFUSED - target is not confirmed OT&E/UAT: environment=%q host=%q", cfg.Environment, cfg.Server.Host)
	}

	// Certificate paths in the config are relative to the repository root.
	root := filepath.Dir(filepath.Dir(filepath.Dir(path)))
	if !filepath.IsAbs(path) {
		root = filepath.Dir(filepath.Dir(path))
	}
	cfg.TLS.CertFile = resolve(root, cfg.TLS.CertFile)
	cfg.TLS.KeyFile = resolve(root, cfg.TLS.KeyFile)
	cfg.TLS.CAFile = resolve(root, cfg.TLS.CAFile)

	return cfg
}

func resolve(root, path string) string {
	if path == "" || filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

// connect opens a session without logging in.
func connect(t *testing.T) *epp.Client {
	t.Helper()

	client, err := epp.Connect(loadConfig(t))
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })
	return client
}

// session opens a session and logs in, logging out on cleanup.
func session(t *testing.T) *epp.Client {
	t.Helper()

	client := connect(t)
	if err := client.Login(); err != nil {
		t.Fatalf("login failed: %v", err)
	}
	t.Cleanup(func() { _ = client.Logout() })
	return client
}

// sharedSession is a single logged-in session reused by read-only tests to
// keep the number of registry logins low.
var (
	sharedOnce   sync.Once
	sharedClient *epp.Client
	sharedErr    error
)

func shared(t *testing.T) *epp.Client {
	t.Helper()

	sharedOnce.Do(func() {
		cfg, err := epp.LoadConfig(configPath())
		if err != nil {
			sharedErr = err
			return
		}
		root := filepath.Dir(filepath.Dir(configPath()))
		cfg.TLS.CertFile = resolve(root, cfg.TLS.CertFile)
		cfg.TLS.KeyFile = resolve(root, cfg.TLS.KeyFile)
		cfg.TLS.CAFile = resolve(root, cfg.TLS.CAFile)
		sharedClient, sharedErr = epp.Connect(cfg)
		if sharedErr != nil {
			return
		}
		sharedErr = sharedClient.Login()
	})
	if sharedErr != nil {
		t.Fatalf("shared session unavailable: %v", sharedErr)
	}
	return sharedClient
}

// requireResultCode asserts that err carries a specific EPP result code.
func requireResultCode(t *testing.T, err error, code int, context string) {
	t.Helper()
	_ = requireResultCodeErr(t, err, code, context)
}

// requireResultCodeErr is requireResultCode for callers that need the error.
func requireResultCodeErr(t *testing.T, err error, code int, context string) *epp.Error {
	t.Helper()

	if err == nil {
		t.Fatalf("%s: expected EPP result %d, got success", context, code)
	}
	eppErr := asEPPError(t, err, context)
	if eppErr.Code != code {
		t.Fatalf("%s: expected EPP result %d, got %d (%s)", context, code, eppErr.Code, eppErr.Message)
	}
	return eppErr
}

func asEPPError(t *testing.T, err error, context string) *epp.Error {
	t.Helper()

	var eppErr *epp.Error
	if !errors.As(err, &eppErr) {
		t.Fatalf("%s: expected *epp.Error, got %T: %v", context, err, err)
	}
	return eppErr
}

// asSDKError reports whether err is an *epp.SDKError and stores it in target.
func asSDKError(err error, target **epp.SDKError) bool {
	return errors.As(err, target)
}

// domainCheckRequest builds a minimal single-name domain check.
func domainCheckRequest(names ...string) types.DomainCheckRequest {
	return types.DomainCheckRequest{Domains: names}
}

// devanagariTag renders the run tag using Devanagari letters, so an IDN test
// label stays within one script. Registries reject mixed-script labels.
func devanagariTag() string {
	letters := []rune("कखगघचछजझटठडढणतथदधनपफबभमयरलवशषसह")
	var builder strings.Builder
	for _, char := range tag() {
		index := strings.IndexRune("0123456789abcdefghijklmnopqrstuvwxyz", char)
		if index < 0 {
			continue
		}
		builder.WriteRune(letters[index%len(letters)])
	}
	return builder.String()
}

// errResultf builds a formatted error for use inside goroutines.
func errResultf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

// loadConfigNoT loads the live configuration outside a *testing.T, for use in
// goroutines that open their own sessions.
func loadConfigNoT() *epp.Config {
	cfg, err := epp.LoadConfig(configPath())
	if err != nil {
		return nil
	}
	root := filepath.Dir(filepath.Dir(configPath()))
	cfg.TLS.CertFile = resolve(root, cfg.TLS.CertFile)
	cfg.TLS.KeyFile = resolve(root, cfg.TLS.KeyFile)
	cfg.TLS.CAFile = resolve(root, cfg.TLS.CAFile)
	return cfg
}
