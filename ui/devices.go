package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tenntenn/natureremo"
)

var dateFormat = "2006/01/02 15:04:05"

type Devices struct {
	*tview.Table
	header []string
}

func NewDevices() *Devices {
	d := &Devices{
		Table: tview.NewTable().SetSelectable(true, false),
	}

	d.SetTitle(" Devices ").SetTitleAlign(tview.AlignLeft)
	d.SetFixed(1, 0).SetBorder(true)
	d.SetBorderColor(tcell.ColorBlue)

	d.header = []string{
		"Name",
		"Mac",
		"Serial",
		"Firmware",
		"Created",
		"Updated",
	}

	return d
}

func (d *Devices) UpdateView(devices []*natureremo.Device) {
	d.Clear()
	for i, h := range d.header {
		d.SetCell(0, i, &tview.TableCell{
			Text:            h,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite,
			BackgroundColor: tcell.ColorDefault,
			Attributes:      tcell.AttrBold | tcell.AttrUnderline,
		})
	}

	for i, dev := range devices {
		cols := []string{
			dev.Name,
			dev.MacAddress,
			dev.SerialNumber,
			dev.FirmwareVersion,
			dev.CreatedAt.Local().Format(dateFormat),
			dev.UpdatedAt.Local().Format(dateFormat),
		}

		for j, col := range cols {
			cell := tview.NewTableCell(col).SetTextColor(tcell.ColorWhite)
			d.SetCell(i+1, j, cell)
		}
	}
}
