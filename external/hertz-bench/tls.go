package main

import (
	"crypto/tls"
	"fmt"
	"strings"

	"gnalloy.org/benchmarks/external/internal/benchtls"
)

func normalizeTLSVersion(value string) (string, error) {
	switch strings.TrimSpace(value) {
	case tlsVersion11:
		return tlsVersion11, nil
	case tlsVersion12:
		return tlsVersion12, nil
	case "", tlsVersion13:
		return tlsVersion13, nil
	default:
		return "", fmt.Errorf("hertz-bench: unsupported TLS version %q", value)
	}
}

func serverTLSConfig(cfg config) (*tls.Config, error) {
	certificates, err := benchtls.SelfSignedCertificates(benchtls.DefaultServerName, benchtls.CertificateKeyECDSA, benchtls.CertificateKeyRSA)
	if err != nil {
		return nil, err
	}
	version, err := cryptoTLSVersion(cfg.TLSVersion)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: certificates,
		MinVersion:   version,
		MaxVersion:   version,
		NextProtos:   []string{cfg.ALPN},
		CipherSuites: cfg.CipherSuiteIDs,
	}, nil
}

func cryptoTLSVersion(value string) (uint16, error) {
	switch value {
	case tlsVersion11:
		return tls.VersionTLS11, nil
	case tlsVersion12:
		return tls.VersionTLS12, nil
	case tlsVersion13:
		return tls.VersionTLS13, nil
	default:
		return 0, fmt.Errorf("hertz-bench: unsupported TLS version %q", value)
	}
}
