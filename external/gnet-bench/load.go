package main

import "context"

func runLoad(ctx context.Context, addr string, cfg config) (benchResult, error) {
	if cfg.Protocol == "udp-echo" {
		return runUDPEchoLoad(ctx, addr, cfg)
	}
	return runTCPEchoLoad(ctx, addr, cfg)
}
