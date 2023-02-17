package tui

import "github.com/rivo/tview"

type formPage struct {
	view *tview.Form
}

func newFormPage() *formPage {
	genericUserInputFormBox := tview.NewForm()
	genericUserInputFormBox.SetBorder(true).SetTitle("Form")

	return &formPage{
		view: genericUserInputFormBox,
	}
}

func (m *formPage) registerEventHandler() error {
	return nil
}
