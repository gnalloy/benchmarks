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
		t.Fatalf("cpu profile is empty")
	}
}

func TestStartCPUProfileAllowsEmptyPath(t *testing.T) {
	stop, err := startCPUProfile("")
	if err != nil {
		t.Fatal(err)
	}
	stop()
}

func TestStartAllocProfileWritesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "alloc.pprof")
	stop, err := startAllocProfile(path)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 1000; i++ {
		_ = make([]byte, i+1)
	}
	stop()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Fatal("allocation profile is empty")
	}
}

func TestStartAllocProfileAllowsEmptyPath(t *testing.T) {
	stop, err := startAllocProfile("")
	if err != nil {
		t.Fatal(err)
	}
	stop()
}

func TestStartRuntimeTraceWritesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "runtime.trace")
	stop, err := startRuntimeTrace(path)
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
		t.Fatalf("runtime trace is empty")
	}
}

func TestStartRuntimeTraceAllowsEmptyPath(t *testing.T) {
	stop, err := startRuntimeTrace("")
	if err != nil {
		t.Fatal(err)
	}
	stop()
}
