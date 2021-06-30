package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/skanehira/remonade/cmd"
	"github.com/skanehira/remonade/config"
	"github.com/skanehira/remonade/util"
)

func main() {
	config.Init()
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
	cmd.Execute()
}
