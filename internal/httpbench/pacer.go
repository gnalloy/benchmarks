package httpbench

import (
	"context"
	"math"
	"math/bits"
	"time"
)

type phasePacer struct {
	start       time.Time
	connections uint64
	targetRate  uint64
}

func newPhasePacer(start time.Time, connections int, targetRate int64) phasePacer {
	return phasePacer{
		start:       start,
		connections: uint64(connections),
		targetRate:  uint64(targetRate),
	}
}

func (p phasePacer) enabled() bool {
	return p.targetRate > 0
}

func (p phasePacer) deadline(clientID int, messageIndex int) time.Time {
	if !p.enabled() {
		return time.Time{}
	}
	ordinal := uint64(messageIndex)*p.connections + uint64(clientID)
	return p.start.Add(scaleDuration(ordinal, p.targetRate))
}

func (p phasePacer) wait(ctx context.Context, timer *time.Timer, clientID int, messageIndex int) (time.Time, error) {
	deadline := p.deadline(clientID, messageIndex)
	delay := time.Until(deadline)
	if delay <= 0 {
		return deadline, nil
	}
	timer.Reset(delay)
	select {
	case <-timer.C:
		return deadline, nil
	case <-ctx.Done():
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		return time.Time{}, ctx.Err()
	}
}

func scaleDuration(ordinal uint64, rate uint64) time.Duration {
	hi, lo := bits.Mul64(ordinal, uint64(time.Second))
	if rate == 0 || hi >= rate {
		return time.Duration(math.MaxInt64)
	}
	value, _ := bits.Div64(hi, lo, rate)
	if value > math.MaxInt64 {
		return time.Duration(math.MaxInt64)
	}
	return time.Duration(value)
}
