package main

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"gnalloy.org/gnalloy/channel"
)

type recordingTCPEchoReadHandler struct {
	mu       sync.Mutex
	messages []any
}

func (h *recordingTCPEchoReadHandler) ChannelRead(_ *channel.HandlerContext, msg any) {
	h.mu.Lock()
	h.messages = append(h.messages, msg)
	h.mu.Unlock()
}

func TestOffloadedTCPEchoHandlerPreservesConnectionOrder(t *testing.T) {
	group, err := newTCPEchoExecutorGroup(1, 128)
	if err != nil {
		t.Fatal(err)
	}
	handler := &recordingTCPEchoReadHandler{}
	offloaded, err := newOffloadedTCPEchoHandler(group, handler)
	if err != nil {
		t.Fatal(err)
	}
	want := make([]any, 100)
	for i := range want {
		want[i] = i
		if err := offloaded.inbound.submit(nil, i); err != nil {
			t.Fatal(err)
		}
	}
	if err := group.shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(handler.messages, want) {
		t.Fatalf("messages=%v, want %v", handler.messages, want)
	}
}
