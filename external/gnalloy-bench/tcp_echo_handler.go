package main

import (
	"fmt"

	"gnalloy.org/gnalloy/buffer"
	"gnalloy.org/gnalloy/channel"
)

func newTCPEchoHandler(cfg config, group *tcpEchoExecutorGroup) (channel.Handler, error) {
	switch cfg.TCPEchoMode {
	case tcpEchoModeDirect:
		return directTCPEchoHandler{}, nil
	case tcpEchoModeReadComplete:
		return readCompleteTCPEchoHandler{}, nil
	case tcpEchoModeOwnerExecutor:
		return newOffloadedTCPEchoHandler(group, ownerLoopTCPEchoHandler{})
	default:
		return nil, fmt.Errorf("%w: tcp-echo-mode %s", errInvalidConfig, cfg.TCPEchoMode)
	}
}

type directTCPEchoHandler struct{}

func (directTCPEchoHandler) ChannelRead(ctx *channel.HandlerContext, msg any) {
	buf, ok := msg.(buffer.ByteBuf)
	if !ok {
		ctx.FireChannelRead(msg)
		return
	}
	if err := ctx.WriteAndFlush(buf); err != nil {
		ctx.FireExceptionCaught(err)
	}
}

func (directTCPEchoHandler) ExceptionCaught(ctx *channel.HandlerContext, _ error) {
	_ = ctx.Close()
}

type readCompleteTCPEchoHandler struct{}

func (readCompleteTCPEchoHandler) ChannelRead(ctx *channel.HandlerContext, msg any) {
	buf, ok := msg.(buffer.ByteBuf)
	if !ok {
		ctx.FireChannelRead(msg)
		return
	}
	if err := ctx.Write(buf); err != nil {
		ctx.FireExceptionCaught(err)
	}
}

func (readCompleteTCPEchoHandler) ChannelReadComplete(ctx *channel.HandlerContext) {
	if err := ctx.Flush(); err != nil {
		ctx.FireExceptionCaught(err)
		return
	}
	ctx.FireChannelReadComplete()
}

func (readCompleteTCPEchoHandler) ExceptionCaught(ctx *channel.HandlerContext, _ error) {
	_ = ctx.Close()
}

type ownerLoopTCPEchoHandler struct{}

func (ownerLoopTCPEchoHandler) ChannelRead(ctx *channel.HandlerContext, msg any) {
	buf, ok := msg.(buffer.ByteBuf)
	if !ok {
		ctx.FireChannelRead(msg)
		return
	}
	if err := ctx.Channel().WriteAndFlush(buf); err != nil {
		ctx.FireExceptionCaught(err)
	}
}

func (ownerLoopTCPEchoHandler) ExceptionCaught(ctx *channel.HandlerContext, _ error) {
	_ = ctx.Close()
}

type offloadedTCPEchoHandler struct {
	worker   *tcpEchoExecutorWorker
	delegate channel.ChannelReadHandler
	inbound  tcpEchoInboundQueue
}

func newOffloadedTCPEchoHandler(group *tcpEchoExecutorGroup, delegate channel.Handler) (*offloadedTCPEchoHandler, error) {
	if group == nil {
		return nil, errInvalidConfig
	}
	read, ok := delegate.(channel.ChannelReadHandler)
	if !ok {
		return nil, errInvalidConfig
	}
	worker := group.bind()
	if worker == nil {
		return nil, errInvalidConfig
	}
	h := &offloadedTCPEchoHandler{worker: worker, delegate: read}
	h.inbound.init(h)
	return h, nil
}

func (h *offloadedTCPEchoHandler) ChannelRead(ctx *channel.HandlerContext, msg any) {
	if err := h.inbound.submit(ctx, msg); err != nil {
		ctx.FireExceptionCaught(err)
	}
}

func (h *offloadedTCPEchoHandler) ExceptionCaught(ctx *channel.HandlerContext, _ error) {
	_ = ctx.Close()
}

func recoverTCPEchoHandler(ctx *channel.HandlerContext) {
	if v := recover(); v != nil {
		ctx.FireExceptionCaught(fmt.Errorf("%w: %v", errTCPEchoHandlerPanic, v))
	}
}
