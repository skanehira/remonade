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

type AppliancePowerOnOff struct {
	Power natureremo.Button
	Row   int
}

func NewAppliances() *Appliances {
	a := &Appliances{
		Table: tview.NewTable().SetSelectable(true, false),
	}
	a.SetTitle(" Appliances ").SetTitleAlign(tview.AlignLeft)
	a.SetFixed(1, 0).SetBorder(true)
	a.SetBorderColor(tcell.ColorYellow).SetBackgroundColor(tcell.ColorDefault)

	a.header = []string{
		"Device",
		"State",
		"NickName",
		"Type",
		"Model",
		"Manufacturer",
		"Country",
	}

	a.SetSelectionChangedFunc(func(row, _ int) {
		ctx := Context{
			Data: row - 1,
		}
		Dispatcher.Dispatch(ActionUpdateSelectIdx, ctx)
	})

	a.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row := a.GetSelect()
		switch event.Rune() {
		case 'u':
			if row == -1 {
				UI.Message("there is not exists any appliances")
				return event
			}
			data := AppliancePowerOnOff{
				Power: natureremo.ButtonPowerOn,
				Row:   row,
			}
			ctx := Context{
				Event: Event{
					Type:  "Appliance",
					Value: "Power On",
				},
				Data: data,
			}
			Dispatcher.Dispatch(ActionAppliancesPower, ctx)
		case 'd':
			if row == -1 {
				UI.Message("there is not exists any appliances")
				return event
			}
			data := AppliancePowerOnOff{
				Power: natureremo.ButtonPowerOff,
				Row:   row,
			}

			ctx := Context{
				Event: Event{
					Type:  "Appliance",
					Value: "Power Off",
				},
				Data: data,
			}
			Dispatcher.Dispatch(ActionAppliancesPower, ctx)
		case 'o':
			if row == -1 {
				UI.Message("there is not exists any appliances")
				return event
			}
			ctx := Context{
				Data: row,
			}
			Dispatcher.Dispatch(ActionOpenUpdateApplianceView, ctx)
		}
		return event
	})

	return a
}

func (a *Appliances) OpenUpdateAirConView(app *natureremo.Appliance) {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" AirCon Settings ").
		SetTitleAlign(tview.AlignLeft).SetBackgroundColor(tcell.ColorDefault)

	viewData := ToUpdateAirConViewData(app)

	row := a.GetSelect()
	if row == -1 {
		UI.Message("there is not any aircon settings")
		return
	}
	dispatcher := make(chan Context)

	addTemp := func() {
		form.AddDropDown("Temperature", viewData.Temp.Values, viewData.Temp.Current,
			func(_ string, idx int) {
				if idx == viewData.Temp.Current {
					return
				}
				viewData.Temp.Current = idx
				ctx := Context{
					Event: Event{
						Type:  "AC Temperature",
						Value: viewData.Temp.Value(),
					},
					Data: viewData,
				}
				dispatcher <- ctx
			})
	}

	addVolume := func() {
		form.AddDropDown("Volume", viewData.Volume.Values, viewData.Volume.Current,
			func(_ string, idx int) {
				if viewData.Volume.Current == idx {
					return
				}
				viewData.Volume.Current = idx
				ctx := Context{
					Event: Event{
						Type:  "AC Volume",
						Value: viewData.Volume.Value(),
					},
					Data: viewData,
				}
				dispatcher <- ctx
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
		func(_ string, idx int) {
			if idx == viewData.Power.Current {
				return
			}
			viewData.Power.Current = idx
			ctx := Context{
				Event: Event{
					Type:  "AC Power",
					Value: viewData.Power.Value(),
				},
				Data: viewData,
			}
			dispatcher <- ctx
		})

	form.AddDropDown("Modes", viewData.Mode.Values, viewData.Mode.Current,
		func(_ string, idx int) {
			if viewData.Mode.Current == idx {
				return
			}
			viewData.Mode.Current = idx
			ctx := Context{
				Event: Event{
					Type:  "AC Mode",
					Value: viewData.Mode.Value(),
				},
				Data: viewData,
			}
			dispatcher <- ctx
			toggleItems()
		})

	toggleItems()

	form.AddDropDown("Direction", viewData.Direction.Values, viewData.Direction.Current,
		func(_ string, idx int) {
			if viewData.Direction.Current == idx {
				return
			}
			viewData.Direction.Current = idx
			ctx := Context{
				Event: Event{
					Type:  "AC Direction",
					Value: viewData.Direction.Value(),
				},
				Data: viewData,
			}
			dispatcher <- ctx
		})
	// update appliance with view data
	go func() {
		for data := range dispatcher {
			Dispatcher.Dispatch(ActionUpdateAirConSettings, data)
		}
		log.Println("aircon settings dispatcher goroutine is closed")
	}()

	close := func() {
		close(dispatcher)
		UI.Pages.RemovePage("form").ShowPage("main")
		UI.App.SetFocus(a)
	}

	form.AddButton("Close", func() {
		close()
	})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlN, tcell.KeyCtrlJ:
			k := tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
			UI.App.QueueEvent(k)
		case tcell.KeyCtrlP, tcell.KeyCtrlK:
			k := tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
			UI.App.QueueEvent(k)
		}

		switch event.Rune() {
		case 'j':
			k := tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
			UI.App.QueueEvent(k)
		case 'k':
			k := tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
			UI.App.QueueEvent(k)
		case 'q', 'c':
			close()
		}
		return event
	})

	UI.Pages.AddPage("form", UI.Modal(form, 50, 15), true, true).SendToFront("form")
}

type LightSetting struct {
	ID    string // appliance id
	Type  string // button or signal
	Value string // button name or signal id
}

func (a *Appliances) OpenUpdateLightView(app *natureremo.Appliance) {
	table := tview.NewTable().SetSelectable(true, false)
	table.SetTitle(" Light Settings ").SetTitleAlign(tview.AlignLeft)
	table.SetFixed(1, 0).SetBorder(true).SetBackgroundColor(tcell.ColorDefault)

	header := []string{
		"Type",
		"Name",
		"Value",
	}

	for i, h := range header {
		table.SetCell(0, i, &tview.TableCell{
			Text:            h,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite,
			BackgroundColor: tcell.ColorDefault,
			Attributes:      tcell.AttrBold | tcell.AttrUnderline,
		})
	}

	list := make([][]string, len(app.Light.Buttons)+len(app.Signals))
	for i, button := range app.Light.Buttons {
		list[i] = []string{
			"button",
			button.Name,
			button.Label,
		}
	}
	for i, signal := range app.Signals {
		list[len(app.Light.Buttons)+i] = []string{
			"signal",
			signal.Name,
			signal.ID,
		}
	}

	for i, row := range list {
		for j, col := range row {
			cell := tview.NewTableCell(col).SetTextColor(tcell.ColorWhite)
			table.SetCell(i+1, j, cell)
		}
	}

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			row, _ := table.GetSelection()
			row--
			if len(list) <= row {
				return event
			}
			selected := list[row]
			setting := LightSetting{
				ID:   app.ID,
				Type: selected[0],
			}
			if setting.Type == "button" {
				setting.Value = selected[1]
			} else {
				setting.Value = selected[2]
			}
			ctx := Context{
				Event: Event{
					Type:  "AC Light",
					Value: selected[1],
				},
				Data: setting,
			}
			Dispatcher.Dispatch(ActionUpdateLight, ctx)
		}

		switch event.Rune() {
		case 'c', 'q':
			UI.Pages.RemovePage("light").ShowPage("main")
			UI.App.SetFocus(a)
		}
		return event
	})

	UI.Pages.AddPage("light", UI.Modal(table, 60, 15), true, true).SendToFront("light")
}

type TVSetting struct {
	ID     string // appliance id
	Button string // TV button
}

func (a *Appliances) OpenUpdateTVView(app *natureremo.Appliance) {
	table := tview.NewTable().SetSelectable(true, false)
	table.SetTitle(" TV Buttons ").SetTitleAlign(tview.AlignLeft)
	table.SetFixed(1, 0).SetBorder(true).SetBackgroundColor(tcell.ColorDefault)
	pageName := "TV"

	header := []string{
		"Label",
		"Name",
	}

	for i, h := range header {
		table.SetCell(0, i, &tview.TableCell{
			Text:            h,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite,
			BackgroundColor: tcell.ColorDefault,
			Attributes:      tcell.AttrBold | tcell.AttrUnderline,
		})
	}

	list := make([][]string, len(app.TV.Buttons))
	for i, button := range app.TV.Buttons {
		list[i] = []string{
			button.Label,
			button.Name,
		}
	}

	for i, row := range list {
		for j, col := range row {
			cell := tview.NewTableCell(col).SetTextColor(tcell.ColorWhite)
			table.SetCell(i+1, j, cell)
		}
	}

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			row, _ := table.GetSelection()
			row--
			if len(list) <= row {
				return event
			}
			selected := list[row]
			setting := TVSetting{
				ID:     app.ID,
				Button: selected[1],
			}
			ctx := Context{
				Event: Event{
					Type:  "TV Button",
					Value: setting.Button,
				},
				Data: setting,
			}
			Dispatcher.Dispatch(ActionSendTVButton, ctx)
		}

		switch event.Rune() {
		case 'c', 'q':
			UI.Pages.RemovePage(pageName).ShowPage("main")
			UI.App.SetFocus(a)
		}
		return event
	})

	UI.Pages.AddPage(pageName, UI.Modal(table, 60, 15), true, true).SendToFront(pageName)
}

func (a *Appliances) OpenUpdateIRView(app *natureremo.Appliance) {
	table := tview.NewTable().SetSelectable(true, false)
	table.SetTitle(" IR Signals ").SetTitleAlign(tview.AlignLeft)
	table.SetFixed(1, 0).SetBorder(true).SetBackgroundColor(tcell.ColorDefault)
	pageName := "IR"

	header := []string{
		"Name",
		"Value",
	}

	for i, h := range header {
		table.SetCell(0, i, &tview.TableCell{
			Text:            h,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite,
			BackgroundColor: tcell.ColorDefault,
			Attributes:      tcell.AttrBold | tcell.AttrUnderline,
		})
	}

	list := make([][]string, len(app.Signals))
	for i, sig := range app.Signals {
		list[i] = []string{
			sig.Name,
			sig.ID,
		}
	}

	for i, row := range list {
		for j, col := range row {
			cell := tview.NewTableCell(col).SetTextColor(tcell.ColorWhite)
			table.SetCell(i+1, j, cell)
		}
	}

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			row, _ := table.GetSelection()
			row--
			if len(list) <= row {
				return event
			}
			selected := list[row]
			signal := selected[1]
			ctx := Context{
				Event: Event{
					Type:  "IR Signal",
					Value: selected[0],
				},
				Data: signal,
			}
			Dispatcher.Dispatch(ActionSendSignal, ctx)
		}

		switch event.Rune() {
		case 'c', 'q':
			UI.Pages.RemovePage(pageName).ShowPage("main")
			UI.App.SetFocus(a)
		}
		return event
	})

	UI.Pages.AddPage(pageName, UI.Modal(table, 60, 15), true, true).SendToFront(pageName)
}

func makeApplianceRow(app *natureremo.Appliance) []string {
	row := []string{
		app.Device.Name,
	}

	switch app.Type {
	case natureremo.ApplianceTypeAirCon:
		if app.AirConSettings.Button == "" {
			row = append(row, "ON")
		} else {
			row = append(row, "OFF")
		}
	case natureremo.ApplianceTypeLight:
		if app.Light.State.Power == "off" {
			row = append(row, "OFF")
		} else {
			row = append(row, "ON")
		}
	case natureremo.ApplianceTypeTV:
		row = append(row, string(app.TV.State.Input))
	default:
		row = append(row, "-")
	}

	name := "-"
	manufacturer := "-"
	country := "-"

	if app.Model != nil {
		name = app.Model.Name
		manufacturer = app.Model.Manufacturer
		country = app.Model.Country
	}

	row = append(row, []string{
		app.Nickname,
		string(app.Type),
		name,
		manufacturer,
		country,
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
