package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type MainView struct {
	view *tview.Table
}

func NewMainView() *MainView {
	table := tview.NewTable().SetFixed(0, 0)
	table.SetBorder(true).SetBorderAttributes(tcell.AttrDim).SetTitle("Table")
	table.SetSelectable(true, false)

	// table.Select(0, 0).SetSelectedFunc(func(row int, column int) {
	// 	if row == 0 {
	// 		return
	// 	}
	// 	c <- model.Event{
	// 		Type:     model.ReadResource,
	// 		RowIndex: row,
	// 	}
	// })

	return &MainView{
		view: table,
	}
}

func (m *MainView) GetView() *tview.Table {
	return m.view
}

func (m *MainView) SetTitle(title string) {
	m.GetView().SetTitle(cases.Title(language.AmericanEnglish).String(title))
}

func (m *MainView) Refresh(data [][]string) {
	m.GetView().Clear()

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

// package main

// import (
// 	"fmt"

// 	"github.com/gdamore/tcell"
// 	"github.com/rivo/tview"
// )

// func main() {
// 	// Create a new table.
// 	table := tview.NewTable().
// 		SetBorders(false).
// 		SetSelectable(true, false)

// 	// Add some rows to the table.
// 	table.SetCell(0, 0, &tview.TableCell{Text: "Row 1"})
// 	table.SetCell(1, 0, &tview.TableCell{Text: "Row 2"})
// 	table.SetCell(2, 0, &tview.TableCell{Text: "Row 3"})

// 	// Set a handler for when the enter key is pressed on a row.
// 	table.SetSelectedFunc(func(row, col int) {
// 		// Show a modal with the table index number.
// 		modal := tview.NewModal().
// 			SetText(fmt.Sprintf("You pressed enter on table index: %d", row)).
// 			AddButtons([]string{"OK"}).
// 			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
// 				// Close the modal when the OK button is pressed.
// 				table.SetModal(nil)
// 			})
// 		table.SetModal(modal)
// 	})

// 	// Create a new primitive which will contain the table.
// 	tablePrimitive := tview.NewFlex().
// 		SetDirection(t
