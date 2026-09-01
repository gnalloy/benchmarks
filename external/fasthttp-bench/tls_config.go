package main

import (
	"crypto/tls"
	"fmt"
	"strings"

	"gnalloy.org/benchmarks/external/internal/benchtls"
)

const (
	tlsVersion11 = "1.1"
	tlsVersion12 = "1.2"
	tlsVersion13 = "1.3"
)

func serverTLSConfig(cfg config) (*tls.Config, error) {
	certs, err := benchtls.SelfSignedCertificates(benchtls.DefaultServerName, benchtls.CertificateKeyECDSA, benchtls.CertificateKeyRSA)
	if err != nil {
		return nil, err
	}
	version, err := cryptoTLSVersion(cfg.TLSVersion)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: certs,
		MinVersion:   version,
		MaxVersion:   version,
		NextProtos:   alpnProtocols(cfg.ALPN),
		CipherSuites: cfg.CipherSuiteIDs,
	}, nil
}

func clientTLSConfig(cfg config) (*tls.Config, error) {
	version, err := cryptoTLSVersion(cfg.TLSVersion)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		ServerName:         benchtls.DefaultServerName,
		InsecureSkipVerify: true,
		MinVersion:         version,
		MaxVersion:         version,
		NextProtos:         alpnProtocols(cfg.ALPN),
		CipherSuites:       cfg.CipherSuiteIDs,
	}, nil
}

func normalizeTLSVersion(value string) (string, error) {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "", "1.3", "tls1.3", "tls13", "tlsv1.3":
		return tlsVersion13, nil
	case "1.2", "tls1.2", "tls12", "tlsv1.2":
		return tlsVersion12, nil
	case "1.1", "tls1.1", "tls11", "tlsv1.1":
		return tlsVersion11, nil
	default:
		return "", fmt.Errorf("%w: unsupported TLS version %q", errInvalidConfig, value)
	}
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
		return 0, fmt.Errorf("%w: unsupported TLS version %q", errInvalidConfig, value)
	}
}

func parseCipherSuites(value string) ([]uint16, []string, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil, nil
	}
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ';' || r == ':'
	})
	ids := make([]uint16, 0, len(parts))
	names := make([]string, 0, len(parts))
	for _, part := range parts {
		name := strings.TrimSpace(part)
		if name == "" {
			continue
		}
		id, canonical, err := lookupCipherSuite(name)
		if err != nil {
			return nil, nil, err
		}
		ids = append(ids, id)
		names = append(names, canonical)
	}
	return ids, names, nil
}

func lookupCipherSuite(name string) (uint16, string, error) {
	switch strings.ToUpper(strings.TrimSpace(name)) {
	case "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":
		return tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA, "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA", nil
	case "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":
		return tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256", nil
	case "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256":
		return tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256", nil
	default:
		return 0, "", fmt.Errorf("%w: unsupported cipher suite %q", errInvalidConfig, name)
	}
}

func alpnProtocols(value string) []string {
	parts := strings.Split(value, ",")
	protocols := make([]string, 0, len(parts))
	for _, part := range parts {
		protocol := strings.TrimSpace(part)
		if protocol != "" {
			protocols = append(protocols, protocol)
		}
	}
	return protocols
}
