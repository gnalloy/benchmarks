package httpbench

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"time"
)

var (
	statusOKHeader      = []byte("HTTP/1.1 200")
	contentLengthHeader = []byte("Content-Length")
)

type client struct {
	conn    net.Conn
	reader  *bufio.Reader
	request []byte
	body    []byte
	reply   []byte
	alpn    string
}

func prepareClients(ctx context.Context, cfg Config) ([]client, error) {
	clients := make([]client, 0, cfg.Connections)
	host := requestHost(cfg)
	for range cfg.Connections {
		conn, alpn, err := dial(ctx, cfg)
		if err != nil {
			closeClients(clients)
			return nil, err
		}
		clients = append(clients, client{
			conn:    conn,
			reader:  bufio.NewReaderSize(conn, 16*1024),
			request: RequestBytes(host),
			body:    ResponseBody(cfg.Payload),
			reply:   make([]byte, cfg.Payload),
			alpn:    alpn,
		})
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

func readResponse(reader *bufio.Reader, reply []byte, expected []byte) error {
	status, err := reader.ReadSlice('\n')
	if err != nil {
		return err
	}
	if !bytes.HasPrefix(status, statusOKHeader) {
		return fmt.Errorf("httpbench: unexpected status %q", bytes.TrimSpace(status))
	}
	contentLength := -1
	for {
		line, err := reader.ReadSlice('\n')
		if err != nil {
			return err
		}
		if bytes.Equal(line, []byte("\r\n")) {
			break
		}
		if value, ok := headerValue(line, contentLengthHeader); ok {
			contentLength, err = parseDecimal(value)
			if err != nil {
				return err
			}
		}
	}
	if contentLength != len(expected) {
		return fmt.Errorf("httpbench: content length %d, want %d", contentLength, len(expected))
	}
	if len(reply) < contentLength {
		return fmt.Errorf("httpbench: reply buffer too small")
	}
	if _, err := io.ReadFull(reader, reply[:contentLength]); err != nil {
		return err
	}
	if !bytes.Equal(reply[:contentLength], expected) {
		return fmt.Errorf("httpbench: response body mismatch")
	}
	return nil
}

func headerValue(line []byte, name []byte) ([]byte, bool) {
	key, value, ok := bytes.Cut(line, []byte(":"))
	if !ok || !bytes.EqualFold(bytes.TrimSpace(key), name) {
		return nil, false
	}
	return bytes.TrimSpace(value), true
}

func parseDecimal(value []byte) (int, error) {
	if len(value) == 0 {
		return 0, fmt.Errorf("httpbench: empty decimal value")
	}
	result := 0
	maxInt := int(^uint(0) >> 1)
	for _, digit := range value {
		if digit < '0' || digit > '9' {
			return 0, fmt.Errorf("httpbench: invalid decimal value %q", value)
		}
		numeric := int(digit - '0')
		if result > (maxInt-numeric)/10 {
			return 0, fmt.Errorf("httpbench: decimal value overflows int")
		}
		result = result*10 + numeric
	}
	return result, nil
}

func closeClients(clients []client) {
	for index := range clients {
		_ = clients[index].conn.Close()
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

func writeFull(writer io.Writer, payload []byte) error {
	for len(payload) > 0 {
		written, err := writer.Write(payload)
		if written > 0 {
			payload = payload[written:]
		}
		if err != nil {
			return err
		}
		if written == 0 {
			return io.ErrShortWrite
		}
	}
	return nil
}
