package cmd

import (
	"github.com/skanehira/remonade/ui"
	"github.com/skanehira/remonade/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "remonade",
}

func Execute() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		ui.Start()
	}

	if err := rootCmd.Execute(); err != nil {
		util.ExitError(err)
	}
}
