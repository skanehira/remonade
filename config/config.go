package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skanehira/remonade/util"
	"gopkg.in/yaml.v3"
)

type config struct {
	Token string `yaml:"token"`
}

var Config config
var Path string

func Init() {
	path, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	Path = filepath.Join(path, "remonade", "config.yaml")

	f, err := os.Open(Path)
	if err != nil {
		util.ExitError(err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&Config); err != nil {
		util.ExitError(err)
	}
}
