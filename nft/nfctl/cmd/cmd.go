package cmd

import (
	"github.com/loveuer/nf/nft/log"
	"github.com/loveuer/nf/nft/nfctl/opt"
	"github.com/spf13/cobra"
)

var (
	Root = &cobra.Command{
		Use:   "nfctl",
		Short: "nfctl: easy start your nf backend work",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if opt.Debug == true {
				log.SetLogLevel(log.LogLevelDebug)
			}
		},
	}
)

func init() {
	initNew()
	Root.PersistentFlags().BoolVar(&opt.Debug, "debug", false, "debug mode")

	Root.AddCommand(
		versionCmd,
		checkCmd,
		cmdNew,
	)
}
