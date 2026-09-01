package main

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"gnalloy.org/transport-quic"
)

type quicStreamClient struct {
	conn    quic.Connection
	payload []byte
	reply   []byte
	alpn    string
}

func runQUICStreamBenchmark(ctx context.Context, cfg config) (benchResult, error) {
	server, err := startQUICStreamServer(ctx, cfg)
	if err != nil {
		return benchResult{}, err
	}
	defer server.stop()

	resourcesBefore := captureResourceSnapshot()
	result, err := runQUICStreamLoad(ctx, server.addr, cfg)
	result.Resources = resourceDeltaSince(resourcesBefore, captureResourceSnapshot())
	if err != nil {
		return result, err
	}
	return result, nil
}

func runQUICStreamLoad(parent context.Context, addr string, cfg config) (benchResult, error) {
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
		samples   [][]int64
	)
	if latencySamplingEnabled(cfg.LatencySampleRate) {
		samples = make([][]int64, cfg.Connections)
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
	for i := range clients {
		clientID := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			var clientSamples *[]int64
			if samples != nil {
				samples[clientID] = newLatencySamples(cfg.Messages, cfg.LatencySampleRate)
				clientSamples = &samples[clientID]
			}
			recordError(runQUICStreamClientMessages(ctx, clients[clientID], cfg, clientID, startCh, cfg.Messages, &successes, clientSamples))
		}()
	}
	start := time.Now()
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
	}
	if samples != nil {
		allSamples := make([]int64, 0, estimateLatencySampleCount(cfg.Connections, cfg.Messages, cfg.LatencySampleRate))
		for _, clientSamples := range samples {
			allSamples = append(allSamples, clientSamples...)
		}
		result.Latency = summarizeLatencySamples(allSamples)
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
			recordError(runQUICStreamClientMessages(ctx, clients[clientID], cfg, clientID, startCh, cfg.WarmupMessages, nil, nil))
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

func runQUICStreamClientMessages(ctx context.Context, client quicStreamClient, cfg config, clientID int, startCh <-chan struct{}, messageCount int, successes *atomic.Int64, latencySamples *[]int64) error {
	select {
	case <-startCh:
	case <-ctx.Done():
		return ctx.Err()
	}
	for i := 0; i < messageCount; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		client.payload[0] = byte(clientID + i)
		recordLatency := latencySamples != nil && shouldRecordLatency(i, cfg.LatencySampleRate)
		var requestStarted time.Time
		if recordLatency {
			requestStarted = time.Now()
		}
		if err := runQUICStreamRequest(ctx, client, cfg); err != nil {
			return err
		}
		if recordLatency && latencySamples != nil {
			*latencySamples = append(*latencySamples, elapsedLatencyNanos(requestStarted))
		}
		if successes != nil {
			successes.Add(1)
		}
	}
	return nil
}

func runQUICStreamRequest(ctx context.Context, client quicStreamClient, cfg config) error {
	stream, err := client.conn.OpenStreamSync(ctx)
	if err != nil {
		return err
	}
	if deadline, ok := ctx.Deadline(); ok {
		_ = stream.SetDeadline(deadline)
	}
	if err := writeQUICStreamFrame(stream, client.payload); err != nil {
		stream.CancelWrite(0)
		return err
	}
	if err := stream.Close(); err != nil {
		return err
	}
	reply, err := readQUICStreamFrame(stream, client.reply, cfg.Payload)
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
