package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FlashView struct {
	view *tview.TextView
}

func NewFlashView() *FlashView {
	t := tview.NewTextView()
	t.SetTextColor(tcell.ColorYellow)

	v := &FlashView{
		view: t,
	}

	return v
}

func (g *FlashView) GetView() *tview.TextView {
	return g.view
}

func (g *FlashView) SetText(text string) {
	g.view.SetText(text)
}
