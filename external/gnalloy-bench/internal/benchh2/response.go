package benchh2

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	http2 "gnalloy.org/codec-http2"
	"gnalloy.org/gnalloy/buffer"
	"gnalloy.org/gnalloy/channel"
	gnalloytls "gnalloy.org/handler-tls"
)

const connectionWindowIncrease = 1<<30 - 1

type streamResponse struct {
	streamID http2.StreamID
	err      error
}

type responseHandler struct {
	expected  []byte
	reply     []byte
	ready     chan struct{}
	responses chan streamResponse
	readyOnce sync.Once

	streamID http2.StreamID
	received int
	statusOK bool
	alpn     string
}

func newResponseHandler(expected []byte) *responseHandler {
	return &responseHandler{
		expected:  expected,
		reply:     make([]byte, len(expected)),
		ready:     make(chan struct{}),
		responses: make(chan streamResponse, 1),
	}
}

func (h *responseHandler) ChannelActive(ctx *channel.HandlerContext) {
	settings := http2.SettingsFrame{Settings: []http2.Setting{{ID: http2.SettingEnablePush, Value: 0}}}
	if err := ctx.Write(settings); err != nil {
		h.fail(ctx, err)
		return
	}
	if err := ctx.Write(http2.WindowUpdateFrame{Increment: connectionWindowIncrease}); err != nil {
		h.fail(ctx, err)
		return
	}
	if err := ctx.Flush(); err != nil {
		h.fail(ctx, err)
		return
	}
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
	switch event := msg.(type) {
	case http2.StreamEvent:
		h.readStreamEvent(ctx, event)
	case http2.TypedFrame:
		event.Release()
	default:
		ctx.FireChannelRead(msg)
	}
}

func (h *responseHandler) readStreamEvent(ctx *channel.HandlerContext, event http2.StreamEvent) {
	defer event.Release()
	if event.Type != http2.StreamEventRead {
		return
	}
	switch frame := event.Frame.(type) {
	case http2.HeadersBlock:
		h.readHeaders(ctx, event.StreamID, frame)
	case http2.DataFrame:
		h.readData(ctx, event.StreamID, frame)
	}
}

func (h *responseHandler) readHeaders(ctx *channel.HandlerContext, streamID http2.StreamID, frame http2.HeadersBlock) {
	if h.streamID != 0 && h.streamID != streamID {
		h.fail(ctx, fmt.Errorf("benchh2: interleaved response stream %d while waiting for %d", streamID, h.streamID))
		return
	}
	h.streamID = streamID
	h.statusOK = responseStatusOK(frame.Fields)
	if frame.EndStream {
		h.complete(ctx, streamID)
	}
}

func (h *responseHandler) readData(ctx *channel.HandlerContext, streamID http2.StreamID, frame http2.DataFrame) {
	if h.streamID == 0 || h.streamID != streamID {
		h.fail(ctx, fmt.Errorf("benchh2: DATA for unexpected stream %d", streamID))
		return
	}
	readable := 0
	if frame.Data != nil {
		readable = frame.Data.ReadableBytes()
	}
	if readable > len(h.reply)-h.received {
		h.fail(ctx, fmt.Errorf("benchh2: response body too large"))
		return
	}
	if readable > 0 {
		h.received += buffer.CopyReadableBytes(h.reply[h.received:], frame.Data)
	}
	if frame.Flags&http2.FlagEndStream != 0 {
		h.complete(ctx, streamID)
	}
}

func (h *responseHandler) complete(ctx *channel.HandlerContext, streamID http2.StreamID) {
	var err error
	switch {
	case !h.statusOK:
		err = fmt.Errorf("benchh2: response status is not 200")
	case h.received != len(h.expected):
		err = fmt.Errorf("benchh2: response body length %d, want %d", h.received, len(h.expected))
	case !bytes.Equal(h.reply[:h.received], h.expected):
		err = fmt.Errorf("benchh2: response body mismatch")
	}
	h.streamID = 0
	h.received = 0
	h.statusOK = false
	select {
	case h.responses <- streamResponse{streamID: streamID, err: err}:
	default:
		h.fail(ctx, fmt.Errorf("benchh2: response queue overflow"))
	}
}

func (h *responseHandler) ChannelInactive(ctx *channel.HandlerContext) {
	h.signalError(io.ErrUnexpectedEOF)
	ctx.FireChannelInactive()
}

func (h *responseHandler) ExceptionCaught(ctx *channel.HandlerContext, err error) {
	h.fail(ctx, err)
}

func (h *responseHandler) fail(ctx *channel.HandlerContext, err error) {
	h.signalError(err)
	h.readyOnce.Do(func() { close(h.ready) })
	_ = ctx.Close()
}

func (h *responseHandler) signalError(err error) {
	select {
	case h.responses <- streamResponse{err: err}:
	default:
	}
}

func responseStatusOK(fields []http2.HeaderField) bool {
	for _, field := range fields {
		if field.Name == ":status" {
			return field.Value == "200"
		}
	}
	return false
}
