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
	App        *tview.Application
	Pages      *tview.Pages
	Primitives []tview.Primitive

	Events     *Events
	Appliances *Appliances
	Devices    *Devices
}

func (ui *ui) Modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func (ui *ui) Message(msg string) {
	oldFocus := ui.App.GetFocus()
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"OK"}).
		SetBackgroundColor(tcell.ColorDefault).
		SetDoneFunc(func(_ int, _ string) {
			ui.Pages.RemovePage("message").ShowPage("main")
			ui.App.SetFocus(oldFocus)
		})
	go UI.App.QueueUpdateDraw(func() {
		ui.Pages.AddPage("message", ui.Modal(modal, 80, 29), true, true).SendToFront("message")
	})
}

func (ui *ui) Confirm(msg, doLabel string, doFunc func() error) {
	oldFocus := ui.App.GetFocus()
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{doLabel, "Cancel"}).
		SetBackgroundColor(tcell.ColorDefault).
		SetDoneFunc(func(_ int, buttonLabel string) {
			ui.Pages.RemovePage("modal").ShowPage("main")
			ui.App.SetFocus(oldFocus)
			if buttonLabel == doLabel {
				if err := doFunc(); err != nil {
					ui.Message(err.Error())
					ui.App.SetFocus(oldFocus)
				}
			}
		})
	go UI.App.QueueUpdateDraw(func() {
		ui.Pages.AddPage("modal", ui.Modal(modal, 80, 29), true, true).SendToFront("modal")
	})
}

func (ui *ui) NextPane() {
	c := ui.App.GetFocus()
	for i, p := range ui.Primitives {
		if c == p {
			idx := (i + 1) % len(ui.Primitives)
			ui.App.SetFocus(ui.Primitives[idx])
			break
		}
	}
}

func (ui *ui) PrevPane() {
	c := ui.App.GetFocus()

	for i, p := range ui.Primitives {
		if c == p {
			var idx int
			if i == 0 {
				idx = len(ui.Primitives) - 1
			} else {
				idx = (i - 1) % len(ui.Primitives)
			}
			ui.App.SetFocus(ui.Primitives[idx])
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
		App: tview.NewApplication(),
	}

	// for readability
	row, col, rowSpan, colSpan := 0, 0, 0, 0

	events := NewEvents()
	header := NewHeader()
	devices := NewDevices()
	apps := NewAppliances()

	UI.Primitives = []tview.Primitive{
		devices,
		apps,
		events,
	}

	UI.Events = events
	UI.Devices = devices
	UI.Appliances = apps

	// nolint gomnd
	grid := tview.NewGrid().SetRows(1, 0, 0).SetColumns(0, 0, 0).
		AddItem(header, row, col, rowSpan+1, colSpan+3, 0, 0, true).
		AddItem(devices, row+1, col, rowSpan+1, colSpan+2, 0, 0, true).
		AddItem(apps, row+2, col, rowSpan+1, colSpan+2, 0, 0, true).
		AddItem(events, row+1, col+2, rowSpan+2, colSpan+1, 0, 0, true)

	UI.Pages = tview.NewPages().
		AddAndSwitchToPage("main", grid, true)

	UI.App.SetRoot(UI.Pages, true)
	UI.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlN:
			UI.NextPane()
		case tcell.KeyCtrlP:
			UI.PrevPane()
		}
		return event
	})

	UI.App.SetFocus(apps)

	ctx := Context{
		Event: Event{
			Type:  "API",
			Value: "get appliances",
		},
	}
	Dispatcher.Dispatch(ActionGetAppliances, ctx)

	ctx = Context{
		Event: Event{
			Type:  "API",
			Value: "get devices",
		},
	}
	Dispatcher.Dispatch(ActionGetDevices, ctx)

	go func() {
		t := time.NewTicker(INTERVAL * time.Hour)
		for range t.C {
			Dispatcher.Dispatch(ActionGetDevices, ctx)
		}
	}()

	if err := UI.App.Run(); err != nil {
		UI.App.Stop()
		util.ExitError(err)
	}
}
