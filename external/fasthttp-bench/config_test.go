package main

import (
	"context"
	"crypto/tls"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestParseConfigAcceptsHTTPS1TLS12Cipher(t *testing.T) {
	cfg, err := parseConfig([]string{
		"-protocol", "https1",
		"-payload", "128",
		"-connections", "2",
		"-messages", "3",
		"-timeout", "1s",
		"-cpuprofile", "cpu.pprof",
		"-tls-version", "1.2",
		"-cipher-suites", "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Protocol != protocolHTTPS1 || cfg.TLSVersion != tlsVersion12 || cfg.CPUProfile != "cpu.pprof" || cfg.CipherSuiteIDs[0] != tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 {
		t.Fatalf("cfg=%+v", cfg)
	}
}

func TestParseConfigRejectsTLS13CipherSuite(t *testing.T) {
	_, err := parseConfig([]string{"-protocol", "https1", "-tls-version", "1.3", "-cipher-suites", "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want errInvalidConfig", err)
	}
}

func TestParseConfigRejectsUnsupportedProtocol(t *testing.T) {
	_, err := parseConfig([]string{"-protocol", "http2"})
	if !errors.Is(err, errUnsupportedProtocol) {
		t.Fatalf("err=%v, want errUnsupportedProtocol", err)
	}
}

func TestParseConfigSupportsHTTP1ServerOnly(t *testing.T) {
	cfg, err := parseConfig([]string{"-protocol", "http1", "-server-only=true"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.ServerOnly {
		t.Fatal("server-only=false, want true")
	}
}

func TestParseConfigSupportsHTTPS1ServerOnly(t *testing.T) {
	cfg, err := parseConfig([]string{"-protocol", "https1", "-server-only=true", "-tls-version", "1.2"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.ServerOnly || cfg.TLSVersion != tlsVersion12 {
		t.Fatalf("cfg=%+v", cfg)
	}
}

func TestRunServerOnlyHTTPS1(t *testing.T) {
	cfg, err := parseConfig([]string{
		"-protocol", "https1",
		"-server-only=true",
		"-addr", "127.0.0.1:0",
		"-payload", "16",
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	output := readyCancelWriter{cancel: cancel}
	if err := runServerOnly(ctx, cfg, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "serverReady=true framework=fasthttp protocol=https1") {
		t.Fatalf("output=%q", output.String())
	}
}

type readyCancelWriter struct {
	strings.Builder
	cancel context.CancelFunc
}

func (w *readyCancelWriter) Write(payload []byte) (int, error) {
	written, err := w.Builder.Write(payload)
	w.cancel()
	return written, err
}

func TestRunCLISmokeHTTP1(t *testing.T) {
	var out strings.Builder
	err := runCLI([]string{
		"-protocol", "http1",
		"-payload", "16",
		"-connections", "1",
		"-messages", "2",
		"-timeout", "3s",
	}, &out)
	if err != nil {
		t.Fatal(err)
	}
	text := out.String()
	for _, want := range []string{"framework=fasthttp", "protocol=http1", "total=2", "errors=0", "BenchmarkFastHTTPHTTP1-"} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in %s", want, text)
		}
	}
}

func TestRunCLISmokeHTTPS1TLS12(t *testing.T) {
	var out strings.Builder
	err := runCLI([]string{
		"-protocol", "https1",
		"-payload", "16",
		"-connections", "1",
		"-messages", "1",
		"-timeout", "3s",
		"-tls-version", "1.2",
		"-cipher-suites", "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
	}, &out)
	if err != nil {
		t.Fatal(err)
	}
	text := out.String()
	for _, want := range []string{"framework=fasthttp", "protocol=https1", "tlsVersion=1.2", "total=1", "errors=0"} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in %s", want, text)
		}
	}
}

func TestConfigRejectsNegativeTimeout(t *testing.T) {
	cfg := config{Protocol: protocolHTTP1, Addr: "127.0.0.1:0", Payload: 1, Connections: 1, Messages: 1, Timeout: -time.Second}
	if !errors.Is(cfg.validate(), errInvalidConfig) {
		t.Fatalf("validate err=%v, want errInvalidConfig", cfg.validate())
	}
}
