package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/rivo/tview"
	"github.com/sharadregoti/devops/model"
)

const (
	rootPage        = "root-page"
	deleteModalPage = "delete-modal-page"
	viewPage        = "view-page"
)

type Application struct {
	rootView           *tview.Flex
	MainView           *MainView
	GeneralInfoView    *GeneralInfo
	SearchView         *SearchView
	IsolatorView       *IsolatorView
	PluginView         *PluginView
	FlashView          *FlashView
	eventChan          chan model.Event
	application        *tview.Application
	page               *tview.Pages
	Ta                 *tview.TextView
	SpecificActionView *SpecificActions
	ActionView         *Actions
	modal              *tview.Modal
	logger             hclog.Logger
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
	flash := NewFlashView(logger, c)

	// Global page info flex view
	global := tview.NewFlex().
		AddItem(g.GetView(), 0, 1, false).
		AddItem(i.GetView(), 0, 1, false).
		AddItem(act.view, 0, 1, false).
		AddItem(sa.GetView(), 0, 1, false).
		AddItem(p.GetView(), 0, 1, false)

	// Main page
	rootFlexContainer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(global, 6, 1, false).
		AddItem(m.GetView(), 0, 1, true).
		AddItem(s.GetView(), 0, 0, false) // Default disable

	// Model page
	// TODO: Preselect no button
	modal := tview.NewModal().
		SetText("Do you want to delete the resource?").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				row, _ := m.view.GetSelection()
				c <- model.Event{
					Type:         model.DeleteResource,
					RowIndex:     row,
					ResourceName: m.view.GetCell(row, 1).Text,
					// This will be NA if resource in not isolator specific
					IsolatorName: m.view.GetCell(row, 0).Text,
					ResourceType: strings.ToLower(m.view.GetTitle()),
				}
				pa.SwitchToPage(rootPage)
				c <- model.Event{
					Type:         model.RefreshResource,
					RowIndex:     row,
					ResourceName: m.view.GetCell(row, 1).Text,
					// This will be NA if resource in not isolator specific
					IsolatorName: m.view.GetCell(row, 0).Text,
					ResourceType: strings.ToLower(m.view.GetTitle()),
				}
			}
			if buttonLabel == "No" {
				pa.SwitchToPage(rootPage)
			}
		})

	// View only page
	ta := tview.NewTextView()
	ta.SetDynamicColors(true)

	pa.AddPage(rootPage, rootFlexContainer, true, true)
	pa.AddPage(deleteModalPage, modal, true, true)
	pa.AddPage(viewPage, ta, true, true)
	pa.SwitchToPage(rootPage)

	ta.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			c <- model.Event{
				Type: model.Close,
			}
			logger.Debug("Sending close event chan")
			pa.SwitchToPage(rootPage)
		}
	})

	a := tview.NewApplication().SetRoot(pa, true)

	// This is here because during the search we if we press any rune key that corresponding function get triggered
	m.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			// Specific Action Checks
			row, _ := m.view.GetSelection()
			for _, action := range sa.actions {
				// TODO: Validate key bindings
				stringToRune := action.KeyBinding[0]
				if event.Rune() == rune(stringToRune) {
					c <- model.Event{
						Type:               model.SpecificActionOccured,
						SpecificActionName: action.Name,
						ResourceName:       m.view.GetCell(row, 1).Text,
					}
				}
			}
		}
		return event
	})

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
		Ta:                 ta,
		ActionView:         act,
		logger:             logger,
		FlashView:          flash,
		modal:              modal,
	}

	r.SetKeyboardShortCuts()

	return r
}

func (a *Application) SetTextAndSwitchView(text string) {
	a.Ta.SetText(text)
	a.page.SwitchToPage(viewPage)
}

func (a *Application) SwitchToMain() {
	a.page.SwitchToPage(rootPage)
}

func (a *Application) GetApp() *tview.Application {
	return a.application
}

func (a *Application) ViewModel(rType, rName string) {
	a.modal.SetText(fmt.Sprintf("Do you want to delete the %s/%s?", rType, rName))
	a.page.SwitchToPage(deleteModalPage)
}

func (a *Application) SetKeyboardShortCuts() {
	a.rootView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		// if tcell.ModAlt != 0 {
		// 	for i, pluginName := range a.PluginView.data {
		// 		numToRune := fmt.Sprintf("%d", i)[0]
		// 		if event.Rune() == rune(numToRune) {
		// 			a.eventChan <- model.Event{
		// 				Type:       model.PluginChanged,
		// 				PluginName: pluginName,
		// 			}
		// 			return event
		// 		}
		// 	}
		// }

		switch event.Key() {

		case tcell.KeyRune:
			if event.Rune() == 'u' {
				row, _ := a.MainView.view.GetSelection()
				a.eventChan <- model.Event{
					Type: model.AddIsolator,
					// +1 because "name" is always the second column that indicates the isolator name
					// TODO: In future, change it to dynamically detected isolator name from table content irrespecitve of where the column
					IsolatorName: a.MainView.view.GetCell(row, 1).Text,
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

			// // Specific Action Checks
			// row, _ := a.MainView.view.GetSelection()
			// for _, action := range a.SpecificActionView.actions {
			// 	// TODO: Validate key bindings
			// 	stringToRune := action.KeyBinding[0]
			// 	if event.Rune() == rune(stringToRune) {
			// 		a.eventChan <- model.Event{
			// 			Type:               model.SpecificActionOccured,
			// 			SpecificActionName: action.Name,
			// 			ResourceName:       a.MainView.view.GetCell(row, 1).Text,
			// 		}
			// 	}
			// }

		case tcell.KeyCtrlR:
			row, _ := a.MainView.view.GetSelection()
			if row == 0 {
				return event
			}
			a.eventChan <- model.Event{
				Type:         model.RefreshResource,
				RowIndex:     row,
				ResourceName: a.MainView.view.GetCell(row, 1).Text,
				ResourceType: strings.ToLower(a.MainView.view.GetTitle()),
				// This will be NA if resource in not isolator specific
				IsolatorName: a.MainView.view.GetCell(row, 0).Text,
			}

		case tcell.KeyCtrlY:
			row, _ := a.MainView.view.GetSelection()
			if row == 0 {
				// Remove header row
				return event
			}
			a.eventChan <- model.Event{
				Type:         model.ReadResource,
				RowIndex:     row,
				ResourceName: a.MainView.view.GetCell(row, 1).Text,
				ResourceType: strings.ToLower(a.MainView.view.GetTitle()),
				// This will be NA if resource in not isolator specific
				IsolatorName: a.MainView.view.GetCell(row, 0).Text,
			}

		case tcell.KeyCtrlD:
			row, _ := a.MainView.view.GetSelection()
			if row == 0 {
				return event
			}
			a.eventChan <- model.Event{
				Type:         model.ShowModal,
				RowIndex:     row,
				ResourceName: a.MainView.view.GetCell(row, 1).Text,
				ResourceType: strings.ToLower(a.MainView.view.GetTitle()),
				// This will be NA if resource in not isolator specific
				IsolatorName: a.MainView.view.GetCell(row, 0).Text,
			}

			// case tcell.ctrl
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

func (a *Application) SetFlashText(text string) {
	a.rootView.AddItem(a.FlashView.GetView(), 2, 1, true)
	a.FlashView.SetText(text)
	go func() {
		<-time.After(5 * time.Second)
		a.rootView.RemoveItem(a.FlashView.GetView())
		a.application.Draw()
	}()
	a.application.Draw()
}
