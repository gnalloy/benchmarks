package main

import (
	"os"
	"runtime/pprof"
	"runtime/trace"
	"strings"
)

func startCPUProfile(path string) (func(), error) {
	if strings.TrimSpace(path) == "" {
		return func() {}, nil
	}
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if err := pprof.StartCPUProfile(file); err != nil {
		_ = file.Close()
		return nil, err
	}
	return func() {
		pprof.StopCPUProfile()
		_ = file.Close()
	}, nil
}

func startRuntimeTrace(path string) (func(), error) {
	if strings.TrimSpace(path) == "" {
		return func() {}, nil
	}
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if err := trace.Start(file); err != nil {
		_ = file.Close()
		return nil, err
	}
	return func() {
		trace.Stop()
		_ = file.Close()
	}, nil
}
