package main

import (
	"testing"

	"gnalloy.org/codec-http1"
	"gnalloy.org/gnalloy/buffer"
	"gnalloy.org/gnalloy/channel"
	"gnalloy.org/gnalloy/message"
	"gnalloy.org/gnalloy/transport"
)

func TestHTTP1HandlerDoesNotResubmitOwnerLoopWrite(t *testing.T) {
	executor := &countingExecutor{}
	sink := &releasingSink{}
	ch := channel.NewLocalChannel(1, buffer.NewHeapAllocator(), sink)
	ch.BindEventExecutor(executor)
	if err := ch.Pipeline().AddLast("encoder", http1.NewResponseEncoder()); err != nil {
		t.Fatal(err)
	}
	if err := ch.Pipeline().AddLast("handler", http1Handler{body: []byte("ok")}); err != nil {
		t.Fatal(err)
	}

	ch.Pipeline().FireChannelRead(http1.Request{})

	if executor.submissions != 0 {
		t.Fatalf("owner-loop submissions=%d, want 0", executor.submissions)
	}
	if sink.writes == 0 || sink.flushes != 1 {
		t.Fatalf("writes=%d flushes=%d, want encoded response and one flush", sink.writes, sink.flushes)
	}
}

type countingExecutor struct {
	submissions int
}

func (e *countingExecutor) Submit(transport.Task) error {
	e.submissions++
	return nil
}

type releasingSink struct {
	writes  int
	flushes int
}

func (s *releasingSink) Write(msg any) error {
	s.writes++
	message.Release(msg)
	return nil
}

func (s *releasingSink) Flush() error {
	s.flushes++
	return nil
}

func (*releasingSink) Close() error { return nil }
