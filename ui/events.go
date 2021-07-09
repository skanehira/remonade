package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Events struct {
	*tview.Table
	header []string
}

func NewEvents() *Events {
	e := &Events{
		Table: tview.NewTable().SetSelectable(true, false),
	}

	e.SetTitle(" Events ").SetTitleAlign(tview.AlignLeft)
	e.SetFixed(1, 0).SetBorder(true)
	e.Clear().SetBorderColor(tcell.ColorGreen)

	e.header = []string{
		"Device",
		"Type",
		"Value",
		"Created",
	}

	return e
}

func (e *Events) UpdateView(events []Event) {
	e.Clear()

	for i, h := range e.header {
		e.SetCell(0, i, &tview.TableCell{
			Text:            h,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite,
			BackgroundColor: tcell.ColorDefault,
			Attributes:      tcell.AttrBold | tcell.AttrUnderline,
		})
	}

	rows := make([][]string, len(events))
	for i, e := range events {
		rows[i] = []string{
			e.Device,
			e.Type,
			e.Value,
			e.Created.Local().Format(dateFormat),
		}
	}

	for i, rows := range rows {
		for j, col := range rows {
			cell := tview.NewTableCell(col).SetTextColor(tcell.ColorWhite)
			e.SetCell(i+1, j, cell)
		}
	}
}
