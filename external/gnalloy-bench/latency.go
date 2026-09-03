package main

import (
	"sort"
	"time"
)

type latencySummary struct {
	Samples int64
	P50     time.Duration
	P95     time.Duration
	P99     time.Duration
	P999    time.Duration
	Max     time.Duration
}

type decomposedLatencySamples struct {
	total         []int64
	scheduleDelay []int64
	roundTrip     []int64
}

func latencySamplingEnabled(rate int) bool {
	return rate > 0
}

func shouldRecordLatency(messageIndex int, rate int) bool {
	return rate > 0 && messageIndex%rate == 0
}

func elapsedLatencyNanos(started time.Time) int64 {
	elapsed := time.Since(started).Nanoseconds()
	if elapsed <= 0 {
		return 1
	}
	return elapsed
}

func positiveLatencyNanos(value time.Duration) int64 {
	if value <= 0 {
		return 1
	}
	return value.Nanoseconds()
}

func nonNegativeLatencyNanos(value time.Duration) int64 {
	if value <= 0 {
		return 0
	}
	return value.Nanoseconds()
}

func newLatencySamples(messages int, rate int) []int64 {
	if !latencySamplingEnabled(rate) {
		return nil
	}
	capacity := messages / rate
	if messages%rate != 0 {
		capacity++
	}
	if capacity < 1 {
		capacity = 1
	}
	return make([]int64, 0, capacity)
}

func summarizeLatencySamples(samples []int64) latencySummary {
	if len(samples) == 0 {
		return latencySummary{}
	}
	sort.Slice(samples, func(i, j int) bool {
		return samples[i] < samples[j]
	})
	return latencySummary{
		Samples: int64(len(samples)),
		P50:     percentileDuration(samples, 0.50),
		P95:     percentileDuration(samples, 0.95),
		P99:     percentileDuration(samples, 0.99),
		P999:    percentileDuration(samples, 0.999),
		Max:     time.Duration(samples[len(samples)-1]),
	}
}

func summarizeDecomposedLatencySamples(samples []decomposedLatencySamples, values func(decomposedLatencySamples) []int64) latencySummary {
	count := 0
	for _, sample := range samples {
		count += len(values(sample))
	}
	all := make([]int64, 0, count)
	for _, sample := range samples {
		all = append(all, values(sample)...)
	}
	return summarizeLatencySamples(all)
}

func percentileDuration(sorted []int64, q float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	if q <= 0 {
		return time.Duration(sorted[0])
	}
	if q >= 1 {
		return time.Duration(sorted[len(sorted)-1])
	}
	index := int(float64(len(sorted)-1) * q)
	return time.Duration(sorted[index])
}
