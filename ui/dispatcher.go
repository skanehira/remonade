package ui

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/google/go-cmp/cmp"
)

type dispatcher struct {
	state *State
}

var Dispatcher = &dispatcher{
	state: &State{},
}

func copyState(state *State) *State {
	// NOTE "github.com/jinzhu/copier" doesn't copy time.Time, using json instead
	// https://github.com/jinzhu/copier/pull/103
	newState := &State{}
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(&state); err != nil {
		log.Println(err)
		return state
	}

	if err := json.NewDecoder(buf).Decode(&newState); err != nil {
		log.Println(err)
		return state
	}
	return newState
}

func (d *dispatcher) Dispatch(action Action, ctx interface{}) {
	old := copyState(d.state)
	err := action(d.state, Client, ctx)
	if err != nil {
		UI.Message(err.Error())
		return
	}

	go UI.app.QueueUpdateDraw(func() {
		if !cmp.Equal(old.Appliances, d.state.Appliances) {
			log.Println("update appliance view")
			d.state.UpdateAppliances()
		}

		if !cmp.Equal(old.Devices, d.state.Devices) {
			log.Println("update devices view")
			d.state.UpdateDevices()
		}

		if !cmp.Equal(old.Events, d.state.Events) {
			log.Println("update events view")
			d.state.UpdateEvents()
		}
	})
}
