package udpbench

import (
	"testing"
	"time"
)

func TestPhasePacerDistributesAggregateRate(t *testing.T) {
	start := time.Unix(123, 0)
	pacer := newPhasePacer(start, 4, 100)

	if got := pacer.deadline(2, 3); !got.Equal(start.Add(140 * time.Millisecond)) {
		t.Fatalf("deadline=%s, want %s", got, start.Add(140*time.Millisecond))
	}
}

func TestScaleDurationSaturatesOnOverflow(t *testing.T) {
	if got := scaleDuration(^uint64(0), 1); got != time.Duration(1<<63-1) {
		t.Fatalf("duration=%s, want maximum duration", got)
	}
}
