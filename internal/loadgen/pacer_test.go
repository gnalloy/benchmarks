package loadgen

import (
	"testing"
	"time"
)

func TestPacerDistributesAggregateRate(t *testing.T) {
	start := time.Unix(123, 0)
	pacer := NewPacer(start, 4, 100)

	if got := pacer.Deadline(2, 3); !got.Equal(start.Add(140 * time.Millisecond)) {
		t.Fatalf("deadline=%s, want %s", got, start.Add(140*time.Millisecond))
	}
}

func TestPacerDisablesInvalidRateOrTopology(t *testing.T) {
	for _, pacer := range []Pacer{
		NewPacer(time.Now(), 1, 0),
		NewPacer(time.Now(), 1, -1),
		NewPacer(time.Now(), 0, 100),
	} {
		if pacer.Enabled() {
			t.Fatal("invalid pacing configuration was enabled")
		}
	}
}

func TestPacerDeadlineSaturatesOnOverflow(t *testing.T) {
	pacer := NewPacer(time.Time{}, int(^uint(0)>>1), 1)
	if got := pacer.Deadline(int(^uint(0)>>1), int(^uint(0)>>1)); got.Sub(time.Time{}) != time.Duration(1<<63-1) {
		t.Fatalf("duration=%s, want maximum duration", got.Sub(time.Time{}))
	}
}
