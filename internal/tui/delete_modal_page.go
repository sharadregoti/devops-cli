package tui

import "github.com/rivo/tview"

type deleteModalPage struct {
	view *tview.Modal
}

func newDeleteModalPage() *deleteModalPage {
	// TODO: Preselect no button
	modal := tview.NewModal().
		SetText("Do you want to delete the resource?").
		AddButtons([]string{"Yes", "No"})

	return &deleteModalPage{
		view: modal,
	}
}
