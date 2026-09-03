package main

import (
	"context"
	cryptotls "crypto/tls"
	"errors"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"gnalloy.org/gnalloy/channel"
	"gnalloy.org/gnalloy/transport"
	handlertls "gnalloy.org/handler-tls"
)

func TestParseConfig(t *testing.T) {
	cfg, err := parseConfig([]string{
		"-protocol", "tcp-echo",
		"-addr", "127.0.0.1:0",
		"-payload", "32",
		"-connections", "2",
		"-messages", "3",
		"-timeout", "2s",
		"-backend", "std",
		"-boss", "1",
		"-workers", "2",
		"-read-buffer-size", "8192",
		"-max-messages-per-read", "5",
		"-event-batch-size", "7",
		"-tcp-echo-mode", "owner-executor",
		"-flush-strategy", "immediate",
		"-tcp-echo-executor-workers", "3",
		"-tcp-echo-executor-queue-size", "11",
		"-boss-cpus", "0",
		"-worker-cpus", "1,2",
		"-reuseport",
		"-warmup-messages", "7",
		"-cpuprofile", "cpu.pprof",
		"-allocprofile", "alloc.pprof",
		"-trace", "runtime.trace",
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Payload != 32 || cfg.Connections != 2 || cfg.Messages != 3 || cfg.Timeout != 2*time.Second {
		t.Fatalf("cfg=%+v", cfg)
	}
	if cfg.Boss != 1 || cfg.Workers != 2 || cfg.ReadBufferSize != 8192 || cfg.MaxMessagesPerRead != 5 || cfg.EventBatchSize != 7 || backendLabel(cfg.Backend) != "std" {
		t.Fatalf("cfg=%+v", cfg)
	}
	if cfg.TCPEchoMode != tcpEchoModeOwnerExecutor || cfg.TCPEchoExecutorWorkers != 3 || cfg.TCPEchoExecutorQueueSize != 11 {
		t.Fatalf("tcp echo config=%+v", cfg)
	}
	if cfg.FlushStrategyName != flushStrategyImmediate || cfg.FlushStrategy != channel.FlushImmediate {
		t.Fatalf("flush strategy=%s/%d, want immediate", cfg.FlushStrategyName, cfg.FlushStrategy)
	}
	if !cfg.ReusePort {
		t.Fatalf("reuseport=%v, want true", cfg.ReusePort)
	}
	if cfg.WarmupMessages != 7 {
		t.Fatalf("warmupMessages=%d, want 7", cfg.WarmupMessages)
	}
	if cfg.CPUProfile != "cpu.pprof" {
		t.Fatalf("cpuProfile=%q, want cpu.pprof", cfg.CPUProfile)
	}
	if cfg.RuntimeTrace != "runtime.trace" {
		t.Fatalf("runtimeTrace=%q, want runtime.trace", cfg.RuntimeTrace)
	}
	if cfg.AllocProfile != "alloc.pprof" {
		t.Fatalf("allocProfile=%q, want alloc.pprof", cfg.AllocProfile)
	}
	if cfg.BossCPUs != "0" || len(cfg.BossCPUSet) != 1 || cfg.BossCPUSet[0] != 0 {
		t.Fatalf("boss CPU set=%q/%v, want 0/[0]", cfg.BossCPUs, cfg.BossCPUSet)
	}
	if cfg.WorkerCPUs != "1,2" || len(cfg.WorkerCPUSet) != 2 || cfg.WorkerCPUSet[0] != 1 || cfg.WorkerCPUSet[1] != 2 {
		t.Fatalf("worker CPU set=%q/%v, want 1,2/[1 2]", cfg.WorkerCPUs, cfg.WorkerCPUSet)
	}
	if cfg.ALPN != "http/1.1" {
		t.Fatalf("alpn=%q, want http/1.1", cfg.ALPN)
	}
}

func TestParseConfigSupportsHTTPS1ALPN(t *testing.T) {
	cfg, err := parseConfig([]string{
		"-protocol", "https1",
		"-alpn", "h2,http/1.1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Protocol != "https1" || cfg.ALPN != "h2,http/1.1" {
		t.Fatalf("cfg=%+v", cfg)
	}
}

func TestParseConfigRejectsHTTP1BenchmarkBypass(t *testing.T) {
	_, err := parseConfig([]string{
		"-protocol", "http1",
		"-http1-mode", "raw",
	})
	if err == nil {
		t.Fatal("expected removed http1-mode flag to be rejected")
	}
}

func TestParseConfigSupportsTCPEchoModes(t *testing.T) {
	for _, mode := range []string{tcpEchoModeDirect, tcpEchoModeReadComplete, tcpEchoModeOwnerExecutor} {
		cfg, err := parseConfig([]string{
			"-protocol", "tcp-echo",
			"-tcp-echo-mode", mode,
		})
		if err != nil {
			t.Fatalf("mode %s: %v", mode, err)
		}
		if cfg.TCPEchoMode != mode {
			t.Fatalf("mode=%q, want %q", cfg.TCPEchoMode, mode)
		}
	}
}

func TestParseConfigSupportsFlushStrategies(t *testing.T) {
	tests := []struct {
		name string
		want channel.FlushStrategy
	}{
		{name: flushStrategyImmediate, want: channel.FlushImmediate},
		{name: flushStrategyReadComplete, want: channel.FlushOnReadComplete},
		{name: flushStrategyEventLoopBatch, want: channel.FlushOnEventLoopBatch},
		{name: "batch", want: channel.FlushOnEventLoopBatch},
	}

	for _, tt := range tests {
		cfg, err := parseConfig([]string{"-flush-strategy", tt.name})
		if err != nil {
			t.Fatalf("strategy %s: %v", tt.name, err)
		}
		if cfg.FlushStrategy != tt.want {
			t.Fatalf("strategy=%s/%d, want %d", cfg.FlushStrategyName, cfg.FlushStrategy, tt.want)
		}
	}
}

func TestParseConfigResolvesDefaultTCPEchoExecutor(t *testing.T) {
	cfg, err := parseConfig([]string{})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.TCPEchoMode != defaultTCPEchoMode {
		t.Fatalf("tcpEchoMode=%q, want %q", cfg.TCPEchoMode, defaultTCPEchoMode)
	}
	if cfg.TCPEchoExecutorWorkers != runtime.GOMAXPROCS(0) {
		t.Fatalf("tcpEchoExecutorWorkers=%d, want GOMAXPROCS", cfg.TCPEchoExecutorWorkers)
	}
	if cfg.TCPEchoExecutorQueueSize != defaultTCPEchoExecutorQueueSize {
		t.Fatalf("tcpEchoExecutorQueueSize=%d, want %d", cfg.TCPEchoExecutorQueueSize, defaultTCPEchoExecutorQueueSize)
	}
	if cfg.FlushStrategyName != defaultFlushStrategyName || cfg.FlushStrategy != channel.FlushOnReadComplete {
		t.Fatalf("flush strategy=%s/%d, want read-complete", cfg.FlushStrategyName, cfg.FlushStrategy)
	}
}

func TestParseConfigRejectsInvalidTCPEchoMode(t *testing.T) {
	_, err := parseConfig([]string{"-tcp-echo-mode", "inline"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigRejectsInvalidFlushStrategy(t *testing.T) {
	_, err := parseConfig([]string{"-flush-strategy", "spin"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigRejectsTCPEchoModeForHTTP1(t *testing.T) {
	_, err := parseConfig([]string{
		"-protocol", "http1",
		"-tcp-echo-mode", tcpEchoModeReadComplete,
	})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigSupportsTLSVersions(t *testing.T) {
	for _, version := range []string{"1.1", "1.2", "1.3"} {
		cfg, err := parseConfig([]string{
			"-protocol", "https1",
			"-tls-version", version,
		})
		if err != nil {
			t.Fatalf("version %s: %v", version, err)
		}
		if cfg.TLSVersion != version {
			t.Fatalf("version=%q, want %s", cfg.TLSVersion, version)
		}
	}
}

func TestParseConfigSupportsCipherSuites(t *testing.T) {
	cfg, err := parseConfig([]string{
		"-protocol", "https1",
		"-tls-version", "1.2",
		"-cipher-suites", "ECDHE-RSA-AES128-GCM-SHA256,0xC030",
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.CipherSuites != "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384" {
		t.Fatalf("cipherSuites=%q", cfg.CipherSuites)
	}
	want := []uint16{cryptotls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, cryptotls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384}
	if len(cfg.CipherSuiteIDs) != len(want) {
		t.Fatalf("cipherSuiteIDs=%x, want %x", cfg.CipherSuiteIDs, want)
	}
	for i := range want {
		if cfg.CipherSuiteIDs[i] != want[i] {
			t.Fatalf("cipherSuiteIDs=%x, want %x", cfg.CipherSuiteIDs, want)
		}
	}
}

func TestParseConfigRejectsInsecureCipherSuiteByDefault(t *testing.T) {
	_, err := parseConfig([]string{
		"-protocol", "https1",
		"-tls-version", "1.1",
		"-cipher-suites", "TLS_RSA_WITH_AES_128_CBC_SHA",
	})
	if !errors.Is(err, handlertls.ErrInsecureCipherSuite) {
		t.Fatalf("err=%v, want %v", err, handlertls.ErrInsecureCipherSuite)
	}
}

func TestParseConfigSupportsInsecureCipherSuiteOptIn(t *testing.T) {
	cfg, err := parseConfig([]string{
		"-protocol", "https1",
		"-tls-version", "1.1",
		"-cipher-suites", "TLS_RSA_WITH_AES_128_CBC_SHA",
		"-allow-insecure-cipher-suites",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.CipherSuiteIDs) != 1 || cfg.CipherSuiteIDs[0] != cryptotls.TLS_RSA_WITH_AES_128_CBC_SHA {
		t.Fatalf("cipherSuiteIDs=%x, want TLS_RSA_WITH_AES_128_CBC_SHA", cfg.CipherSuiteIDs)
	}
}

func TestParseConfigRejectsCipherSuitesForTLS13(t *testing.T) {
	_, err := parseConfig([]string{
		"-protocol", "https1",
		"-tls-version", "1.3",
		"-cipher-suites", "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
	})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigRejectsHTTP2TLS11(t *testing.T) {
	_, err := parseConfig([]string{"-protocol", "https2", "-tls-version", "1.1"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigSupportsHTTP2Family(t *testing.T) {
	cfg, err := parseConfig([]string{"-protocol", "http2"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Protocol != "http2" {
		t.Fatalf("protocol=%q, want http2", cfg.Protocol)
	}

	cfg, err = parseConfig([]string{"-protocol", "https2"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Protocol != "https2" || cfg.ALPN != "h2" {
		t.Fatalf("cfg=%+v, want https2 with h2 ALPN", cfg)
	}
}

func TestParseConfigSupportsHTTP3(t *testing.T) {
	cfg, err := parseConfig([]string{"-protocol", "http3"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Protocol != "http3" || cfg.ALPN != "h3" || cfg.TLSVersion != "1.3" {
		t.Fatalf("cfg=%+v, want http3 with h3/TLS1.3", cfg)
	}
}

func TestParseConfigSupportsQUICStream(t *testing.T) {
	cfg, err := parseConfig([]string{"-protocol", "quic-stream"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Protocol != "quic-stream" || cfg.ALPN != "gnalloy-quic" || cfg.TLSVersion != "1.3" {
		t.Fatalf("cfg=%+v, want quic-stream with gnalloy-quic/TLS1.3", cfg)
	}
}

func TestParseConfigRejectsHTTP3TLS12(t *testing.T) {
	_, err := parseConfig([]string{"-protocol", "http3", "-tls-version", "1.2"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigRejectsQUICStreamTLS12(t *testing.T) {
	_, err := parseConfig([]string{"-protocol", "quic-stream", "-tls-version", "1.2"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigSupportsUDPEcho(t *testing.T) {
	cfg, err := parseConfig([]string{
		"-protocol", "udp-echo",
		"-payload", "128",
		"-connections", "2",
		"-messages", "3",
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Protocol != "udp-echo" || cfg.Payload != 128 || cfg.Connections != 2 || cfg.Messages != 3 {
		t.Fatalf("cfg=%+v", cfg)
	}
}

func TestParseConfigSupportsUDPServerOnly(t *testing.T) {
	cfg, err := parseConfig([]string{"-protocol", "udp-echo", "-server-only=true"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.ServerOnly {
		t.Fatal("server-only=false, want true")
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

func TestParseConfigSupportsHTTP2ServerOnly(t *testing.T) {
	for _, protocol := range []string{"http2", "https2"} {
		t.Run(protocol, func(t *testing.T) {
			cfg, err := parseConfig([]string{"-protocol", protocol, "-server-only=true"})
			if err != nil {
				t.Fatal(err)
			}
			if !cfg.ServerOnly {
				t.Fatal("server-only=false, want true")
			}
		})
	}
}

func TestParseConfigSupportsHTTP2ClientOnly(t *testing.T) {
	for _, protocol := range []string{"http2", "https2"} {
		t.Run(protocol, func(t *testing.T) {
			cfg, err := parseConfig([]string{"-protocol", protocol, "-client-only=true"})
			if err != nil {
				t.Fatal(err)
			}
			if !cfg.ClientOnly {
				t.Fatal("client-only=false, want true")
			}
		})
	}
}

func TestParseConfigSupportsHTTP3StandaloneModes(t *testing.T) {
	for _, flag := range []string{"-server-only=true", "-client-only=true"} {
		t.Run(flag, func(t *testing.T) {
			cfg, err := parseConfig([]string{"-protocol", "http3", flag})
			if err != nil {
				t.Fatal(err)
			}
			if !cfg.ServerOnly && !cfg.ClientOnly {
				t.Fatalf("cfg=%+v, want standalone mode", cfg)
			}
		})
	}
}

func TestParseConfigRejectsConflictingStandaloneModes(t *testing.T) {
	_, err := parseConfig([]string{"-protocol", "http2", "-server-only=true", "-client-only=true"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestRunServerOnlyHTTPS1(t *testing.T) {
	cfg, err := parseConfig([]string{
		"-protocol", "https1",
		"-server-only=true",
		"-addr", "127.0.0.1:0",
		"-payload", "16",
		"-boss", "1",
		"-workers", "1",
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
	if !strings.Contains(output.String(), "serverReady=true framework=gnalloy protocol=https1") {
		t.Fatalf("output=%q", output.String())
	}
}

func TestRunServerOnlyHTTP2(t *testing.T) {
	cfg, err := parseConfig([]string{
		"-protocol", "http2",
		"-server-only=true",
		"-addr", "127.0.0.1:0",
		"-payload", "16",
		"-boss", "1",
		"-workers", "1",
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
	if !strings.Contains(output.String(), "serverReady=true framework=gnalloy protocol=http2") {
		t.Fatalf("output=%q", output.String())
	}
}

func TestRunServerOnlyHTTP3(t *testing.T) {
	cfg, err := parseConfig([]string{
		"-protocol", "http3",
		"-server-only=true",
		"-addr", "127.0.0.1:0",
		"-payload", "16",
		"-connections", "1",
		"-messages", "1",
		"-boss", "1",
		"-workers", "1",
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
	if !strings.Contains(output.String(), "serverReady=true framework=gnalloy protocol=http3") {
		t.Fatalf("output=%q", output.String())
	}
}

func TestRunClientOnlyHTTP2(t *testing.T) {
	serverConfig, err := parseConfig([]string{
		"-protocol", "http2",
		"-addr", "127.0.0.1:0",
		"-payload", "16",
		"-connections", "1",
		"-messages", "1",
		"-boss", "1",
		"-workers", "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	server, err := startHTTP2Server(context.Background(), serverConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer server.stop()

	clientConfig := serverConfig
	clientConfig.ClientOnly = true
	clientConfig.Addr = server.addr
	result, err := runClientOnly(context.Background(), clientConfig)
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalRequests != 1 || result.Errors != 0 {
		t.Fatalf("result=%+v, want one successful request", result)
	}
}

func TestRunClientOnlyHTTP3(t *testing.T) {
	serverConfig, err := parseConfig([]string{
		"-protocol", "http3",
		"-addr", "127.0.0.1:0",
		"-payload", "16",
		"-connections", "1",
		"-messages", "1",
		"-boss", "1",
		"-workers", "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	server, err := startHTTP3Server(context.Background(), serverConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer server.stop()

	clientConfig := serverConfig
	clientConfig.ClientOnly = true
	clientConfig.Addr = server.addr
	result, err := runClientOnly(context.Background(), clientConfig)
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalRequests != 1 || result.Errors != 0 || result.Protocol != http3ALPNValue {
		t.Fatalf("result=%+v, want one successful HTTP/3 request", result)
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

func TestParseConfigRejectsServerOnlyForOtherProtocols(t *testing.T) {
	_, err := parseConfig([]string{"-protocol", "tcp-echo", "-server-only=true"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigResolvesNativePerformanceFlags(t *testing.T) {
	cfg, err := parseConfig([]string{
		"-backend", "iouring",
		"-mmap",
		"-mmap-block-size", "8192",
		"-mmap-blocks", "1024",
		"-iouring-fixed-buffers",
		"-iouring-multishot-accept",
		"-iouring-sqpoll",
		"-latency-sample-rate", "32",
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Backend != transport.BackendIOUring || !cfg.Mmap || cfg.MmapBlockSize != 8192 || cfg.MmapBlocks != 1024 {
		t.Fatalf("cfg=%+v", cfg)
	}
	if !cfg.IOUringFixedBuffers || !cfg.IOUringMultishotAccept || !cfg.IOUringSQPoll {
		t.Fatalf("iouring flags=%+v", cfg)
	}
	if cfg.LatencySampleRate != 32 {
		t.Fatalf("latencySampleRate=%d, want 32", cfg.LatencySampleRate)
	}
}

func TestParseConfigResolvesAutoWorkers(t *testing.T) {
	cfg, err := parseConfig([]string{
		"-backend", "iocp",
		"-workers", "0",
	})
	if err != nil {
		t.Fatal(err)
	}
	want := defaultWorkerCount(workerSizingInput{
		GOOS:       runtime.GOOS,
		Backend:    transport.BackendIOCP,
		GOMAXPROCS: runtime.GOMAXPROCS(0),
		Bosses:     1,
	})
	if cfg.Workers != want {
		t.Fatalf("workers=%d, want %d", cfg.Workers, want)
	}
}

func TestParseConfigResolvesAutoReadBufferSize(t *testing.T) {
	cfg, err := parseConfig([]string{
		"-payload", "16384",
		"-read-buffer-size", "0",
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ReadBufferSize != 16384 {
		t.Fatalf("readBufferSize=%d, want 16384", cfg.ReadBufferSize)
	}
}

func TestParseConfigKeepsMinimumAutoReadBufferSize(t *testing.T) {
	cfg, err := parseConfig([]string{
		"-payload", "64",
		"-read-buffer-size", "0",
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ReadBufferSize != 4096 {
		t.Fatalf("readBufferSize=%d, want 4096", cfg.ReadBufferSize)
	}
}

func TestParseConfigResolvesDefaultMaxMessagesPerRead(t *testing.T) {
	cfg, err := parseConfig([]string{"-max-messages-per-read", "0"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.MaxMessagesPerRead != channel.OptionMaxMessagesPerRead.Default() {
		t.Fatalf("maxMessagesPerRead=%d, want %d", cfg.MaxMessagesPerRead, channel.OptionMaxMessagesPerRead.Default())
	}
}

func TestParseConfigRejectsUnsupportedProtocol(t *testing.T) {
	_, err := parseConfig([]string{"-protocol", "sctp-echo"})
	if !errors.Is(err, errUnsupportedProtocol) {
		t.Fatalf("err=%v, want %v", err, errUnsupportedProtocol)
	}
}

func TestParseConfigRejectsNegativeWorkers(t *testing.T) {
	_, err := parseConfig([]string{"-workers", "-1"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigRejectsNegativeReadBufferSize(t *testing.T) {
	_, err := parseConfig([]string{"-read-buffer-size", "-1"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigRejectsNegativeMaxMessagesPerRead(t *testing.T) {
	_, err := parseConfig([]string{"-max-messages-per-read", "-1"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigRejectsNegativeEventBatchSize(t *testing.T) {
	_, err := parseConfig([]string{"-event-batch-size", "-1"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigRejectsInvalidCPUSet(t *testing.T) {
	for _, args := range [][]string{
		{"-boss-cpus", "-1"},
		{"-worker-cpus", "0,,1"},
		{"-worker-cpus", "x"},
	} {
		_, err := parseConfig(args)
		if !errors.Is(err, errInvalidConfig) {
			t.Fatalf("args=%v err=%v, want %v", args, err, errInvalidConfig)
		}
	}
}

func TestParseConfigRejectsNegativeLatencySampleRate(t *testing.T) {
	_, err := parseConfig([]string{"-latency-sample-rate", "-1"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigSupportsTargetRate(t *testing.T) {
	cfg, err := parseConfig([]string{"-protocol", "http2", "-target-rate", "60000"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.TargetRate != 60000 {
		t.Fatalf("targetRate=%d, want 60000", cfg.TargetRate)
	}
}

func TestParseConfigSupportsHTTP3TargetRate(t *testing.T) {
	cfg, err := parseConfig([]string{"-protocol", "http3", "-target-rate", "30000"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.TargetRate != 30000 {
		t.Fatalf("targetRate=%d, want 30000", cfg.TargetRate)
	}
}

func TestParseConfigRejectsNegativeTargetRate(t *testing.T) {
	_, err := parseConfig([]string{"-protocol", "http2", "-target-rate", "-1"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigRejectsUnsupportedTargetRateProtocol(t *testing.T) {
	_, err := parseConfig([]string{"-protocol", "udp-echo", "-target-rate", "100"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigRejectsNegativeWarmupMessages(t *testing.T) {
	_, err := parseConfig([]string{"-warmup-messages", "-1"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigRejectsInvalidTLSVersion(t *testing.T) {
	_, err := parseConfig([]string{"-protocol", "https1", "-tls-version", "1.0"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigRejectsFixedBuffersWithoutMmap(t *testing.T) {
	_, err := parseConfig([]string{"-backend", "iouring", "-iouring-fixed-buffers"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigRejectsMmapBlockSmallerThanReadBuffer(t *testing.T) {
	_, err := parseConfig([]string{"-mmap", "-mmap-block-size", "1024", "-read-buffer-size", "4096"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigRejectsMmapSizeOverflow(t *testing.T) {
	_, err := parseConfig([]string{"-mmap", "-mmap-block-size", strconv.Itoa(maxInt), "-mmap-blocks", "2"})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigRejectsInvalidBackend(t *testing.T) {
	_, err := parseConfig([]string{"-backend", "select"})
	if !errors.Is(err, errInvalidBackend) {
		t.Fatalf("err=%v, want %v", err, errInvalidBackend)
	}
}
