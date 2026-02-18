package httputil

import (
	"crypto/tls"
	"fmt"
	"os"
)

// ServerTLSConfig holds TLS configuration for an HTTP server.
type ServerTLSConfig struct {
	TLSCert string
	TLSKey  string
}

// TLSConfig creates a tls.Config from the ServerTLSConfig.
func (c *ServerTLSConfig) TLSConfig() (*tls.Config, error) {
	if c.TLSCert == "" || c.TLSKey == "" {
		return nil, fmt.Errorf("TLS cert and key must be specified")
	}
	cert, err := tls.LoadX509KeyPair(c.TLSCert, c.TLSKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS cert/key: %w", err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}

// CLIFlags returns the TLS-related CLI flags.
func (c *ServerTLSConfig) Check() error {
	if c.TLSCert != "" {
		if _, err := os.Stat(c.TLSCert); err != nil {
			return fmt.Errorf("TLS cert file not found: %w", err)
		}
	}
	if c.TLSKey != "" {
		if _, err := os.Stat(c.TLSKey); err != nil {
			return fmt.Errorf("TLS key file not found: %w", err)
		}
	}
	return nil
}
