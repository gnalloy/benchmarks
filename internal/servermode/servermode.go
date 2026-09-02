package servermode

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

type Info struct {
	Framework string
	Protocol  string
	Addr      string
}

func WriteReady(writer io.Writer, info Info) error {
	_, err := fmt.Fprintf(writer, "serverReady=true framework=%s protocol=%s addr=%s\n", info.Framework, info.Protocol, info.Addr)
	return err
}

func Wait(parent context.Context) {
	if parent == nil {
		parent = context.Background()
	}
	ctx, stop := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
}
