package ui

import (
	"context"
	"fmt"

	"github.com/tenntenn/natureremo"
)

type (
	// process action, must return new state
	ActionFunc func(state State, action Action, ctx interface{}) (State, error)
)

// Action type
type Action string

var (
	GetAppliances Action = "get appliances"
	GetDevices    Action = "get devices"
	PowerON       Action = "power on"
	PowerOFF      Action = "power off"
)

func ActionGetAppliances(state State, action Action, ctx interface{}) (State, error) {
	apps, err := Client.ApplianceService.GetAll(context.Background())
	if err != nil {
		return state, err
	}
	state.Appliances = apps
	return state, nil
}

func ActionGetDevices(state State, action Action, ctx interface{}) (State, error) {
	devices, err := Client.DeviceService.GetAll(context.Background())
	if err != nil {
		return state, err
	}
	state.Devices = devices
	return state, nil
}

func ActionAppliancesPower(state State, action Action, ctx interface{}) (State, error) {
	row, ok := ctx.(int)
	if !ok {
		return state, fmt.Errorf("ctx is not int: %#+v", ctx)
	}

	if row >= len(state.Appliances) {
		return state, fmt.Errorf("index out of range, row: %v, state: %#+v", row, state)
	}

	app := state.Appliances[row]
	on := action == PowerON

	var (
		err error
	)

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
		return state, fmt.Errorf("unsupported appliance: %v", app.Type)
	}

	return state, err
}
