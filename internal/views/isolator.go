package views

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/rivo/tview"
	"github.com/sharadregoti/devops/model"
)

type IsolatorView struct {
	view          *tview.TextView
	currentKeyMap []string
	logger        hclog.Logger
}

func NewIsolatorView(logger hclog.Logger, c chan model.Event) *IsolatorView {
	t := tview.NewTextView()
	t.SetBorder(true)
	t.SetTitle("Isolator")

	v := &IsolatorView{
		view:          t,
		logger:        logger.Named("isolator"),
		currentKeyMap: make([]string, 0),
	}

	// t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	// 	if event.Key() == tcell.KeyRune {
	// 		for i := range v.currentKeyMap {
	// 			fmt.Println("Here")
	// 			numToRune := fmt.Sprintf("%d", i)[0]
	// 			if event.Rune() == rune(numToRune) {
	// 				c <- model.Event{
	// 					Type:         model.IsolatorChanged,
	// 					IsolatorName: v.currentKeyMap[i],
	// 				}
	// 			}
	// 		}
	// 	}
	// 	return event
	// })

	return v
}

func (g *IsolatorView) GetView() *tview.TextView {
	return g.view
}

func (g *IsolatorView) SetTitle(data string) {
	g.view.SetTitle(data)
}

func (g *IsolatorView) SetDefault(data string) {
	g.view.Clear()
	temp := []string{data}
	g.currentKeyMap = temp
	g.view.SetText(createKeyValuePairsIsolator(temp))
}

func (g *IsolatorView) AddAndRefreshView(isolatorName string) {
	if isolatorName == "" {
		return
	}

	// Don't add key if already exists
	for _, v := range g.currentKeyMap {
		if v == isolatorName {
			g.logger.Debug(fmt.Sprintf("trying to add the same isolator key <%s>, ignoring it", isolatorName))
			return
		}
	}

	// Insert the element at specific index, shift remaining by 1
	g.currentKeyMap = append(g.currentKeyMap[:1], append([]string{isolatorName}, g.currentKeyMap[1:]...)...)

	limit := 3
	if len(g.currentKeyMap) > limit {
		// Cut off extra keys
		g.currentKeyMap = g.currentKeyMap[:limit]
	}

	g.view.SetText(createKeyValuePairsIsolator(g.currentKeyMap))
}

func createKeyValuePairsIsolator(m []string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%d: %s\n", key, value)
	}
	return b.String()
}
