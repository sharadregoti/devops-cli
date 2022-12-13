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
	searchBox.SetDoneFunc(func(key tcell.Key) {
		// searchBox.Autocomplete().Blur()
		c <- model.Event{
			Type:         "search-complete",
			ResourceType: searchBox.GetText(),
		}
		// searchBox.SetText("")
		// searchBox.Blur()
	})
	// searchBox.auto(tcell.Color100)
	searchBox.SetDoneFunc(func(key tcell.Key) {
		c <- model.Event{
			Type:         "search-complete",
			ResourceType: searchBox.GetText(),
		}
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

// func (s *Search) Refresh(data map[string]string) {
// 	s.view.Clear()

// 	for k, v := range data {
// 		s.view.AddItem(fmt.Sprintf("%s:%s", k, v), "", ' ', nil)
// 	}
// }
