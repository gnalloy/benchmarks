package main

import "errors"

var (
	errInvalidConfig       = errors.New("gnalloy-bench: invalid config")
	errUnsupportedProtocol = errors.New("gnalloy-bench: unsupported protocol")
	errInvalidBackend      = errors.New("gnalloy-bench: invalid backend")
	errTCPEchoExecutorFull = errors.New("gnalloy-bench: tcp echo executor queue full")
	errTCPEchoExecutorDone = errors.New("gnalloy-bench: tcp echo executor closed")
	errTCPEchoHandlerPanic = errors.New("gnalloy-bench: tcp echo handler panic")
)
