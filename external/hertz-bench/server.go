package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	hertzconfig "github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	http2config "github.com/hertz-contrib/http2/config"
	"github.com/hertz-contrib/http2/factory"
)

type benchmarkServer struct {
	addr     string
	listener net.Listener
	engine   *server.Hertz
	errCh    chan error
}

func startServer(cfg config) (*benchmarkServer, error) {
	listener, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	body := responseBody(cfg.Payload)
	options := []hertzconfig.Option{
		server.WithListener(listener),
		server.WithH2C(true),
	}
	engine := server.New(options...)
	engine.AddProtocol("h2", factory.NewServerFactory(
		http2config.WithReadTimeout(cfg.Timeout),
		http2config.WithDisableKeepAlive(false),
	))
	engine.GET("/bench", func(_ context.Context, ctx *app.RequestContext) {
		ctx.SetStatusCode(consts.StatusOK)
		ctx.SetContentType("application/octet-stream")
		ctx.Response.SetBodyRaw(body)
	})

	result := &benchmarkServer{
		addr:     listener.Addr().String(),
		listener: listener,
		engine:   engine,
		errCh:    make(chan error, 1),
	}
	go func() {
		result.errCh <- engine.Run()
	}()
	select {
	case err := <-result.errCh:
		_ = listener.Close()
		return nil, fmt.Errorf("hertz-bench: start server: %w", err)
	case <-time.After(20 * time.Millisecond):
		return result, nil
	}
}

func (s *benchmarkServer) stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	shutdownErr := s.engine.Shutdown(ctx)
	_ = s.listener.Close()
	select {
	case runErr := <-s.errCh:
		if runErr != nil && !errors.Is(runErr, net.ErrClosed) && shutdownErr == nil {
			shutdownErr = runErr
		}
	case <-ctx.Done():
		if shutdownErr == nil {
			shutdownErr = ctx.Err()
		}
	}
	return shutdownErr
}
