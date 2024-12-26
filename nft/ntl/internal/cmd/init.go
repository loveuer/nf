package cmd

import (
	"context"
	"fmt"
	"os"
	"time"
)

func Init() {
	initRoot(
		initUpdate(),
		initNew(),
	)
}

func Run(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Printf("❌ %s\n", err.Error())
		os.Exit(1)
	}

	time.Sleep(300 * time.Millisecond)
}
