package config

type Config struct {
	Server ServerConfig
	Authentication AuthenticationConfig
	TLS            TLSConfig
	Timeout        TimeoutConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type AuthenticationConfig struct {
	Username string
	Password string
}

type TLSConfig struct {
    CertFile string
    KeyFile  string
}


type TimeoutConfig struct {
	Connect int
	Read    int
	Write   int
}
