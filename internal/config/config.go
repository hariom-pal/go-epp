package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config contains all settings required to open an EPP session.
type Config struct {
	Environment    string `yaml:"environment"`
	Server         ServerConfig
	Authentication AuthenticationConfig
	TLS            TLSConfig
	Timeout        TimeoutConfig
	Transport      TransportConfig
	Login          LoginConfig
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
	CertFile   string `yaml:"cert_file"`
	KeyFile    string `yaml:"key_file"`
	CAFile     string `yaml:"ca_file"`
	ServerName string `yaml:"server_name"`

	// CAOnly restricts trust to CAFile alone. By default the CA file is added
	// to the host's public roots, so any publicly trusted CA can still vouch
	// for the registry; a registrar issued a private registry CA usually wants
	// that CA to be the only one accepted. Requires CAFile.
	CAOnly bool `yaml:"ca_only"`

	InsecureSkipVerify bool `yaml:"insecure_skip_verify"`
	AllowInsecure      bool `yaml:"allow_insecure"`
}

// TimeoutConfig contains timeout settings in seconds.
type TimeoutConfig struct {
	Connect int `yaml:"connect"`
	Read    int `yaml:"read"`
	Write   int `yaml:"write"`
}

// TransportConfig contains wire-level safety limits.
type TransportConfig struct {
	MaxFrameSize int `yaml:"max_frame_size"`
}

// LoginConfig controls which services are advertised during login.
type LoginConfig struct {
	ObjectURIs                 []string `yaml:"object_uris"`
	RequireSupportedObjects    bool     `yaml:"require_supported_objects"`
	ExtensionURIs              []string `yaml:"extension_uris"`
	RequireSupportedExtensions bool     `yaml:"require_supported_extensions"`
}

// LoadFromFile reads and parses a YAML configuration file.
// This function is intended for use by CLI tools, examples, and integration tests
// that need to load configuration from a file path.
// The SDK itself should not use this function directly for client instantiation.
func LoadFromFile(path string) (*Config, error) {
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
