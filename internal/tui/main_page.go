package tui

import (
	"github.com/rivo/tview"
)

type mainPage struct {
	flexView          *tview.Flex
	isolatorBox       *IsolatorView
	generalInfoBox    *GeneralInfo
	normalActionBox   *Actions
	specificActionBox *SpecificActions
	tableBox          *MainView
	searchBox         *SearchView
}

func newMainPage() *mainPage {
	// Section1: Top level boxes that shows various info
	isolatorBox := NewIsolatorView()
	generalInfoBox := NewGeneralInfo()
	normalActionBox := NewAction()
	specificActionBox := NewSpecificAction()
	global := tview.NewFlex().
		AddItem(generalInfoBox.GetView(), 0, 1, false).
		AddItem(isolatorBox.GetView(), 0, 1, false).
		AddItem(normalActionBox.view, 0, 1, false).
		AddItem(specificActionBox.GetView(), 0, 2, false)

	// Section2: Table
	tableBox := NewTableView()

	// Section3: Search
	searchBox := NewSearchView()

	rootFlexContainer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(global, 5, 1, false).             // Top level boxes that shows various info
		AddItem(tableBox.view, 0, 1, true).       // Table
		AddItem(searchBox.GetView(), 0, 0, false) // Default disable

	return &mainPage{
		flexView:          rootFlexContainer,
		isolatorBox:       isolatorBox,
		generalInfoBox:    generalInfoBox,
		normalActionBox:   normalActionBox,
		specificActionBox: specificActionBox,
		tableBox:          tableBox,
		searchBox:         searchBox,
	}
}
