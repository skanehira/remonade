package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Events struct {
	*tview.Table
}

func NewEvents() *Events {
	e := &Events{
		Table: tview.NewTable().SetSelectable(true, false),
	}

	e.SetTitle(" Events ").SetTitleAlign(tview.AlignLeft)
	e.SetFixed(1, 0).SetBorder(true)
	e.Clear().SetBorderColor(tcell.ColorGreen)

	headers := []string{
		"Type",
		"Value",
		"Created",
	}

	for i, h := range headers {
		e.SetCell(0, i, &tview.TableCell{
			Text:            h,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite,
			BackgroundColor: tcell.ColorDefault,
			Attributes:      tcell.AttrBold | tcell.AttrUnderline,
		})
	}

	devices, err := Client.DeviceService.GetAll(context.Background())
	if err != nil {
		return e
	}

	update := func() {
		row := e.GetRowCount()
		var lines [][]string
		for _, dev := range devices {
			var cols []string
			for st, v := range dev.NewestEvents {
				cols = append(cols, []string{
					string(st),
					fmt.Sprintf("%v", v.Value),
					v.CreatedAt.Local().Format(dateFormat),
				}...)
			}
			lines = append(lines, cols)
		}

		go UI.app.QueueUpdateDraw(func() {
			for i, rows := range lines {
				for j, col := range rows {
					cell := tview.NewTableCell(col).SetTextColor(tcell.ColorWhite)
					e.SetCell(i+row, j, cell)
				}
			}
		})
	}

	go func() {
		t := time.NewTicker(3 * time.Minute)
		for range t.C {
			update()
		}
	}()

	update()

	return e
}

func (e *Events) AppendRow(cols []string) {
	go UI.app.QueueUpdateDraw(func() {
		row := e.GetRowCount()
		for i, col := range cols {
			cell := tview.NewTableCell(col).SetTextColor(tcell.ColorWhite)
			e.SetCell(row, i, cell)
		}
	})
}
