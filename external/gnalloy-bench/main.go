package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"gnalloy.org/benchmarks/internal/servermode"
)

func main() {
	if err := runCLI(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runCLI(args []string, stdout io.Writer) error {
	cfg, err := parseConfig(args)
	if err != nil {
		return err
	}
	stopAllocProfile, err := startAllocProfile(cfg.AllocProfile)
	if err != nil {
		return err
	}
	defer stopAllocProfile()
	stopProfile, err := startCPUProfile(cfg.CPUProfile)
	if err != nil {
		return err
	}
	defer stopProfile()
	stopTrace, err := startRuntimeTrace(cfg.RuntimeTrace)
	if err != nil {
		return err
	}
	defer stopTrace()
	if cfg.ServerOnly {
		return runServerOnly(context.Background(), cfg, stdout)
	}
	result, err := runBenchmark(context.Background(), cfg)
	if result.TotalRequests > 0 {
		writeBenchmarkResult(stdout, cfg, result)
	}
	return err
}

func runServerOnly(ctx context.Context, cfg config, stdout io.Writer) error {
	switch cfg.Protocol {
	case "udp-echo":
		server, err := startUDPEchoServer(ctx, cfg)
		if err != nil {
			return err
		}
		defer server.stop()
		return waitForServerStop(ctx, stdout, cfg, server.addr)
	case "http1", "https1":
		server, err := startHTTP1Server(ctx, cfg)
		if err != nil {
			return err
		}
		defer server.stop()
		return waitForServerStop(ctx, stdout, cfg, server.addr)
	case "http2", "https2":
		server, err := startHTTP2Server(ctx, cfg)
		if err != nil {
			return err
		}
		defer server.stop()
		return waitForServerStop(ctx, stdout, cfg, server.addr)
	default:
		return fmt.Errorf("%w: server-only does not support %s", errInvalidConfig, cfg.Protocol)
	}
}

func waitForServerStop(ctx context.Context, stdout io.Writer, cfg config, addr string) error {
	if err := servermode.WriteReady(stdout, servermode.Info{Framework: "gnalloy", Protocol: cfg.Protocol, Addr: addr}); err != nil {
		return err
	}
	servermode.Wait(ctx)
	return nil
}

func benchmarkName(protocol string) string {
	switch protocol {
	case "http1":
		return "BenchmarkGnalloyHTTP1"
	case "https1":
		return "BenchmarkGnalloyHTTPS1"
	case "http2":
		return "BenchmarkGnalloyHTTP2"
	case "https2":
		return "BenchmarkGnalloyHTTPS2"
	case "http3":
		return "BenchmarkGnalloyHTTP3"
	case "quic-stream":
		return "BenchmarkGnalloyQUICStream"
	case "udp-echo":
		return "BenchmarkGnalloyUDPEcho"
	default:
		return "BenchmarkGnalloyTCPEcho"
	}
}
