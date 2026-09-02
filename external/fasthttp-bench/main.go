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
	stopProfile, err := startCPUProfile(cfg.CPUProfile)
	if err != nil {
		return err
	}
	defer stopProfile()
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
	server, err := startHTTPServer(ctx, cfg)
	if err != nil {
		return err
	}
	defer server.stop()
	if err := servermode.WriteReady(stdout, servermode.Info{Framework: "fasthttp", Protocol: cfg.Protocol, Addr: server.addr}); err != nil {
		return err
	}
	servermode.Wait(ctx)
	return nil
}
