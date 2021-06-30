package ui

import (
	"context"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tenntenn/natureremo"
)

type Appliances struct {
	*tview.Table
	apps []*natureremo.Appliance
	rows [][]string
}

func NewAppliances() *Appliances {
	a := &Appliances{
		Table: tview.NewTable().SetSelectable(true, false),
	}
	a.SetTitle(" Appliances ").SetTitleAlign(tview.AlignLeft)
	a.SetFixed(1, 0).SetBorder(true)
	a.SetBorderColor(tcell.ColorYellow)

	headers := []string{
		"State",
		"NickName",
		"Type",
		"Model",
		"Manufacturer",
		"Country",
	}

	for i, h := range headers {
		a.SetCell(0, i, &tview.TableCell{
			Text:            h,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite,
			BackgroundColor: tcell.ColorDefault,
			Attributes:      tcell.AttrBold | tcell.AttrUnderline,
		})
	}

	apps, err := Client.ApplianceService.GetAll(context.Background())
	if err != nil {
		return a
	}
	a.apps = apps

	for i, app := range apps {
		row := a.makeRow(app)
		for j, col := range row {
			cell := tview.NewTableCell(col).SetTextColor(tcell.ColorWhite)
			a.SetCell(i+1, j, cell)
		}
		a.rows = append(a.rows, row)
	}

	a.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'u':
			a.Power(true)
		case 'd':
			a.Power(false)
		case 'e':
			a.Update()
		}
		return event
	})

	return a
}

func (a *Appliances) makeRow(app *natureremo.Appliance) []string {
	var row []string

	if app.Type == natureremo.ApplianceTypeAirCon {
		if app.AirConSettings.Button == "" {
			row = []string{"ON"}
		} else {
			row = []string{"OFF"}
		}
	} else if app.Type == natureremo.ApplianceTypeLight {
		if app.Light.State.Power == "off" {
			row = []string{"OFF"}
		} else {
			row = []string{"ON"}
		}
	} else if app.Type == natureremo.ApplianceTypeTV {
		row = []string{string(app.TV.State.Input)}
	} else {
		row = []string{"-"}
	}

	row = append(row, []string{
		app.Nickname,
		string(app.Type),
		app.Model.Name,
		app.Model.Manufacturer,
		app.Model.Country,
	}...)
	return row
}

func (a *Appliances) SelectedApp() *natureremo.Appliance {
	row, _ := a.GetSelection()
	idx := row - 1
	if len(a.apps) <= idx {
		return nil
	}

	return a.apps[idx]
}

func (a *Appliances) Power(on bool) {
	app := a.SelectedApp()
	if app == nil {
		return
	}

	var (
		err  error
		cols []string
	)

	state := "ON"
	if !on {
		state = "OFF"
	}

	switch app.Type {
	case natureremo.ApplianceTypeAirCon:
		cols = []string{"Aircon", "Power-ON", time.Now().Local().Format(dateFormat)}
		settings := &natureremo.AirConSettings{
			Button: natureremo.ButtonPowerOn,
		}
		if !on {
			cols[1] = "Power-OFF"
			settings.Button = natureremo.ButtonPowerOff
		}
		err = Client.ApplianceService.
			UpdateAirConSettings(context.Background(), app, settings)
	case natureremo.ApplianceTypeLight:
		cols = []string{"Light", "Power-ON", time.Now().Local().Format(dateFormat)}
		btn := "on"
		if !on {
			cols[1] = "Power-OFF"
			btn = "off"
		}
		_, err = Client.ApplianceService.SendLightSignal(context.Background(), app, btn)
	case natureremo.ApplianceTypeTV:
		cols = []string{"TV", "Power-ON", time.Now().Local().Format(dateFormat)}
		btn := "power"
		if !on {
			cols[1] = "Power-OFF"
		}
		_, err = Client.ApplianceService.SendTVSignal(context.Background(), app, btn)
		state = ""
	default:
		return
	}

	if err != nil {
		UI.Message(err.Error(), func() {
			UI.app.SetFocus(a)
		})
		return
	}

	// update table columns
	idx, _ := a.GetSelection()
	if state != "" {
		a.SetCellSimple(idx, 0, state)
	}

	// insert row to event panel
	UI.events.AppendRow(cols)
}

func (a *Appliances) Update() {
	app := a.SelectedApp()
	if app == nil {
		return
	}

	switch app.Type {
	case natureremo.ApplianceTypeAirCon:
		a.UpdateAirCon(app)
	case natureremo.ApplianceTypeLight:
	}
}

func (a *Appliances) UpdateLight(app *natureremo.Appliance) {

}

func (a *Appliances) UpdateAirCon(app *natureremo.Appliance) {
	var (
		currentPower int
	)

	if app.AirConSettings.Button != "" {
		currentPower = 1
	}

	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" AirCon Settings ")
	form.SetTitleAlign(tview.AlignLeft)

	form.AddDropDown("Power", []string{
		"ON",
		"OFF",
	}, currentPower, func(text string, idx int) {
		if idx == currentPower {
			return
		}
		currentPower = idx

		// TODO
	})

	var modes []string
	var currentMode int
	for m := range app.AirCon.Range.Modes {
		modes = append(modes, string(m))
	}
	for i, m := range modes {
		if string(app.AirConSettings.OperationMode) == m {
			currentMode = i
			break
		}
	}

	form.AddDropDown("Modes", modes, currentMode, func(opt string, idx int) {
		if currentMode == idx {
			return
		}
		idx = currentMode
		// TODO
	})

	opeMode := modes[currentMode]
	modeInfo := app.AirCon.Range.Modes[natureremo.OperationMode(opeMode)]
	if opeMode != "below" {
		var currentTemp int
		temps := modeInfo.Temperature
		temp := app.AirConSettings.Temperature

		for i, m := range temps {
			if m == temp {
				currentTemp = i
				break
			}
		}

		form.AddDropDown("Temperature", temps, currentTemp, func(opt string, idx int) {
			if idx == currentTemp {
				return
			}
			currentTemp = idx
			// TODO
		})
	}

	if opeMode != "dry" {
		vol := string(app.AirConSettings.AirVolume)
		var vols []string
		var currentVol int
		for i, v := range modeInfo.AirVolume {
			if vol == string(v) {
				currentVol = i
			}
			vols = append(vols, string(v))
		}

		form.AddDropDown("Volume", vols, currentVol, func(opt string, idx int) {
			if currentVol == idx {
				return
			}

			idx = currentVol
			// TODO
		})
	}

	dir := string(app.AirConSettings.AirDirection)
	var dirs []string
	var currentDir int
	for i, d := range modeInfo.AirDirection {
		if string(d) == dir {
			currentDir = i
		}
		dirs = append(dirs, string(d))
	}

	form.AddDropDown("Direction", dirs, currentDir, func(opt string, idx int) {
		if currentDir == idx {
			return
		}
		idx = currentDir
		// TODO
	})

	close := func() {
		UI.pages.RemovePage("form").ShowPage("main")
		UI.app.SetFocus(a)
	}

	form.AddButton("Close", func() {
		close()
	})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlN:
			k := tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
			UI.app.QueueEvent(k)
		case tcell.KeyCtrlP:
			k := tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
			UI.app.QueueEvent(k)
		}

		switch event.Rune() {
		case 'q', 'c':
			close()
		}
		return event
	})

	UI.pages.AddAndSwitchToPage("form", UI.Modal(form, 50, 15), true).ShowPage("main")
}
