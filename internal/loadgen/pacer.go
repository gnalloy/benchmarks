package loadgen

import (
	"context"
	"math"
	"math/bits"
	"time"
)

// Pacer 把多连接请求映射到统一的全局发送时间线。
type Pacer struct {
	start       time.Time
	connections uint64
	targetRate  uint64
}

// NewPacer 创建聚合目标速率控制器；targetRate 为 0 时禁用节流。
func NewPacer(start time.Time, connections int, targetRate int64) Pacer {
	if connections <= 0 || targetRate <= 0 {
		return Pacer{}
	}
	return Pacer{
		start:       start,
		connections: uint64(connections),
		targetRate:  uint64(targetRate),
	}
}

// SetStart 设置所有连接共享的测量起点。
func (p *Pacer) SetStart(start time.Time) {
	p.start = start
}

// Enabled 报告是否启用目标速率控制。
func (p Pacer) Enabled() bool {
	return p.targetRate > 0
}

// Deadline 返回指定连接和消息在全局时间线上的计划发送时间。
func (p Pacer) Deadline(clientID int, messageIndex int) time.Time {
	if !p.Enabled() {
		return time.Time{}
	}
	hi, ordinal := bits.Mul64(uint64(messageIndex), p.connections)
	if hi != 0 || math.MaxUint64-ordinal < uint64(clientID) {
		return p.start.Add(time.Duration(math.MaxInt64))
	}
	ordinal += uint64(clientID)
	return p.start.Add(scaleDuration(ordinal, p.targetRate))
}

// Wait 等待计划发送时间，并在上下文取消时可靠回收 timer 状态。
func (p Pacer) Wait(ctx context.Context, timer *time.Timer, clientID int, messageIndex int) (time.Time, error) {
	deadline := p.Deadline(clientID, messageIndex)
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
