package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "remonade",
}

func Execute() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		_ = rootCmd.Help()
	}

	if err := rootCmd.Execute(); err != nil {
		util.ExitError(err)
	}
}
