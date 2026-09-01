package main

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestTCPEchoExecutorGroupRunsBoundWorkersConcurrently(t *testing.T) {
	group, err := newTCPEchoExecutorGroup(4, 8)
	if err != nil {
		t.Fatal(err)
	}
	release := make(chan struct{})
	started := make(chan struct{}, 4)
	for range 4 {
		worker := group.bind()
		if worker == nil {
			t.Fatal("expected bound worker")
		}
		if err := worker.submit(func() {
			started <- struct{}{}
			<-release
		}); err != nil {
			t.Fatal(err)
		}
	}
	for range 4 {
		select {
		case <-started:
		case <-time.After(time.Second):
			t.Fatal("fixed workers did not run concurrently")
		}
	}
	close(release)
	if err := group.shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestTCPEchoExecutorGroupShutdownDrainsAndRejects(t *testing.T) {
	group, err := newTCPEchoExecutorGroup(2, 8)
	if err != nil {
		t.Fatal(err)
	}
	var completed atomic.Int32
	workers := make([]*tcpEchoExecutorWorker, 0, 8)
	for range 8 {
		worker := group.bind()
		workers = append(workers, worker)
		if err := worker.submit(func() {
			completed.Add(1)
		}); err != nil {
			t.Fatal(err)
		}
	}
	if err := group.shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := completed.Load(); got != 8 {
		t.Fatalf("completed=%d, want 8", got)
	}
	for _, worker := range workers {
		if err := worker.submit(func() {}); !errors.Is(err, errTCPEchoExecutorDone) {
			t.Fatalf("submit err=%v, want %v", err, errTCPEchoExecutorDone)
		}
	}
}
