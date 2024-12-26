package cmd

import (
	"github.com/loveuer/nf/nft/ntl/internal/opt"
	"github.com/loveuer/nf/pkg/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nfctl",
	Short: "nfctl is a tool for quick start a nf projects",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if opt.Cfg.Debug {
			log.SetLogLevel(log.LogLevelDebug)
		}

		if !opt.Cfg.DisableUpdate {
			doUpdate(cmd.Context())
		}

		return nil
	},
	DisableSuggestions: true,
	SilenceUsage:       true,

	Run: func(cmd *cobra.Command, args []string) {},
}

func initRoot(cmds ...*cobra.Command) {
	rootCmd.PersistentFlags().BoolVar(&opt.Cfg.Debug, "debug", false, "debug mode")
	rootCmd.PersistentFlags().BoolVar(&opt.Cfg.DisableUpdate, "disable-update", false, "disable self update")
	rootCmd.AddCommand(cmds...)
}
