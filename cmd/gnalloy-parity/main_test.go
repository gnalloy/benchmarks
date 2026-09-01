package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestRunDumpLinuxFullSpec(t *testing.T) {
	var stdout bytes.Buffer
	if err := run([]string{"-matrix", "linux-full", "-dump-spec"}, &stdout, ioDiscard{}); err != nil {
		t.Fatal(err)
	}
	text := stdout.String()
	for _, want := range []string{`"name": "gnalloy linux full parity matrix"`, `"framework": "fasthttp"`, `"protocol": "quic-stream"`} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in %s", want, text)
		}
	}
}

func TestParseFlagsRejectsAmbiguousSource(t *testing.T) {
	_, err := parseFlags([]string{"-spec", "a.json", "-matrix", "linux-full"}, ioDiscard{})
	if err == nil || !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("err=%v, want mutually exclusive", err)
	}
}

func TestParseFlagsRejectsNonPositiveTimeout(t *testing.T) {
	_, err := parseFlags([]string{"-matrix", "linux-full", "-timeout", "0s"}, ioDiscard{})
	if err == nil || !strings.Contains(err.Error(), "timeout") {
		t.Fatalf("err=%v, want timeout error", err)
	}
}

func TestFormatValue(t *testing.T) {
	var f formatValue
	if err := f.Set("json"); err != nil {
		t.Fatal(err)
	}
	if got := f.String(); got != "json" {
		t.Fatalf("format=%q, want json", got)
	}
}

func TestLoadSpecRejectsUnknownMatrix(t *testing.T) {
	_, err := loadSpec(config{matrixName: "missing", timeout: time.Second})
	if err == nil || !strings.Contains(err.Error(), "unknown matrix") {
		t.Fatalf("err=%v, want unknown matrix", err)
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}
