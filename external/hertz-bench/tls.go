package main

import (
	"crypto/tls"
	"fmt"
	"strings"

	"gnalloy.org/benchmarks/external/internal/benchtls"
)

func normalizeTLSVersion(value string) (string, error) {
	switch strings.TrimSpace(value) {
	case tlsVersion12:
		return tlsVersion12, nil
	case "", tlsVersion13:
		return tlsVersion13, nil
	case "1.1":
		return "", fmt.Errorf("hertz-bench: HTTP/2 over TLS requires TLS 1.2 or newer")
	default:
		return "", fmt.Errorf("hertz-bench: unsupported TLS version %q", value)
	}
}

func serverTLSConfig(cfg config) (*tls.Config, error) {
	certificate, err := benchtls.SelfSignedCertificate(benchtls.DefaultServerName)
	if err != nil {
		return nil, err
	}
	version := uint16(tls.VersionTLS13)
	if cfg.TLSVersion == tlsVersion12 {
		version = tls.VersionTLS12
	}
	return &tls.Config{
		Certificates: []tls.Certificate{certificate},
		MinVersion:   version,
		MaxVersion:   version,
		NextProtos:   []string{"h2"},
	}, nil
}
