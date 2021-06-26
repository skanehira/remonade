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
	a.Clear().SetBorderColor(tcell.ColorYellow)

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

	apps, err := UI.cli.ApplianceService.GetAll(context.Background())
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

func (a *Appliances) GetSelected() *natureremo.Appliance {
	row, _ := a.GetSelection()
	idx := row - 1
	if len(a.apps) <= idx {
		return nil
	}

	return a.apps[idx]
}

func (a *Appliances) Power(on bool) {
	app := a.GetSelected()
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
		cols = []string{"Aircon", "Power-ON", time.Now().Format(dateFormat)}
		settings := &natureremo.AirConSettings{
			Button: natureremo.ButtonPowerOn,
		}
		if !on {
			cols[1] = "Power-OFF"
			settings.Button = natureremo.ButtonPowerOff
		}
		err = UI.cli.ApplianceService.
			UpdateAirConSettings(context.Background(), app, settings)
	case natureremo.ApplianceTypeLight:
		cols = []string{"Light", "Power-ON", time.Now().Format(dateFormat)}
		btn := "on"
		if !on {
			cols[1] = "Power-OFF"
			btn = "off"
		}
		_, err = UI.cli.ApplianceService.SendLightSignal(context.Background(), app, btn)
	case natureremo.ApplianceTypeTV:
		cols = []string{"TV", "Power-ON", time.Now().Format(dateFormat)}
		btn := "power"
		if !on {
			cols[1] = "Power-OFF"
		}
		_, err = UI.cli.ApplianceService.SendTVSignal(context.Background(), app, btn)
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
