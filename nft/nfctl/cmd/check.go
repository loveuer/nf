package cmd

import (
	"github.com/loveuer/nf/nft/nfctl/version"
	"github.com/spf13/cobra"
)

var (
	checkCmd = &cobra.Command{
		Use:   "check",
		Short: "nfctl new version check",
		Run: func(cmd *cobra.Command, args []string) {
			version.Check(true, true, 30)
		},
	}
)
