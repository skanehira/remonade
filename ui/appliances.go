package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tenntenn/natureremo"
)

type Appliances struct {
	*tview.Table
	header []string
}

func (a *Appliances) UpdateView(apps []*natureremo.Appliance) {
	a.Clear()
	for i, h := range a.header {
		a.SetCell(0, i, &tview.TableCell{
			Text:            h,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite,
			BackgroundColor: tcell.ColorDefault,
			Attributes:      tcell.AttrBold | tcell.AttrUnderline,
		})
	}

	for i, row := range makeApplianceRows(apps) {
		for j, col := range row {
			cell := tview.NewTableCell(col).SetTextColor(tcell.ColorWhite)
			a.SetCell(i+1, j, cell)
		}
	}
}

func (a *Appliances) GetSelect() int {
	row, _ := a.GetSelection()
	row--
	return row
}

func NewAppliances() *Appliances {
	a := &Appliances{
		Table: tview.NewTable().SetSelectable(true, false),
	}
	a.SetTitle(" Appliances ").SetTitleAlign(tview.AlignLeft)
	a.SetFixed(1, 0).SetBorder(true)
	a.SetBorderColor(tcell.ColorYellow)

	a.header = []string{
		"State",
		"NickName",
		"Type",
		"Model",
		"Manufacturer",
		"Country",
	}

	a.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row := a.GetSelect()
		switch event.Rune() {
		case 'u':
			Dispatcher.Dispatch(PowerON, row)
		case 'd':
			Dispatcher.Dispatch(PowerOFF, row)
		case 'e':
			Dispatcher.Dispatch(OpenUpdateApplianceView, row)
		}
		return event
	})

	return a
}

func (a *Appliances) OpenUpdateAirConView(app *natureremo.Appliance) {
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

	// nolint prealloc
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
		// nolint ineffassign
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

			// nolint ineffassign
			idx = currentVol
			// TODO
		})
	}

	// nolint prealloc
	var dirs []string
	var currentDir int
	dir := string(app.AirConSettings.AirDirection)
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
		// nolint ineffassign
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

func makeApplianceRow(app *natureremo.Appliance) []string {
	var row []string

	switch app.Type {
	case natureremo.ApplianceTypeAirCon:
		if app.AirConSettings.Button == "" {
			row = []string{"ON"}
		} else {
			row = []string{"OFF"}
		}
	case natureremo.ApplianceTypeLight:
		if app.Light.State.Power == "off" {
			row = []string{"OFF"}
		} else {
			row = []string{"ON"}
		}
	case natureremo.ApplianceTypeTV:
		row = []string{string(app.TV.State.Input)}
	default:
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

func makeApplianceRows(apps []*natureremo.Appliance) [][]string {
	rows := make([][]string, len(apps))

	for i, app := range apps {
		row := makeApplianceRow(app)
		rows[i] = row
	}

	return rows
}
