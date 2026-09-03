package main

import (
	"crypto/tls"
	"testing"
)

func TestParseConfigHTTP1(t *testing.T) {
	cfg, err := parseConfig([]string{"-protocol", "http1", "-addr", "127.0.0.1:8080", "-payload", "128"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Protocol != protocolHTTP1 || cfg.Addr != "127.0.0.1:8080" || cfg.Payload != 128 {
		t.Fatalf("cfg=%+v", cfg)
	}
}

func TestParseConfigHTTPS1TLSVersions(t *testing.T) {
	tests := []struct {
		version string
		cipher  string
		wantID  uint16
	}{
		{version: tlsVersion11, cipher: "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA", wantID: tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA},
		{version: tlsVersion12, cipher: "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256", wantID: tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		{version: tlsVersion13},
	}
	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			cfg, err := parseConfig([]string{"-protocol", "https1", "-tls-version", tt.version, "-alpn", "http/1.1", "-cipher-suites", tt.cipher})
			if err != nil {
				t.Fatal(err)
			}
			if cfg.TLSVersion != tt.version {
				t.Fatalf("tlsVersion=%q, want %q", cfg.TLSVersion, tt.version)
			}
			if tt.wantID != 0 && (len(cfg.CipherSuiteIDs) != 1 || cfg.CipherSuiteIDs[0] != tt.wantID) {
				t.Fatalf("cipherSuiteIDs=%v, want [%d]", cfg.CipherSuiteIDs, tt.wantID)
			}
		})
	}
}

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

func TestParseConfigRejectsHTTPS1WithoutHTTP11ALPN(t *testing.T) {
	if _, err := parseConfig([]string{"-protocol", "https1", "-alpn", "h2"}); err == nil {
		t.Fatal("expected HTTP/1 ALPN rejection")
	}
}
