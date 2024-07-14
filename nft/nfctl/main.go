package main

import (
	"context"
	"github.com/loveuer/nf/nft/nfctl/cmd"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	//if !(len(os.Args) >= 2 && os.Args[1] == "version") {
	//	version.Check(5)
	//}

	_ = cmd.Root.ExecuteContext(ctx)
}
