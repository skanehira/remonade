package ui

import (
	"context"
	"fmt"

	"github.com/tenntenn/natureremo"
)

type (
	// process action, must return new state
	ActionFunc func(state *State, action Action, ctx interface{}) error
)

// Action type
type Action string

var (
	GetAppliances           Action = "get appliances"
	GetDevices              Action = "get devices"
	PowerON                 Action = "power on"
	PowerOFF                Action = "power off"
	OpenUpdateApplianceView Action = "open update appliance view"
	UpdateAirConSettings    Action = "update aircon settings"
)

func ActionGetAppliances(state *State, action Action, ctx interface{}) error {
	apps, err := Client.ApplianceService.GetAll(context.Background())
	if err != nil {
		return err
	}
	state.Appliances = apps
	return nil
}

func ActionGetDevices(state *State, action Action, ctx interface{}) error {
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

func ActionAppliancesPower(state *State, action Action, ctx interface{}) error {
	app, err := getAppliance(state, ctx)
	if err != nil {
		return err
	}

	on := action == PowerON

	switch app.Type {
	case natureremo.ApplianceTypeAirCon:
		app.AirConSettings.Button = natureremo.ButtonPowerOn

		if !on {
			app.AirConSettings.Button = natureremo.ButtonPowerOff
		}
		err = Client.ApplianceService.
			UpdateAirConSettings(context.Background(), app, app.AirConSettings)
	case natureremo.ApplianceTypeLight:
		power := "on"
		if !on {
			power = "off"
		}
		app.Light.State.Power = power
		_, err = Client.ApplianceService.SendLightSignal(context.Background(), app, power)
	case natureremo.ApplianceTypeTV:
		btn := "power"
		_, err = Client.ApplianceService.SendTVSignal(context.Background(), app, btn)
	default:
		return fmt.Errorf("unsupported appliance: %v", app.Type)
	}

	return err
}

func ActionOpenUpdateApplianceView(state *State, action Action, ctx interface{}) error {
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

func ActionOpenUpdateAirConSettings(state *State, action Action, ctx interface{}) error {
	data, ok := ctx.(map[string]UpdateAirConFormData)
	if !ok {
		return fmt.Errorf(`ctx type is not valid type: %T`, ctx)
	}

	var (
		id   string
		form UpdateAirConFormData
	)

	for id, form = range data {
		break
	}

	app := &natureremo.Appliance{ID: id}
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

	newApps, err := Client.ApplianceService.GetAll(context.Background())
	if err != nil {
		return err
	}
	state.Appliances = newApps

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
