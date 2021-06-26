package ui

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var dateFormat = "2006/01/02 15:04:05"

type Devices struct {
	*tview.Table
}

func NewDevices() *Devices {
	d := &Devices{
		Table: tview.NewTable().SetSelectable(true, false),
	}

	d.SetTitle(" Devices ").SetTitleAlign(tview.AlignLeft)
	d.SetFixed(1, 0).SetBorder(true)
	d.Clear().SetBorderColor(tcell.ColorBlue)

	headers := []string{
		"Name",
		"Mac",
		"Serial",
		"Firmware",
		"Created",
		"Updated",
	}

	for i, h := range headers {
		d.SetCell(0, i, &tview.TableCell{
			Text:            h,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite,
			BackgroundColor: tcell.ColorDefault,
			Attributes:      tcell.AttrBold | tcell.AttrUnderline,
		})
	}

	devices, err := UI.cli.DeviceService.GetAll(context.Background())
	if err != nil {
		return d
	}

	for i, dev := range devices {
		cols := []string{
			dev.Name,
			dev.MacAddress,
			dev.SerialNumber,
			dev.FirmwareVersion,
			dev.CreatedAt.Format(dateFormat),
			dev.UpdatedAt.Format(dateFormat),
		}

		for j, col := range cols {
			cell := tview.NewTableCell(col).SetTextColor(tcell.ColorWhite)
			d.SetCell(i+1, j, cell)
		}
	}

	return d
}
