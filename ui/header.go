package ui

import (
	"context"
	"fmt"

	"github.com/rivo/tview"
)

type Header struct {
	*tview.TextView
}

func NewHeader() *Header {
	h := &Header{
		TextView: tview.NewTextView(),
	}
	h.SetDynamicColors(true)

	me, err := Client.UserService.Me(context.Background())
	if err != nil {
		h.SetText("unknown")
		return h
	}

	text := fmt.Sprintf(" [yellow::bl]ID: [-:-:-]%s [yellow::bl]Name: [-:-:-]%s", me.ID, me.Nickname)
	h.SetText(text)
	return h
}
