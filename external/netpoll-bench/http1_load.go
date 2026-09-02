package main

import (
	"context"

	"gnalloy.org/benchmarks/internal/httpbench"
)

func runHTTP1Load(ctx context.Context, addr string, cfg config) (benchResult, error) {
	result, err := httpbench.Run(ctx, httpbench.Config{
		Addr:              addr,
		Payload:           cfg.Payload,
		Connections:       cfg.Connections,
		Messages:          cfg.Messages,
		Timeout:           cfg.Timeout,
		LatencySampleRate: cfg.LatencySampleRate,
		WarmupMessages:    cfg.WarmupMessages,
	})
	return benchResult{
		TotalRequests: result.TotalRequests,
		Errors:        result.Errors,
		Elapsed:       result.Elapsed,
		Throughput:    result.Throughput,
		NsPerOp:       result.NsPerOp,
		Latency: latencySummary{
			Samples: result.Latency.Samples,
			P50:     result.Latency.P50,
			P95:     result.Latency.P95,
			P99:     result.Latency.P99,
			P999:    result.Latency.P999,
			Max:     result.Latency.Max,
		},
	}, err
}
