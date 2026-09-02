package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"gnalloy.org/benchmarks/internal/udpbench"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	cfg := udpbench.DefaultConfig()
	serverFramework := "unknown"
	flags := flag.NewFlagSet("udp-load", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&cfg.Addr, "addr", cfg.Addr, "UDP echo server address")
	flags.IntVar(&cfg.Payload, "payload", cfg.Payload, "payload bytes")
	flags.IntVar(&cfg.Connections, "connections", cfg.Connections, "concurrent UDP clients")
	flags.IntVar(&cfg.Messages, "messages", cfg.Messages, "messages per client")
	flags.IntVar(&cfg.WarmupMessages, "warmup-messages", cfg.WarmupMessages, "warmup messages per client")
	flags.IntVar(&cfg.LatencySampleRate, "latency-sample-rate", cfg.LatencySampleRate, "record one latency every N messages")
	flags.Int64Var(&cfg.TargetRate, "target-rate", cfg.TargetRate, "aggregate target requests per second; 0 runs without pacing")
	flags.DurationVar(&cfg.Timeout, "timeout", cfg.Timeout, "overall timeout")
	flags.StringVar(&serverFramework, "server-framework", serverFramework, "server implementation label")
	if err := flags.Parse(args); err != nil {
		return err
	}
	serverFramework = strings.TrimSpace(serverFramework)
	if serverFramework == "" {
		return fmt.Errorf("udp-load: empty server framework")
	}
	result, err := udpbench.Run(context.Background(), cfg)
	if result.TotalRequests > 0 {
		fmt.Fprintf(stdout, "framework=common-udp-client serverFramework=%s protocol=udp-echo payload=%d connections=%d messages=%d warmupMessages=%d targetRate=%d latencySampleRate=%d latencySamples=%d p50LatencyNs=%d p95LatencyNs=%d p99LatencyNs=%d p999LatencyNs=%d maxLatencyNs=%d total=%d errors=%d elapsed=%s throughput=%.2f ops/s\n",
			serverFramework, cfg.Payload, cfg.Connections, cfg.Messages, cfg.WarmupMessages, cfg.TargetRate, cfg.LatencySampleRate,
			result.Latency.Samples, result.Latency.P50.Nanoseconds(), result.Latency.P95.Nanoseconds(), result.Latency.P99.Nanoseconds(), result.Latency.P999.Nanoseconds(), result.Latency.Max.Nanoseconds(),
			result.TotalRequests, result.Errors, result.Elapsed, result.Throughput)
	}
	return err
}
