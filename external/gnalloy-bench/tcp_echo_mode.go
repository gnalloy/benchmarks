package main

import (
	"fmt"
	"strings"
)

const (
	defaultTCPEchoMode = tcpEchoModeDirect

	tcpEchoModeDirect        = "direct"
	tcpEchoModeReadComplete  = "read-complete"
	tcpEchoModeOwnerExecutor = "owner-executor"

	defaultTCPEchoExecutorQueueSize = 4096
)

func normalizeTCPEchoMode(mode string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "", tcpEchoModeDirect:
		return tcpEchoModeDirect, nil
	case tcpEchoModeReadComplete:
		return tcpEchoModeReadComplete, nil
	case tcpEchoModeOwnerExecutor:
		return tcpEchoModeOwnerExecutor, nil
	default:
		return "", fmt.Errorf("%w: tcp-echo-mode %s", errInvalidConfig, mode)
	}
}

func tcpEchoModeUsesExecutor(mode string) bool {
	return mode == tcpEchoModeOwnerExecutor
}
