package views

import (
	"github.com/rivo/tview"
	"github.com/sharadregoti/devops/shared"
)

type SpecificActions struct {
	view    *tview.TextView
	actions []shared.SpecificAction
}

func NewSpecificAction() *SpecificActions {
	t := tview.NewTextView()
	t.SetBorder(true)
	t.SetTitle("Specific Actions")

	return &SpecificActions{
		view: t,
	}
}

func (g *SpecificActions) GetView() *tview.TextView {
	return g.view
}

func (g *SpecificActions) RefreshActions(data []shared.SpecificAction) {
	temp := map[string]string{}
	for _, sa := range data {
		temp[sa.KeyBinding] = sa.Name
	}
	g.view.SetText(createKeyValuePairsWithBrackets(temp))
	g.actions = data
}
