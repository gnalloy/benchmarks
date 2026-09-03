package main

import (
	"testing"
	"time"

	http2config "github.com/hertz-contrib/http2/config"
)

func TestHTTP2ServerTimeoutsCoverConnectionPreparation(t *testing.T) {
	timeout := 5 * time.Minute
	cfg := http2config.NewConfig(http2ServerOptions(timeout)...)
	if cfg.ReadTimeout != timeout || cfg.IdleTimeout != timeout {
		t.Fatalf("readTimeout=%v idleTimeout=%v, want %v", cfg.ReadTimeout, cfg.IdleTimeout, timeout)
	}
	if cfg.DisableKeepalive {
		t.Fatal("HTTP/2 benchmark requires persistent connections")
	}
}
