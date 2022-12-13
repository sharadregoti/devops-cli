package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sharadregoti/devops/model"
)

type Application struct {
	rootView        *tview.Flex
	MainView        *MainView
	GeneralInfoView *GeneralInfo
	SearchView      *SearchView
	IsolatorView    *IsolatorView
	PluginView      *PluginView
	eventChan       chan model.Event
	application     *tview.Application
}

func NewApplication(c chan model.Event) *Application {
	m := NewMainView()
	g := NewGeneralInfo()
	s := NewSearchView(c)
	i := NewIsolatorView()
	p := NewPluginView()

	global := tview.NewFlex().
		AddItem(g.GetView(), 0, 1, false).
		AddItem(i.GetView(), 0, 1, false).
		AddItem(p.GetView(), 0, 1, false)

	rootFlexContainer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(global, 6, 1, false).
		AddItem(m.GetView(), 0, 1, false).
		AddItem(s.GetView(), 0, 0, false) // Default disable

	a := tview.NewApplication().SetRoot(rootFlexContainer, true)

	r := &Application{
		rootView:        rootFlexContainer,
		MainView:        m,
		GeneralInfoView: g,
		SearchView:      s,
		IsolatorView:    i,
		PluginView:      p,
		eventChan:       c,
		application:     a,
	}

	r.SetKeyboardShortCuts()

	return r
}

func (a *Application) GetApp() *tview.Application {
	return a.application
}

func (a *Application) SetKeyboardShortCuts() {
	a.rootView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyCtrlA {
			if a.rootView.GetItemCount() == 2 {
				a.rootView.AddItem(a.SearchView.GetView(), 2, 1, true)
				a.application.SetFocus(a.SearchView.GetView())
			} else {
				a.rootView.RemoveItem(a.SearchView.GetView())
				a.application.SetFocus(a.MainView.GetView())
			}
			return event
		}

		return event
	})
}

func (a *Application) Start() error {
	return a.application.EnableMouse(true).Run()
}
