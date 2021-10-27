package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/skanehira/remonade/config"
	"github.com/skanehira/remonade/util"
	"github.com/spf13/cobra"
	"github.com/tenntenn/natureremo"
)

var appCmd = &cobra.Command{
	Use: "app",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
	},
}

var appActionCmd = &cobra.Command{
	Use:  "action",
	Long: "executing any actions",
	Run:  runAppAction,
}

var appList = &cobra.Command{
	Use:  "list",
	Long: "appliance list",
	Run: func(cmd *cobra.Command, args []string) {
		config.Load()

		printJSON, err := cmd.Flags().GetBool("json")
		if err != nil {
			util.ExitError(err)
			return
		}

		if printJSON {
			if err := json.NewEncoder(os.Stdout).Encode(config.Config.Apps); err != nil {
				util.ExitError(err)
			}
			return
		}
	},
}

func actionOptions(actions map[string]Action) []string {
	var opts []string
	for a := range actions {
		opts = append(opts, a)
	}
	return opts
}

func appOptions(apps []*natureremo.Appliance) []string {
	var opts []string
	for _, app := range apps {
		opts = append(opts, app.Nickname)
	}
	return opts
}

func sigOpts(app *natureremo.Appliance) []string {
	var sigs []string
	for _, sig := range app.Signals {
		sigs = append(sigs, sig.Name)
	}
	return sigs
}

func btnOpts(buttons []natureremo.DefaultButton) []string {
	var bts []string
	for _, bt := range buttons {
		bts = append(bts, bt.Name)
	}
	return bts
}

func runAppAction(cmd *cobra.Command, args []string) {
	config.Load()
	//cli := natureremo.NewClient(config.Config.Token)

	var appNames []string

	qs := &survey.MultiSelect{
		Message: "Choose an appliance:",
		Options: appOptions(config.Config.Apps),
	}

	err := survey.AskOne(qs, &appNames)
	if err != nil {
		util.ExitError(err)
	}

	var selectApps []*natureremo.Appliance

	for _, name := range appNames {
		for _, app := range config.Config.Apps {
			if name == app.Nickname {
				selectApps = append(selectApps, app)
			}
		}
	}

	qss := []*survey.Question{}
	for _, app := range selectApps {
		switch app.Type {
		case natureremo.ApplianceTypeAirCon:
			// select power or modes
			var selected string
			survey.AskOne(&survey.Select{
				Message: "Choose an option: ",
				Options: []string{"Power", "Mode"},
			}, &selected)

			if selected == "Power" {
				var power string
				survey.AskOne(&survey.Select{
					Message: "ON/OFF",
					Options: []string{"ON", "OFF"},
				}, &power)

				setting := &natureremo.AirConSettings{}
				if power == "ON" {
					setting.Button = natureremo.ButtonPowerOn
				} else {
					setting.Button = natureremo.ButtonPowerOff
				}
				//if err := cli.ApplianceService.UpdateAirConSettings(context.Background(), app, setting); err != nil {
				//	util.ExitError(err)
				//}
				return
			}

			var (
				i    int
				mode string
			)
			modes := make([]string, len(app.AirCon.Range.Modes))
			for m := range app.AirCon.Range.Modes {
				modes[i] = string(m)
				i++
			}
			survey.AskOne(&survey.Select{
				Message: "Mode",
				Options: modes,
			}, &mode)
			// directions

			// temps
			// volumes

			switch natureremo.OperationMode(mode) {
			case natureremo.OperationModeAuto:
			}
		case natureremo.ApplianceTypeLight:
			opts := btnOpts(app.Light.Buttons)
			for _, sig := range app.Signals {
				opts = append(opts, sig.Name)
			}
			qss = append(qss, &survey.Question{
				Name: app.Nickname,
				Prompt: &survey.Select{
					Message: app.Nickname,
					Options: opts,
				},
				Validate: survey.Required,
			})
		case natureremo.ApplianceTypeTV:
			qss = append(qss, &survey.Question{
				Name: app.Nickname,
				Prompt: &survey.Select{
					Message: app.Nickname,
					Options: btnOpts(app.TV.Buttons),
				},
				Validate: survey.Required,
			})
		case natureremo.ApplianceTypeIR:
			// add signals
			qss = append(qss, &survey.Question{
				Name: app.Nickname,
				Prompt: &survey.Select{
					Message: app.Nickname,
					Options: sigOpts(app),
				},
				Validate: survey.Required,
			})
		}
	}

	actions := Action{
		apps: config.Config.Apps,
	}

	err = survey.Ask(qss, &actions)
	if err != nil {
		util.ExitError(err)
	}
}

type Action struct {
	apps []*natureremo.Appliance
}

func (a Action) WriteAnswer(name string, value interface{}) error {
	ans := value.(core.OptionAnswer)
	for _, app := range a.apps {
		if app.Nickname == name {
			switch app.Type {
			case natureremo.ApplianceTypeAirCon:
				//app.AirConSettings
			case natureremo.ApplianceTypeIR:
				sig := app.Signals[ans.Index]
				fmt.Printf("%#+v\n", sig)
			case natureremo.ApplianceTypeLight:
				btn := app.Light.Buttons[ans.Index]
				fmt.Printf("%#+v\n", btn)
			case natureremo.ApplianceTypeTV:
				btn := app.TV.Buttons[ans.Index]
				fmt.Printf("%#+v\n", btn)
			}
		}
	}
	return nil
}

func init() {
	appList.Flags().Bool("json", false, "print json")
	appCmd.AddCommand(appActionCmd)
	appCmd.AddCommand(appList)
	rootCmd.AddCommand(appCmd)
}
