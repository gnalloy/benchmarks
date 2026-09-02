package udpbench

import (
	"sort"
	"time"
)

type Result struct {
	TotalRequests int64
	Errors        int64
	Elapsed       time.Duration
	Throughput    float64
	Latency       Latency
	ScheduleDelay Latency
	RoundTrip     Latency
}

type Latency struct {
	Samples int64
	P50     time.Duration
	P95     time.Duration
	P99     time.Duration
	P999    time.Duration
	Max     time.Duration
}

func summarizeLatency(samples []int64) Latency {
	if len(samples) == 0 {
		return Latency{}
	}
	sort.Slice(samples, func(i int, j int) bool {
		return samples[i] < samples[j]
	})
	return Latency{
		Samples: int64(len(samples)),
		P50:     percentile(samples, 0.50),
		P95:     percentile(samples, 0.95),
		P99:     percentile(samples, 0.99),
		P999:    percentile(samples, 0.999),
		Max:     time.Duration(samples[len(samples)-1]),
	}
}

func percentile(sorted []int64, quantile float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	if quantile <= 0 {
		return time.Duration(sorted[0])
	}
	if quantile >= 1 {
		return time.Duration(sorted[len(sorted)-1])
	}
	index := int(float64(len(sorted)-1) * quantile)
	return time.Duration(sorted[index])
}
