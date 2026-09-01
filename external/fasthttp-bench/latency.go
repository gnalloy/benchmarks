package main

import "time"

type latencySummary struct {
	Samples int64
	P50     time.Duration
	P95     time.Duration
	P99     time.Duration
	P999    time.Duration
	Max     time.Duration
}
