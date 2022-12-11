package views

import (
	"github.com/rivo/tview"
)

type SearchView struct {
	view *tview.Box
}

func NewSearchView() *SearchView {
	t := tview.NewBox().SetBorder(true).SetTitle("Search")

	return &SearchView{
		view: t,
	}
}

func (s *SearchView) GetView() *tview.Box {
	return s.view
}

// func (s *Search) Refresh(data map[string]string) {
// 	s.view.Clear()

// 	for k, v := range data {
// 		s.view.AddItem(fmt.Sprintf("%s:%s", k, v), "", ' ', nil)
// 	}
// }
