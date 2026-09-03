package benchh2

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"gnalloy.org/benchmarks/internal/loadgen"
	"gnalloy.org/gnalloy/transport"
)

// RunLoad 通过 Gnalloy TCP、TLS Handler 和 HTTP/2 codec 执行负载。
func RunLoad(parent context.Context, cfg Config) (Result, error) {
	if err := cfg.Validate(); err != nil {
		return Result{}, err
	}
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithTimeout(parent, cfg.Timeout)
	defer cancel()

	group, err := newClientGroup(cfg.Connections)
	if err != nil {
		return Result{Errors: 1}, err
	}
	defer shutdownGroup(group)
	clients, err := prepareClients(ctx, cfg, group)
	if err != nil {
		return Result{Errors: 1}, err
	}
	stopClosing := closeClientsOnContext(ctx, clients)
	defer stopClosing()
	defer closeClients(clients)
	if err := warmupClients(ctx, clients, cfg); err != nil {
		return Result{Errors: 1}, err
	}

	var (
		successes atomic.Int64
		errorsN   atomic.Int64
		firstErr  error
		once      sync.Once
		wg        sync.WaitGroup
		samples   []clientSamples
	)
	if latencySamplingEnabled(cfg.LatencySampleRate) {
		samples = make([]clientSamples, cfg.Connections)
	}
	recordError := func(err error) {
		if err == nil {
			return
		}
		errorsN.Add(1)
		once.Do(func() {
			firstErr = err
			cancel()
		})
	}

	startCh := make(chan struct{})
	pacer := loadgen.NewPacer(time.Time{}, len(clients), cfg.TargetRate)
	pacingEnabled := pacer.Enabled()
	for i := range clients {
		clientID := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			var clientSample *clientSamples
			if samples != nil {
				capacity := estimateLatencySampleCount(1, cfg.Messages, cfg.LatencySampleRate)
				samples[clientID].total = make([]int64, 0, capacity)
				samples[clientID].roundTrip = make([]int64, 0, capacity)
				if pacingEnabled {
					samples[clientID].scheduleDelay = make([]int64, 0, capacity)
				}
				clientSample = &samples[clientID]
			}
			recordError(runClientMessages(ctx, &clients[clientID], clientID, cfg.Messages, cfg.LatencySampleRate, startCh, &pacer, &successes, clientSample))
		}()
	}
	start := time.Now()
	pacer.SetStart(start)
	close(startCh)
	wg.Wait()
	elapsed := time.Since(start)
	if elapsed <= 0 {
		elapsed = time.Nanosecond
	}
	return loadResult(cfg, clients, samples, successes.Load(), errorsN.Load(), elapsed, firstErr)
}

func newClientGroup(connections int) (*transport.EventLoopGroup, error) {
	size := runtime.GOMAXPROCS(0)
	if size > connections {
		size = connections
	}
	return transport.NewEventLoopGroup(transport.EventLoopGroupConfig{
		Size:         size,
		PollerConfig: transport.Config{Backend: transport.DefaultBackend()},
	})
}

func loadResult(cfg Config, clients []client, samples []clientSamples, total int64, errorsN int64, elapsed time.Duration, firstErr error) (Result, error) {
	result := Result{
		TotalRequests:      total,
		Errors:             errorsN,
		Elapsed:            elapsed,
		Throughput:         float64(total) / elapsed.Seconds(),
		NegotiatedProtocol: firstNegotiatedProtocol(clients),
	}
	if total > 0 {
		result.NsPerOp = float64(elapsed.Nanoseconds()) / float64(total)
	}
	if samples != nil {
		result.Latency = summarizeClientSamples(samples, func(sample clientSamples) []int64 { return sample.total })
		result.ScheduleDelay = summarizeClientSamples(samples, func(sample clientSamples) []int64 { return sample.scheduleDelay })
		result.RoundTrip = summarizeClientSamples(samples, func(sample clientSamples) []int64 { return sample.roundTrip })
	}
	if firstErr != nil {
		return result, firstErr
	}
	expected := int64(cfg.Connections * cfg.Messages)
	if total != expected {
		return result, fmt.Errorf("benchh2: completed %d requests, want %d", total, expected)
	}
	return result, nil
}
