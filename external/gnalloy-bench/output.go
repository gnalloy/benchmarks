package main

import (
	"fmt"
	"io"
	"runtime"
)

func writeBenchmarkResult(w io.Writer, cfg config, result benchResult) {
	if w == nil {
		return
	}
	fmt.Fprintf(w, "framework=gnalloy protocol=%s backend=%s tcpEchoMode=%s flushStrategy=%s tcpEchoExecutorWorkers=%d tcpEchoExecutorQueueSize=%d boss=%d workers=%d readBufferSize=%d maxMessagesPerRead=%d eventBatchSize=%d bossCPUs=%s workerCPUs=%s reuseport=%t mmap=%t mmapBlockSize=%d mmapBlocks=%d iouringFixedBuffers=%t iouringMultishotAccept=%t iouringSQPoll=%t tlsVersion=%s cipherSuites=%s negotiatedProtocol=%s targetRate=%d latencySampleRate=%d warmupMessages=%d latencySamples=%d p50LatencyNs=%d p95LatencyNs=%d p99LatencyNs=%d p999LatencyNs=%d maxLatencyNs=%d scheduleDelaySamples=%d p50ScheduleDelayNs=%d p95ScheduleDelayNs=%d p99ScheduleDelayNs=%d p999ScheduleDelayNs=%d maxScheduleDelayNs=%d roundTripLatencySamples=%d p50RoundTripLatencyNs=%d p95RoundTripLatencyNs=%d p99RoundTripLatencyNs=%d p999RoundTripLatencyNs=%d maxRoundTripLatencyNs=%d rssBytes=%d heapAllocBytes=%d heapSysBytes=%d heapObjects=%d gcCount=%d gcPauseNs=%d goroutines=%d payload=%d connections=%d messages=%d total=%d errors=%d elapsed=%s throughput=%.2f ops/s\n",
		cfg.Protocol, benchmarkBackendLabel(cfg), cfg.TCPEchoMode, cfg.FlushStrategyName, cfg.TCPEchoExecutorWorkers, cfg.TCPEchoExecutorQueueSize, cfg.Boss, cfg.Workers, cfg.ReadBufferSize, cfg.MaxMessagesPerRead, cfg.EventBatchSize, cfg.BossCPUs, cfg.WorkerCPUs, cfg.ReusePort, cfg.Mmap, cfg.MmapBlockSize, cfg.MmapBlocks, cfg.IOUringFixedBuffers, cfg.IOUringMultishotAccept, cfg.IOUringSQPoll, cfg.TLSVersion, cfg.CipherSuites, result.Protocol, cfg.TargetRate, cfg.LatencySampleRate, cfg.WarmupMessages, result.Latency.Samples, result.Latency.P50.Nanoseconds(), result.Latency.P95.Nanoseconds(), result.Latency.P99.Nanoseconds(), result.Latency.P999.Nanoseconds(), result.Latency.Max.Nanoseconds(), result.ScheduleDelay.Samples, result.ScheduleDelay.P50.Nanoseconds(), result.ScheduleDelay.P95.Nanoseconds(), result.ScheduleDelay.P99.Nanoseconds(), result.ScheduleDelay.P999.Nanoseconds(), result.ScheduleDelay.Max.Nanoseconds(), result.RoundTrip.Samples, result.RoundTrip.P50.Nanoseconds(), result.RoundTrip.P95.Nanoseconds(), result.RoundTrip.P99.Nanoseconds(), result.RoundTrip.P999.Nanoseconds(), result.RoundTrip.Max.Nanoseconds(), result.Resources.RSSBytes, result.Resources.HeapAllocBytes, result.Resources.HeapSysBytes, result.Resources.HeapObjects, result.Resources.GCCount, result.Resources.GCPauseNanos, result.Resources.Goroutines, cfg.Payload, cfg.Connections, cfg.Messages, result.TotalRequests, result.Errors, result.Elapsed, result.Throughput)
	fmt.Fprintf(w, "%s-%d %d %.0f ns/op\n", benchmarkName(cfg.Protocol), runtime.GOMAXPROCS(0), result.TotalRequests, result.NsPerOp)
}

func benchmarkBackendLabel(cfg config) string {
	if cfg.Protocol == "http3" || cfg.Protocol == "quic-stream" {
		return "rfc9000"
	}
	return backendLabel(cfg.Backend)
}
