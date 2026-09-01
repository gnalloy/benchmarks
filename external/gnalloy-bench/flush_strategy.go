package main

import (
	"fmt"
	"strings"

	"gnalloy.org/gnalloy/channel"
)

const (
	defaultFlushStrategyName = flushStrategyEventLoopBatch

	flushStrategyImmediate      = "immediate"
	flushStrategyReadComplete   = "read-complete"
	flushStrategyEventLoopBatch = "event-loop-batch"
)

func normalizeFlushStrategy(name string) (string, channel.FlushStrategy, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "", flushStrategyEventLoopBatch, "batch":
		return flushStrategyEventLoopBatch, channel.FlushOnEventLoopBatch, nil
	case flushStrategyImmediate:
		return flushStrategyImmediate, channel.FlushImmediate, nil
	case flushStrategyReadComplete:
		return flushStrategyReadComplete, channel.FlushOnReadComplete, nil
	default:
		return "", 0, fmt.Errorf("%w: flush-strategy %s", errInvalidConfig, name)
	}
}
