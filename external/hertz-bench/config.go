package main

import (
	"flag"
	"fmt"
	"io"
	"strings"
	"time"
)

const protocolHTTP2 = "http2"

type config struct {
	Protocol string
	Addr     string
	Payload  int
	Timeout  time.Duration
}

func parseConfig(args []string) (config, error) {
	cfg := config{
		Protocol: protocolHTTP2,
		Addr:     "127.0.0.1:0",
		Payload:  1024,
		Timeout:  5 * time.Minute,
	}
	fs := flag.NewFlagSet("hertz-bench", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&cfg.Protocol, "protocol", cfg.Protocol, "benchmark protocol")
	fs.StringVar(&cfg.Addr, "addr", cfg.Addr, "listen address")
	fs.IntVar(&cfg.Payload, "payload", cfg.Payload, "payload bytes")
	fs.DurationVar(&cfg.Timeout, "timeout", cfg.Timeout, "server lifetime")
	if err := fs.Parse(args); err != nil {
		return config{}, err
	}
	return cfg, cfg.validate()
}

func (c config) validate() error {
	if c.Protocol != protocolHTTP2 {
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
	return nil
}
