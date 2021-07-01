package ui

import (
	"encoding/json"

	"github.com/tenntenn/natureremo"
)

type UpdateAirConFormItem struct {
	Current int
	Values  []string
}

func (u UpdateAirConFormItem) Value() string {
	if len(u.Values) == 0 {
		return ""
	}
	return u.Values[u.Current]
}

type UpdateAirConFormData struct {
	Power     UpdateAirConFormItem
	Mode      UpdateAirConFormItem
	Temp      UpdateAirConFormItem
	Volume    UpdateAirConFormItem
	Direction UpdateAirConFormItem
}

func (u UpdateAirConFormData) String() string {
	m := map[string]interface{}{
		"Power":     u.Power.Values[u.Power.Current],
		"Mode":      u.Mode.Values[u.Mode.Current],
		"Temp":      u.Temp.Values[u.Temp.Current],
		"Volume":    u.Volume.Values[u.Volume.Current],
		"Direction": u.Direction.Values[u.Direction.Current],
	}

	b, err := json.Marshal(&m)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func ToUpdateAirConViewData(app *natureremo.Appliance) UpdateAirConFormData {
	form := UpdateAirConFormData{}

	// power
	// if button is not empty, the power is off
	form.Power.Values = []string{"ON", "OFF"}
	if app.AirConSettings.Button != natureremo.ButtonPowerOn {
		form.Power.Current = 1
	}

	// modes
	var (
		i int
	)
	form.Mode.Values = make([]string, len(app.AirCon.Range.Modes))
	for m := range app.AirCon.Range.Modes {
		if app.AirConSettings.OperationMode == m {
			form.Mode.Current = i
		}
		form.Mode.Values[i] = string(m)
		i++
	}

	// temps
	opeMode := form.Mode.Values[form.Mode.Current]
	modeInfo := app.AirCon.Range.Modes[natureremo.OperationMode(opeMode)]
	temp := app.AirConSettings.Temperature
	form.Temp.Values = make([]string, len(modeInfo.Temperature))
	for i, value := range modeInfo.Temperature {
		if value == temp {
			form.Temp.Current = i
		}
		form.Temp.Values[i] = value
	}

	// volume
	vol := string(app.AirConSettings.AirVolume)
	form.Volume.Values = make([]string, len(modeInfo.AirVolume))
	for i, v := range modeInfo.AirVolume {
		if vol == string(v) {
			form.Volume.Current = i
		}
		form.Volume.Values[i] = string(v)
	}

	// direction
	dir := string(app.AirConSettings.AirDirection)
	form.Direction.Values = make([]string, len(modeInfo.AirDirection))
	for i, d := range modeInfo.AirDirection {
		if string(d) == dir {
			form.Direction.Current = i
		}
		form.Direction.Values[i] = string(d)
	}

	return form
}
