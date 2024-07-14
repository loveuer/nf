package cmd

import (
	"github.com/fatih/color"
	"github.com/loveuer/nf/nft/nfctl/version"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print nfctl version and exit",
		Run: func(cmd *cobra.Command, args []string) {
			color.Cyan("nfctl - version: %s", version.Version)
			version.Check(true, false, 5)
		},
	}
)
