package loading

import (
	"context"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestLoadingPrint(t *testing.T) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	Do(ctx, func(ctx context.Context, print func(msg string, types ...Type)) error {
		print("start task 1...")
		time.Sleep(3 * time.Second)

		print("warning...1", TypeWarning)

		time.Sleep(2 * time.Second)

		return nil
	})
}
