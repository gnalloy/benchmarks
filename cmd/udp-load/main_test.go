package main

import (
	"bytes"
	"net"
	"strings"
	"testing"
	"time"
)

func TestRunRejectsEmptyServerFramework(t *testing.T) {
	var output bytes.Buffer
	if err := run([]string{"-server-framework", " "}, &output); err == nil {
		t.Fatal("expected error")
	}
}

func TestRunWritesPacedLatencyComponents(t *testing.T) {
	server, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = server.Close() })
	go func() {
		payload := make([]byte, 64)
		for {
			n, addr, readErr := server.ReadFrom(payload)
			if readErr != nil {
				return
			}
			if _, writeErr := server.WriteTo(payload[:n], addr); writeErr != nil {
				return
			}
		}
	}()

	var output bytes.Buffer
	err = run([]string{
		"-server-framework", "test",
		"-addr", server.LocalAddr().String(),
		"-payload", "16",
		"-connections", "1",
		"-messages", "2",
		"-warmup-messages", "0",
		"-latency-sample-rate", "1",
		"-target-rate", "1000",
		"-timeout", time.Second.String(),
	}, &output)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"scheduleDelaySamples=2", "p99ScheduleDelayNs=", "roundTripLatencySamples=2", "p99RoundTripLatencyNs="} {
		if !strings.Contains(output.String(), field) {
			t.Fatalf("output %q does not contain %q", output.String(), field)
		}
	}
}
