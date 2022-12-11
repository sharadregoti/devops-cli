package views

import (
	"bytes"
	"fmt"

	"github.com/rivo/tview"
)

type IsolatorView struct {
	view *tview.TextView
}

func NewIsolatorView() *IsolatorView {
	t := tview.NewTextView()
	t.SetBorder(true)
	t.SetTitle("Isolator")

	return &IsolatorView{
		view: t,
	}
}

func (g *IsolatorView) GetView() *tview.TextView {
	return g.view
}

func (g *IsolatorView) Refresh(data map[string]string) {
	g.view.Clear()
	g.view.SetText(createKeyValuePairsIsolator(data))
}

func createKeyValuePairsIsolator(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s: \"%s\"\n", key, value)
	}
	return b.String()
}
