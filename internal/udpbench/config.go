package udpbench

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrInvalidConfig = errors.New("udpbench: invalid config")

type Config struct {
	Addr              string
	Payload           int
	Connections       int
	Messages          int
	WarmupMessages    int
	LatencySampleRate int
	TargetRate        int64
	Timeout           time.Duration
}

func DefaultConfig() Config {
	return Config{
		Addr:              "127.0.0.1:19090",
		Payload:           128,
		Connections:       64,
		Messages:          20000,
		WarmupMessages:    1000,
		LatencySampleRate: 64,
		Timeout:           5 * time.Minute,
	}
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.Addr) == "" {
		return fmt.Errorf("%w: empty address", ErrInvalidConfig)
	}
	if c.Payload <= 0 || c.Connections <= 0 || c.Messages <= 0 {
		return fmt.Errorf("%w: payload, connections and messages must be positive", ErrInvalidConfig)
	}
	if c.Connections > int(^uint(0)>>1)/c.Messages {
		return fmt.Errorf("%w: connections * messages overflows int", ErrInvalidConfig)
	}
	if c.WarmupMessages < 0 || c.LatencySampleRate < 0 || c.TargetRate < 0 {
		return fmt.Errorf("%w: warmup, latency sample rate and target rate must not be negative", ErrInvalidConfig)
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("%w: timeout must be positive", ErrInvalidConfig)
	}
	return nil
}
