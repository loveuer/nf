package cmd

import (
	"os"

	"github.com/loveuer/ursa/ursatool/log"
	"github.com/loveuer/ursa/ursatool/ursactl/internal/opt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ursactl",
	Short: "ursactl is a tool for quick start a nf projects",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if opt.Cfg.Debug {
			log.SetLogLevel(log.LogLevelDebug)
		}

		if opt.Cfg.Version {
			doVersion(cmd, args)
			os.Exit(0)
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
	rootCmd.PersistentFlags().BoolVarP(&opt.Cfg.Version, "version", "v", false, "print ursactl version")
	rootCmd.AddCommand(cmds...)
}
