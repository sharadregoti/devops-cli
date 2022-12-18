package views

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sharadregoti/devops/model"
)

type SearchView struct {
	view *tview.InputField
}

func NewSearchView(c chan model.Event) *SearchView {
	searchBox := tview.NewInputField()
	searchBox.SetFieldBackgroundColor(tcell.Color100)
	searchBox.Autocomplete().SetDoneFunc(func(key tcell.Key) {
		c <- model.Event{
			Type:         model.ResourceTypeChanged,
			ResourceType: searchBox.GetText(),
		}
		searchBox.SetText("")
	})

	// searchBox.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

	// 	// Check if the user pressed Shift + ":".
	// 	if event.Key() == ':' && event.Modifiers() == tcell.ModShift {
	// 		// Toggle the visibility of the root primitive.
	// 		// searchBox.
	// 		// searchBox.Setvi(!searchBox.IsVisible())
	// 		println("Key pressed")
	// 		// Return without handling the event.
	// 		return event
	// 	}

	// 	// Check if the user pressed the enter key.
	// 	if event.Key() == tcell.KeyEnter {
	// 		// Get the search query.
	// 		query := searchBox.GetText()

	// 		items := []string{"apple", "banana", "orange", "grape", "strawberry"}

	// 		// Search the array of items for matches.
	// 		matches := []string{}
	// 		for _, item := range items {
	// 			if strings.Contains(strings.ToLower(item), strings.ToLower(query)) {
	// 				matches = append(matches, item)
	// 			}
	// 		}

	// 		// Print the results.
	// 		if len(matches) > 0 {
	// 			println("Matches:")
	// 			for _, match := range matches {
	// 				println(match)
	// 			}
	// 		} else {
	// 			println("No matches found.")
	// 		}
	// 	}

	// 	return event
	// })

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
		for _, word := range arr {
			if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) {
				entries = append(entries, word)
			}
		}
		if len(entries) <= 1 {
			entries = nil
		}
		return
	})
}
