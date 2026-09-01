package main

import (
	"context"
	"sync"
	"sync/atomic"
)

type tcpEchoTask func()

type tcpEchoExecutorGroup struct {
	workers []*tcpEchoExecutorWorker
	next    atomic.Uint64
	closed  atomic.Bool
}

func newTCPEchoExecutorGroup(size int, queueSize int) (*tcpEchoExecutorGroup, error) {
	if size <= 0 || queueSize <= 0 {
		return nil, errInvalidConfig
	}
	group := &tcpEchoExecutorGroup{
		workers: make([]*tcpEchoExecutorWorker, 0, size),
	}
	for i := 0; i < size; i++ {
		group.workers = append(group.workers, newTCPEchoExecutorWorker(queueSize))
	}
	return group, nil
}

func (g *tcpEchoExecutorGroup) bind() *tcpEchoExecutorWorker {
	if g == nil || len(g.workers) == 0 {
		return nil
	}
	idx := g.next.Add(1) - 1
	return g.workers[idx%uint64(len(g.workers))]
}

func (g *tcpEchoExecutorGroup) shutdown(ctx context.Context) error {
	if g == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if !g.closed.Swap(true) {
		for _, worker := range g.workers {
			worker.close()
		}
	}
	done := make(chan struct{})
	go func() {
		for _, worker := range g.workers {
			worker.await()
		}
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type tcpEchoExecutorWorker struct {
	mu     sync.RWMutex
	tasks  chan tcpEchoTask
	closed bool
	done   chan struct{}
}

func newTCPEchoExecutorWorker(queueSize int) *tcpEchoExecutorWorker {
	worker := &tcpEchoExecutorWorker{
		tasks: make(chan tcpEchoTask, queueSize),
		done:  make(chan struct{}),
	}
	go worker.run()
	return worker
}

func (w *tcpEchoExecutorWorker) submit(task tcpEchoTask) error {
	if task == nil {
		return nil
	}
	w.mu.RLock()
	defer w.mu.RUnlock()
	if w.closed {
		return errTCPEchoExecutorDone
	}
	select {
	case w.tasks <- task:
		return nil
	default:
		return errTCPEchoExecutorFull
	}
}

func (w *tcpEchoExecutorWorker) close() {
	w.mu.Lock()
	if !w.closed {
		w.closed = true
		close(w.tasks)
	}
	w.mu.Unlock()
}

func (w *tcpEchoExecutorWorker) await() {
	<-w.done
}

func (w *tcpEchoExecutorWorker) run() {
	defer close(w.done)
	for task := range w.tasks {
		task()
	}
}
