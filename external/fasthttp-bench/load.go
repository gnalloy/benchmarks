package main

import (
	"context"
	"time"

	"gnalloy.org/benchmarks/external/internal/benchhttp"
)

type benchResult struct {
	TotalRequests int64
	Errors        int64
	Elapsed       time.Duration
	Throughput    float64
	NsPerOp       float64
	Protocol      string
	Latency       latencySummary
	Resources     resourceDelta
}

func runBenchmark(parent context.Context, cfg config) (benchResult, error) {
	if err := (&cfg).resolve(); err != nil {
		return benchResult{}, err
	}
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithTimeout(parent, cfg.Timeout)
	defer cancel()
	server, err := startHTTPServer(ctx, cfg)
	if err != nil {
		return benchResult{}, err
	}
	defer server.stop()

	resourcesBefore := captureResourceSnapshot()
	httpConfig := benchhttp.Config{
		Addr:              server.addr,
		ServerName:        benchtlsServerName(),
		Payload:           cfg.Payload,
		Connections:       cfg.Connections,
		Messages:          cfg.Messages,
		Timeout:           cfg.Timeout,
		LatencySampleRate: cfg.LatencySampleRate,
		WarmupMessages:    cfg.WarmupMessages,
	}
	if cfg.Protocol == protocolHTTPS1 {
		tlsConfig, err := clientTLSConfig(cfg)
		if err != nil {
			return benchResult{}, err
		}
		httpConfig.TLS = tlsConfig
	}
	result, err := benchhttp.RunLoad(ctx, httpConfig)
	return benchResult{
		TotalRequests: result.TotalRequests,
		Errors:        result.Errors,
		Elapsed:       result.Elapsed,
		Throughput:    result.Throughput,
		NsPerOp:       result.NsPerOp,
		Protocol:      result.NegotiatedProtocol,
		Latency: latencySummary{
			Samples: result.Latency.Samples,
			P50:     result.Latency.P50,
			P95:     result.Latency.P95,
			P99:     result.Latency.P99,
			P999:    result.Latency.P999,
			Max:     result.Latency.Max,
		},
		Resources: resourceDeltaSince(resourcesBefore, captureResourceSnapshot()),
	}, err
}

func benchtlsServerName() string {
	return "gnalloy.local"
}
