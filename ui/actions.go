package ui

import (
	"context"
	"fmt"
	"log"

	"github.com/jinzhu/copier"
	"github.com/tenntenn/natureremo"
)

type (
	Action func(state *State, cli *natureremo.Client, ctx interface{}) error
)

func ActionGetAppliances(state *State, cli *natureremo.Client, ctx interface{}) error {
	apps, err := Client.ApplianceService.GetAll(context.Background())
	if err != nil {
		return err
	}
	state.Appliances = apps
	return nil
}

func ActionGetDevices(state *State, cli *natureremo.Client, ctx interface{}) error {
	devices, err := Client.DeviceService.GetAll(context.Background())
	if err != nil {
		return err
	}
	state.Devices = devices
	for _, dev := range devices {
		for t, e := range dev.NewestEvents {
			ev := Event{
				Device:  dev.Name,
				Type:    parseEventType(string(t)),
				Value:   fmt.Sprintf("%v", e.Value),
				Created: e.CreatedAt,
			}
			state.PushEvent(ev)
		}
	}
	return nil
}

func ActionAppliancesPower(state *State, cli *natureremo.Client, ctx interface{}) error {
	data, ok := ctx.(AppliancePowerOnOff)
	if !ok {
		return fmt.Errorf(`ctx is not "AppliancePowerOnOff": %T`, ctx)
	}

	var row interface{} = data.Row
	app, err := getAppliance(state, row)
	if err != nil {
		return err
	}

	on := data.Power == natureremo.ButtonPowerOn

	switch app.Type {
	case natureremo.ApplianceTypeAirCon:
		btn := natureremo.ButtonPowerOn
		if !on {
			btn = natureremo.ButtonPowerOff
		}
		settings := &natureremo.AirConSettings{
			Button: btn,
		}
		app.AirConSettings.Button = btn
		err = cli.ApplianceService.
			UpdateAirConSettings(context.Background(), app, settings)
	case natureremo.ApplianceTypeLight:
		power := "on"
		if !on {
			power = "off"
		}
		app.Light.State.Power = power
		_, err = cli.ApplianceService.SendLightSignal(context.Background(), app, power)
	case natureremo.ApplianceTypeTV:
		btn := "power"
		_, err = cli.ApplianceService.SendTVSignal(context.Background(), app, btn)
	case natureremo.ApplianceTypeIR:
		// DO NOTHING
		return nil
	default:
		return fmt.Errorf("unsupported appliance: %v", app.Type)
	}

	return err
}

func ActionOpenUpdateApplianceView(state *State, cli *natureremo.Client, ctx interface{}) error {
	app, err := getAppliance(state, ctx)
	if err != nil {
		return err
	}

	switch app.Type {
	case natureremo.ApplianceTypeAirCon:
		UI.Appliances.OpenUpdateAirConView(app)
	case natureremo.ApplianceTypeLight:
		UI.Appliances.OpenUpdateLightView(app)
	case natureremo.ApplianceTypeTV:
		UI.Appliances.OpenUpdateTVView(app)
	case natureremo.ApplianceTypeIR:
		UI.Appliances.OpenUpdateIRView(app)
	default:
		return fmt.Errorf("unsupported appliance type: %s", app.Type)
	}

	return nil
}

func ActionUpdateAirConSettings(state *State, cli *natureremo.Client, ctx interface{}) error {
	form, ok := ctx.(UpdateAirConFormData)
	if !ok {
		return fmt.Errorf(`ctx is invalid type: %T`, ctx)
	}

	oldapp, err := state.SelectAppliance()
	if err != nil {
		return err
	}

	app := &natureremo.Appliance{ID: oldapp.ID}
	settings := &natureremo.AirConSettings{}

	if form.Power.Value() == "ON" {
		settings.Button = natureremo.ButtonPowerOn
	} else {
		settings.Button = natureremo.ButtonPowerOff
	}
	settings.OperationMode = natureremo.OperationMode(form.Mode.Value())
	settings.AirDirection = natureremo.AirDirection(form.Direction.Value())

	mode := form.Mode.Value()
	switch mode {
	case "below":
		settings.AirVolume = natureremo.AirVolume(form.Volume.Value())
	case "cool", "warm":
		settings.Temperature = form.Temp.Value()
		settings.AirVolume = natureremo.AirVolume(form.Volume.Value())
	case "dry":
		settings.Temperature = form.Temp.Value()
	}

	if err := cli.ApplianceService.UpdateAirConSettings(context.Background(), app, settings); err != nil {
		return err
	}

	err = copier.CopyWithOption(oldapp.AirConSettings, settings,
		copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		log.Printf("destination: %#+v, source: %#+v", oldapp.AirConSettings, settings)
		return err
	}
	// NOTE copier option is ignore empty, but when power is on, the value is empty
	// so copier doesn't copy button
	oldapp.AirConSettings.Button = settings.Button

	return nil
}

func ActionUpdateLight(state *State, cli *natureremo.Client, ctx interface{}) error {
	setting, ok := ctx.(LightSetting)
	if !ok {
		return fmt.Errorf("ctx is invalid type: %T", ctx)
	}

	log.Printf("action update light: ctx: %#+v", setting)
	if setting.Type == "button" {
		app := &natureremo.Appliance{
			ID: setting.ID,
		}
		// TODO update light state
		if _, err := cli.ApplianceService.SendLightSignal(context.Background(), app, setting.Value); err != nil {
			return err
		}
	} else {
		signal := &natureremo.Signal{
			ID: setting.Value,
		}
		if err := cli.SignalService.Send(context.Background(), signal); err != nil {
			return err
		}
	}
	return nil
}

func ActionSendTVButton(state *State, cli *natureremo.Client, ctx interface{}) error {
	setting, ok := ctx.(TVSetting)
	if !ok {
		return fmt.Errorf(`ctx is not "TVSetting": %T`, ctx)
	}
	app := &natureremo.Appliance{ID: setting.ID}
	tvState, err := cli.ApplianceService.SendTVSignal(context.Background(), app, setting.Button)
	if err != nil {
		return err
	}

	for _, app := range state.Appliances {
		if app.ID == setting.ID {
			app.TV.State = tvState
		}
	}
	return nil
}

func ActionSendSignal(state *State, cli *natureremo.Client, ctx interface{}) error {
	id, ok := ctx.(string)
	if !ok {
		return fmt.Errorf("ctx is not string: %T", ctx)
	}

	signal := &natureremo.Signal{
		ID: id,
	}
	if err := cli.SignalService.Send(context.Background(), signal); err != nil {
		return err
	}
	return nil
}

func ActionUpdateSelectIdx(state *State, cli *natureremo.Client, ctx interface{}) error {
	idx, ok := ctx.(int)
	if !ok {
		log.Printf(`ctx is not "int": %T\n`, ctx)
		return nil
	}

	state.SelectApplianceIdx = idx
	return nil
}

func getAppliance(state *State, ctx interface{}) (*natureremo.Appliance, error) {
	row, ok := ctx.(int)
	if !ok {
		return nil, fmt.Errorf("ctx is not int: %#+v", ctx)
	}

	if row >= len(state.Appliances) {
		return nil, fmt.Errorf("index out of range, row: %v, state: %#+v", row, state)
	}

	return state.Appliances[row], nil
}
