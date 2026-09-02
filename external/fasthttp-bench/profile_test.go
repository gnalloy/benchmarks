package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStartCPUProfileWritesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cpu.pprof")
	stop, err := startCPUProfile(path)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 100000; i++ {
		_ = i * i
	}
	stop()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Fatal("cpu profile is empty")
	}
}

func TestStartCPUProfileAllowsEmptyPath(t *testing.T) {
	stop, err := startCPUProfile("")
	if err != nil {
		t.Fatal(err)
	}
	stop()
}
