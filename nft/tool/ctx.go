package tool

import (
	"context"
	"time"
)

func Timeout(seconds ...int) (context.Context, context.CancelFunc) {
	var duration time.Duration

	if len(seconds) > 0 && seconds[0] > 0 {
		duration = time.Duration(seconds[0]) * time.Second
	} else {
		duration = time.Duration(30) * time.Second
	}

	return context.WithTimeout(context.Background(), duration)
}

func TimeoutCtx(ctx context.Context, seconds ...int) (context.Context, context.CancelFunc) {
	var duration time.Duration

	if len(seconds) > 0 && seconds[0] > 0 {
		duration = time.Duration(seconds[0]) * time.Second
	} else {
		duration = time.Duration(30) * time.Second
	}

	return context.WithTimeout(ctx, duration)
}
