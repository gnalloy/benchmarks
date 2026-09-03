package httpbench

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type clientResult struct {
	completed int64
	err       error
}

// Run 执行 HTTP/1.1 keep-alive 负载。
func Run(parent context.Context, cfg Config) (Result, error) {
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
	defer shutdownClientGroup(group)
	clients, err := prepareClients(ctx, cfg, group)
	if err != nil {
		return Result{Errors: 1}, err
	}
	defer closeClients(clients)
	if _, _, _, err := runPhase(ctx, clients, cfg, cfg.WarmupMessages, false); err != nil {
		return Result{Errors: 1}, err
	}

	total, samples, elapsed, err := runPhase(ctx, clients, cfg, cfg.Messages, true)
	if elapsed <= 0 {
		elapsed = time.Nanosecond
	}
	result := Result{
		TotalRequests:      total,
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
	if err != nil {
		result.Errors = 1
		return result, err
	}
	expected := int64(cfg.Connections * cfg.Messages)
	if total != expected {
		result.Errors = 1
		return result, fmt.Errorf("httpbench: completed %d requests, want %d", total, expected)
	}
	return result, nil
}

func runPhase(ctx context.Context, clients []client, cfg Config, messages int, measured bool) (int64, []clientSamples, time.Duration, error) {
	if messages == 0 {
		return 0, nil, 0, nil
	}
	results := make([]clientResult, len(clients))
	var samples []clientSamples
	if measured && cfg.LatencySampleRate > 0 {
		samples = make([]clientSamples, len(clients))
	}
	phaseCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	targetRate := cfg.TargetRate
	if !measured {
		targetRate = 0
	}
	pacer := newPhasePacer(time.Time{}, len(clients), targetRate)
	start := make(chan struct{})
	var wait sync.WaitGroup
	for clientID := range clients {
		clientID := clientID
		var clientSample *clientSamples
		if samples != nil {
			capacity := sampleCapacity(messages, cfg.LatencySampleRate)
			samples[clientID].total = make([]int64, 0, capacity)
			samples[clientID].roundTrip = make([]int64, 0, capacity)
			if pacer.Enabled() {
				samples[clientID].scheduleDelay = make([]int64, 0, capacity)
			}
			clientSample = &samples[clientID]
		}
		wait.Add(1)
		go func(sample *clientSamples) {
			defer wait.Done()
			results[clientID].completed, results[clientID].err = runClient(
				phaseCtx, &clients[clientID], clientID, messages, cfg.LatencySampleRate, start, &pacer, sample,
			)
			if results[clientID].err != nil {
				cancel()
			}
		}(clientSample)
	}
	phaseStart := time.Now()
	pacer.SetStart(phaseStart)
	close(start)
	wait.Wait()
	elapsed := time.Since(phaseStart)
	var total int64
	var firstErr error
	for index := range results {
		total += results[index].completed
		if firstErr == nil && results[index].err != nil {
			firstErr = results[index].err
		}
	}
	return total, samples, elapsed, firstErr
}

func runClient(
	ctx context.Context,
	client *client,
	clientID int,
	messages int,
	sampleRate int,
	start <-chan struct{},
	pacer *phasePacer,
	samples *clientSamples,
) (int64, error) {
	select {
	case <-start:
	case <-ctx.Done():
		return 0, ctx.Err()
	}
	localPacer := *pacer
	if localPacer.Enabled() {
		return runPacedClient(ctx, client, clientID, messages, sampleRate, localPacer, samples)
	}
	return runSaturatedClient(ctx, client, messages, sampleRate, samples)
}

func runSaturatedClient(ctx context.Context, client *client, messages int, sampleRate int, samples *clientSamples) (int64, error) {
	var completed int64
	for index := 0; index < messages; index++ {
		if err := ctx.Err(); err != nil {
			return completed, err
		}
		record := samples != nil && index%sampleRate == 0
		var started time.Time
		if record {
			started = time.Now()
		}
		if err := exchange(ctx, client); err != nil {
			return completed, err
		}
		if record {
			elapsed := positiveNanoseconds(time.Since(started))
			samples.total = append(samples.total, elapsed)
			samples.roundTrip = append(samples.roundTrip, elapsed)
		}
		completed++
	}
	return completed, nil
}

func runPacedClient(ctx context.Context, client *client, clientID int, messages int, sampleRate int, pacer phasePacer, samples *clientSamples) (int64, error) {
	timer := time.NewTimer(time.Hour)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()
	var completed int64
	for index := 0; index < messages; index++ {
		if err := ctx.Err(); err != nil {
			return completed, err
		}
		deadline, err := pacer.Wait(ctx, timer, clientID, index)
		if err != nil {
			return completed, err
		}
		record := samples != nil && index%sampleRate == 0
		var sendStarted time.Time
		if record {
			sendStarted = time.Now()
		}
		if err := exchange(ctx, client); err != nil {
			return completed, err
		}
		if record {
			completedAt := time.Now()
			samples.total = append(samples.total, positiveNanoseconds(completedAt.Sub(deadline)))
			samples.scheduleDelay = append(samples.scheduleDelay, nonNegativeNanoseconds(sendStarted.Sub(deadline)))
			samples.roundTrip = append(samples.roundTrip, positiveNanoseconds(completedAt.Sub(sendStarted)))
		}
		completed++
	}
	return completed, nil
}

func exchange(ctx context.Context, client *client) error {
	return client.protocol.exchange(ctx)
}
