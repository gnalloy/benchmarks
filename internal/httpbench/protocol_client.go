package httpbench

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"gnalloy.org/codec-http1"
	"gnalloy.org/gnalloy/buffer"
	"gnalloy.org/gnalloy/channel"
	"gnalloy.org/gnalloy/message"
)

const clientReadBufferSize = 16 * 1024

type protocolClient struct {
	conn     net.Conn
	channel  *channel.LocalChannel
	alloc    buffer.Allocator
	response *responseHandler
	request  http1.Request
}

func newProtocolClient(conn net.Conn, host string, expectedBody []byte) (*protocolClient, error) {
	alloc := buffer.NewHeapAllocator()
	ch := channel.NewLocalChannel(1, alloc, &connectionSink{conn: conn})
	decoder, err := http1.NewResponseDecoder(16*1024, len(expectedBody))
	if err != nil {
		_ = alloc.Close()
		return nil, err
	}
	response := &responseHandler{expectedBody: expectedBody}
	if err := ch.Pipeline().AddLast("httpEncoder", http1.NewRequestEncoder()); err != nil {
		_ = alloc.Close()
		return nil, err
	}
	if err := ch.Pipeline().AddLast("httpDecoder", decoder); err != nil {
		_ = alloc.Close()
		return nil, err
	}
	if err := ch.Pipeline().AddLast("response", response); err != nil {
		_ = alloc.Close()
		return nil, err
	}
	return &protocolClient{
		conn:     conn,
		channel:  ch,
		alloc:    alloc,
		response: response,
		request: http1.Request{
			Method: "GET",
			URI:    "/bench",
			Headers: http1.Headers{
				"Host":       host,
				"User-Agent": "gnalloy-bench",
				"Accept":     "*/*",
				"Connection": "keep-alive",
			},
		},
	}, nil
}

func (c *protocolClient) exchange() error {
	c.response.reset()
	if err := c.channel.WriteAndFlush(c.request); err != nil {
		return err
	}
	return c.readResponse()
}

func (c *protocolClient) close() {
	c.channel.Pipeline().FireChannelInactive()
	_ = c.alloc.Close()
}

func (c *protocolClient) readResponse() error {
	for !c.response.received && c.response.err == nil {
		in, err := c.alloc.Acquire(clientReadBufferSize)
		if err != nil {
			return err
		}
		n, readErr := c.conn.Read(in.WritableBytesView())
		if n > 0 {
			if err := in.AdvanceWriter(n); err != nil {
				in.Release()
				return err
			}
			c.channel.Pipeline().FireChannelRead(in)
		} else {
			in.Release()
		}
		if c.response.err != nil {
			return c.response.err
		}
		if c.response.received {
			return nil
		}
		if readErr != nil {
			return readErr
		}
		if n == 0 {
			return io.ErrNoProgress
		}
	}
	return c.response.err
}

type connectionSink struct {
	conn net.Conn
}

func (s *connectionSink) Write(msg any) error {
	buf, ok := msg.(buffer.ByteBuf)
	if !ok {
		message.Release(msg)
		return fmt.Errorf("httpbench: codec emitted unsupported outbound type %T", msg)
	}
	defer buf.Release()
	for _, part := range buf.ReadableSlices(nil) {
		if err := writeFull(s.conn, part); err != nil {
			return err
		}
	}
	return nil
}

func (*connectionSink) Flush() error {
	return nil
}

func (s *connectionSink) Close() error {
	return s.conn.Close()
}

type responseHandler struct {
	expectedBody []byte
	received     bool
	err          error
}

func (h *responseHandler) reset() {
	h.received = false
	h.err = nil
}

func (h *responseHandler) ChannelRead(_ *channel.HandlerContext, msg any) {
	response, ok := msg.(http1.Response)
	if !ok {
		message.Release(msg)
		h.err = fmt.Errorf("httpbench: codec emitted unsupported inbound type %T", msg)
		return
	}
	defer response.Release()
	if h.received {
		h.err = fmt.Errorf("httpbench: received multiple responses for one request")
		return
	}
	h.received = true
	if response.StatusCode != 200 {
		h.err = fmt.Errorf("httpbench: unexpected status %d", response.StatusCode)
		return
	}
	if response.Body == nil || !bytes.Equal(response.Body.Bytes(), h.expectedBody) {
		h.err = fmt.Errorf("httpbench: response body mismatch")
	}
}

func (h *responseHandler) ExceptionCaught(_ *channel.HandlerContext, err error) {
	if h.err == nil {
		h.err = err
	}
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
