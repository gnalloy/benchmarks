package httpbench

import (
	"context"
	"crypto/tls"
	"net"
	"time"
)

type client struct {
	conn     net.Conn
	protocol *protocolClient
	alpn     string
}

func prepareClients(ctx context.Context, cfg Config) ([]client, error) {
	clients := make([]client, 0, cfg.Connections)
	host := requestHost(cfg)
	expectedBody := ResponseBody(cfg.Payload)
	for range cfg.Connections {
		conn, alpn, err := dial(ctx, cfg)
		if err != nil {
			closeClients(clients)
			return nil, err
		}
		protocol, err := newProtocolClient(conn, host, expectedBody)
		if err != nil {
			_ = conn.Close()
			closeClients(clients)
			return nil, err
		}
		clients = append(clients, client{conn: conn, protocol: protocol, alpn: alpn})
	}
	return clients, nil
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

func dial(ctx context.Context, cfg Config) (net.Conn, string, error) {
	dialer := net.Dialer{Timeout: cfg.Timeout}
	conn, err := dialer.DialContext(ctx, "tcp", cfg.Addr)
	if err != nil {
		return nil, "", err
	}
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		_ = tcpConn.SetNoDelay(true)
	}
	if err := conn.SetDeadline(time.Now().Add(cfg.Timeout)); err != nil {
		_ = conn.Close()
		return nil, "", err
	}
	if cfg.TLS == nil {
		return conn, "", nil
	}
	tlsConfig := cfg.TLS.Clone()
	if tlsConfig.ServerName == "" {
		tlsConfig.ServerName = cfg.ServerName
	}
	tlsConn := tls.Client(conn, tlsConfig)
	if err := tlsConn.HandshakeContext(ctx); err != nil {
		_ = conn.Close()
		return nil, "", err
	}
	return tlsConn, tlsConn.ConnectionState().NegotiatedProtocol, nil
}

func closeClients(clients []client) {
	for index := range clients {
		clients[index].close()
	}
}

func (c *client) close() {
	if c.protocol != nil {
		c.protocol.close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

func firstNegotiatedProtocol(clients []client) string {
	for index := range clients {
		if clients[index].alpn != "" {
			return clients[index].alpn
		}
	}
	return ""
}
