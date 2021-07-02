package ui

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/tenntenn/natureremo"
)

type (
	// process action, must return new state
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
			state.Events = append(state.Events, Event{
				Type:    string(t),
				Value:   fmt.Sprintf("%v", e.Value),
				Created: e.CreatedAt,
			})
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
		UI.appliances.OpenUpdateAirConView(app)
	case natureremo.ApplianceTypeLight:
		// TODO
	}

	return nil
}

func ActionUpdateAirConSettings(state *State, cli *natureremo.Client, ctx interface{}) error {
	data, ok := ctx.(map[int]UpdateAirConFormData)
	if !ok {
		return fmt.Errorf(`ctx is invalid type: %T`, ctx)
	}

	var (
		idx  int
		form UpdateAirConFormData
	)

	for idx, form = range data {
		break
	}

	oldapp := state.Appliances[idx]

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
	case "cool", "swarm":
		settings.Temperature = form.Temp.Value()
		settings.AirVolume = natureremo.AirVolume(form.Volume.Value())
	case "dry":
		settings.Temperature = form.Temp.Value()
	}

	if err := Client.ApplianceService.UpdateAirConSettings(context.Background(), app, settings); err != nil {
		return err
	}

	err := copier.CopyWithOption(oldapp.AirConSettings, settings,
		copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return err
	}
	// NOTE copier option is ignore empty, but when power is on, the value is empty
	// so copier doesn't copy button
	oldapp.AirConSettings.Button = settings.Button

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
