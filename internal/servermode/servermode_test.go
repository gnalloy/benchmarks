package servermode

import (
	"bytes"
	"context"
	"testing"
)

func TestWriteReady(t *testing.T) {
	var output bytes.Buffer
	err := WriteReady(&output, Info{Framework: "gnalloy", Protocol: "udp-echo", Addr: "127.0.0.1:19090"})
	if err != nil {
		t.Fatal(err)
	}
	want := "serverReady=true framework=gnalloy protocol=udp-echo addr=127.0.0.1:19090\n"
	if output.String() != want {
		t.Fatalf("output=%q, want %q", output.String(), want)
	}
}

func TestWaitReturnsWhenContextIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	Wait(ctx)
}
