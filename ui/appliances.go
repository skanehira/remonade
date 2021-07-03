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
			if row == -1 {
				UI.Message("there is not exists any appliances")
				return event
			}
			ctx := AppliancePowerOnOff{
				Power: natureremo.ButtonPowerOn,
				Row:   row,
			}
			Dispatcher.Dispatch(ActionAppliancesPower, ctx)
		case 'd':
			if row == -1 {
				UI.Message("there is not exists any appliances")
				return event
			}
			ctx := AppliancePowerOnOff{
				Power: natureremo.ButtonPowerOff,
				Row:   row,
			}
			Dispatcher.Dispatch(ActionAppliancesPower, ctx)
		case 'o':
			if row == -1 {
				UI.Message("there is not exists any appliances")
				return event
			}
			Dispatcher.Dispatch(ActionOpenUpdateApplianceView, row)
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

	row := a.GetSelect()
	if row == -1 {
		UI.Message("there is not any aircon settings")
		return
	}
	dispatcher := make(chan map[int]UpdateAirConFormData)

	addTemp := func() {
		form.AddDropDown("Temperature", viewData.Temp.Values, viewData.Temp.Current,
			func(opt string, idx int) {
				if idx == viewData.Temp.Current {
					return
				}
				viewData.Temp.Current = idx
				updateData := map[int]UpdateAirConFormData{row: viewData}
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
				updateData := map[int]UpdateAirConFormData{row: viewData}
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
			updateData := map[int]UpdateAirConFormData{row: viewData}
			dispatcher <- updateData
		})

	form.AddDropDown("Modes", viewData.Mode.Values, viewData.Mode.Current,
		func(opt string, idx int) {
			if viewData.Mode.Current == idx {
				return
			}
			viewData.Mode.Current = idx
			updateData := map[int]UpdateAirConFormData{row: viewData}
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
			updateData := map[int]UpdateAirConFormData{row: viewData}
			dispatcher <- updateData
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

type LightSetting struct {
	ID    string // appliance id
	Type  string // button or signal
	Value string // button name or signal id
}

func (a *Appliances) OpenUpdateLightView(app *natureremo.Appliance) {
	table := tview.NewTable().SetSelectable(true, false)
	table.SetTitle(" Light Settings ").SetTitleAlign(tview.AlignLeft)
	table.SetFixed(1, 0).SetBorder(true)

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
			Dispatcher.Dispatch(ActionUpdateLight, setting)
		}

		switch event.Rune() {
		case 'c', 'q':
			UI.pages.RemovePage("light").ShowPage("main")
			UI.app.SetFocus(a)
		}
		return event
	})

	UI.pages.AddPage("light", UI.Modal(table, 60, 15), true, true).SendToFront("light")
}

type TVSetting struct {
	ID     string // appliance id
	Button string // TV button
}

func (a *Appliances) OpenUpdateTVView(app *natureremo.Appliance) {
	table := tview.NewTable().SetSelectable(true, false)
	table.SetTitle(" TV Buttons ").SetTitleAlign(tview.AlignLeft)
	table.SetFixed(1, 0).SetBorder(true)
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
			Dispatcher.Dispatch(ActionSendTVButton, setting)
		}

		switch event.Rune() {
		case 'c', 'q':
			UI.pages.RemovePage(pageName).ShowPage("main")
			UI.app.SetFocus(a)
		}
		return event
	})

	UI.pages.AddPage(pageName, UI.Modal(table, 60, 15), true, true).SendToFront(pageName)
}

func (a *Appliances) OpenUpdateIRView(app *natureremo.Appliance) {
	table := tview.NewTable().SetSelectable(true, false)
	table.SetTitle(" IR Signals ").SetTitleAlign(tview.AlignLeft)
	table.SetFixed(1, 0).SetBorder(true)
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
			Dispatcher.Dispatch(ActionSendSignal, signal)
		}

		switch event.Rune() {
		case 'c', 'q':
			UI.pages.RemovePage(pageName).ShowPage("main")
			UI.app.SetFocus(a)
		}
		return event
	})

	UI.pages.AddPage(pageName, UI.Modal(table, 60, 15), true, true).SendToFront(pageName)
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
