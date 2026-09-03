package benchh2

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	http2 "gnalloy.org/codec-http2"
	"gnalloy.org/gnalloy/bootstrap"
	"gnalloy.org/gnalloy/channel"
	"gnalloy.org/gnalloy/transport"
	"gnalloy.org/transport-tcp"
)

const closeTimeout = time.Second

type client struct {
	channel    channel.Channel
	response   *responseHandler
	fields     []http2.HeaderField
	nextStream http2.StreamID
}

func prepareClients(ctx context.Context, cfg Config, group *transport.EventLoopGroup) ([]client, error) {
	clients := make([]client, 0, cfg.Connections)
	clientTransport := newClientTransport(cfg)
	for range cfg.Connections {
		prepared, err := dialClient(ctx, cfg, group, clientTransport)
		if err != nil {
			closeClients(clients)
			return nil, err
		}
		clients = append(clients, prepared)
	}
	return clients, nil
}

func newClientTransport(cfg Config) *tcp.Transport {
	tcpConfig := tcp.DefaultConfig()
	tcpConfig.ConnectTimeoutMillis = int(cfg.Timeout / time.Millisecond)
	tcpConfig.ReadBufferSize = max(cfg.Payload+http2.FrameHeaderSize, http2.DefaultMaxFrameSize)
	return tcp.NewTransport(tcpConfig)
}

func dialClient(ctx context.Context, cfg Config, group *transport.EventLoopGroup, clientTransport *tcp.Transport) (client, error) {
	expected := ResponseBody(cfg.Payload)
	response := newResponseHandler(expected)
	ch, err := bootstrap.NewDialer().
		Group(group).
		Transport(clientTransport).
		Initializer(func(ch channel.Channel) error {
			return addClientPipeline(ch, cfg, response)
		}).
		DialContext(ctx, cfg.Addr)
	if err != nil {
		return client{}, err
	}
	select {
	case <-response.ready:
	case <-ctx.Done():
		_ = ch.Close()
		return client{}, ctx.Err()
	}
	if err := responseReadyError(response); err != nil {
		_ = ch.Close()
		return client{}, err
	}
	return client{
		channel:    ch,
		response:   response,
		fields:     requestFields(cfg.ServerName, cfg.TLS != nil),
		nextStream: 1,
	}, nil
}

func responseReadyError(response *responseHandler) error {
	select {
	case result := <-response.responses:
		return result.err
	default:
		return nil
	}
}

func requestFields(authority string, tlsEnabled bool) []http2.HeaderField {
	if authority == "" {
		authority = "127.0.0.1"
	}
	scheme := "http"
	if tlsEnabled {
		scheme = "https"
	}
	return []http2.HeaderField{
		{Name: ":method", Value: "GET"},
		{Name: ":scheme", Value: scheme},
		{Name: ":path", Value: "/bench"},
		{Name: ":authority", Value: authority},
	}
}

func warmupClients(ctx context.Context, clients []client, cfg Config) error {
	if cfg.WarmupMessages <= 0 {
		return nil
	}
	var (
		firstErr error
		once     sync.Once
		wg       sync.WaitGroup
	)
	startCh := make(chan struct{})
	for i := range clients {
		clientID := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := runClientMessages(ctx, &clients[clientID], cfg.WarmupMessages, 0, startCh, nil, nil)
			if err != nil {
				once.Do(func() { firstErr = err })
			}
		}()
	}
	close(startCh)
	wg.Wait()
	return firstErr
}

func runClientMessages(ctx context.Context, c *client, messageCount int, latencySampleRate int, startCh <-chan struct{}, successes *atomic.Int64, latencySamples *[]int64) error {
	select {
	case <-startCh:
	case <-ctx.Done():
		return ctx.Err()
	}
	for i := 0; i < messageCount; i++ {
		streamID := c.nextStream
		c.nextStream += 2
		recordLatency := latencySamples != nil && shouldRecordLatency(i, latencySampleRate)
		var requestStarted time.Time
		if recordLatency {
			requestStarted = time.Now()
		}
		request := http2.HeadersBlock{StreamID: streamID, Fields: c.fields, EndStream: true}
		if err := c.channel.WriteAndFlush(request); err != nil {
			return err
		}
		select {
		case response := <-c.response.responses:
			if response.err != nil {
				return response.err
			}
			if response.streamID != streamID {
				return fmt.Errorf("benchh2: response stream %d, want %d", response.streamID, streamID)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
		if recordLatency {
			*latencySamples = append(*latencySamples, elapsedLatencyNanos(requestStarted))
		}
		if successes != nil {
			successes.Add(1)
		}
	}
	return nil
}

func closeClients(clients []client) {
	for i := range clients {
		if clients[i].channel == nil {
			continue
		}
		future := clients[i].channel.CloseFuture()
		_, _ = future.AwaitTimeout(closeTimeout)
	}
}

func closeClientsOnContext(ctx context.Context, clients []client) func() {
	done := make(chan struct{})
	var once sync.Once
	go func() {
		select {
		case <-ctx.Done():
			closeClients(clients)
		case <-done:
		}
	}()
	return func() {
		once.Do(func() { close(done) })
	}
}

func shutdownGroup(group *transport.EventLoopGroup) {
	ctx, cancel := context.WithTimeout(context.Background(), closeTimeout)
	defer cancel()
	_ = group.Shutdown(ctx)
}

func firstNegotiatedProtocol(clients []client) string {
	for i := range clients {
		if clients[i].response.alpn != "" {
			return clients[i].response.alpn
		}
	}
	return ""
}
