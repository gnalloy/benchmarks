package httpbench

import "time"

// Result 是一次 HTTP/1.1 负载的汇总指标。
type Result struct {
	TotalRequests      int64
	Errors             int64
	Elapsed            time.Duration
	Throughput         float64
	NsPerOp            float64
	NegotiatedProtocol string
	Latency            Latency
	ScheduleDelay      Latency
	RoundTrip          Latency
}

// Latency 汇总采样到的延迟分位数。
type Latency struct {
	Samples int64
	P50     time.Duration
	P95     time.Duration
	P99     time.Duration
	P999    time.Duration
	Max     time.Duration
}
