package udpbench

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type client struct {
	conn    net.Conn
	payload []byte
	reply   []byte
}

func Run(parent context.Context, cfg Config) (Result, error) {
	if err := cfg.Validate(); err != nil {
		return Result{}, err
	}
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithTimeout(parent, cfg.Timeout)
	defer cancel()
	clients, err := prepareClients(ctx, cfg)
	if err != nil {
		return Result{Errors: 1}, err
	}
	defer closeClients(clients)
	if _, err := runPhase(ctx, clients, cfg, cfg.WarmupMessages, nil); err != nil {
		return Result{Errors: 1}, err
	}

	samples := make([][]int64, cfg.Connections)
	if cfg.LatencySampleRate == 0 {
		samples = nil
	}
	started := time.Now()
	total, err := runPhase(ctx, clients, cfg, cfg.Messages, samples)
	elapsed := time.Since(started)
	if elapsed <= 0 {
		elapsed = time.Nanosecond
	}
	result := Result{
		TotalRequests: total,
		Elapsed:       elapsed,
		Throughput:    float64(total) / elapsed.Seconds(),
	}
	if samples != nil {
		all := make([]int64, 0, sampleCapacity(cfg.Connections*cfg.Messages, cfg.LatencySampleRate))
		for _, values := range samples {
			all = append(all, values...)
		}
		result.Latency = summarizeLatency(all)
	}
	if err != nil {
		result.Errors = 1
		return result, err
	}
	expected := int64(cfg.Connections * cfg.Messages)
	if result.TotalRequests != expected {
		result.Errors = 1
		return result, fmt.Errorf("udpbench: completed %d requests, want %d", result.TotalRequests, expected)
	}
	return result, nil
}

func prepareClients(ctx context.Context, cfg Config) ([]client, error) {
	clients := make([]client, 0, cfg.Connections)
	for id := 0; id < cfg.Connections; id++ {
		dialer := net.Dialer{Timeout: cfg.Timeout}
		conn, err := dialer.DialContext(ctx, "udp", cfg.Addr)
		if err != nil {
			closeClients(clients)
			return nil, err
		}
		if err := conn.SetDeadline(time.Now().Add(cfg.Timeout)); err != nil {
			_ = conn.Close()
			closeClients(clients)
			return nil, err
		}
		clients = append(clients, client{
			conn:    conn,
			payload: makePayload(cfg.Payload, id),
			reply:   make([]byte, cfg.Payload),
		})
	}
	return clients, nil
}

type clientResult struct {
	completed int64
	err       error
}

func runPhase(ctx context.Context, clients []client, cfg Config, messages int, samples [][]int64) (int64, error) {
	if messages == 0 {
		return 0, nil
	}
	results := make([]clientResult, len(clients))
	var wait sync.WaitGroup
	start := make(chan struct{})
	pacer := newPhasePacer(time.Time{}, len(clients), cfg.TargetRate)
	for id := range clients {
		id := id
		wait.Add(1)
		go func() {
			defer wait.Done()
			var clientSamples *[]int64
			if samples != nil {
				samples[id] = make([]int64, 0, sampleCapacity(messages, cfg.LatencySampleRate))
				clientSamples = &samples[id]
			}
			results[id].completed, results[id].err = runClient(ctx, clients[id], id, messages, cfg.LatencySampleRate, start, &pacer, clientSamples)
		}()
	}
	pacer.start = time.Now()
	close(start)
	wait.Wait()
	var total int64
	for id := range results {
		total += results[id].completed
		if results[id].err != nil {
			return total, results[id].err
		}
	}
	return total, nil
}

func runClient(ctx context.Context, client client, id int, messages int, sampleRate int, start <-chan struct{}, pacer *phasePacer, samples *[]int64) (int64, error) {
	select {
	case <-start:
	case <-ctx.Done():
		return 0, ctx.Err()
	}
	if pacer.enabled() {
		return runPacedClient(ctx, client, id, messages, sampleRate, pacer, samples)
	}
	return runSaturatedClient(ctx, client, id, messages, sampleRate, samples)
}

func runSaturatedClient(ctx context.Context, client client, id int, messages int, sampleRate int, samples *[]int64) (int64, error) {
	var completed int64
	for index := 0; index < messages; index++ {
		if err := ctx.Err(); err != nil {
			return completed, err
		}
		client.payload[0] = byte(id + index)
		record := samples != nil && sampleRate > 0 && index%sampleRate == 0
		var started time.Time
		if record {
			started = time.Now()
		}
		if err := writeFull(client.conn, client.payload); err != nil {
			return completed, err
		}
		n, err := client.conn.Read(client.reply)
		if err != nil {
			return completed, err
		}
		if n != len(client.payload) || !bytes.Equal(client.reply[:n], client.payload) {
			return completed, fmt.Errorf("udpbench: echo mismatch")
		}
		if record {
			elapsed := time.Since(started).Nanoseconds()
			if elapsed <= 0 {
				elapsed = 1
			}
			*samples = append(*samples, elapsed)
		}
		completed++
	}
	return completed, nil
}

func runPacedClient(ctx context.Context, client client, id int, messages int, sampleRate int, pacer *phasePacer, samples *[]int64) (int64, error) {
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
		client.payload[0] = byte(id + index)
		deadline, err := pacer.wait(ctx, timer, id, index)
		if err != nil {
			return completed, err
		}
		record := samples != nil && sampleRate > 0 && index%sampleRate == 0
		if err := writeFull(client.conn, client.payload); err != nil {
			return completed, err
		}
		n, err := client.conn.Read(client.reply)
		if err != nil {
			return completed, err
		}
		if n != len(client.payload) || !bytes.Equal(client.reply[:n], client.payload) {
			return completed, fmt.Errorf("udpbench: echo mismatch")
		}
		if record {
			elapsed := time.Since(deadline).Nanoseconds()
			if elapsed <= 0 {
				elapsed = 1
			}
			*samples = append(*samples, elapsed)
		}
		completed++
	}
	return completed, nil
}

func writeFull(writer io.Writer, payload []byte) error {
	n, err := writer.Write(payload)
	if err != nil {
		return err
	}
	if n != len(payload) {
		return io.ErrShortWrite
	}
	return nil
}

func closeClients(clients []client) {
	for index := range clients {
		_ = clients[index].conn.Close()
	}
}

func makePayload(size int, seed int) []byte {
	payload := make([]byte, size)
	for index := range payload {
		payload[index] = byte(seed + index)
	}
	return payload
}

func sampleCapacity(messages int, rate int) int {
	if rate <= 0 || messages <= 0 {
		return 0
	}
	return (messages + rate - 1) / rate
}
