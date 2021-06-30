package ui

import (
	"fmt"
)

type dispatcher struct {
	state   *State
	actions map[Action]ActionFunc
}

func (d *dispatcher) Dispatch(action Action, ctx interface{}) {
	f, ok := d.actions[action]
	if !ok {
		msg := fmt.Sprintf("doesn't register action: %v\n", action)
		UI.Message(msg)
		return
	}

	newState, err := f(*d.state, action, ctx)
	if err != nil {
		UI.Message(err.Error())
		return
	}

	d.state = &newState

	go UI.app.QueueUpdateDraw(func() {
		d.state.UpdateAppliances()
		d.state.UpdateDevices()
		d.state.UpdateEvents()
	})
}

var Dispatcher *dispatcher

func init() {
	Dispatcher = &dispatcher{
		state: &State{},
		actions: map[Action]ActionFunc{
			GetAppliances:           ActionGetAppliances,
			GetDevices:              ActionGetDevices,
			PowerON:                 ActionAppliancesPower,
			PowerOFF:                ActionAppliancesPower,
			OpenUpdateApplianceView: ActionOpenUpdateApplianceView,
		},
	}
}
