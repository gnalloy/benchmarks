package benchh3

import (
	"context"
	"errors"
	"io"
	"sync/atomic"
	"time"

	"gnalloy.org/benchmarks/internal/loadgen"
	codechttp3 "gnalloy.org/codec-http3"
	h3transport "gnalloy.org/transport-http3"
	"gnalloy.org/transport-quic"
)

type client struct {
	conn     quic.Connection
	session  *h3transport.Session
	headers  codechttp3.HeadersBlock
	expected []byte
	reply    []byte
	alpn     string
}

func runClientMessages(ctx context.Context, c *client, clientID int, messageCount int, latencySampleRate int, startCh <-chan struct{}, sharedPacer *loadgen.Pacer, successes *atomic.Int64, samples *clientSamples) error {
	select {
	case <-startCh:
	case <-ctx.Done():
		return ctx.Err()
	}
	pacer := *sharedPacer
	var timer *time.Timer
	if pacer.Enabled() {
		timer = time.NewTimer(time.Hour)
		if !timer.Stop() {
			<-timer.C
		}
		defer timer.Stop()
	}
	for i := 0; i < messageCount; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		deadline := time.Time{}
		if pacer.Enabled() {
			var err error
			deadline, err = pacer.Wait(ctx, timer, clientID, i)
			if err != nil {
				return err
			}
		}
		recordLatency := samples != nil && shouldRecordLatency(i, latencySampleRate)
		var sendStarted time.Time
		if recordLatency {
			sendStarted = time.Now()
		}
		if err := runRequest(ctx, c); err != nil {
			return err
		}
		if recordLatency {
			completedAt := time.Now()
			roundTrip := positiveLatencyNanos(completedAt.Sub(sendStarted))
			samples.roundTrip = append(samples.roundTrip, roundTrip)
			if pacer.Enabled() {
				samples.total = append(samples.total, positiveLatencyNanos(completedAt.Sub(deadline)))
				samples.scheduleDelay = append(samples.scheduleDelay, nonNegativeLatencyNanos(sendStarted.Sub(deadline)))
			} else {
				samples.total = append(samples.total, roundTrip)
			}
		}
		if successes != nil {
			successes.Add(1)
		}
	}
	return nil
}

func runRequest(ctx context.Context, c *client) error {
	streamCh, err := c.session.OpenRequestStream(ctx)
	if err != nil {
		return err
	}
	defer streamCh.Close()
	capture := &responseCapture{expected: c.expected, reply: c.reply[:0]}
	if err := streamCh.Channel().Pipeline().AddLast("capture", capture); err != nil {
		return err
	}
	if err := streamCh.Channel().WriteAndFlush(c.headers); err != nil {
		return err
	}
	if err := readResponse(ctx, streamCh, capture); err != nil {
		return err
	}
	return nil
}

func readResponse(ctx context.Context, streamCh *h3transport.StreamChannel, capture *responseCapture) error {
	for {
		if capture.complete() {
			return nil
		}
		_, err := streamCh.ReadOnce(ctx)
		if capture.err != nil {
			return capture.err
		}
		if err == nil {
			continue
		}
		if errors.Is(err, io.EOF) && capture.complete() {
			return nil
		}
		return err
	}
}
