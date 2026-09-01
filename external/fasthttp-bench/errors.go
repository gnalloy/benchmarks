package main

import "errors"

var (
	errInvalidConfig       = errors.New("fasthttp-bench: invalid config")
	errUnsupportedProtocol = errors.New("fasthttp-bench: unsupported protocol")
)
