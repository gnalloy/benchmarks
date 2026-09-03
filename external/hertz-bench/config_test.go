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

func TestParseConfigHTTPS2TLSVersions(t *testing.T) {
	for _, version := range []string{tlsVersion12, tlsVersion13} {
		t.Run(version, func(t *testing.T) {
			cfg, err := parseConfig([]string{"-protocol", "https2", "-tls-version", version})
			if err != nil {
				t.Fatal(err)
			}
			if cfg.TLSVersion != version {
				t.Fatalf("tlsVersion=%q, want %q", cfg.TLSVersion, version)
			}
		})
	}
}

func TestParseConfigRejectsHTTP2TLS11(t *testing.T) {
	if _, err := parseConfig([]string{"-protocol", "https2", "-tls-version", "1.1"}); err == nil {
		t.Fatal("expected HTTP/2 TLS 1.1 rejection")
	}
}

func TestParseConfigRejectsHTTPS2WithoutH2ALPN(t *testing.T) {
	if _, err := parseConfig([]string{"-protocol", "https2", "-alpn", "http/1.1"}); err == nil {
		t.Fatal("expected HTTP/2 ALPN rejection")
	}
}
