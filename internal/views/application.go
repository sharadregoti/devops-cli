package views

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/rivo/tview"
	"github.com/sharadregoti/devops/model"
)

type Application struct {
	rootView           *tview.Flex
	MainView           *MainView
	GeneralInfoView    *GeneralInfo
	SearchView         *SearchView
	IsolatorView       *IsolatorView
	PluginView         *PluginView
	eventChan          chan model.Event
	application        *tview.Application
	page               *tview.Pages
	ta                 *tview.TextView
	SpecificActionView *SpecificActions
	ActionView         *Actions
}

func NewApplication(logger hclog.Logger, c chan model.Event) *Application {

	pa := tview.NewPages()

	i := NewIsolatorView(logger, c)
	m := NewMainView()
	s := NewSearchView(c)
	g := NewGeneralInfo()
	p := NewPluginView()
	act := NewAction()
	sa := NewSpecificAction()

	global := tview.NewFlex().
		AddItem(g.GetView(), 0, 1, false).
		AddItem(i.GetView(), 0, 1, false).
		AddItem(act.view, 0, 1, false).
		AddItem(sa.GetView(), 0, 1, false).
		AddItem(p.GetView(), 0, 1, false)

	rootFlexContainer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(global, 6, 1, false).
		AddItem(m.GetView(), 0, 1, true).
		AddItem(s.GetView(), 0, 0, false) // Default disable

	modal := tview.NewModal().
		SetText("Do you want to delete the resource?").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				row, col := m.view.GetSelection()
				c <- model.Event{
					Type:         model.DeleteResource,
					ResourceName: m.view.GetCell(row, col+1).Text,
					IsolatorName: "",
					ResourceType: strings.ToLower(m.view.GetTitle()),
				}
				pa.SwitchToPage("main")
				c <- model.Event{
					Type: model.RefreshResource,
				}
			}
			if buttonLabel == "No" {
				pa.SwitchToPage("main")
			}

		})

	ta := tview.NewTextView()
	ta.SetDynamicColors(true)

	pa.AddPage("main", rootFlexContainer, true, true)
	pa.AddPage("modal", modal, true, true)
	pa.AddPage("ta", ta, true, true)
	pa.SwitchToPage("main")

	ta.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			pa.SwitchToPage("main")
		}
	})

	a := tview.NewApplication().SetRoot(pa, true)

	r := &Application{
		rootView:           rootFlexContainer,
		MainView:           m,
		GeneralInfoView:    g,
		SearchView:         s,
		IsolatorView:       i,
		PluginView:         p,
		SpecificActionView: sa,
		eventChan:          c,
		application:        a,
		page:               pa,
		ta:                 ta,
		ActionView:         act,
	}

	r.SetKeyboardShortCuts()

	return r
}

func (a *Application) SetText(text string) {
	a.ta.SetText(text)
	a.page.SwitchToPage("ta")
}

func (a *Application) SwitchToMain() {
	a.page.SwitchToPage("main")
}

func (a *Application) GetApp() *tview.Application {
	return a.application
}

func (a *Application) ViewModel() {
	a.page.SwitchToPage("modal")
}

func (a *Application) SetKeyboardShortCuts() {
	a.rootView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		switch event.Key() {
		case tcell.KeyRune:
			if event.Rune() == 'u' {
				row, col := a.MainView.view.GetSelection()
				a.eventChan <- model.Event{
					Type: model.AddIsolator,
					// +1 because "name" is always the second column that indicates the isolator name
					// TODO: In future, change it to dynamically detected isolator name from table content irrespecitve of where the column
					IsolatorName: a.MainView.view.GetCell(row, col+1).Text,
					ResourceType: strings.ToLower(a.MainView.view.GetTitle()),
				}
				return event
			}

			// Single digit numbers only
			for i := range a.IsolatorView.currentKeyMap {
				numToRune := fmt.Sprintf("%d", i)[0]
				if event.Rune() == rune(numToRune) {
					a.eventChan <- model.Event{
						Type:         model.IsolatorChanged,
						IsolatorName: a.IsolatorView.currentKeyMap[i],
					}
				}
			}

			// Specific Action Checks
			row, col := a.MainView.view.GetSelection()
			for _, action := range a.SpecificActionView.actions {
				stringToRune := action.KeyBinding[0]
				if event.Rune() == rune(stringToRune) {
					a.eventChan <- model.Event{
						Type:               model.SpecificActionOccured,
						SpecificActionName: action.Name,
						ResourceName:       a.MainView.view.GetCell(row, col+1).Text,
					}
				}
			}

		case tcell.KeyCtrlY:
			row, _ := a.MainView.view.GetSelection()
			if row == 0 {
				return event
			}
			a.eventChan <- model.Event{
				Type:     model.ReadResource,
				RowIndex: row,
			}

		case tcell.KeyCtrlD:
			row, _ := a.MainView.view.GetSelection()
			if row == 0 {
				return event
			}
			a.eventChan <- model.Event{
				Type:     model.ShowModal,
				RowIndex: row,
			}

		case tcell.KeyCtrlA:
			if a.rootView.GetItemCount() == 2 {
				a.rootView.AddItem(a.SearchView.GetView(), 2, 1, true)
				a.application.SetFocus(a.SearchView.GetView())
			} else {
				a.SearchView.GetView().SetText("")
				a.rootView.RemoveItem(a.SearchView.GetView())
				a.application.SetFocus(a.MainView.GetView())
			}

		}

		return event
	})
}

func (a *Application) RemoveSearchView() {
	a.rootView.RemoveItem(a.SearchView.GetView())
	a.application.SetFocus(a.MainView.GetView())
}

func (a *Application) Start() error {
	return a.application.EnableMouse(true).Run()
}
