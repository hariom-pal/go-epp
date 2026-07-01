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

// LoadConfig reads and parses a YAML configuration file.
func LoadConfig(path string) (*Config, error) {
	return config.Load(path)
}
