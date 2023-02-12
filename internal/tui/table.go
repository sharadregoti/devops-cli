package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sharadregoti/devops/model"
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

func (m *MainView) Refresh(data []*model.TableRow, rowNum int) {
	m.GetView().Clear()

	for r, cols := range data {
		for c, col := range cols.Data {
			// Set header
			if r < 1 {
				m.SetHeaderCell(r, c, col)
				continue
			}

			m.SetCell(r, c, col, getColor(cols.Color))
		}
	}
	m.GetView().Select(rowNum, 0)
}

func (m *MainView) SetHeaderCell(x, y int, text string) {
	m.view.SetCell(x, y,
		tview.NewTableCell(text).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft).SetExpansion(1))
}

func (m *MainView) SetCell(x, y int, text string, color tcell.Color) {
	m.view.SetCell(x, y,
		tview.NewTableCell(text).
			SetTextColor(color).
			SetAlign(tview.AlignLeft).SetExpansion(1))
}

func getColor(color string) tcell.Color {

	switch color {
	case "darkorange":
		return tcell.ColorDarkOrange
	case "gray":
		return tcell.ColorGray
	case "white":
		return tcell.ColorWhite
	case "lightskyblue":
		return tcell.ColorLightSkyBlue
	case "mediumpurple":
		return tcell.ColorMediumPurple
	case "red":
		return tcell.ColorRed
	case "yellow":
		return tcell.ColorYellow
	case "blue":
		return tcell.ColorBlue
	case "orange":
		return tcell.ColorOrange
	case "green":
		return tcell.ColorGreen
	case "aqua":
		return tcell.ColorAqua
	default:
		return tcell.ColorWhite
	}
}
