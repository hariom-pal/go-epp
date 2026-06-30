package config
import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig
	Authentication AuthenticationConfig
	TLS            TLSConfig
	Timeout        TimeoutConfig
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
}

type AuthenticationConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type TLSConfig struct {
    CertFile string `yaml:"cert_file"`
    KeyFile  string `yaml:"key_file"`
    CAFile   string `yaml:"ca_file"`
    InsecureSkipVerify  bool   `yaml:"insecure_skip_verify"`
}


type TimeoutConfig struct {
	Connect int `yaml:"connect"`
	Read    int `yaml:"read"`
	Write   int `yaml:"write"`
}


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
