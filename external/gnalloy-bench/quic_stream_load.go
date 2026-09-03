package main

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"gnalloy.org/benchmarks/internal/loadgen"
	"gnalloy.org/transport-quic"
	"gnalloy.org/transport-quic/application"
)

type quicStreamClient struct {
	conn    quic.Connection
	payload []byte
	reply   []byte
	alpn    string
	codec   application.LengthPrefixedCodec
}

func runQUICStreamBenchmark(ctx context.Context, cfg config) (benchResult, error) {
	server, err := startQUICStreamServer(ctx, cfg)
	if err != nil {
		return benchResult{}, err
	}
	defer server.stop()

	return runQUICStreamLoad(ctx, server.addr, cfg)
}

func runQUICStreamLoad(parent context.Context, addr string, cfg config) (benchResult, error) {
	resourcesBefore := captureResourceSnapshot()
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	clients, err := prepareQUICStreamClients(ctx, addr, cfg)
	if err != nil {
		return benchResult{Errors: 1}, err
	}
	stopClosing := closeQUICStreamClientsOnContext(ctx, clients)
	defer stopClosing()
	defer closeQUICStreamClients(clients)
	if err := warmupQUICStreamClients(ctx, clients, cfg); err != nil {
		return benchResult{Errors: 1}, err
	}

	var (
		successes atomic.Int64
		errorsN   atomic.Int64
		firstErr  error
		once      sync.Once
		wg        sync.WaitGroup
		samples   []decomposedLatencySamples
	)
	if latencySamplingEnabled(cfg.LatencySampleRate) {
		samples = make([]decomposedLatencySamples, cfg.Connections)
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
			var clientSamples *decomposedLatencySamples
			if samples != nil {
				capacity := estimatePerClientLatencySampleCount(cfg.Messages, cfg.LatencySampleRate)
				samples[clientID].total = make([]int64, 0, capacity)
				samples[clientID].roundTrip = make([]int64, 0, capacity)
				if pacingEnabled {
					samples[clientID].scheduleDelay = make([]int64, 0, capacity)
				}
				clientSamples = &samples[clientID]
			}
			recordError(runQUICStreamClientMessages(ctx, clients[clientID], cfg, clientID, startCh, cfg.Messages, &pacer, &successes, clientSamples))
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

	total := successes.Load()
	result := benchResult{
		TotalRequests: total,
		Errors:        errorsN.Load(),
		Elapsed:       elapsed,
		Protocol:      firstQUICStreamALPN(clients),
		Resources:     resourceDeltaSince(resourcesBefore, captureResourceSnapshot()),
	}
	if samples != nil {
		result.Latency = summarizeDecomposedLatencySamples(samples, func(sample decomposedLatencySamples) []int64 { return sample.total })
		result.ScheduleDelay = summarizeDecomposedLatencySamples(samples, func(sample decomposedLatencySamples) []int64 { return sample.scheduleDelay })
		result.RoundTrip = summarizeDecomposedLatencySamples(samples, func(sample decomposedLatencySamples) []int64 { return sample.roundTrip })
	}
	result.Throughput = float64(total) / elapsed.Seconds()
	if total > 0 {
		result.NsPerOp = float64(elapsed.Nanoseconds()) / float64(total)
	}
	if firstErr != nil {
		return result, firstErr
	}
	expected := int64(cfg.Connections * cfg.Messages)
	if total != expected {
		return result, fmt.Errorf("gnalloy-bench: completed %d requests, want %d", total, expected)
	}
	return result, nil
}

func prepareQUICStreamClients(ctx context.Context, addr string, cfg config) ([]quicStreamClient, error) {
	clients := make([]quicStreamClient, 0, cfg.Connections)
	for i := 0; i < cfg.Connections; i++ {
		conn, err := dialQUICStreamClient(ctx, addr, cfg)
		if err != nil {
			closeQUICStreamClients(clients)
			return nil, err
		}
		clients = append(clients, quicStreamClient{
			conn:    conn,
			payload: makePayload(cfg.Payload, i),
			reply:   make([]byte, cfg.Payload),
			alpn:    conn.ConnectionState().TLS.NegotiatedProtocol,
			codec:   application.LengthPrefixedCodec{MaxFrameSize: cfg.Payload},
		})
	}
	return clients, nil
}

func dialQUICStreamClient(ctx context.Context, addr string, cfg config) (quic.Connection, error) {
	tlsConfig, err := clientTLSConfig(cfg)
	if err != nil {
		return nil, err
	}
	return quic.DialAddr(ctx, addr, quic.Config{
		TLS:                   tlsConfig,
		NextProtos:            []string{quicStreamALPN(cfg)},
		MaxIncomingStreams:    -1,
		MaxIncomingUniStreams: -1,
		MaxIdleTimeout:        cfg.Timeout,
		KeepAlivePeriod:       cfg.Timeout / 4,
		InitialPacketSize:     quic.MinInitialPacketSize,
	})
}

func warmupQUICStreamClients(ctx context.Context, clients []quicStreamClient, cfg config) error {
	if cfg.WarmupMessages <= 0 {
		return nil
	}
	var (
		firstErr error
		once     sync.Once
		wg       sync.WaitGroup
	)
	recordError := func(err error) {
		if err == nil {
			return
		}
		once.Do(func() {
			firstErr = err
		})
	}
	startCh := make(chan struct{})
	for i := range clients {
		clientID := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			disabledPacer := loadgen.Pacer{}
			recordError(runQUICStreamClientMessages(ctx, clients[clientID], cfg, clientID, startCh, cfg.WarmupMessages, &disabledPacer, nil, nil))
		}()
	}
	close(startCh)
	wg.Wait()
	return firstErr
}

func closeQUICStreamClients(clients []quicStreamClient) {
	for i := range clients {
		if clients[i].conn != nil {
			_ = clients[i].conn.CloseWithError(0, "benchmark done")
		}
	}
}

func closeQUICStreamClientsOnContext(ctx context.Context, clients []quicStreamClient) func() {
	done := make(chan struct{})
	var once sync.Once
	go func() {
		select {
		case <-ctx.Done():
			closeQUICStreamClients(clients)
		case <-done:
		}
	}()
	return func() {
		once.Do(func() {
			close(done)
		})
	}
}

func runQUICStreamClientMessages(ctx context.Context, client quicStreamClient, cfg config, clientID int, startCh <-chan struct{}, messageCount int, sharedPacer *loadgen.Pacer, successes *atomic.Int64, samples *decomposedLatencySamples) error {
	select {
	case <-startCh:
	case <-ctx.Done():
		return ctx.Err()
	}
	pacer := *sharedPacer
	var timer *time.Timer
	if pacer.Enabled() {
		timer = time.NewTimer(time.Hour)
		if !timer.Stop() {
			<-timer.C
		}
		defer timer.Stop()
	}
	for i := 0; i < messageCount; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		deadline := time.Time{}
		if pacer.Enabled() {
			var err error
			deadline, err = pacer.Wait(ctx, timer, clientID, i)
			if err != nil {
				return err
			}
		}
		client.payload[0] = byte(clientID + i)
		recordLatency := samples != nil && shouldRecordLatency(i, cfg.LatencySampleRate)
		var sendStarted time.Time
		if recordLatency {
			sendStarted = time.Now()
		}
		if err := runQUICStreamRequest(ctx, client, cfg); err != nil {
			return err
		}
		if recordLatency {
			completedAt := time.Now()
			roundTrip := positiveLatencyNanos(completedAt.Sub(sendStarted))
			samples.roundTrip = append(samples.roundTrip, roundTrip)
			if pacer.Enabled() {
				samples.total = append(samples.total, positiveLatencyNanos(completedAt.Sub(deadline)))
				samples.scheduleDelay = append(samples.scheduleDelay, nonNegativeLatencyNanos(sendStarted.Sub(deadline)))
			} else {
				samples.total = append(samples.total, roundTrip)
			}
		}
		if successes != nil {
			successes.Add(1)
		}
	}
	return nil
}

func estimatePerClientLatencySampleCount(messages int, rate int) int {
	if messages <= 0 || rate <= 0 {
		return 0
	}
	count := messages / rate
	if messages%rate != 0 {
		count++
	}
	return count
}

func runQUICStreamRequest(ctx context.Context, client quicStreamClient, cfg config) error {
	stream, err := client.conn.OpenStreamSync(ctx)
	if err != nil {
		return err
	}
	if deadline, ok := ctx.Deadline(); ok {
		_ = stream.SetDeadline(deadline)
	}
	if err := client.codec.WriteFrame(stream, client.payload); err != nil {
		stream.CancelWrite(0)
		return err
	}
	if err := stream.Close(); err != nil {
		return err
	}
	reply, err := client.codec.ReadFrameInto(stream, client.reply)
	if err != nil {
		stream.CancelRead(0)
		return err
	}
	if !bytes.Equal(reply, client.payload) {
		return fmt.Errorf("gnalloy-bench: quic stream echo mismatch")
	}
	return nil
}

func firstQUICStreamALPN(clients []quicStreamClient) string {
	for i := range clients {
		if clients[i].alpn != "" {
			return clients[i].alpn
		}
	}
	return ""
}
