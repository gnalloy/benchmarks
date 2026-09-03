package main

import (
	"testing"

	"gnalloy.org/codec-http1"
	"gnalloy.org/gnalloy/buffer"
	"gnalloy.org/gnalloy/channel"
	"gnalloy.org/gnalloy/message"
	"gnalloy.org/gnalloy/transport"
)

func TestHTTP1ServerAppliesBenchmarkChildOptions(t *testing.T) {
	cfg := config{
		MaxMessagesPerRead: 7,
		FlushStrategy:      channel.FlushOnEventLoopBatch,
	}
	options := channel.NewChannelOptions()
	options.Apply(benchmarkChildOptions(cfg)...)

	if got := channel.OptionMaxMessagesPerRead.Get(options); got != cfg.MaxMessagesPerRead {
		t.Fatalf("max messages per read=%d, want %d", got, cfg.MaxMessagesPerRead)
	}
	if got := channel.OptionFlushStrategy.Get(options); got != cfg.FlushStrategy {
		t.Fatalf("flush strategy=%d, want %d", got, cfg.FlushStrategy)
	}
}

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

	ch.Pipeline().FireChannelRead(&http1.Request{})

	if executor.submissions != 0 {
		t.Fatalf("owner-loop submissions=%d, want 0", executor.submissions)
	}
	if sink.writes == 0 || sink.flushes != 1 {
		t.Fatalf("writes=%d flushes=%d, want encoded response and one flush", sink.writes, sink.flushes)
	}
}

func TestHTTPS1ResponseEncoderCoalescesSmallBody(t *testing.T) {
	sink := &releasingSink{}
	ch := channel.NewLocalChannel(1, buffer.NewHeapAllocator(), sink)
	if err := ch.Pipeline().AddLast("encoder", newHTTP1ResponseEncoder("https1", 128)); err != nil {
		t.Fatal(err)
	}

	if err := ch.Write(http1.Response{StatusCode: 200, Body: buffer.NewSharedBuffer(make([]byte, 128))}); err != nil {
		t.Fatal(err)
	}
	if sink.writes != 1 {
		t.Fatalf("writes=%d, want one coalesced buffer", sink.writes)
	}
}

func TestHTTP1ResponseEncoderSplitsSmallBody(t *testing.T) {
	sink := &releasingSink{}
	ch := channel.NewLocalChannel(1, buffer.NewHeapAllocator(), sink)
	if err := ch.Pipeline().AddLast("encoder", newHTTP1ResponseEncoder("http1", 128)); err != nil {
		t.Fatal(err)
	}

	if err := ch.Write(http1.Response{StatusCode: 200, Body: buffer.NewSharedBuffer(make([]byte, 128))}); err != nil {
		t.Fatal(err)
	}
	if sink.writes != 2 {
		t.Fatalf("writes=%d, want split header and body", sink.writes)
	}
}

func TestHTTPS1ResponseEncoderSplitsLargeBody(t *testing.T) {
	sink := &releasingSink{}
	ch := channel.NewLocalChannel(1, buffer.NewHeapAllocator(), sink)
	if err := ch.Pipeline().AddLast("encoder", newHTTP1ResponseEncoder("https1", maxCoalescedHTTP1BodyBytes+1)); err != nil {
		t.Fatal(err)
	}

	body := make([]byte, maxCoalescedHTTP1BodyBytes+1)
	if err := ch.Write(http1.Response{StatusCode: 200, Body: buffer.NewSharedBuffer(body)}); err != nil {
		t.Fatal(err)
	}
	if sink.writes != 2 {
		t.Fatalf("writes=%d, want split header and body", sink.writes)
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
