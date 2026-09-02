package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"gnalloy.org/benchmarks/internal/httpbench"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	cfg := httpbench.DefaultConfig()
	serverFramework := "unknown"
	flags := flag.NewFlagSet("http1-load", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&cfg.Addr, "addr", cfg.Addr, "HTTP/1 server address")
	flags.StringVar(&cfg.ServerName, "server-name", cfg.ServerName, "HTTP Host header and TLS server name")
	flags.IntVar(&cfg.Payload, "payload", cfg.Payload, "response payload bytes")
	flags.IntVar(&cfg.Connections, "connections", cfg.Connections, "concurrent keep-alive connections")
	flags.IntVar(&cfg.Messages, "messages", cfg.Messages, "requests per connection")
	flags.IntVar(&cfg.WarmupMessages, "warmup-messages", cfg.WarmupMessages, "warmup requests per connection")
	flags.IntVar(&cfg.LatencySampleRate, "latency-sample-rate", cfg.LatencySampleRate, "record one latency every N requests")
	flags.Int64Var(&cfg.TargetRate, "target-rate", cfg.TargetRate, "aggregate target requests per second; 0 runs without pacing")
	flags.DurationVar(&cfg.Timeout, "timeout", cfg.Timeout, "overall timeout")
	flags.StringVar(&serverFramework, "server-framework", serverFramework, "server implementation label")
	if err := flags.Parse(args); err != nil {
		return err
	}
	serverFramework = strings.TrimSpace(serverFramework)
	if serverFramework == "" {
		return fmt.Errorf("http1-load: empty server framework")
	}
	result, err := httpbench.Run(context.Background(), cfg)
	if result.TotalRequests > 0 {
		fmt.Fprintf(stdout, "framework=common-http1-client serverFramework=%s protocol=http1 payload=%d connections=%d messages=%d warmupMessages=%d targetRate=%d latencySampleRate=%d latencySamples=%d p50LatencyNs=%d p95LatencyNs=%d p99LatencyNs=%d p999LatencyNs=%d maxLatencyNs=%d scheduleDelaySamples=%d p50ScheduleDelayNs=%d p95ScheduleDelayNs=%d p99ScheduleDelayNs=%d p999ScheduleDelayNs=%d maxScheduleDelayNs=%d roundTripLatencySamples=%d p50RoundTripLatencyNs=%d p95RoundTripLatencyNs=%d p99RoundTripLatencyNs=%d p999RoundTripLatencyNs=%d maxRoundTripLatencyNs=%d total=%d errors=%d elapsed=%s throughput=%.2f ops/s\n",
			serverFramework, cfg.Payload, cfg.Connections, cfg.Messages, cfg.WarmupMessages, cfg.TargetRate, cfg.LatencySampleRate,
			result.Latency.Samples, result.Latency.P50.Nanoseconds(), result.Latency.P95.Nanoseconds(), result.Latency.P99.Nanoseconds(), result.Latency.P999.Nanoseconds(), result.Latency.Max.Nanoseconds(),
			result.ScheduleDelay.Samples, result.ScheduleDelay.P50.Nanoseconds(), result.ScheduleDelay.P95.Nanoseconds(), result.ScheduleDelay.P99.Nanoseconds(), result.ScheduleDelay.P999.Nanoseconds(), result.ScheduleDelay.Max.Nanoseconds(),
			result.RoundTrip.Samples, result.RoundTrip.P50.Nanoseconds(), result.RoundTrip.P95.Nanoseconds(), result.RoundTrip.P99.Nanoseconds(), result.RoundTrip.P999.Nanoseconds(), result.RoundTrip.Max.Nanoseconds(),
			result.TotalRequests, result.Errors, result.Elapsed, result.Throughput)
	}
	return err
}
