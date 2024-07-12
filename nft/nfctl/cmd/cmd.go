package cmd

import "github.com/spf13/cobra"

var (
	Root = &cobra.Command{
		Use:   "nfctl",
		Short: "nfctl: easy start your nf backend work",
	}
)

func init() {
	initNew()

	Root.AddCommand(
		versionCmd,
		cmdNew,
	)
}
