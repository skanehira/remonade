package ui

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/tenntenn/natureremo"
)

var (
	IndexOutOfAppliances = errors.New("index out of appliances")
)

type Event struct {
	Device  string
	Type    string
	Value   string
	Created time.Time
}

type State struct {
	SelectApplianceIdx int
	Devices            []*natureremo.Device
	Appliances         []*natureremo.Appliance
	Events             []Event
}

func (s *State) String() string {
	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(&s); err != nil {
		return err.Error()
	}
	return b.String()
}

func (s *State) UpdateDevices() {
	UI.Devices.UpdateView(s.Devices)
}

func (s *State) UpdateAppliances() {
	UI.Appliances.UpdateView(s.Appliances)
}

func (s *State) UpdateEvents() {
	UI.Events.UpdateView(s.Events)
}

func (s *State) PushEvent(ev Event) {
	ev.Created = time.Now().Local()
	s.Events = append(s.Events, ev)
}

func (s *State) SelectAppliance() (*natureremo.Appliance, error) {
	if s.SelectApplianceIdx >= len(s.Appliances) {
		log.Println("select idx is greater than state.Appliances's length")
		return nil, IndexOutOfAppliances
	}
	return s.Appliances[s.SelectApplianceIdx], nil
}

func parseEventType(ev string) string {
	switch ev {
	case "te":
		return "temperature"
	case "hu":
		return "humidity"
	case "il":
		return "illumination"
	}
	return ev
}
