package udpbench

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"
)

func TestRunUsesCommonUDPEchoClient(t *testing.T) {
	conn := startEchoServer(t)
	cfg := DefaultConfig()
	cfg.Addr = conn.LocalAddr().String()
	cfg.Payload = 32
	cfg.Connections = 4
	cfg.Messages = 20
	cfg.WarmupMessages = 2
	cfg.LatencySampleRate = 2
	cfg.Timeout = 5 * time.Second

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalRequests != 80 || result.Errors != 0 {
		t.Fatalf("total=%d errors=%d", result.TotalRequests, result.Errors)
	}
	if result.Latency.Samples != 40 || result.Latency.P99 <= 0 || result.Throughput <= 0 {
		t.Fatalf("latency=%+v throughput=%f", result.Latency, result.Throughput)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []Config{
		{},
		{Addr: "127.0.0.1:1", Payload: 1, Connections: 1, Messages: 1},
	}
	for _, cfg := range tests {
		if err := cfg.Validate(); !errors.Is(err, ErrInvalidConfig) {
			t.Fatalf("config=%+v error=%v", cfg, err)
		}
	}
}

func startEchoServer(t *testing.T) net.PacketConn {
	t.Helper()
	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	go func() {
		payload := make([]byte, 2048)
		for {
			n, addr, err := conn.ReadFrom(payload)
			if err != nil {
				return
			}
			if _, err := conn.WriteTo(payload[:n], addr); err != nil {
				return
			}
		}
	}()
	return conn
}
