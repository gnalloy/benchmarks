package main

import (
	"errors"
	"runtime"
	"strconv"
	"testing"
)

func TestParseCPUSet(t *testing.T) {
	got, err := parseCPUSet("0, 2,4")
	if err != nil {
		t.Fatal(err)
	}
	want := []int{0, 2, 4}
	if len(got) != len(want) {
		t.Fatalf("cpus=%v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("cpus=%v, want %v", got, want)
		}
	}
}

func TestParseCPUSetEmpty(t *testing.T) {
	got, err := parseCPUSet("")
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Fatalf("cpus=%v, want nil", got)
	}
}

func TestParseCPUSetRejectsInvalidValue(t *testing.T) {
	_, err := parseCPUSet("0,-1")
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigRejectsOversubscribedEventLoops(t *testing.T) {
	workers := runtime.GOMAXPROCS(0)
	_, err := parseConfig([]string{"-boss", "1", "-workers", strconv.Itoa(workers)})
	if !errors.Is(err, errInvalidConfig) {
		t.Fatalf("err=%v, want %v", err, errInvalidConfig)
	}
}

func TestParseConfigRejectsUnsafeCPUAffinity(t *testing.T) {
	tests := [][]string{
		{"-boss", "1", "-workers", "1", "-boss-cpus", "0", "-worker-cpus", "0"},
		{"-boss", "1", "-workers", "2", "-boss-cpus", "0", "-worker-cpus", "1"},
		{"-boss", "1", "-workers", "2", "-boss-cpus", "0", "-worker-cpus", "1,1"},
	}
	for _, args := range tests {
		_, err := parseConfig(args)
		if !errors.Is(err, errInvalidConfig) {
			t.Fatalf("args=%v err=%v, want %v", args, err, errInvalidConfig)
		}
	}
}
