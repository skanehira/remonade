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
		if err := Create(Path); err != nil {
			util.ExitError(fmt.Errorf("cannot create file %s: %w", Path, err))
		}
		return
	}

	f, err := os.Open(Path)
	if err != nil {
		util.ExitError(err)
	}
	defer func() {
		_ = f.Close()
	}()

	if err := yaml.NewDecoder(f).Decode(&Config); err != nil {
		util.ExitError(fmt.Errorf("cannot decode %s: %w", Path, err))
	}
}

func Create(path string) error {
	base := filepath.Dir(path)

	if util.NotExist(base) {
		if err := os.Mkdir(base, os.ModePerm); err != nil {
			return err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	return nil
}
