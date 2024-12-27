package cmd

import (
	"context"
	"fmt"
	"os"
	"time"
)

func Init() {
	initRoot(
		initVersion(),
		initUpdate(),
		initNew(),
	)
}

func Run(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Printf("❌ %s\n", err.Error())
		time.Sleep(300 * time.Millisecond)
		os.Exit(1)
	}

	time.Sleep(300 * time.Millisecond)
}
