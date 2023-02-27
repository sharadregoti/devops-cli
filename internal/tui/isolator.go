package tui

import (
	"bytes"
	"fmt"

	"github.com/rivo/tview"
	"github.com/sharadregoti/devops/utils/logger"
)

type IsolatorView struct {
	view *tview.TextView
	// currentKeyMap map[string]string
	currentKeyMap []string
	// requiredIs    []string
}

func NewIsolatorView() *IsolatorView {
	t := tview.NewTextView()
	// t.SetBorder(true)
	t.SetTitle("Isolator")

	v := &IsolatorView{
		view:          t,
		currentKeyMap: []string{},
		// currentKeyMap: make(map[string]string),
		// requiredIs:    []string{},
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
	// tempMap := make(map[string]string)
	// for i, v := range data {
	// 	tempMap[fmt.Sprintf("%d", i)] = v
	// }
	g.currentKeyMap = data
	// g.currentKeyMap = tempMap
	// g.requiredIs = data
	g.view.SetText(createKeyValuePairsIsolator(data))
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

	// if len(g.currentKeyMap) == 4 {
	// 	// Remove the first element
	// 	index := len(g.requiredIs)
	// 	newMap := map[string]string{}

	// 	for k, v := range g.currentKeyMap {
	// 		i, _ := strconv.Atoi(k)
	// 		if i < index {
	// 			continue
	// 		} else if i == index {
	// 			newMap[fmt.Sprintf("%d", i)] = isolatorName
	// 		} else {
	// 			newMap[fmt.Sprintf("%d", i)] = v
	// 		}
	// 	}

	// 	g.currentKeyMap = newMap
	// } else {
	// 	g.currentKeyMap[fmt.Sprintf("%d", len(g.currentKeyMap))] = isolatorName
	// }

	// Insert the element at specific index, shift remaining by 1
	g.currentKeyMap = append(g.currentKeyMap[:1], append([]string{isolatorName}, g.currentKeyMap[1:]...)...)

	logger.LogDebug("Current Key Map: %v", len(g.currentKeyMap))
	limit := 8
	if len(g.currentKeyMap) >= limit {
		// Cut off extra keys
		g.currentKeyMap = g.currentKeyMap[:limit-1]
	}

	tempMap := make(map[string]string)
	for i, v := range g.currentKeyMap {
		tempMap[fmt.Sprintf("%d", i)] = v
	}

	g.view.SetText(getNiceFormat(tempMap))
}

func createKeyValuePairsIsolator(m []string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%d: %s\n", key, value)
	}
	return b.String()
}
