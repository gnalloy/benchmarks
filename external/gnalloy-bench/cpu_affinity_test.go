package main

import (
	"errors"
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
