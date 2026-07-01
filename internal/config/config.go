package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config contains all settings required to open an EPP session.
type Config struct {
	Server         ServerConfig
	Authentication AuthenticationConfig
	TLS            TLSConfig
	Timeout        TimeoutConfig
}

// ServerConfig contains the EPP server endpoint.
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// AuthenticationConfig contains EPP login credentials.
type AuthenticationConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// TLSConfig contains client certificate and trust settings.
type TLSConfig struct {
	CertFile           string `yaml:"cert_file"`
	KeyFile            string `yaml:"key_file"`
	CAFile             string `yaml:"ca_file"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
}

// TimeoutConfig contains timeout settings in seconds.
type TimeoutConfig struct {
	Connect int `yaml:"connect"`
	Read    int `yaml:"read"`
	Write   int `yaml:"write"`
}

// Load reads and parses a YAML configuration file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
