package main

import "testing"

func TestParseConfigHTTP2(t *testing.T) {
	cfg, err := parseConfig([]string{"-protocol", "http2", "-addr", "127.0.0.1:8080", "-payload", "128"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Protocol != protocolHTTP2 || cfg.Addr != "127.0.0.1:8080" || cfg.Payload != 128 {
		t.Fatalf("cfg=%+v", cfg)
	}
}

func TestParseConfigRejectsUnsupportedProtocol(t *testing.T) {
	if _, err := parseConfig([]string{"-protocol", "http3"}); err == nil {
		t.Fatal("expected unsupported protocol error")
	}
}
