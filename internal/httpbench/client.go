package httpbench

import (
	"context"
	"net"
	"runtime"
	"time"

	"gnalloy.org/gnalloy/bootstrap"
	"gnalloy.org/gnalloy/channel"
	"gnalloy.org/gnalloy/transport"
	"gnalloy.org/transport-tcp"
)

const clientCloseTimeout = time.Second

type client struct {
	channel  channel.Channel
	protocol *protocolClient
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

func prepareClients(ctx context.Context, cfg Config, group *transport.EventLoopGroup) ([]client, error) {
	clients := make([]client, 0, cfg.Connections)
	host := requestHost(cfg)
	expectedBody := ResponseBody(cfg.Payload)
	for range cfg.Connections {
		prepared, err := dialClient(ctx, cfg, group, host, expectedBody)
		if err != nil {
			closeClients(clients)
			return nil, err
		}
		clients = append(clients, prepared)
	}
	return clients, nil
}

func dialClient(ctx context.Context, cfg Config, group *transport.EventLoopGroup, host string, expectedBody []byte) (client, error) {
	protocol := newProtocolClient(host, expectedBody)
	tcpConfig := tcp.DefaultConfig()
	tcpConfig.ConnectTimeoutMillis = connectTimeoutMillis(cfg.Timeout)
	tcpConfig.ReadBufferSize = max(clientReadBufferSize, cfg.Payload)
	ch, err := bootstrap.NewDialer().
		Group(group).
		Transport(tcp.NewTransport(tcpConfig)).
		Initializer(func(ch channel.Channel) error {
			return protocol.addPipeline(ch, cfg)
		}).
		DialContext(ctx, cfg.Addr)
	if err != nil {
		return client{}, err
	}
	select {
	case <-protocol.response.ready:
	case <-ctx.Done():
		_ = ch.Close()
		return client{}, ctx.Err()
	}
	if err := protocol.response.readyError(); err != nil {
		_ = ch.Close()
		return client{}, err
	}
	protocol.channel = ch
	return client{channel: ch, protocol: protocol}, nil
}

func requestHost(cfg Config) string {
	if cfg.ServerName != "" {
		return cfg.ServerName
	}
	host, _, err := net.SplitHostPort(cfg.Addr)
	if err == nil && host != "" {
		return host
	}
	return cfg.Addr
}

func connectTimeoutMillis(timeout time.Duration) int {
	millis := timeout.Milliseconds()
	maxInt := int64(^uint(0) >> 1)
	if millis > maxInt {
		return int(maxInt)
	}
	return int(millis)
}

func closeClients(clients []client) {
	for index := range clients {
		clients[index].close()
	}
}

func (c *client) close() {
	if c.channel == nil {
		return
	}
	future := c.channel.CloseFuture()
	_, _ = future.AwaitTimeout(clientCloseTimeout)
}

func shutdownClientGroup(group *transport.EventLoopGroup) {
	ctx, cancel := context.WithTimeout(context.Background(), clientCloseTimeout)
	defer cancel()
	_ = group.Shutdown(ctx)
}

func firstNegotiatedProtocol(clients []client) string {
	for index := range clients {
		if clients[index].protocol != nil && clients[index].protocol.response.alpn != "" {
			return clients[index].protocol.response.alpn
		}
	}
	return ""
}
