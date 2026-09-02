package main

import (
	"net"
	"net/http"
	"strings"
	"testing"

	"gnalloy.org/benchmarks/internal/httpbench"
)

func TestRun(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	body := httpbench.ResponseBody(8)
	server := &http.Server{Handler: http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		_, _ = writer.Write(body)
	})}
	go func() {
		_ = server.Serve(listener)
	}()
	t.Cleanup(func() { _ = server.Close() })

	var output strings.Builder
	err = run([]string{
		"-addr", listener.Addr().String(),
		"-payload", "8",
		"-connections", "1",
		"-messages", "1",
		"-warmup-messages", "1",
		"-latency-sample-rate", "1",
		"-server-framework", "test",
	}, &output)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"framework=common-http1-client",
		"serverFramework=test",
		"protocol=http1",
		"roundTripLatencySamples=1",
		"total=1",
		"errors=0",
	} {
		if !strings.Contains(output.String(), want) {
			t.Fatalf("missing %q in %s", want, output.String())
		}
	}
}

func TestRunRejectsEmptyServerFramework(t *testing.T) {
	err := run([]string{"-server-framework", " "}, &strings.Builder{})
	if err == nil {
		t.Fatal("expected empty server framework error")
	}
}
