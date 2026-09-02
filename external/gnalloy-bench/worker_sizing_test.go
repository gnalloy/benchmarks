package main

import (
	"testing"

	"gnalloy.org/gnalloy/transport"
)

func TestDefaultWorkerCount(t *testing.T) {
	tests := []struct {
		name  string
		input workerSizingInput
		want  int
	}{
		{
			name:  "windows iocp cap",
			input: workerSizingInput{GOOS: "windows", Backend: transport.BackendIOCP, GOMAXPROCS: 16, Bosses: 1},
			want:  8,
		},
		{
			name:  "linux epoll cap",
			input: workerSizingInput{GOOS: "linux", Backend: transport.BackendEpoll, GOMAXPROCS: 8, Bosses: 1},
			want:  4,
		},
		{
			name:  "linux io uring cap",
			input: workerSizingInput{GOOS: "linux", Backend: transport.BackendIOUring, GOMAXPROCS: 8, Bosses: 1},
			want:  4,
		},
		{
			name:  "standard backend reserves boss",
			input: workerSizingInput{GOOS: "linux", Backend: transport.BackendStd, GOMAXPROCS: 16, Bosses: 1},
			want:  15,
		},
		{
			name:  "invalid cpu count",
			input: workerSizingInput{GOOS: "windows", Backend: transport.BackendIOCP, GOMAXPROCS: 0, Bosses: 1},
			want:  1,
		},
		{
			name:  "four schedulers reserve one boss",
			input: workerSizingInput{GOOS: "linux", Backend: transport.BackendEpoll, GOMAXPROCS: 4, Bosses: 1},
			want:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultWorkerCount(tt.input); got != tt.want {
				t.Fatalf("workers=%d, want %d", got, tt.want)
			}
		})
	}
}
