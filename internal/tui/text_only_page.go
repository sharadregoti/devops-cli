package tui

import (
	"github.com/rivo/tview"
)

type textOnlyPage struct {
	view *tview.TextView
}

func newTextOnlyPage() *textOnlyPage {
	textOnlyBox := tview.NewTextView()
	textOnlyBox.SetDynamicColors(true)

	return &textOnlyPage{
		view: textOnlyBox,
	}
}
