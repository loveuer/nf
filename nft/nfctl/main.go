package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/loveuer/nf/nft/nfctl/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	_ = cmd.Root.ExecuteContext(ctx)
}
