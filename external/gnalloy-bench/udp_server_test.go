package main

import (
	"net"
	"testing"

	"gnalloy.org/gnalloy/buffer"
	"gnalloy.org/gnalloy/channel/embedded"
	"gnalloy.org/transport-udp"
)

func TestUDPEchoHandlerFlushesReadBatchOnce(t *testing.T) {
	ch, err := embedded.New(udpEchoHandler{})
	if err != nil {
		t.Fatal(err)
	}
	defer ch.FinishAndReleaseAll()

	for _, payload := range []string{"one", "two"} {
		ch.Pipeline().FireChannelRead(testUDPDatagram(payload))
	}
	if ch.Flushes() != 0 {
		t.Fatalf("flushes=%d before read complete, want 0", ch.Flushes())
	}
	ch.Pipeline().FireChannelReadComplete()
	if ch.Flushes() != 1 {
		t.Fatalf("flushes=%d after read complete, want 1", ch.Flushes())
	}
	for _, want := range []string{"one", "two"} {
		msg, ok := ch.ReadOutbound()
		if !ok {
			t.Fatalf("missing outbound payload %q", want)
		}
		datagram, ok := msg.(udp.Datagram)
		if !ok {
			t.Fatalf("message=%T, want udp.Datagram", msg)
		}
		if got := string(datagram.Payload.Bytes()); got != want {
			t.Fatalf("payload=%q, want %q", got, want)
		}
		datagram.Release()
	}
}

func testUDPDatagram(payload string) udp.Datagram {
	buf := buffer.NewHeapBuffer(len(payload))
	_, _ = buf.WriteBytes([]byte(payload))
	return udp.Datagram{
		Payload: buf,
		Addr:    udp.Address{IP: net.IPv4(127, 0, 0, 1), Port: 19090},
	}
}
