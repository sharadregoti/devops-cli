package views

import (
	"github.com/rivo/tview"
)

type Application struct {
	rootView        *tview.Flex
	MainView        *MainView
	GeneralInfoView *GeneralInfo
	SearchView      *SearchView
	IsolatorView    *IsolatorView
	PluginView      *PluginView
}

func NewApplication() *Application {
	m := NewMainView()
	g := NewGeneralInfo()
	s := NewSearchView()
	i := NewIsolatorView()
	p := NewPluginView()

	global := tview.NewFlex().
		AddItem(g.GetView(), 0, 1, false).
		AddItem(i.GetView(), 0, 1, false).
		AddItem(p.GetView(), 0, 1, false)

	rootFlexContainer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(global, 6, 1, false).
		AddItem(m.GetView(), 0, 1, true).
		AddItem(s.GetView(), 2, 1, false)

	// rootFlexContainer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

	// 	tcell.NewEventKey(tcell.Key(tcell.ModShift), ':', tcell.ModAlt)

	// 	if event.Key() == tcell.KeyCtrlB {
	// 		rootFlexContainer.RemoveItem(fixedSearchBox)
	// 	}

	// 	if event.Key() == tcell.KeyCtrlA {
	// 		rootFlexContainer.AddItem(fixedSearchBox, 2, 1, false)
	// 	}

	// 	return event
	// })

	return &Application{
		rootView:        rootFlexContainer,
		MainView:        m,
		GeneralInfoView: g,
		SearchView:      s,
		IsolatorView:    i,
		PluginView:      p,
	}
}

func (a *Application) Start() error {
	return tview.NewApplication().SetRoot(a.rootView, true).Run()
}
