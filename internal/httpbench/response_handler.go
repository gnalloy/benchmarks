package httpbench

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"gnalloy.org/codec-http1"
	"gnalloy.org/gnalloy/channel"
	"gnalloy.org/gnalloy/message"
	gnalloytls "gnalloy.org/handler-tls"
)

type responseHandler struct {
	expectedBody []byte
	responses    chan error
	ready        chan struct{}
	readyOnce    sync.Once
	alpn         string
}

func newResponseHandler(expectedBody []byte) *responseHandler {
	return &responseHandler{
		expectedBody: expectedBody,
		responses:    make(chan error, 1),
		ready:        make(chan struct{}),
	}
}

func (h *responseHandler) ChannelActive(ctx *channel.HandlerContext) {
	h.readyOnce.Do(func() { close(h.ready) })
	ctx.FireChannelActive()
}

func (h *responseHandler) UserEventTriggered(ctx *channel.HandlerContext, event any) {
	if handshake, ok := event.(gnalloytls.HandshakeEvent); ok {
		h.alpn = handshake.NegotiatedProtocol
	}
	ctx.FireUserEventTriggered(event)
}

func (h *responseHandler) ChannelRead(ctx *channel.HandlerContext, msg any) {
	response, ok := msg.(http1.Response)
	if !ok {
		message.Release(msg)
		h.fail(ctx, fmt.Errorf("httpbench: codec emitted unsupported inbound type %T", msg))
		return
	}
	err := h.validate(response)
	response.Release()
	select {
	case h.responses <- err:
	default:
		h.fail(ctx, fmt.Errorf("httpbench: response queue overflow"))
	}
}

func (h *responseHandler) ChannelInactive(ctx *channel.HandlerContext) {
	h.signalError(io.ErrUnexpectedEOF)
	ctx.FireChannelInactive()
}

func (h *responseHandler) ExceptionCaught(ctx *channel.HandlerContext, err error) {
	h.fail(ctx, err)
}

func (h *responseHandler) readyError() error {
	select {
	case err := <-h.responses:
		return err
	default:
		return nil
	}
}

func (h *responseHandler) validate(response http1.Response) error {
	if response.StatusCode != 200 {
		return fmt.Errorf("httpbench: unexpected status %d", response.StatusCode)
	}
	if response.Body == nil || !bytes.Equal(response.Body.Bytes(), h.expectedBody) {
		return fmt.Errorf("httpbench: response body mismatch")
	}
	return nil
}

func (h *responseHandler) fail(ctx *channel.HandlerContext, err error) {
	h.signalError(err)
	h.readyOnce.Do(func() { close(h.ready) })
	_ = ctx.Close()
}

func (h *responseHandler) signalError(err error) {
	select {
	case h.responses <- err:
	default:
	}
}
