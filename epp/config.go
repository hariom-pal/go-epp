package epp

import "github.com/hariom-pal/go-epp/internal/config"

// Config contains all settings required to open an EPP session.
type Config = config.Config

// ServerConfig contains the EPP server endpoint.
type ServerConfig = config.ServerConfig

// AuthenticationConfig contains EPP login credentials.
type AuthenticationConfig = config.AuthenticationConfig

// TLSConfig contains client certificate and trust settings.
type TLSConfig = config.TLSConfig

// TimeoutConfig contains timeout settings in seconds.
type TimeoutConfig = config.TimeoutConfig

// TransportConfig contains wire-level safety limits.
type TransportConfig = config.TransportConfig

// LoginConfig controls which services are advertised during login.
type LoginConfig = config.LoginConfig

// LoadConfig reads and parses a YAML configuration file.
// This function is intended for use by CLI tools, examples, and integration tests
// that need to load configuration from a file path.
// The SDK itself should not use this function directly for client instantiation.
func LoadConfig(path string) (*Config, error) {
	return config.LoadFromFile(path)
}
