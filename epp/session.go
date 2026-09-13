package epp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/hariom-pal/go-epp/internal/config"
)

// Connect opens a TLS EPP session and reads the server greeting.
func Connect(cfg *config.Config) (*Client, error) {
	return ConnectContext(context.Background(), cfg)
}

// ConnectContext opens a TLS EPP session and reads the server greeting.
func ConnectContext(ctx context.Context, cfg *config.Config) (*Client, error) {
	if cfg == nil {
		return nil, newSDKError(ErrorKindConfiguration, "config is required", nil)
	}
	if err := validateTLSConfig(cfg); err != nil {
		return nil, err
	}

	// Load client certificate
	cert, err := tls.LoadX509KeyPair(
		cfg.TLS.CertFile,
		cfg.TLS.KeyFile,
	)
	if err != nil {
		return nil, newSDKError(ErrorKindTLS, "failed to load client certificate", err)
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
			return nil, newSDKError(ErrorKindTLS, "failed to read root CA", err)
		}

		if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
			return nil, newSDKError(ErrorKindTLS, "failed to load root CA", nil)
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

	dialer := &net.Dialer{
		Timeout: timeoutDuration(cfg.Timeout.Connect),
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", address, tlsConfig)
	if err != nil {
		return nil, newSDKError(ErrorKindTLS, "TLS connection failed", err)
	}

	// Read Greeting immediately after connect
	if err := setContextDeadline(ctx, conn.SetReadDeadline, cfg.Timeout.Read); err != nil {
		conn.Close()
		return nil, err
	}
	greeting, err := ReadFrameWithMax(conn, maxFrameSize(cfg))
	if err != nil {
		conn.Close()
		return nil, err
	}
	_ = conn.SetReadDeadline(time.Time{})

	greetingInfo, err := parseGreetingXML(greeting)
	if err != nil {
		conn.Close()
		return nil, newSDKError(ErrorKindProtocol, "failed to parse EPP greeting", err)
	}

	client := &Client{
		conn:         conn,
		config:       cfg,
		greeting:     greeting,
		greetingInfo: greetingInfo,
	}
	client.emit(Event{Type: EventConnect})

	return client, nil
}

func validateTLSConfig(cfg *config.Config) error {
	if cfg.Server.Host == "" || cfg.Server.Port == 0 {
		return newSDKError(ErrorKindConfiguration, "server host and port are required", nil)
	}
	if cfg.TLS.CertFile == "" || cfg.TLS.KeyFile == "" {
		return newSDKError(ErrorKindConfiguration, "client certificate and key are required", nil)
	}
	if cfg.TLS.InsecureSkipVerify && !cfg.TLS.AllowInsecure {
		return newSDKError(
			ErrorKindConfiguration,
			"insecure TLS verification requires explicit allow_insecure opt-in for local testing",
			nil,
		)
	}
	return nil
}

func timeoutDuration(seconds int) time.Duration {
	if seconds <= 0 {
		return 0
	}
	return time.Duration(seconds) * time.Second
}

func maxFrameSize(cfg *config.Config) int {
	if cfg == nil || cfg.Transport.MaxFrameSize <= 0 {
		return DefaultMaxFrameSize
	}
	return cfg.Transport.MaxFrameSize
}

func setContextDeadline(ctx context.Context, setter func(time.Time) error, seconds int) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return contextError(err)
	}

	var deadline time.Time
	if contextDeadline, ok := ctx.Deadline(); ok {
		deadline = contextDeadline
	}
	if seconds > 0 {
		timeoutDeadline := time.Now().Add(time.Duration(seconds) * time.Second)
		if deadline.IsZero() || timeoutDeadline.Before(deadline) {
			deadline = timeoutDeadline
		}
	}
	if err := setter(deadline); err != nil {
		return classifyTransportError(err)
	}
	return nil
}

func contextError(err error) error {
	switch {
	case errors.Is(err, context.Canceled):
		return newSDKError(ErrorKindCancellation, "operation canceled", err)
	case errors.Is(err, context.DeadlineExceeded):
		return newSDKError(ErrorKindTimeout, "operation deadline exceeded", err)
	default:
		return newSDKError(ErrorKindCancellation, "context error", err)
	}
}
