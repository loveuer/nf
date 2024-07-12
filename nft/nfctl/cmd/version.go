package cmd

import (
	"github.com/loveuer/nf/nft/log"
	"github.com/loveuer/nf/nft/nfctl/version"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print nfctl version and exit",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info("version: %s", version.Version)
		},
	}
)
