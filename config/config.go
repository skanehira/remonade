package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skanehira/remonade/util"
	"github.com/tenntenn/natureremo"
	"gopkg.in/yaml.v3"
)

type config struct {
	Token string                  `yaml:"token"`
	Apps  []*natureremo.Appliance `yaml:"apps"`
}

var (
	Config config
	Path   string
)

func init() {
	path, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("cannot get user config dir: %w", err))
		os.Exit(1)
	}
	Path = filepath.Join(path, "remonade", "config.yaml")
}

func Load() {
	if util.NotExist(Path) {
		if err := util.Create(Path); err != nil {
			util.ExitError(fmt.Errorf("cannot create file %s: %w", Path, err))
		}
		return
	}

	f, err := os.Open(Path)
	if err != nil {
		util.ExitError(err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&Config); err != nil {
		util.ExitError(fmt.Errorf("cannot decode %s: %w", Path, err))
	}
}
