package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/skanehira/remonade/config"
	"github.com/skanehira/remonade/util"
	"github.com/tenntenn/natureremo"
)

var UI *ui

type ui struct {
	cli        *natureremo.Client
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

func (ui *ui) Message(msg string, focusFunc func()) {
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			ui.pages.RemovePage("message").ShowPage("main")
			focusFunc()
		})
	ui.pages.AddAndSwitchToPage("message", ui.Modal(modal, 80, 29), true).ShowPage("main")
}

func (ui *ui) Confirm(msg, doLabel string, doFunc func() error, focusFunc func()) {
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{doLabel, "Cancel"}).
		SetDoneFunc(func(_ int, buttonLabel string) {
			ui.pages.RemovePage("modal").ShowPage("main")
			focusFunc()
			if buttonLabel == doLabel {
				if err := doFunc(); err != nil {
					ui.Message(err.Error(), func() {
						focusFunc()
					})
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

func (ui *ui) Start() {
	// for readability
	row, col, rowSpan, colSpan := 0, 0, 0, 0

	events := NewEvents()
	user := NewHeader()
	devices := NewDevices()
	apps := NewAppliances()

	ui.primitives = []tview.Primitive{
		devices,
		apps,
		events,
	}

	ui.events = events
	ui.devices = devices
	ui.appliances = apps

	grid := tview.NewGrid().SetRows(1, 0, 0).SetColumns(0, 0, 0).
		AddItem(user, row, col, rowSpan+1, colSpan+3, 0, 0, true).
		AddItem(devices, row+1, col, rowSpan+1, colSpan+2, 0, 0, true).
		AddItem(apps, row+2, col, rowSpan+1, colSpan+2, 0, 0, true).
		AddItem(events, row+1, col+2, rowSpan+2, colSpan+1, 0, 0, true)

	ui.pages = tview.NewPages().
		AddAndSwitchToPage("main", grid, true)

	ui.app.SetRoot(ui.pages, true)
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlN:
			ui.next()
		case tcell.KeyCtrlP:
			ui.prev()
		}
		return event
	})

	ui.app.SetFocus(devices)

	if err := ui.app.Run(); err != nil {
		ui.app.Stop()
		util.ExitError(err)
	}
}

func NewUI() {
	ui := &ui{
		cli: natureremo.NewClient(config.Config.Token),
		app: tview.NewApplication(),
	}

	UI = ui
}

func Run() {
	NewUI()
	UI.Start()
}
