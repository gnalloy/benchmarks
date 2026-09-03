package main

import (
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"gnalloy.org/benchmarks/external/internal/benchtls"
)

const (
	protocolHTTP1  = "http1"
	protocolHTTPS1 = "https1"
	protocolHTTP2  = "http2"
	protocolHTTPS2 = "https2"
	tlsVersion11   = "1.1"
	tlsVersion12   = "1.2"
	tlsVersion13   = "1.3"
)

type config struct {
	Protocol       string
	Addr           string
	Payload        int
	Timeout        time.Duration
	TLSVersion     string
	ALPN           string
	CipherSuites   string
	CipherSuiteIDs []uint16
}

func parseConfig(args []string) (config, error) {
	cfg := config{
		Protocol:   protocolHTTP2,
		Addr:       "127.0.0.1:0",
		Payload:    1024,
		Timeout:    5 * time.Minute,
		TLSVersion: tlsVersion13,
		ALPN:       "h2",
	}
	fs := flag.NewFlagSet("hertz-bench", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&cfg.Protocol, "protocol", cfg.Protocol, "benchmark protocol")
	fs.StringVar(&cfg.Addr, "addr", cfg.Addr, "listen address")
	fs.IntVar(&cfg.Payload, "payload", cfg.Payload, "payload bytes")
	fs.DurationVar(&cfg.Timeout, "timeout", cfg.Timeout, "server lifetime")
	fs.StringVar(&cfg.TLSVersion, "tls-version", cfg.TLSVersion, "TLS protocol version: 1.1, 1.2 or 1.3")
	fs.StringVar(&cfg.ALPN, "alpn", cfg.ALPN, "TLS ALPN protocol")
	fs.StringVar(&cfg.CipherSuites, "cipher-suites", cfg.CipherSuites, "comma-separated TLS cipher suites")
	if err := fs.Parse(args); err != nil {
		return config{}, err
	}
	version, err := normalizeTLSVersion(cfg.TLSVersion)
	if err != nil {
		return config{}, err
	}
	cfg.TLSVersion = version
	ids, names, err := benchtls.ParseCipherSuites(cfg.CipherSuites)
	if err != nil {
		return config{}, err
	}
	if len(ids) > 0 && cfg.TLSVersion == tlsVersion13 {
		return config{}, fmt.Errorf("hertz-bench: cipher suites are configurable only for TLS 1.1 and TLS 1.2")
	}
	cfg.CipherSuiteIDs = ids
	cfg.CipherSuites = strings.Join(names, ",")
	return cfg, cfg.validate()
}

func (c config) validate() error {
	switch c.Protocol {
	case protocolHTTP1, protocolHTTPS1, protocolHTTP2, protocolHTTPS2:
	default:
		return fmt.Errorf("hertz-bench: unsupported protocol %q", c.Protocol)
	}
	if strings.TrimSpace(c.Addr) == "" {
		return fmt.Errorf("hertz-bench: empty addr")
	}
	if c.Payload <= 0 {
		return fmt.Errorf("hertz-bench: payload must be positive")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("hertz-bench: timeout must be positive")
	}
	if c.Protocol == protocolHTTPS2 && c.ALPN != "h2" {
		return fmt.Errorf("hertz-bench: HTTP/2 over TLS requires ALPN h2")
	}
	if c.Protocol == protocolHTTPS2 && c.TLSVersion == tlsVersion11 {
		return fmt.Errorf("hertz-bench: HTTP/2 over TLS requires TLS 1.2 or newer")
	}
	if c.Protocol == protocolHTTPS1 && c.ALPN != "http/1.1" {
		return fmt.Errorf("hertz-bench: HTTP/1 over TLS requires ALPN http/1.1")
	}
	if (c.Protocol == protocolHTTP1 || c.Protocol == protocolHTTP2) && len(c.CipherSuiteIDs) > 0 {
		return fmt.Errorf("hertz-bench: cipher suites require TLS")
	}
	return nil
}
