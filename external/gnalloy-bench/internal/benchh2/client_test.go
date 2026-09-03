package benchh2

import (
	"context"
	cryptotls "crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	xhttp2 "golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func TestRequestFieldsUseHTTP2PseudoHeaders(t *testing.T) {
	fields := requestFields("bench.local", true)
	want := []http2Field{
		{name: ":method", value: "GET"},
		{name: ":scheme", value: "https"},
		{name: ":path", value: "/bench"},
		{name: ":authority", value: "bench.local"},
	}
	if len(fields) != len(want) {
		t.Fatalf("fields=%v", fields)
	}
	for i := range want {
		if fields[i].Name != want[i].name || fields[i].Value != want[i].value {
			t.Fatalf("field[%d]=%+v, want=%+v", i, fields[i], want[i])
		}
	}
}

func TestRunLoadHTTP2Cleartext(t *testing.T) {
	var requests atomic.Int64
	server := httptest.NewServer(h2c.NewHandler(benchmarkHandler(t, 16, &requests), &xhttp2.Server{}))
	defer server.Close()

	result, err := RunLoad(context.Background(), Config{
		Addr:              serverAddress(t, server.URL),
		Payload:           16,
		Connections:       1,
		Messages:          2,
		WarmupMessages:    1,
		LatencySampleRate: 1,
		Timeout:           5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalRequests != 2 || result.Errors != 0 || result.Latency.Samples != 2 {
		t.Fatalf("result=%+v", result)
	}
	if requests.Load() != 3 {
		t.Fatalf("requests=%d, want=3", requests.Load())
	}
}

func TestRunLoadHTTP2PacesAggregateTargetRate(t *testing.T) {
	var requests atomic.Int64
	server := httptest.NewServer(h2c.NewHandler(benchmarkHandler(t, 16, &requests), &xhttp2.Server{}))
	defer server.Close()

	result, err := RunLoad(context.Background(), Config{
		Addr:              serverAddress(t, server.URL),
		Payload:           16,
		Connections:       1,
		Messages:          3,
		LatencySampleRate: 1,
		TargetRate:        100,
		Timeout:           5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Elapsed < 20*time.Millisecond {
		t.Fatalf("elapsed=%s, want at least 20ms", result.Elapsed)
	}
	if result.Latency.Samples != 3 || result.ScheduleDelay.Samples != 3 || result.RoundTrip.Samples != 3 {
		t.Fatalf("result=%+v", result)
	}
	if result.Latency.P50 <= 0 || result.RoundTrip.P50 <= 0 || result.ScheduleDelay.P50 < 0 {
		t.Fatalf("result=%+v", result)
	}
}

func TestRunLoadHTTP2MoreConnectionsThanEventLoops(t *testing.T) {
	previous := runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(previous)

	var requests atomic.Int64
	server := httptest.NewServer(h2c.NewHandler(benchmarkHandler(t, 16, &requests), &xhttp2.Server{}))
	defer server.Close()

	result, err := RunLoad(context.Background(), Config{
		Addr:           serverAddress(t, server.URL),
		Payload:        16,
		Connections:    3,
		Messages:       2,
		WarmupMessages: 1,
		Timeout:        2 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalRequests != 6 || result.Errors != 0 {
		t.Fatalf("result=%+v", result)
	}
	if requests.Load() != 9 {
		t.Fatalf("requests=%d, want=9", requests.Load())
	}
}

func TestRunLoadHTTP2TLSALPN(t *testing.T) {
	var requests atomic.Int64
	server := httptest.NewUnstartedServer(benchmarkHandler(t, 8, &requests))
	server.EnableHTTP2 = true
	server.StartTLS()
	defer server.Close()

	result, err := RunLoad(context.Background(), Config{
		Addr:              serverAddress(t, server.URL),
		ServerName:        "localhost",
		Payload:           8,
		Connections:       1,
		Messages:          2,
		WarmupMessages:    1,
		LatencySampleRate: 1,
		Timeout:           5 * time.Second,
		TLS: &cryptotls.Config{
			ServerName:         "localhost",
			InsecureSkipVerify: true,
			MinVersion:         cryptotls.VersionTLS13,
			NextProtos:         []string{"h2"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalRequests != 2 || result.Errors != 0 || result.NegotiatedProtocol != "h2" {
		t.Fatalf("result=%+v", result)
	}
	if requests.Load() != 3 {
		t.Fatalf("requests=%d, want=3", requests.Load())
	}
}

func TestRunLoadTimeoutClosesBlockedClients(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			return
		}
		defer conn.Close()
		<-stop
	}()
	defer func() {
		close(stop)
		_ = listener.Close()
		<-done
	}()

	started := time.Now()
	_, err = RunLoad(context.Background(), Config{
		Addr:        listener.Addr().String(),
		Payload:     8,
		Connections: 1,
		Messages:    1,
		Timeout:     150 * time.Millisecond,
	})
	if err == nil {
		t.Fatal("expected timeout or closed connection error")
	}
	if elapsed := time.Since(started); elapsed > 2*time.Second {
		t.Fatalf("RunLoad elapsed=%v, want bounded timeout", elapsed)
	}
}

type http2Field struct {
	name  string
	value string
}

func benchmarkHandler(t *testing.T, payload int, requests *atomic.Int64) http.Handler {
	t.Helper()
	body := ResponseBody(payload)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor != 2 || r.Method != http.MethodGet || r.URL.Path != "/bench" {
			t.Errorf("request=%s %s proto=%s", r.Method, r.URL.Path, r.Proto)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		requests.Add(1)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(body); err != nil {
			t.Errorf("write response: %v", err)
		}
	})
}

func serverAddress(t *testing.T, rawURL string) string {
	t.Helper()
	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatal(err)
	}
	return parsed.Host
}
