package udpbench

import (
	"time"

	"gnalloy.org/benchmarks/internal/loadgen"
)

type phasePacer = loadgen.Pacer

func newPhasePacer(start time.Time, connections int, targetRate int64) phasePacer {
	return loadgen.NewPacer(start, connections, targetRate)
}
