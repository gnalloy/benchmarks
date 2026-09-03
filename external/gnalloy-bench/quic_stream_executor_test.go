package main

import (
	"context"
	"sync/atomic"
	"testing"

	"gnalloy.org/transport-quic"
)

func TestQUICStreamExecutorProcessesSubmittedStreams(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{}, 4)
	var handled atomic.Int64
	executor := newQUICStreamExecutor(ctx, 2, 4, func(quic.Stream) {
		handled.Add(1)
		done <- struct{}{}
	})
	for range 4 {
		if !executor.submit(nil) {
			t.Fatal("submit rejected before cancellation")
		}
	}
	for range 4 {
		<-done
	}
	cancel()
	executor.stop()
	if got := handled.Load(); got != 4 {
		t.Fatalf("handled=%d, want 4", got)
	}
	if executor.submit(nil) {
		t.Fatal("submit accepted after cancellation")
	}
}

func TestQUICStreamQueueSize(t *testing.T) {
	tests := []struct {
		name        string
		workerCount int
		want        int
	}{
		{name: "single worker", workerCount: 1, want: 4},
		{name: "regular load", workerCount: 64, want: 256},
		{name: "queue limit", workerCount: 1024, want: maxQUICStreamQueueSize},
		{name: "overflow resistant", workerCount: maxInt, want: maxQUICStreamQueueSize},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := quicStreamQueueSize(test.workerCount); got != test.want {
				t.Fatalf("quicStreamQueueSize(%d)=%d, want %d", test.workerCount, got, test.want)
			}
		})
	}
}
