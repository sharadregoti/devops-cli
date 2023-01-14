package views

import (
	"github.com/rivo/tview"
	"github.com/sharadregoti/devops/shared"
)

type Actions struct {
	view           *tview.TextView
	nestingEnabled bool
}

func NewAction() *Actions {
	t := tview.NewTextView()
	t.SetBorder(true)
	t.SetTitle("Actions")

	return &Actions{
		view: t,
	}
}

func (g *Actions) EnableNesting(v bool) {
	g.nestingEnabled = v
}

func (g *Actions) RefreshActions(data shared.GenericActions) {
	tempMap := map[string]string{"ctrl-y": "read", "ctrl-r": "refresh", "ctrl-a": "toggle search bar"}
	if data.IsCreate {
		tempMap["ctrl-c"] = "create"
	}
	if data.IsUpdate {
		tempMap["ctrl-u"] = "update"
	}
	if data.IsDelete {
		tempMap["ctrl-d"] = "delete"
	}
	g.view.SetText(createKeyValuePairsWithBrackets(tempMap))
}