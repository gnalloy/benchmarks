package main

import (
	"context"

	"gnalloy.org/benchmarks/internal/httpbench"
	"gnalloy.org/codec-http1"
	"gnalloy.org/gnalloy/bootstrap"
	"gnalloy.org/gnalloy/buffer"
	"gnalloy.org/gnalloy/channel"
	"gnalloy.org/gnalloy/transport"
	gnalloytls "gnalloy.org/handler-tls"
	"gnalloy.org/transport-tcp"
)

const maxCoalescedHTTP1BodyBytes = 16 * 1024

func startHTTP1Server(ctx context.Context, cfg config) (*echoServer, error) {
	boss, workers, err := newGroups(cfg)
	if err != nil {
		return nil, err
	}
	server, err := bindHTTP1Server(ctx, cfg, boss, workers)
	if err != nil {
		shutdownGroups(boss, workers)
		return nil, err
	}
	return &echoServer{addr: server.Addr(), server: server, boss: boss, workers: workers}, nil
}

func bindHTTP1Server(ctx context.Context, cfg config, boss *transport.EventLoopGroup, workers *transport.EventLoopGroup) (bootstrap.Server, error) {
	tcpConfig := tcp.DefaultConfig()
	tcpConfig.ReadBufferSize = cfg.ReadBufferSize
	tcpConfig.ReusePort = cfg.ReusePort
	tcpConfig.IOUringFixedBuffers = cfg.IOUringFixedBuffers
	if cfg.Mmap {
		tcpConfig.AllocatorFactory = tcp.NewMmapAllocatorFactory(buffer.MmapAllocatorConfig{
			BlockSize: cfg.MmapBlockSize,
			Blocks:    cfg.MmapBlocks,
		}, false)
	}
	return bootstrap.NewServerBootstrap().
		Group(boss, workers).
		Transport(tcp.NewTransport(tcpConfig)).
		ChildOption(benchmarkChildOptions(cfg)...).
		ChildInitializer(func(ch channel.Channel) error {
			if cfg.Protocol == "https1" {
				tlsConfig, err := serverTLSConfig(cfg)
				if err != nil {
					return err
				}
				if err := ch.Pipeline().AddLast("tls", gnalloytls.Server(gnalloytls.Config{TLS: tlsConfig})); err != nil {
					return err
				}
			}
			return addHTTP1CodecPipeline(ch, cfg)
		}).
		BindContext(ctx, cfg.Addr)
}

func addHTTP1CodecPipeline(ch channel.Channel, cfg config) error {
	decoder, err := http1.NewRequestDecoder(16*1024, 0)
	if err != nil {
		return err
	}
	encoder := newHTTP1ResponseEncoder(cfg.Payload)
	if err := ch.Pipeline().AddLast("httpEncoder", encoder); err != nil {
		return err
	}
	if err := ch.Pipeline().AddLast("httpDecoder", decoder); err != nil {
		return err
	}
	if err := ch.Pipeline().AddLast("continue", http1.NewContinueHandler()); err != nil {
		return err
	}
	return ch.Pipeline().AddLast("handler", http1Handler{
		body: httpbench.ResponseBody(cfg.Payload),
		headers: http1.Headers{
			"Content-Type": "application/octet-stream",
		},
	})
}

func newHTTP1ResponseEncoder(payload int) *http1.ResponseEncoder {
	if payload > 0 && payload <= maxCoalescedHTTP1BodyBytes {
		return http1.NewResponseEncoderWithOptions(http1.ResponseEncoderOptions{
			CoalesceBodyBytes: maxCoalescedHTTP1BodyBytes,
		})
	}
	return http1.NewResponseEncoder()
}

type http1Handler struct {
	body    []byte
	headers http1.Headers
}

func (h http1Handler) ChannelRead(ctx *channel.HandlerContext, msg any) {
	req, ok := msg.(*http1.Request)
	if !ok {
		ctx.FireChannelRead(msg)
		return
	}
	defer req.Release()
	resp := http1.AcquireResponse()
	resp.StatusCode = 200
	resp.Headers = h.headers
	resp.Body = buffer.NewSharedBuffer(h.body)
	if err := ctx.WriteAndFlush(resp); err != nil {
		ctx.FireExceptionCaught(err)
	}
}

func (http1Handler) ExceptionCaught(ctx *channel.HandlerContext, _ error) {
	_ = ctx.Close()
}
