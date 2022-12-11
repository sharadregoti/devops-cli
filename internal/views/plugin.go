package views

import (
	"bytes"
	"fmt"

	"github.com/rivo/tview"
)

type PluginView struct {
	view *tview.TextView
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

func (g *PluginView) Refresh(data map[string]string) {
	g.view.Clear()
	g.view.SetText(createKeyValuePairsPlugin(data))
}

func createKeyValuePairsPlugin(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s: \"%s\"\n", key, value)
	}
	return b.String()
}
