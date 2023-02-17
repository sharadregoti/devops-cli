package tui

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/rivo/tview"
)

type IsolatorView struct {
	view          *tview.TextView
	currentKeyMap map[string]string
	requiredIs    []string
}

func NewIsolatorView() *IsolatorView {
	t := tview.NewTextView()
	// t.SetBorder(true)
	t.SetTitle("Isolator")

	v := &IsolatorView{
		view:          t,
		currentKeyMap: make(map[string]string),
		requiredIs:    []string{},
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

func (g *IsolatorView) SetDefault(data []string) {
	g.view.Clear()
	// convert data to map
	tempMap := make(map[string]string)
	for i, v := range data {
		tempMap[fmt.Sprintf("%d", i)] = v
	}
	g.currentKeyMap = tempMap
	g.view.SetText(getNiceFormat(tempMap))
}

func (g *IsolatorView) AddAndRefreshView(isolatorName string) {
	if isolatorName == "" {
		return
	}

	// Don't add key if already exists
	for _, v := range g.currentKeyMap {
		if v == isolatorName {
			return
		}
	}

	if len(g.currentKeyMap) == 8 {
		// Remove the first element
		index := len(g.requiredIs)
		newMap := map[string]string{}

		for k, v := range g.currentKeyMap {
			i, _ := strconv.Atoi(k)
			if i < index {
				continue
			}
			newMap[fmt.Sprintf("%d", i+1)] = v
		}
		newMap[fmt.Sprintf("%d", index)] = isolatorName
	} else {
		g.currentKeyMap[fmt.Sprintf("%d", len(g.currentKeyMap))] = isolatorName
	}

	// Insert the element at specific index, shift remaining by 1
	// g.currentKeyMap = append(g.currentKeyMap[:1], append([]string{isolatorName}, g.currentKeyMap[1:]...)...)

	// limit := 5
	// if len(g.currentKeyMap) > limit {
	// 	// Cut off extra keys
	// 	g.currentKeyMap = g.currentKeyMap[:limit]
	// }

	g.view.SetText(getNiceFormat(g.currentKeyMap))
}

func createKeyValuePairsIsolator(m []string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%d: %s\n", key, value)
	}
	return b.String()
}
