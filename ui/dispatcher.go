package ui

import (
	"log"
)

type dispatcher struct {
	state   *State
	actions map[Action]ActionFunc
}

func (d *dispatcher) Dispatch(action Action, ctx interface{}) {
	f, ok := d.actions[action]
	if !ok {
		log.Printf("doesn't register action: %v\n", action)
		return
	}

	newState, err := f(*d.state, action, ctx)
	if err != nil {
		log.Println(err)
		return
	}

	d.state = &newState

	d.state.UpdateAppliances()
	d.state.UpdateDevices()
}

var Dispatcher *dispatcher

func init() {
	Dispatcher = &dispatcher{
		state: &State{},
		actions: map[Action]ActionFunc{
			GetAppliances: ActionGetAppliances,
			GetDevices:    ActionGetDevices,
			PowerON:       ActionAppliancesPower,
			PowerOFF:      ActionAppliancesPower,
		},
	}
}
