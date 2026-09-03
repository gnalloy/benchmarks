package main

import (
	"context"
	"sync"

	"gnalloy.org/transport-quic"
)

const maxQUICStreamQueueSize = 4096

// quicStreamExecutor 对 stream 处理施加固定并发上限，避免高频创建短生命周期 goroutine。
type quicStreamExecutor struct {
	ctx     context.Context
	queue   chan quic.Stream
	handle  func(quic.Stream)
	workers sync.WaitGroup
}

func newQUICStreamExecutor(ctx context.Context, workerCount, queueSize int, handle func(quic.Stream)) *quicStreamExecutor {
	if workerCount < 1 {
		workerCount = 1
	}
	if queueSize < workerCount {
		queueSize = workerCount
	}
	executor := &quicStreamExecutor{
		ctx:    ctx,
		queue:  make(chan quic.Stream, queueSize),
		handle: handle,
	}
	executor.workers.Add(workerCount)
	for range workerCount {
		go executor.run()
	}
	return executor
}

func quicStreamQueueSize(workerCount int) int {
	if workerCount >= maxQUICStreamQueueSize/4 {
		return maxQUICStreamQueueSize
	}
	return workerCount * 4
}

func (e *quicStreamExecutor) submit(stream quic.Stream) bool {
	if e.ctx.Err() != nil {
		return false
	}
	select {
	case e.queue <- stream:
		return true
	case <-e.ctx.Done():
		return false
	}
}

func (e *quicStreamExecutor) run() {
	defer e.workers.Done()
	for {
		select {
		case stream := <-e.queue:
			e.handle(stream)
		case <-e.ctx.Done():
			return
		}
	}
}

func (e *quicStreamExecutor) stop() {
	e.workers.Wait()
	for {
		select {
		case stream := <-e.queue:
			if stream != nil {
				stream.CancelRead(0)
				stream.CancelWrite(0)
			}
		default:
			return
		}
	}
}
