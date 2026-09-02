package httpbench

import (
	"context"
	"errors"
	"net"
	"net/http"
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
			listener := startTestServer(t, 16)
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
		})
	}
}

func TestPhasePacerDistributesAggregateRate(t *testing.T) {
	start := time.Unix(123, 0)
	pacer := newPhasePacer(start, 4, 100)
	if got := pacer.deadline(2, 3); !got.Equal(start.Add(140 * time.Millisecond)) {
		t.Fatalf("deadline=%s, want %s", got, start.Add(140*time.Millisecond))
	}
}

func startTestServer(t *testing.T, payload int) net.Listener {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	body := ResponseBody(payload)
	server := &http.Server{
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/octet-stream")
			_, _ = writer.Write(body)
		}),
	}
	go func() {
		defer close(done)
		_ = server.Serve(listener)
	}()
	t.Cleanup(func() {
		_ = server.Close()
		<-done
	})
	return listener
}

func rateName(rate int64) string {
	if rate == 0 {
		return "saturated"
	}
	return "paced"
}
