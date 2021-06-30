package ui

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/tenntenn/natureremo"
)

type Event struct {
	Type    string
	Value   string
	Created time.Time
}

type State struct {
	Devices    []*natureremo.Device
	Appliances []*natureremo.Appliance
	Events     []Event
}

func (s *State) String() string {
	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(&s); err != nil {
		return err.Error()
	}
	return b.String()
}

func (s *State) UpdateDevices() {
	UI.devices.UpdateView(s.Devices)
}

func (s *State) UpdateAppliances() {
	UI.appliances.UpdateView(s.Appliances)
}

func (s *State) UpdateEvents() {
	UI.events.UpdateView(s.Events)
}
