package views

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sahilm/fuzzy"
)

type SearchView struct {
	view *tview.InputField
}

func NewSearchView() *SearchView {
	searchBox := tview.NewInputField()
	searchBox.SetFieldBackgroundColor(tcell.Color100)
	searchBox.Autocomplete().SetDoneFunc(func(key tcell.Key) {
		// c <- model.Event{
		// 	Type:         string(model.ResourceTypeChanged),
		// 	ResourceType: searchBox.GetText(),
		// }
		// searchBox.SetText("")
	})

	return &SearchView{
		view: searchBox,
	}
}

func (s *SearchView) GetView() *tview.InputField {
	return s.view
}

func (s *SearchView) SetResourceTypes(arr []string) {
	s.view.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if len(currentText) == 0 {
			return
		}

		matches := fuzzy.Find(strings.ToLower(currentText), arr)

		for _, v := range matches {
			entries = append(entries, v.Str)
		}
		// for _, word := range arr {

		// 	if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) {
		// 		entries = append(entries, word)
		// 	}
		// }
		if len(entries) == 0 {
			entries = nil
		}
		return
	})
}
