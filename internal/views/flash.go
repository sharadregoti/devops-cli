package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/rivo/tview"
	"github.com/sharadregoti/devops/model"
)

type FlashView struct {
	view   *tview.TextView
	logger hclog.Logger
}

func NewFlashView(logger hclog.Logger, c chan model.Event) *FlashView {
	t := tview.NewTextView()
	t.SetTextColor(tcell.ColorYellow)

	v := &FlashView{
		view:   t,
		logger: logger.Named("flash-view"),
	}

	return v
}

func (g *FlashView) GetView() *tview.TextView {
	return g.view
}

func (g *FlashView) SetText(text string) {
	g.view.SetText(text)
}
