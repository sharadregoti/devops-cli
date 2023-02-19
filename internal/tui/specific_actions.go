package tui

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/sharadregoti/devops-plugin-sdk/proto"
	"github.com/sharadregoti/devops/utils/logger"
)

type SpecificActions struct {
	view    *tview.TextView
	actions []*proto.Action
}

func NewSpecificAction() *SpecificActions {
	t := tview.NewTextView()
	// t.SetBorder(true)
	t.SetTitle("Specific Actions")

	return &SpecificActions{
		view: t,
	}
}

func (g *SpecificActions) GetView() *tview.TextView {
	return g.view
}

func (g *SpecificActions) RefreshActions(data []*proto.Action) {
	temp := map[string]string{}
	for _, sa := range data {
		temp[fmt.Sprintf("<%s>", sa.KeyBinding)] = sa.Name
	}
	// g.view.SetText(createKeyValuePairsWithBrackets(temp))
	// TODO: This can be expensive
	logger.LogDebug("Nice format")
	logger.LogDebug("%s", getNiceFormat(temp))
	g.view.SetText(getNiceFormat(temp))
	g.actions = data
}
