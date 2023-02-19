package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/utils/logger"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type MainView struct {
	view *tview.Table
}

func NewTableView() *MainView {
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

func (m *MainView) SetTitle(title string) {
	m.view.SetTitle(cases.Title(language.AmericanEnglish).String(title))
}

func (m *MainView) Refresh(data []*model.TableRow, rowNum int) {
	m.view.Clear()

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
	m.view.Select(rowNum, 0)
}

func (m *MainView) GetRowNum(id string) int {
	for r := 0; r < m.view.GetRowCount(); r++ {
		// 1 is always the id column
		if m.view.GetCell(r, 1).Text == id {
			return r
		}
	}
	return -1
}

func (m *MainView) SetHeader(d *model.TableRow) {
	for i, colValue := range d.Data {
		m.SetCell(0, i, colValue, getColor("white"))
	}
}

func (m *MainView) AddRow(d []*model.TableRow) {
	currentRowCount := m.view.GetRowCount()
	for rowNumber, row := range d {
		for i, colValue := range row.Data {
			m.SetCell(currentRowCount+rowNumber, i, colValue, getColor(row.Color))
		}
	}
}

func (m *MainView) UpdateRow(d []*model.TableRow) {
	for _, rowInfo := range d {
		row := m.GetRowNum(rowInfo.Data[1])
		if row <= 0 {
			logger.LogDebug("Update row, row not found for id: (%s)", rowInfo.Data[1])
			return
		}
		logger.LogDebug("Updating row number (%d) with id (%s)", row, rowInfo.Data[1])
		for i, colValue := range rowInfo.Data {
			m.SetCell(row, i, colValue, getColor(rowInfo.Color))
		}
	}
}

func (m *MainView) DeleteRow(d *model.TableRow) {
	row := m.GetRowNum(d.Data[1])
	if row <= 0 {
		logger.LogDebug("Delete row, row not found for id: (%s)", d.Data[1])
		return
	}
	m.view.RemoveRow(row)
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
