package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/skanehira/remonade/config"
	"github.com/skanehira/remonade/util"
	"github.com/spf13/cobra"
	"github.com/tenntenn/natureremo"
	"gopkg.in/yaml.v3"
)

var (
	ErrEmptyToken     = errors.New("access token is empty")
	ErrNotExistConfig = errors.New("config file is not exists")
	ErrEmptyEDITOR    = errors.New("$EDITOR is empty")
)

func runInit(path string) error {
	var (
		err error
	)

	if util.NotExist(path) {
		if err = util.Create(path); err != nil {
			return err
		}
	}

	fmt.Print("Your access token: ")
	sc := bufio.NewScanner(os.Stdin)

	sc.Scan()
	token := sc.Text()
	if token == "" {
		return ErrEmptyToken
	}

	fmt.Println("Initializing...")
	return updateConfig(path, token)
}

func updateConfig(path string, tokens ...string) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&config.Config); err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}

	var token string
	// use config's token
	if len(tokens) == 0 {
		if config.Config.Token == "" {
			return ErrEmptyToken
		}
		token = config.Config.Token
	} else {
		token = tokens[0]
	}

	cli := natureremo.NewClient(token)
	apps, err := cli.ApplianceService.GetAll(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get appliances: %w", err)
	}

	config.Config.Token = token
	config.Config.Apps = apps

	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate config.yaml: %w", err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek config.yaml: %w", err)
	}
	return yaml.NewEncoder(f).Encode(config.Config)
}

func runUpdate(path string) error {
	fmt.Println("Updating...")
	return updateConfig(path)
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

	switch args[0] {
	case "init", "edit", "update":
	default:
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

	if args[0] == "update" {
		if err := runUpdate(config.Path); err != nil {
			util.ExitError(err)
		}
	}
}

var configCmd = &cobra.Command{
	Use: "config",
	Run: run,
}

func init() {
	configCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Print(`Config command

Usage:
  remonade config [edit|init|update]

Command:
  init                Setup config
  edit                Edit config
  update              Update config

Flags:
  -h, --help          help for config
`)
	})
	rootCmd.AddCommand(configCmd)
}
