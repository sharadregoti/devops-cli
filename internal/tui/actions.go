package tui

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/sharadregoti/devops-plugin-sdk/proto"
)

type Actions struct {
	view           *tview.TextView
	nestingEnabled bool
	actions        []*proto.Action
}

func NewAction() *Actions {
	t := tview.NewTextView()
	// t.SetBorder(true)
	t.SetTitle("Actions")

	return &Actions{
		view: t,
	}
}

func (g *Actions) EnableNesting(v bool) {
	g.nestingEnabled = v
}

func (g *Actions) SetDefault(arr []*proto.Action) {
	g.actions = arr
	g.RefreshActions(arr)
}

func (g *Actions) ResetDefault() {
	g.RefreshActions(g.actions)
}

func (g *Actions) RefreshActions(arr []*proto.Action) {
	tempMap := map[string]string{}

	for _, a := range arr {
		tempMap[fmt.Sprintf("<%s>", a.KeyBinding)] = a.Name
		// tempMap[a.KeyBinding] = a.Name
	}

	// tempMap := map[string]string{"ctrl-y": "read", "ctrl-r": "refresh", "ctrl-a": "toggle search bar"}
	// if data.IsCreate {
	// 	tempMap["ctrl-c"] = "create"
	// }
	// if data.IsUpdate {
	// 	tempMap["ctrl-u"] = "update"
	// }
	// if data.IsDelete {
	// 	tempMap["ctrl-d"] = "delete"
	// }
	// g.view.SetText(createKeyValuePairsWithBrackets(tempMap))
	g.view.SetText(getNiceFormat(tempMap))

}
