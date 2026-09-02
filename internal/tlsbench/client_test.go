package tlsbench

import (
	"crypto/tls"
	"errors"
	"reflect"
	"testing"
)

func TestClientOptionsBuild(t *testing.T) {
	cfg, err := (ClientOptions{
		Enabled:            true,
		Version:            Version12,
		ALPN:               "h2,http/1.1",
		CipherSuites:       "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		InsecureSkipVerify: true,
	}).Build("127.0.0.1:443")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.MinVersion != tls.VersionTLS12 || cfg.MaxVersion != tls.VersionTLS12 {
		t.Fatalf("version=%d..%d", cfg.MinVersion, cfg.MaxVersion)
	}
	if cfg.ServerName != "127.0.0.1" || !cfg.InsecureSkipVerify {
		t.Fatalf("config=%+v", cfg)
	}
	if !reflect.DeepEqual(cfg.NextProtos, []string{"h2", "http/1.1"}) {
		t.Fatalf("nextProtos=%v", cfg.NextProtos)
	}
	if !reflect.DeepEqual(cfg.CipherSuites, []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256}) {
		t.Fatalf("cipherSuites=%v", cfg.CipherSuites)
	}
}

func TestClientOptionsBuildDisabled(t *testing.T) {
	cfg, err := (ClientOptions{}).Build("invalid")
	if err != nil || cfg != nil {
		t.Fatalf("config=%v err=%v", cfg, err)
	}
}

func TestClientOptionsBuildRejectsInvalidConfig(t *testing.T) {
	tests := []ClientOptions{
		{Enabled: true, Version: "1.0", ALPN: "http/1.1"},
		{Enabled: true, Version: Version13, ALPN: "http/1.1", CipherSuites: "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"},
		{Enabled: true, Version: Version13},
	}
	for _, options := range tests {
		if _, err := options.Build("127.0.0.1:443"); !errors.Is(err, ErrInvalidConfig) {
			t.Fatalf("options=%+v err=%v", options, err)
		}
	}
}
