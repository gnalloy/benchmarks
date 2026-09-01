package main

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/valyala/fasthttp"
)

type echoServer struct {
	addr     string
	listener net.Listener
	server   *fasthttp.Server
	errCh    chan error
}

func startHTTPServer(ctx context.Context, cfg config) (*echoServer, error) {
	listener, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	if tcpAddr := listener.Addr(); tcpAddr != nil {
		cfg.Addr = tcpAddr.String()
	}
	if cfg.Protocol == protocolHTTPS1 {
		tlsConfig, err := serverTLSConfig(cfg)
		if err != nil {
			_ = listener.Close()
			return nil, err
		}
		listener = tls.NewListener(listener, tlsConfig)
	}
	body := responseBody(cfg.Payload)
	server := &fasthttp.Server{
		Handler:               responseHandler(body),
		NoDefaultServerHeader: true,
		NoDefaultDate:         true,
		DisableKeepalive:      false,
	}
	out := &echoServer{
		addr:     cfg.Addr,
		listener: listener,
		server:   server,
		errCh:    make(chan error, 1),
	}
	go func() {
		out.errCh <- server.Serve(listener)
	}()
	select {
	case err := <-out.errCh:
		_ = listener.Close()
		return nil, err
	case <-time.After(20 * time.Millisecond):
		return out, nil
	case <-ctx.Done():
		_ = listener.Close()
		return nil, ctx.Err()
	}
}

func (s *echoServer) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = s.server.ShutdownWithContext(ctx)
	_ = s.listener.Close()
	select {
	case <-s.errCh:
	case <-ctx.Done():
	}
}

func responseHandler(body []byte) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetContentType("application/octet-stream")
		ctx.Response.SetBodyRaw(body)
	}
}

func responseBody(size int) []byte {
	body := make([]byte, size)
	for i := range body {
		body[i] = byte(i)
	}
	return body
}
