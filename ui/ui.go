package ui

import (
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/skanehira/remonade/config"
	"github.com/skanehira/remonade/util"
	"github.com/tenntenn/natureremo"
)

const INTERVAL = 1

var (
	UI     *ui
	Client *natureremo.Client
)

type ui struct {
	app        *tview.Application
	pages      *tview.Pages
	primitives []tview.Primitive

	events     *Events
	appliances *Appliances
	devices    *Devices
}

func (ui *ui) Modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func (ui *ui) Message(msg string) {
	oldFocus := ui.app.GetFocus()
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			ui.pages.RemovePage("message").ShowPage("main")
			ui.app.SetFocus(oldFocus)
		})
	ui.pages.AddAndSwitchToPage("message", ui.Modal(modal, 80, 29), true).ShowPage("main")
}

func (ui *ui) Confirm(msg, doLabel string, doFunc func() error) {
	oldFocus := ui.app.GetFocus()
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{doLabel, "Cancel"}).
		SetDoneFunc(func(_ int, buttonLabel string) {
			ui.pages.RemovePage("modal").ShowPage("main")
			ui.app.SetFocus(oldFocus)
			if buttonLabel == doLabel {
				if err := doFunc(); err != nil {
					ui.Message(err.Error())
					ui.app.SetFocus(oldFocus)
				}
			}
		})
	ui.pages.AddAndSwitchToPage("modal", ui.Modal(modal, 80, 29), true).ShowPage("main")
}

func (ui *ui) next() {
	c := ui.app.GetFocus()

	for i, p := range ui.primitives {
		if c == p {
			idx := (i + 1) % len(ui.primitives)
			ui.app.SetFocus(ui.primitives[idx])
			break
		}
	}
}

func (ui *ui) prev() {
	c := ui.app.GetFocus()

	for i, p := range ui.primitives {
		if c == p {
			var idx int
			if i == 0 {
				idx = len(ui.primitives) - 1
			} else {
				idx = (i - 1) % len(ui.primitives)
			}
			ui.app.SetFocus(ui.primitives[idx])
			break
		}
	}

}

func Start() {
	Client = natureremo.NewClient(config.Config.Token)
	endpoint := os.Getenv("NATURE_REMO_ENDPOINT")
	if endpoint != "" {
		Client.BaseURL = endpoint
	}

	UI = &ui{
		app: tview.NewApplication(),
	}

	// for readability
	row, col, rowSpan, colSpan := 0, 0, 0, 0

	events := NewEvents()
	header := NewHeader()
	devices := NewDevices()
	apps := NewAppliances()

	UI.primitives = []tview.Primitive{
		devices,
		apps,
		events,
	}

	UI.events = events
	UI.devices = devices
	UI.appliances = apps

	// nolint gomnd
	grid := tview.NewGrid().SetRows(1, 0, 0).SetColumns(0, 0, 0).
		AddItem(header, row, col, rowSpan+1, colSpan+3, 0, 0, true).
		AddItem(devices, row+1, col, rowSpan+1, colSpan+2, 0, 0, true).
		AddItem(apps, row+2, col, rowSpan+1, colSpan+2, 0, 0, true).
		AddItem(events, row+1, col+2, rowSpan+2, colSpan+1, 0, 0, true)

	UI.pages = tview.NewPages().
		AddAndSwitchToPage("main", grid, true)

	UI.app.SetRoot(UI.pages, true)
	UI.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlN:
			UI.next()
		case tcell.KeyCtrlP:
			UI.prev()
		}
		return event
	})

	UI.app.SetFocus(apps)

	Dispatcher.Dispatch(GetAppliances, nil)
	Dispatcher.Dispatch(GetDevices, nil)

	go func() {
		t := time.NewTicker(INTERVAL * time.Hour)
		for range t.C {
			Dispatcher.Dispatch(GetDevices, nil)
		}
	}()

	if err := UI.app.Run(); err != nil {
		UI.app.Stop()
		util.ExitError(err)
	}
}
