package main

import (
	"flag"
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	protocolHTTP1  = "http1"
	protocolHTTPS1 = "https1"
	defaultALPN    = "http/1.1"
)

type config struct {
	Protocol          string
	Addr              string
	Payload           int
	Connections       int
	Messages          int
	Timeout           time.Duration
	LatencySampleRate int
	WarmupMessages    int
	TLSVersion        string
	ALPN              string
	CipherSuites      string
	CipherSuiteIDs    []uint16
	ServerOnly        bool
}

func parseConfig(args []string) (config, error) {
	cfg := config{
		Protocol:    protocolHTTP1,
		Addr:        "127.0.0.1:0",
		Payload:     1024,
		Connections: 256,
		Messages:    100000,
		Timeout:     5 * time.Minute,
		TLSVersion:  tlsVersion13,
		ALPN:        defaultALPN,
	}
	fs := flag.NewFlagSet("fasthttp-bench", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&cfg.Protocol, "protocol", cfg.Protocol, "benchmark protocol")
	fs.StringVar(&cfg.Addr, "addr", cfg.Addr, "listen address")
	fs.IntVar(&cfg.Payload, "payload", cfg.Payload, "payload bytes")
	fs.IntVar(&cfg.Connections, "connections", cfg.Connections, "concurrent connections")
	fs.IntVar(&cfg.Messages, "messages", cfg.Messages, "messages per connection")
	fs.DurationVar(&cfg.Timeout, "timeout", cfg.Timeout, "overall timeout")
	fs.IntVar(&cfg.LatencySampleRate, "latency-sample-rate", cfg.LatencySampleRate, "latency sampling interval")
	fs.IntVar(&cfg.WarmupMessages, "warmup-messages", cfg.WarmupMessages, "messages per connection sent before timed measurement")
	fs.StringVar(&cfg.TLSVersion, "tls-version", cfg.TLSVersion, "TLS protocol version: 1.1, 1.2 or 1.3")
	fs.StringVar(&cfg.ALPN, "alpn", cfg.ALPN, "comma-separated TLS ALPN protocols")
	fs.StringVar(&cfg.CipherSuites, "cipher-suites", cfg.CipherSuites, "comma-separated TLS cipher suites")
	fs.BoolVar(&cfg.ServerOnly, "server-only", cfg.ServerOnly, "run only the benchmark server")
	if err := fs.Parse(args); err != nil {
		return config{}, err
	}
	return cfg, cfg.resolve()
}

func (c *config) resolve() error {
	version, err := normalizeTLSVersion(c.TLSVersion)
	if err != nil {
		return err
	}
	c.TLSVersion = version
	ids, names, err := parseCipherSuites(c.CipherSuites)
	if err != nil {
		return err
	}
	if len(ids) > 0 && c.TLSVersion == tlsVersion13 {
		return fmt.Errorf("%w: cipher suites are configurable only for TLS 1.1 and TLS 1.2", errInvalidConfig)
	}
	c.CipherSuiteIDs = ids
	c.CipherSuites = strings.Join(names, ",")
	return c.validate()
}

func (c config) validate() error {
	switch strings.TrimSpace(c.Protocol) {
	case protocolHTTP1, protocolHTTPS1:
	default:
		return fmt.Errorf("%w: %s", errUnsupportedProtocol, c.Protocol)
	}
	if strings.TrimSpace(c.Addr) == "" {
		return fmt.Errorf("%w: empty addr", errInvalidConfig)
	}
	if c.Payload <= 0 || c.Connections <= 0 || c.Messages <= 0 {
		return fmt.Errorf("%w: payload, connections and messages must be positive", errInvalidConfig)
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("%w: timeout must be positive", errInvalidConfig)
	}
	if c.LatencySampleRate < 0 {
		return fmt.Errorf("%w: latency-sample-rate must not be negative", errInvalidConfig)
	}
	if c.WarmupMessages < 0 {
		return fmt.Errorf("%w: warmup-messages must not be negative", errInvalidConfig)
	}
	if c.Protocol == protocolHTTP1 && len(c.CipherSuiteIDs) > 0 {
		return fmt.Errorf("%w: cipher suites require HTTPS", errInvalidConfig)
	}
	if c.ServerOnly && c.Protocol != protocolHTTP1 && c.Protocol != protocolHTTPS1 {
		return fmt.Errorf("%w: server-only requires http1 or https1 protocol", errInvalidConfig)
	}
	return nil
}
