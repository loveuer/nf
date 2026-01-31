package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/loveuer/ursa/ursatool/ursactl/internal/cmd"
)

func init() {
	time.Local = time.FixedZone("CST", 8*3600)
	cmd.Init()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	cmd.Run(ctx)
}
