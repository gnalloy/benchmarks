package main

import (
	"sync"

	"gnalloy.org/gnalloy/channel"
	"gnalloy.org/gnalloy/message"
)

type tcpEchoInboundEntry struct {
	ctx  *channel.HandlerContext
	msg  any
	next *tcpEchoInboundEntry
}

type tcpEchoInboundQueue struct {
	mu        sync.Mutex
	owner     *offloadedTCPEchoHandler
	head      *tcpEchoInboundEntry
	tail      *tcpEchoInboundEntry
	free      *tcpEchoInboundEntry
	scheduled bool
	task      tcpEchoTask
}

func (q *tcpEchoInboundQueue) init(owner *offloadedTCPEchoHandler) {
	q.owner = owner
	q.task = q.drain
}

func (q *tcpEchoInboundQueue) submit(ctx *channel.HandlerContext, msg any) error {
	q.mu.Lock()
	entry := q.acquire(ctx, msg)
	if q.tail == nil {
		q.head = entry
		q.tail = entry
	} else {
		q.tail.next = entry
		q.tail = entry
	}
	if q.scheduled {
		q.mu.Unlock()
		return nil
	}
	q.scheduled = true
	task := q.task
	q.mu.Unlock()

	if err := q.owner.worker.submit(task); err != nil {
		q.reject()
		return err
	}
	return nil
}

func (q *tcpEchoInboundQueue) drain() {
	for {
		q.mu.Lock()
		entry := q.head
		if entry == nil {
			q.scheduled = false
			q.mu.Unlock()
			return
		}
		q.head = entry.next
		if q.head == nil {
			q.tail = nil
		}
		entry.next = nil
		q.mu.Unlock()

		q.dispatch(entry)
		q.mu.Lock()
		q.release(entry)
		q.mu.Unlock()
	}
}

func (q *tcpEchoInboundQueue) dispatch(entry *tcpEchoInboundEntry) {
	defer recoverTCPEchoHandler(entry.ctx)
	q.owner.delegate.ChannelRead(entry.ctx, entry.msg)
}

func (q *tcpEchoInboundQueue) reject() {
	q.mu.Lock()
	head := q.head
	q.head = nil
	q.tail = nil
	q.scheduled = false
	q.mu.Unlock()
	for entry := head; entry != nil; {
		next := entry.next
		message.Release(entry.msg)
		q.mu.Lock()
		q.release(entry)
		q.mu.Unlock()
		entry = next
	}
}

func (q *tcpEchoInboundQueue) acquire(ctx *channel.HandlerContext, msg any) *tcpEchoInboundEntry {
	if q.free == nil {
		return &tcpEchoInboundEntry{ctx: ctx, msg: msg}
	}
	entry := q.free
	q.free = entry.next
	entry.ctx = ctx
	entry.msg = msg
	entry.next = nil
	return entry
}

func (q *tcpEchoInboundQueue) release(entry *tcpEchoInboundEntry) {
	entry.ctx = nil
	entry.msg = nil
	entry.next = q.free
	q.free = entry
}
