package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MainView struct {
	view *tview.Table
}

func NewMainView() *MainView {
	table := tview.NewTable().SetFixed(0, 0)
	table.SetBorder(true).SetBorderAttributes(tcell.AttrDim).SetTitle("Table")
	table.SetSelectable(true, false)
	table.SetFocusFunc(func() {
		table.Select(1, 1)
	})

	return &MainView{
		view: table,
	}
}

func (m *MainView) GetView() *tview.Table {
	return m.view
}

func (m *MainView) Refresh(data [][]string) {
	m.GetView().Clear()
	// m.GetView().getre

	for r, cols := range data {
		for c, col := range cols {
			// Set header
			if r < 1 {
				m.SetHeaderCell(r, c, col)
				continue
			}

			m.SetCell(r, c, col)
		}
	}
}

func (m *MainView) SetHeaderCell(x, y int, text string) {
	m.view.SetCell(x, y,
		tview.NewTableCell(text).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignLeft))
}

func (m *MainView) SetCell(x, y int, text string) {
	m.view.SetCell(x, y,
		tview.NewTableCell(text).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft))
}
