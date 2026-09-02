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
	if cfg.ServerOnly {
		return runUDPServerOnly(context.Background(), cfg, stdout)
	}
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
	result, err := runBenchmark(context.Background(), cfg)
	if result.TotalRequests > 0 {
		writeBenchmarkResult(stdout, cfg, result)
	}
	return err
}

func runUDPServerOnly(ctx context.Context, cfg config, stdout io.Writer) error {
	server, err := startUDPEchoServer(ctx, cfg)
	if err != nil {
		return err
	}
	defer server.stop()
	if err := servermode.WriteReady(stdout, servermode.Info{Framework: "gnalloy", Protocol: cfg.Protocol, Addr: server.addr}); err != nil {
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
