package main

import (
	"context"
	"time"

	"github.com/panjf2000/gnet/v2"
)

type echoServer struct {
	addr  string
	eng   gnet.Engine
	errCh chan error
}

func startEchoServer(ctx context.Context, cfg config) (*echoServer, error) {
	network := listenNetwork(cfg.Protocol)
	addr, err := resolveListenAddress(network, cfg.Addr)
	if err != nil {
		return nil, err
	}
	ready := make(chan gnet.Engine, 1)
	handler := &gnetEchoHandler{ready: ready}
	errCh := make(chan error, 1)
	opts := []gnet.Option{
		gnet.WithReuseAddr(true),
		gnet.WithLogger(noopLogger{}),
	}
	if network == "tcp" {
		opts = append(opts, gnet.WithTCPNoDelay(gnet.TCPNoDelay))
	}
	if cfg.EventLoops > 0 {
		opts = append(opts, gnet.WithNumEventLoop(cfg.EventLoops))
	} else if cfg.Multicore {
		opts = append(opts, gnet.WithMulticore(true))
	}

	go func() {
		errCh <- gnet.Run(handler, network+"://"+addr, opts...)
	}()

	select {
	case eng := <-ready:
		return &echoServer{addr: addr, eng: eng, errCh: errCh}, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func listenNetwork(protocol string) string {
	if protocol == "udp-echo" {
		return "udp"
	}
	return "tcp"
}

func (s *echoServer) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = s.eng.Stop(ctx)
	select {
	case <-s.errCh:
	case <-ctx.Done():
	}
}

type gnetEchoHandler struct {
	gnet.BuiltinEventEngine
	ready chan<- gnet.Engine
}

type noopLogger struct{}

func (noopLogger) Debugf(string, ...any) {}
func (noopLogger) Infof(string, ...any)  {}
func (noopLogger) Warnf(string, ...any)  {}
func (noopLogger) Errorf(string, ...any) {}
func (noopLogger) Fatalf(string, ...any) {}

func (h *gnetEchoHandler) OnBoot(eng gnet.Engine) gnet.Action {
	select {
	case h.ready <- eng:
	default:
	}
	return gnet.None
}

func (h *gnetEchoHandler) OnTraffic(c gnet.Conn) gnet.Action {
	buf, err := c.Next(-1)
	if err != nil {
		return gnet.Close
	}
	if len(buf) == 0 {
		return gnet.None
	}
	n, err := c.Write(buf)
	if err != nil || n != len(buf) {
		return gnet.Close
	}
	if err := c.Flush(); err != nil {
		return gnet.Close
	}
	return gnet.None
}
