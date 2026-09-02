package main

import "context"

func runLoad(ctx context.Context, addr string, cfg config) (benchResult, error) {
	return runTCPEchoLoad(ctx, addr, cfg)
}
