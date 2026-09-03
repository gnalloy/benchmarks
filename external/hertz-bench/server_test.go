package main

import (
	"crypto/tls"
	"io"
	"net/http"
	"testing"
	"time"

	http2config "github.com/hertz-contrib/http2/config"
)

func TestServerHandlesHTTP1(t *testing.T) {
	server, err := startServer(config{
		Protocol: protocolHTTP1,
		Addr:     "127.0.0.1:0",
		Payload:  128,
		Timeout:  time.Minute,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := server.stop(); err != nil {
			t.Errorf("stop server: %v", err)
		}
	}()
	transport := &http.Transport{}
	defer transport.CloseIdleConnections()
	assertHTTP1Response(t, &http.Client{Transport: transport, Timeout: time.Second}, "http://"+server.addr+"/bench", 128, 0)
}

func TestServerHandlesHTTPS1TLSVersions(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    uint16
		cipher  uint16
	}{
		{name: "tls11", version: tlsVersion11, want: tls.VersionTLS11, cipher: tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA},
		{name: "tls12", version: tlsVersion12, want: tls.VersionTLS12, cipher: tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		{name: "tls13", version: tlsVersion13, want: tls.VersionTLS13},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config{
				Protocol:       protocolHTTPS1,
				Addr:           "127.0.0.1:0",
				Payload:        128,
				Timeout:        time.Minute,
				TLSVersion:     tt.version,
				ALPN:           "http/1.1",
				CipherSuiteIDs: nil,
			}
			if tt.cipher != 0 {
				cfg.CipherSuiteIDs = []uint16{tt.cipher}
			}
			server, err := startServer(cfg)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := server.stop(); err != nil {
					t.Errorf("stop server: %v", err)
				}
			}()
			transport := &http.Transport{TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				MinVersion:         tt.want,
				MaxVersion:         tt.want,
				NextProtos:         []string{"http/1.1"},
				CipherSuites:       cfg.CipherSuiteIDs,
			}}
			defer transport.CloseIdleConnections()
			assertHTTP1Response(t, &http.Client{Transport: transport, Timeout: time.Second}, "https://"+server.addr+"/bench", 128, tt.want)
		})
	}
}

func assertHTTP1Response(t *testing.T, client *http.Client, url string, payload int, tlsVersion uint16) {
	t.Helper()
	response, err := client.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK || len(body) != payload {
		t.Fatalf("status=%d body=%d, want 200/%d", response.StatusCode, len(body), payload)
	}
	if tlsVersion != 0 && (response.TLS == nil || response.TLS.Version != tlsVersion || response.TLS.NegotiatedProtocol != "http/1.1") {
		t.Fatalf("tls=%+v, want version=%x alpn=http/1.1", response.TLS, tlsVersion)
	}
}

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
