package ui

import (
	"log"

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
		case 'o':
			Dispatcher.Dispatch(OpenUpdateApplianceView, row)
		}
		return event
	})

	return a
}

func (a *Appliances) OpenUpdateAirConView(app *natureremo.Appliance) {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" AirCon Settings ").
		SetTitleAlign(tview.AlignLeft)

	viewData := ToUpdateAirConViewData(app)

	dispatcher := make(chan map[string]UpdateAirConFormData)

	addTemp := func() {
		form.AddDropDown("Temperature", viewData.Temp.Values, viewData.Temp.Current,
			func(opt string, idx int) {
				if idx == viewData.Temp.Current {
					return
				}
				viewData.Temp.Current = idx
				updateData := map[string]UpdateAirConFormData{app.ID: viewData}
				dispatcher <- updateData
			})
	}

	addVolume := func() {
		form.AddDropDown("Volume", viewData.Volume.Values, viewData.Volume.Current,
			func(opt string, idx int) {
				if viewData.Volume.Current == idx {
					return
				}
				viewData.Volume.Current = idx
				updateData := map[string]UpdateAirConFormData{app.ID: viewData}
				dispatcher <- updateData
			})
	}

	toggleItems := func() {
		labels := []string{
			"Temperature", "Volume",
		}
		for _, label := range labels {
			idx := form.GetFormItemIndex(label)
			if idx != -1 {
				form.RemoveFormItem(idx)
			}
		}

		switch viewData.Mode.Value() {
		case "below":
			addVolume()
		case "cool", "warm":
			addTemp()
			addVolume()
		case "dry":
			addTemp()
		}

	}

	form.AddDropDown("Power", viewData.Power.Values, viewData.Power.Current,
		func(text string, idx int) {
			if idx == viewData.Power.Current {
				return
			}
			viewData.Power.Current = idx
			updateData := map[string]UpdateAirConFormData{app.ID: viewData}
			dispatcher <- updateData
		})

	form.AddDropDown("Modes", viewData.Mode.Values, viewData.Mode.Current,
		func(opt string, idx int) {
			if viewData.Mode.Current == idx {
				return
			}
			viewData.Mode.Current = idx
			updateData := map[string]UpdateAirConFormData{app.ID: viewData}
			dispatcher <- updateData
			toggleItems()
		})

	toggleItems()

	form.AddDropDown("Direction", viewData.Direction.Values, viewData.Direction.Current,
		func(opt string, idx int) {
			if viewData.Direction.Current == idx {
				return
			}
			viewData.Direction.Current = idx
			updateData := map[string]UpdateAirConFormData{app.ID: viewData}
			dispatcher <- updateData
		})
	// update appliance with view data
	go func() {
		for data := range dispatcher {
			Dispatcher.Dispatch(UpdateAirConSettings, data)
		}
		log.Println("aircon settings dispatcher goroutine is closed")
	}()

	close := func() {
		close(dispatcher)
		UI.pages.RemovePage("form").ShowPage("main")
		UI.app.SetFocus(a)
	}

	form.AddButton("Close", func() {
		close()
	})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlN, tcell.KeyCtrlJ:
			k := tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
			UI.app.QueueEvent(k)
		case tcell.KeyCtrlP, tcell.KeyCtrlK:
			k := tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
			UI.app.QueueEvent(k)
		}

		switch event.Rune() {
		case 'j':
			k := tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
			UI.app.QueueEvent(k)
		case 'k':
			k := tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
			UI.app.QueueEvent(k)
		case 'q', 'c':
			close()
		}
		return event
	})

	UI.pages.AddPage("form", UI.Modal(form, 50, 15), true, true).SendToFront("form")
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
