package tui

import (
	"bytes"
	"fmt"

	"github.com/rivo/tview"
)

type GeneralInfo struct {
	view *tview.TextView
}

func NewGeneralInfo() *GeneralInfo {
	t := tview.NewTextView()
	// t.SetBorder(true)
	t.SetTitle("General Info")

	return &GeneralInfo{
		view: t,
	}
}

func (g *GeneralInfo) GetView() *tview.TextView {
	return g.view
}

func (g *GeneralInfo) Refresh(data map[string]string) {
	g.view.Clear()
	// g.view.SetText(createKeyValuePairsWithoutBrackets(data))
	g.view.SetText(getNiceFormat(data))
}

func createKeyValuePairsWithoutBrackets(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s: %s\n", key, value)
	}
	return b.String()
}

func createKeyValuePairsWithBrackets(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "<%s> %s\n", key, value)
	}
	return b.String()
}
