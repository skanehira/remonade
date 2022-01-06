package main

import (
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
	"github.com/skanehira/remonade/cmd"
	"runtime"
)

func init() {
	if runtime.GOOS == "windows" && runewidth.IsEastAsian() {
		tview.Borders.Horizontal = '-'
		tview.Borders.Vertical = '|'
		tview.Borders.TopLeft = '+'
		tview.Borders.TopRight = '+'
		tview.Borders.BottomLeft = '+'
		tview.Borders.BottomRight = '+'
		tview.Borders.LeftT = '|'
		tview.Borders.RightT = '|'
		tview.Borders.TopT = '-'
		tview.Borders.BottomT = '-'
		tview.Borders.Cross = '+'
		tview.Borders.HorizontalFocus = '='
		tview.Borders.VerticalFocus = '|'
		tview.Borders.TopLeftFocus = '+'
		tview.Borders.TopRightFocus = '+'
		tview.Borders.BottomLeftFocus = '+'
		tview.Borders.BottomRightFocus = '+'
	}
}

func main() {
	cmd.Execute()
}
