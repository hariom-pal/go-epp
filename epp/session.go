package epp

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/hariom-pal/go-epp/internal/config"
)

func Connect(cfg *config.Config) (*Client, error) {

	// Load client certificate
	cert, err := tls.LoadX509KeyPair(
		cfg.TLS.CertFile,
		cfg.TLS.KeyFile,
	)
	if err != nil {
		return nil, err
	}

	// Load system CA pool
	rootCAs, err := x509.SystemCertPool()
	if err != nil || rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Load custom Root CA (optional)
	if cfg.TLS.CAFile != "" {

		caCert, err := os.ReadFile(cfg.TLS.CAFile)
		if err != nil {
			return nil, err
		}

		if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
			return nil, fmt.Errorf("failed to load root CA")
		}
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            rootCAs,
		ServerName:         cfg.Server.Host,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: cfg.TLS.InsecureSkipVerify,
	}

	address := fmt.Sprintf("%s:%d",
		cfg.Server.Host,
		cfg.Server.Port,
	)

	conn, err := tls.Dial("tcp", address, tlsConfig)
	if err != nil {
		return nil, err
	}

	// Read Greeting immediately after connect
	greeting, err := ReadFrame(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	client := &Client{
		conn:     conn,
		config:   cfg,
		greeting: greeting,
	}

	return client, nil
}
