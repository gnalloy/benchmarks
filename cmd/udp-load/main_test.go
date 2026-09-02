package main

import (
	"bytes"
	"testing"
)

func TestRunRejectsEmptyServerFramework(t *testing.T) {
	var output bytes.Buffer
	if err := run([]string{"-server-framework", " "}, &output); err == nil {
		t.Fatal("expected error")
	}
}
