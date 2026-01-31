package cmd

import (
	"fmt"

	"github.com/loveuer/ursa/ursatool/ursactl/internal/opt"
	"github.com/spf13/cobra"
)

func initVersion() *cobra.Command {
	return &cobra.Command{
		Use: "version",
		Run: doVersion,
	}
}

func doVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("%s\nursactl: %s\n\n", opt.Banner, opt.Version)
}
