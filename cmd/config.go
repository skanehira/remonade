package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/skanehira/remonade/config"
	"github.com/skanehira/remonade/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	ErrEmptyToken     = errors.New("access token is empty")
	ErrNotExistConfig = errors.New("config file is not exists")
	ErrEmptyEDITOR    = errors.New("$EDITOR is empty")
)

func runInit(path string) error {
	var (
		f   *os.File
		err error
	)

	if util.NotExist(path) {
		if err = config.Create(path); err != nil {
			return err
		}
	}

	f, err = os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Print("Your access token: ")
	sc := bufio.NewScanner(os.Stdin)

	sc.Scan()
	token := sc.Text()
	if token == "" {
		return ErrEmptyToken
	}

	config.Config.Token = strings.Trim(token, "\r\n")
	return yaml.NewEncoder(f).Encode(config.Config)
}

func runEdit(path string) error {
	if util.NotExist(path) {
		return ErrNotExistConfig
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		return ErrEmptyEDITOR
	}
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		_ = cmd.Help()
		os.Exit(1)
	}

	if args[0] != "edit" && args[0] != "init" {
		_ = cmd.Help()
		os.Exit(1)
	}

	if args[0] == "init" {
		if err := runInit(config.Path); err != nil {
			util.ExitError(err)
		}
	}

	if args[0] == "edit" {
		if err := runEdit(config.Path); err != nil {
			util.ExitError(err)
		}
	}
}

var configCmd = &cobra.Command{
	Use:   "config",
	Run:   run,
	Short: "Edit or setup config",
}

func init() {
	configCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Print(`Config command

Usage:
  remonade config [edit|init]

Command:
  init                Setup config
  edit                Edit config

Flags:
  -h, --help          help for config
`)
	})
	rootCmd.AddCommand(configCmd)
}
