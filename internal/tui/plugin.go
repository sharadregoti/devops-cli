package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

type PluginView struct {
	view *tview.TextView
	data []string
}

func NewPluginView() *PluginView {
	t := tview.NewTextView()
	t.SetBorder(true)
	t.SetTitle("Plugin")

	return &PluginView{
		view: t,
	}
}

func (g *PluginView) GetView() *tview.TextView {
	return g.view
}

func (g *PluginView) Refresh(plugins []string) {
	g.data = plugins
	pluginMap := map[string]string{}
	for i, p := range plugins {
		pluginMap[fmt.Sprintf("alt-%d", i)] = p
	}
	g.view.Clear()
	g.view.SetText(createKeyValuePairsWithBrackets(pluginMap))
}
