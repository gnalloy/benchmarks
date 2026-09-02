package httpbench

import (
	"sort"
	"time"
)

type clientSamples struct {
	total         []int64
	scheduleDelay []int64
	roundTrip     []int64
}

func summarizeClientSamples(samples []clientSamples, values func(clientSamples) []int64) Latency {
	total := 0
	for _, sample := range samples {
		total += len(values(sample))
	}
	all := make([]int64, 0, total)
	for _, sample := range samples {
		all = append(all, values(sample)...)
	}
	return summarizeLatency(all)
}

func summarizeLatency(samples []int64) Latency {
	if len(samples) == 0 {
		return Latency{}
	}
	sort.Slice(samples, func(left int, right int) bool {
		return samples[left] < samples[right]
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

func sampleCapacity(messages int, rate int) int {
	if rate <= 0 || messages <= 0 {
		return 0
	}
	return (messages + rate - 1) / rate
}

func positiveNanoseconds(value time.Duration) int64 {
	if value <= 0 {
		return 1
	}
	return value.Nanoseconds()
}

func nonNegativeNanoseconds(value time.Duration) int64 {
	if value <= 0 {
		return 0
	}
	return value.Nanoseconds()
}
