package benchh3

import (
	"testing"
	"time"
)

func TestPositiveLatencyNanosKeepsPositiveFloor(t *testing.T) {
	if got := positiveLatencyNanos(0); got != 1 {
		t.Fatalf("latency=%d, want positive lower bound", got)
	}
	if got := positiveLatencyNanos(time.Nanosecond); got != 1 {
		t.Fatalf("latency=%d, want positive lower bound", got)
	}
}

func TestNonNegativeLatencyNanosClampsEarlySend(t *testing.T) {
	if got := nonNegativeLatencyNanos(-time.Nanosecond); got != 0 {
		t.Fatalf("latency=%d, want zero", got)
	}
}
