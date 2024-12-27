package cmd

import (
	"fmt"

	"github.com/loveuer/nf/nft/nfctl/internal/opt"
	"github.com/spf13/cobra"
)

func initVersion() *cobra.Command {
	return &cobra.Command{
		Use: "version",
		Run: doVersion,
	}
}

func doVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("%s\nnfctl: %s\n\n", opt.Banner, opt.Version)
}
