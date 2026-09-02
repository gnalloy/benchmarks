package main

import (
	"context"

	"gnalloy.org/gnalloy/bootstrap"
	"gnalloy.org/gnalloy/buffer"
	"gnalloy.org/gnalloy/channel"
	"gnalloy.org/gnalloy/transport"
	"gnalloy.org/transport-udp"
)

func startUDPEchoServer(ctx context.Context, cfg config) (*echoServer, error) {
	boss, workers, err := newGroups(cfg)
	if err != nil {
		return nil, err
	}
	server, err := bindUDPEchoServer(ctx, cfg, boss, workers)
	if err != nil {
		shutdownGroups(boss, workers)
		return nil, err
	}
	return &echoServer{addr: server.Addr(), server: server, boss: boss, workers: workers}, nil
}

func bindUDPEchoServer(ctx context.Context, cfg config, boss *transport.EventLoopGroup, workers *transport.EventLoopGroup) (bootstrap.Server, error) {
	udpConfig := udp.DefaultConfig()
	udpConfig.ReadBufferSize = cfg.ReadBufferSize
	udpConfig.ReusePort = cfg.ReusePort
	udpConfig.PooledInboundDatagrams = true
	if cfg.Mmap {
		udpConfig.AllocatorFactory = udp.NewMmapAllocatorFactory(buffer.MmapAllocatorConfig{
			BlockSize: cfg.MmapBlockSize,
			Blocks:    cfg.MmapBlocks,
		}, false)
	}
	return bootstrap.NewServerBootstrap().
		Group(boss, workers).
		Transport(udp.NewTransport(udpConfig)).
		ChildInitializer(func(ch channel.Channel) error {
			return ch.Pipeline().AddLast("echo", udpEchoHandler{})
		}).
		BindContext(ctx, cfg.Addr)
}

type udpEchoHandler struct{}

func (udpEchoHandler) ChannelRead(ctx *channel.HandlerContext, msg any) {
	switch datagram := msg.(type) {
	case udp.Datagram:
	case *udp.Datagram:
		if datagram == nil {
			ctx.FireChannelRead(msg)
			return
		}
	default:
		ctx.FireChannelRead(msg)
		return
	}
	if err := ctx.Write(msg); err != nil {
		ctx.FireExceptionCaught(err)
	}
}

func (udpEchoHandler) ChannelReadComplete(ctx *channel.HandlerContext) {
	if err := ctx.Flush(); err != nil {
		ctx.FireExceptionCaught(err)
		return
	}
	ctx.FireChannelReadComplete()
}

func (udpEchoHandler) ExceptionCaught(ctx *channel.HandlerContext, _ error) {
	_ = ctx.Close()
}
