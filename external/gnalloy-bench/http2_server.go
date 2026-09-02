package main

import (
	"context"
	"strconv"

	"gnalloy.org/benchmarks/external/gnalloy-bench/internal/benchh2"
	"gnalloy.org/codec-http2"
	"gnalloy.org/gnalloy/bootstrap"
	"gnalloy.org/gnalloy/buffer"
	"gnalloy.org/gnalloy/channel"
	"gnalloy.org/gnalloy/transport"
	gnalloytls "gnalloy.org/handler-tls"
	"gnalloy.org/transport-tcp"
)

const http2PrefaceHandlerName = "http2-preface"

func startHTTP2Server(ctx context.Context, cfg config) (*echoServer, error) {
	boss, workers, err := newGroups(cfg)
	if err != nil {
		return nil, err
	}
	server, err := bindHTTP2Server(ctx, cfg, boss, workers)
	if err != nil {
		shutdownGroups(boss, workers)
		return nil, err
	}
	return &echoServer{addr: server.Addr(), server: server, boss: boss, workers: workers}, nil
}

func bindHTTP2Server(ctx context.Context, cfg config, boss *transport.EventLoopGroup, workers *transport.EventLoopGroup) (bootstrap.Server, error) {
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
		ChildInitializer(func(ch channel.Channel) error {
			if cfg.Protocol == "https2" {
				tlsConfig, err := serverTLSConfig(cfg)
				if err != nil {
					return err
				}
				if err := ch.Pipeline().AddLast("tls", gnalloytls.Server(gnalloytls.Config{TLS: tlsConfig})); err != nil {
					return err
				}
			}
			return addHTTP2Pipeline(ch, cfg)
		}).
		BindContext(ctx, cfg.Addr)
}

func addHTTP2Pipeline(ch channel.Channel, cfg config) error {
	frameDecoder, err := http2.NewFrameDecoder(http2.DefaultMaxFrameSize)
	if err != nil {
		return err
	}
	headerDecoder, err := http2.NewHeaderDecoder(http2.HeaderCodecConfig{})
	if err != nil {
		return err
	}
	headerEncoder, err := http2.NewHeaderEncoder(http2.HeaderCodecConfig{})
	if err != nil {
		return err
	}
	mux, err := http2.NewStreamMultiplexer(http2.MultiplexerConfig{Server: true})
	if err != nil {
		return err
	}
	pipeline := ch.Pipeline()
	if err := pipeline.AddLast(http2PrefaceHandlerName, http2.NewPrefaceDecoder()); err != nil {
		return err
	}
	if err := pipeline.AddLast("http2-frame-decoder", frameDecoder); err != nil {
		return err
	}
	if err := pipeline.AddLast("http2-typed-decoder", http2.NewTypedFrameDecoder()); err != nil {
		return err
	}
	frameEncoder := http2.NewFrameEncoder()
	if cfg.Protocol == "https2" {
		frameEncoder = http2.NewFrameEncoderWithConfig(http2.FrameEncoderConfig{CoalescePayload: true})
	}
	if err := pipeline.AddLast("http2-frame-encoder", frameEncoder); err != nil {
		return err
	}
	if err := pipeline.AddLast("http2-typed-encoder", http2.NewTypedFrameEncoder()); err != nil {
		return err
	}
	if err := pipeline.AddLast("http2-header-decoder", headerDecoder); err != nil {
		return err
	}
	if err := pipeline.AddLast("http2-header-encoder", headerEncoder); err != nil {
		return err
	}
	if err := pipeline.AddLast("http2-settings-ack", http2.NewSettingsAckHandler()); err != nil {
		return err
	}
	if err := pipeline.AddLast("http2-ping-ack", http2.NewPingAckHandler()); err != nil {
		return err
	}
	if err := pipeline.AddLast("http2-mux", mux); err != nil {
		return err
	}
	return pipeline.AddLast("http2-handler", newHTTP2BenchmarkHandler(cfg.Payload))
}

type http2BenchmarkHandler struct {
	body   []byte
	fields []http2.HeaderField
}

func newHTTP2BenchmarkHandler(payload int) http2BenchmarkHandler {
	body := benchh2.ResponseBody(payload)
	return http2BenchmarkHandler{
		body: body,
		fields: []http2.HeaderField{
			{Name: ":status", Value: "200"},
			{Name: "content-type", Value: "application/octet-stream"},
			{Name: "content-length", Value: strconv.Itoa(len(body))},
		},
	}
}

func (h http2BenchmarkHandler) ChannelActive(ctx *channel.HandlerContext) {
	if err := ctx.WriteAndFlush(http2.SettingsFrame{}); err != nil {
		ctx.FireExceptionCaught(err)
		return
	}
	ctx.FireChannelActive()
}

func (h http2BenchmarkHandler) ChannelRead(ctx *channel.HandlerContext, msg any) {
	switch frame := msg.(type) {
	case http2.StreamEvent:
		defer frame.Release()
		if frame.Type != http2.StreamEventRead {
			return
		}
		if _, ok := frame.Frame.(http2.HeadersBlock); !ok {
			return
		}
		if err := h.writeResponse(ctx, frame.StreamID); err != nil {
			ctx.FireExceptionCaught(err)
		}
	default:
		ctx.FireChannelRead(msg)
	}
}

func (h http2BenchmarkHandler) writeResponse(ctx *channel.HandlerContext, streamID http2.StreamID) error {
	if err := ctx.Write(http2.HeadersBlock{StreamID: streamID, Fields: h.fields}); err != nil {
		return err
	}
	body, err := ctx.Channel().Allocator().Acquire(len(h.body))
	if err != nil {
		return err
	}
	if _, err := body.WriteBytes(h.body); err != nil {
		body.Release()
		return err
	}
	if err := ctx.Write(http2.DataFrame{StreamID: streamID, Flags: http2.FlagEndStream, Data: body}); err != nil {
		return err
	}
	return ctx.Flush()
}

func (http2BenchmarkHandler) ExceptionCaught(ctx *channel.HandlerContext, _ error) {
	_ = ctx.Close()
}
