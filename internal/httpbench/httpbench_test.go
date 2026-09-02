package httpbench

import (
	"bufio"
	"context"
	"errors"
	"net"
	"testing"
	"time"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name string
		edit func(*Config)
	}{
		{name: "empty address", edit: func(cfg *Config) { cfg.Addr = "" }},
		{name: "zero payload", edit: func(cfg *Config) { cfg.Payload = 0 }},
		{name: "negative warmup", edit: func(cfg *Config) { cfg.WarmupMessages = -1 }},
		{name: "negative sample rate", edit: func(cfg *Config) { cfg.LatencySampleRate = -1 }},
		{name: "negative target rate", edit: func(cfg *Config) { cfg.TargetRate = -1 }},
		{name: "zero timeout", edit: func(cfg *Config) { cfg.Timeout = 0 }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := DefaultConfig()
			test.edit(&cfg)
			if err := cfg.Validate(); !errors.Is(err, ErrInvalidConfig) {
				t.Fatalf("err=%v, want ErrInvalidConfig", err)
			}
		})
	}
}

func TestServerStateCountsSplitRequests(t *testing.T) {
	var state ServerState
	request := RequestBytes("localhost")
	if got := state.AppendAndCountRequests(request[:5]); got != 0 {
		t.Fatalf("requests=%d, want 0", got)
	}
	if got := state.AppendAndCountRequests(append(request[5:], request...)); got != 2 {
		t.Fatalf("requests=%d, want 2", got)
	}
}

func TestRunSaturatedAndPaced(t *testing.T) {
	for _, targetRate := range []int64{0, 1000} {
		t.Run(rateName(targetRate), func(t *testing.T) {
			listener, done := startTestServer(t, 3, 16)
			result, err := Run(context.Background(), Config{
				Addr:              listener.Addr().String(),
				Payload:           16,
				Connections:       1,
				Messages:          2,
				WarmupMessages:    1,
				LatencySampleRate: 1,
				TargetRate:        targetRate,
				Timeout:           5 * time.Second,
			})
			if err != nil {
				t.Fatal(err)
			}
			if result.TotalRequests != 2 || result.Errors != 0 || result.Latency.Samples != 2 || result.RoundTrip.Samples != 2 {
				t.Fatalf("result=%+v", result)
			}
			if targetRate > 0 && result.ScheduleDelay.Samples != 2 {
				t.Fatalf("scheduleDelay=%+v", result.ScheduleDelay)
			}
			<-done
		})
	}
}

func TestParseDecimal(t *testing.T) {
	if got, err := parseDecimal([]byte("1024")); err != nil || got != 1024 {
		t.Fatalf("got=%d err=%v", got, err)
	}
	if _, err := parseDecimal([]byte("1KiB")); err == nil {
		t.Fatal("expected invalid decimal error")
	}
}

func TestPhasePacerDistributesAggregateRate(t *testing.T) {
	start := time.Unix(123, 0)
	pacer := newPhasePacer(start, 4, 100)
	if got := pacer.deadline(2, 3); !got.Equal(start.Add(140 * time.Millisecond)) {
		t.Fatalf("deadline=%s, want %s", got, start.Add(140*time.Millisecond))
	}
}

func startTestServer(t *testing.T, exchanges int, payload int) (net.Listener, <-chan struct{}) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer listener.Close()
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		reader := bufio.NewReader(conn)
		response := ResponseBytes(payload)
		for range exchanges {
			for {
				line, readErr := reader.ReadString('\n')
				if readErr != nil || line == "\r\n" {
					break
				}
			}
			if _, err := conn.Write(response); err != nil {
				return
			}
		}
	}()
	return listener, done
}

func rateName(rate int64) string {
	if rate == 0 {
		return "saturated"
	}
	return "paced"
}
