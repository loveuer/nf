package main

import (
	"context"
	"github.com/loveuer/nf/nft/nfctl/cmd"
	"github.com/loveuer/nf/nft/nfctl/version"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	version.Check()
	defer version.Fn()

	_ = cmd.Root.ExecuteContext(ctx)
}
