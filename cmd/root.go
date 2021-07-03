package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/skanehira/remonade/config"
	"github.com/skanehira/remonade/ui"
	"github.com/skanehira/remonade/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "remonade",
}

func Execute() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		config.Load()
		debug := os.Getenv("DEBUG")
		if debug != "" {
			path := filepath.Join(filepath.Dir(config.Path), "debug.log")
			f, err := os.OpenFile(path,
				os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				util.ExitError(err)
			}
			defer f.Close()
			log.SetOutput(f)
		} else {
			log.SetOutput(ioutil.Discard)
		}
		ui.Start()
	}

	if err := rootCmd.Execute(); err != nil {
		util.ExitError(err)
	}
}
