package views

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/ghodss/yaml"
	"github.com/rivo/tview"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/utils"
	"github.com/sharadregoti/devops/utils/logger"
)

const (
	rootPage        = "root-page"
	deleteModalPage = "delete-modal-page"
	viewPage        = "view-page"
	formPage        = "form-page"
)

type Application struct {
	rootView           *tview.Flex
	MainView           *MainView
	GeneralInfoView    *GeneralInfo
	SearchView         *SearchView
	IsolatorView       *IsolatorView
	PluginView         *PluginView
	FlashView          *FlashView
	application        *tview.Application
	page               *tview.Pages
	Ta                 *tview.TextView
	SpecificActionView *SpecificActions
	ActionView         *Actions
	modal              *tview.Modal
	addr               string
	form               *tview.Form

	connectionID string

	// Communication Channel
	wsdata chan model.WebsocketResponse

	// Application state
	currentIsolator string
}

func NewApplication(addr string) (*Application, error) {
	// c := make(chan model.Event, 1)
	pa := tview.NewPages()

	i := NewIsolatorView()
	m := NewMainView()
	s := NewSearchView() // Need chan
	g := NewGeneralInfo()
	p := NewPluginView()
	act := NewAction()
	sa := NewSpecificAction()
	flash := NewFlashView()
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Form")

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
		AddItem(global, 7, 1, false).
		AddItem(m.GetView(), 0, 1, true).
		AddItem(s.GetView(), 0, 0, false) // Default disable

	r := &Application{
		rootView:           rootFlexContainer,
		MainView:           m,
		GeneralInfoView:    g,
		SearchView:         s,
		IsolatorView:       i,
		PluginView:         p,
		SpecificActionView: sa,
		form:               form,
		// application:        a,
		page: pa,
		// Ta:                 ta,
		ActionView: act,
		FlashView:  flash,
		// modal:      modal,
		addr: addr,
		// wsdata:     wsdata,
	}

	// Model page
	// TODO: Preselect no button
	modal := tview.NewModal().
		SetText("Do you want to delete the resource?").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				row, _ := m.view.GetSelection()
				resourceType, _ := utils.ParseTableTitle(m.view.GetTitle())

				_, err := r.sendEvent(model.FrontendEvent{
					EventType:    model.NormalAction,
					ActionName:   string(model.DeleteResource),
					ResourceName: m.view.GetCell(row, 1).Text,
					ResourceType: strings.ToLower(resourceType),
					IsolatorName: m.view.GetCell(row, 0).Text,
				})
				if err != nil {
					return
				}
				// c <- model.Event{
				// 	Type:         model.DeleteResource,
				// 	RowIndex:     row,
				// 	ResourceName: m.view.GetCell(row, 1).Text,
				// 	// This will be NA if resource in not isolator specific
				// 	IsolatorName: m.view.GetCell(row, 0).Text,
				// 	ResourceType: strings.ToLower(resourceType),
				// }
				pa.SwitchToPage(rootPage)
				// a.SwitchToMainc()
				// go func() {
				// Give time for refresh
				// time.Sleep(1 * time.Second)
				// r.sendEvent(model.FrontendEvent{
				// EventType:  model.InternalAction,
				// ActionName: string(model.RefreshResource),
				// RowIndex:   row,
				// })
				// c <- model.Event{
				// 	Type:         string(model.RefreshResource),
				// 	RowIndex:     row,
				// 	ResourceName: m.view.GetCell(row, 1).Text,
				// 	// This will be NA if resource in not isolator specific
				// 	IsolatorName: m.view.GetCell(row, 0).Text,
				// 	ResourceType: strings.ToLower(resourceType),
				// }
				// }()
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
	pa.AddPage(formPage, form, true, true)
	pa.SwitchToPage(rootPage)

	ta.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			// r.sendEvent(model.FrontendEvent{
			// 	EventType:  model.InternalAction,
			// 	ActionName: string(model.Close),
			// 	// RowIndex:   row,
			// })
			// c <- model.Event{
			// 	Type: model.Close,
			// }
			// logger.LogDebug("Sending close event chan")
			pa.SwitchToPage(rootPage)
		}
	})

	wsdata := make(chan model.WebsocketResponse, 1)
	go func() {
		for v := range wsdata {
			r.SpecificActionView.RefreshActions(v.SpecificActions)
			// c.appView.ActionView.EnableNesting(rs.currentSchema.Nesting.IsNested)
			r.MainView.SetTitle(v.TableName)
			r.RemoveSearchView()
			r.MainView.Refresh(v.Data, 0)
			r.GetApp().SetFocus(r.MainView.GetView())
			r.GetApp().Draw()
			logger.LogDebug("Websocket: received data from server, total length (%v)", len(v.Data))
		}
	}()

	a := tview.NewApplication().SetRoot(pa, true)
	r.application = a
	r.Ta = ta
	r.modal = modal
	r.wsdata = wsdata

	// This is here because during the search we if we press any rune key that corresponding function get triggered
	m.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter && act.nestingEnabled {
			row, _ := m.view.GetSelection()
			if row == 0 {
				return nil
			}
			// c <- model.Event{
			// 	Type:         string(model.ViewNestedResource),
			// 	RowIndex:     row,
			// 	ResourceName: m.view.GetCell(row, 1).Text,
			// }
			return event
		}

		if event.Key() == tcell.KeyRune {
			// Specific Action Checks
			row, _ := m.view.GetSelection()
			if row == 0 {
				return nil
			}
			for _, action := range sa.actions {

				// TODO: Validate key bindings
				stringToRune := action.KeyBinding[0]
				if event.Rune() == rune(stringToRune) {
					resourceType, _ := utils.ParseTableTitle(m.view.GetTitle())

					fe := model.FrontendEvent{
						EventType:    model.SpecificAction,
						ActionName:   action.Name,
						ResourceType: strings.ToLower(resourceType),
						ResourceName: m.view.GetCell(row, 1).Text,
						IsolatorName: m.view.GetCell(row, 0).Text,
					}
					if action.Execution.UserInput.Required {
						go func() {
							r.ShowForm(action.Execution.UserInput.Args, fe)
							r.GetApp().Draw()
						}()
						return event
					}

					eRes, err := r.sendEvent(fe)
					if err != nil {
						return event
					}

					if action.OutputType == model.OutputTypeString {
						stringData := eRes.Result.(string)
						go func() {
							r.SetTextAndSwitchView(stringData)
							r.GetApp().Draw()
						}()
					}

					switch action.OutputType {
					case model.OutputTypeBidrectional, model.OutputTypeStream:
						logger.LogDebug("Executing action command: %v", eRes.Result.(string))
						a.Suspend(func() {
							fmt.Fprintf(os.Stdout, "\033[2J\033[H")
							logger.LogDebug("Starting suspension:")
							if err := utils.ExecuteCMD(eRes.Result.(string)); err != nil {
								logger.LogError("Failed to execute command: %v", err.Error())
								return
							}
							logger.LogDebug("Existing suspension:")

							go a.Draw()
							logger.LogDebug("Redrawing:")
						})

						logger.LogDebug("Suspension function ended")
					}

					return event
				}
			}
		}
		return event
	})

	r.SetKeyboardShortCuts()

	infoRes, err := r.getInfo()
	if err != nil {
		return nil, err
	}

	if err := r.tableWebsocket(infoRes.SessionID, wsdata); err != nil {
		return nil, err
	}

	r.connectionID = infoRes.SessionID
	r.ActionView.RefreshActions(infoRes.Actions)
	r.SearchView.SetResourceTypes(infoRes.ResourceTypes)
	r.GeneralInfoView.Refresh(infoRes.General)
	r.IsolatorView.SetDefault(infoRes.DefaultIsolator)
	r.currentIsolator = infoRes.DefaultIsolator
	r.IsolatorView.SetTitle(strings.Title(infoRes.IsolatorType))

	r.SearchView.GetView().Autocomplete().SetDoneFunc(func(key tcell.Key) {
		r.sendEvent(model.FrontendEvent{
			EventType:    model.InternalAction,
			ActionName:   string(model.ResourceTypeChanged),
			ResourceType: r.SearchView.GetView().GetText(),
			ResourceName: "",
			IsolatorName: r.currentIsolator,
		})
		r.SearchView.GetView().SetText("")
	})
	return r, nil
}

func (a *Application) SetTextAndSwitchView(text string) {
	a.Ta.SetText(text)
	a.page.SwitchToPage(viewPage)
}

func (a *Application) ShowForm(formData map[string]interface{}, fe model.FrontendEvent) {
	a.form.Clear(true)

	for key, value := range formData {
		a.form.AddInputField(key, value.(string), 0, nil, nil)
	}

	a.form.AddButton("OK", func() {
		a.page.SwitchToPage(rootPage)
		args := map[string]interface{}{}
		for key := range formData {
			fi := a.form.GetFormItemByLabel(key)
			args[key] = fi.(*tview.InputField).GetText()
		}
		fe.Args = args
		_, err := a.sendEvent(fe)
		if err != nil {
			return
		}
	})
	a.form.AddButton("Cancel", func() {
		a.page.SwitchToPage(rootPage)
	})

	a.page.SwitchToPage(formPage)
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
		if event.Key() == tcell.KeyEscape {
			// row, _ := a.MainView.view.GetSelection()
			// if row == 0 {
			// 	return event
			// }
			// a.eventChan <- model.Event{
			// 	Type:     model.NestBack,
			// 	RowIndex: row,
			// }
			return event
		}

		switch event.Key() {

		// case tcell.KeyEnter:
		// 	// Specific action chec
		// 	a.eventChan <- model.Event{
		// 		Type:         model.ResourceTypeChanged,
		// 		RowIndex:     0,
		// 		ResourceType: "pods",
		// 		IsolatorName: "default",
		// 	}

		case tcell.KeyRune:
			row, _ := a.MainView.view.GetSelection()
			if event.Rune() == 'u' {
				if row == 0 {
					return event
				}
				go func() {
					// pCtx.clearNestedResource()
					a.IsolatorView.AddAndRefreshView(a.MainView.view.GetCell(row, 1).Text)
					a.GetApp().Draw()
				}()

				// resourceType, _ := utils.ParseTableTitle(a.MainView.view.GetTitle())
				// a.sendEvent(model.FrontendEvent{
				// EventType:  model.SpecificAction,
				// ActionName: string(model.AddIsolator),
				// RowIndex:   row,
				// })
				// a.eventChan <- model.Event{
				// 	Type: string(model.AddIsolator),
				// 	// +1 because "name" is always the second column that indicates the isolator name
				// 	// TODO: In future, change it to dynamically detected isolator name from table content irrespecitve of where the column
				// 	IsolatorName: a.MainView.view.GetCell(row, 1).Text,
				// 	ResourceType: strings.ToLower(resourceType),
				// }
				return event
			}

			// Single digit numbers only
			// Isolator View
			for i := range a.IsolatorView.currentKeyMap {
				numToRune := fmt.Sprintf("%d", i)[0]
				if event.Rune() == rune(numToRune) {
					resourceType, _ := utils.ParseTableTitle(a.MainView.view.GetTitle())

					_, err := a.sendEvent(model.FrontendEvent{
						EventType:    model.NormalAction,
						ActionName:   string(model.IsolatorChanged),
						ResourceName: a.MainView.view.GetCell(row, 1).Text,
						ResourceType: strings.ToLower(resourceType),
						IsolatorName: a.IsolatorView.currentKeyMap[i],
					})
					if err != nil {
						return event
					}
					a.currentIsolator = a.IsolatorView.currentKeyMap[i]
					// a.eventChan <- model.Event{
					// 	Type:         string(model.IsolatorChanged),
					// 	IsolatorName: a.IsolatorView.currentKeyMap[i],
					// }
				}
			}

		case tcell.KeyCtrlR:
			row, _ := a.MainView.view.GetSelection()
			resourceType, _ := utils.ParseTableTitle(a.MainView.view.GetTitle())
			_, err := a.sendEvent(model.FrontendEvent{
				EventType:    model.InternalAction,
				ActionName:   string(model.RefreshResource),
				ResourceName: a.MainView.view.GetCell(row, 1).Text,
				ResourceType: strings.ToLower(resourceType),
				IsolatorName: a.MainView.view.GetCell(row, 0).Text,
			})
			if err != nil {
				return nil
			}

		case tcell.KeyCtrlE:
			row, _ := a.MainView.view.GetSelection()
			if row == 0 {
				// Remove header row
				return event
			}
			resourceType, _ := utils.ParseTableTitle(a.MainView.view.GetTitle())

			a.GetApp().Suspend(func() {
				fmt.Fprintf(os.Stdout, "\033[2J\033[H")
				logger.LogDebug("Starting suspension:")

				f, err := os.CreateTemp("", "*")
				if err != nil {
					logger.LogError("Failed to create temp file: %v", err.Error())
					return
				}
				defer os.Remove(f.Name())

				fileName := f.Name()
				logger.LogDebug("Temp file name is: %v", fileName)

				res, err := a.sendEvent(model.FrontendEvent{
					EventType:    model.NormalAction,
					ActionName:   string(model.ReadResource),
					ResourceName: a.MainView.view.GetCell(row, 1).Text,
					ResourceType: strings.ToLower(resourceType),
					IsolatorName: a.MainView.view.GetCell(row, 0).Text,
				})
				if err != nil {
					return
				}
				dd, _ := yaml.Marshal(res.Result)

				_, err = f.WriteString(string(dd))
				if err != nil {
					logger.LogError("Failed to write to temp file: %v", err.Error())
					return
				}
				f.Close()

				if err := utils.ExecuteCMD("vi " + fileName); err != nil {
					logger.LogError("Failed to execute command: %v", err.Error())
					return
				}

				index := 0
				for {

					ff, err := os.ReadFile(fileName)
					if err != nil {
						logger.LogError("Failed to open temp file: %v", err.Error())
						return
					}

					yamlContent := map[string]interface{}{}
					if err := yaml.Unmarshal(ff, &yamlContent); err != nil {
						logger.LogError("Failed to marshal yaml: %v", err.Error())
						return
					}

					_, errs := a.sendEvent(model.FrontendEvent{
						EventType:    model.NormalAction,
						ActionName:   string(model.EditResource),
						ResourceName: a.MainView.view.GetCell(row, 1).Text,
						ResourceType: strings.ToLower(resourceType),
						IsolatorName: a.MainView.view.GetCell(row, 0).Text,
						Args:         yamlContent,
					})
					if errs == nil {
						break
					}

					f, err = os.OpenFile(fileName, os.O_RDWR, 0644)
					if err != nil {
						logger.LogError("Failed to open temp file: %v", err.Error())
						return
					}

					ff, err = os.ReadFile(fileName)
					if err != nil {
						logger.LogError("Failed to open temp file: %v", err.Error())
						return
					}

					newContent := ""
					if index == 0 {
						newContent = fmt.Sprintf("# %s\n%s", errs.Error(), string(ff))
						index++
					} else {
						index := strings.IndexByte(string(ff), byte('\n'))
						newContent = fmt.Sprintf("# %s\n%s", errs.Error(), string(ff)[index+1:])
					}
					_, err = f.WriteString(newContent)
					if err != nil {
						logger.LogError("Failed to write to temp file: %v", err.Error())
						return
					}
					f.Close()

					if err := utils.ExecuteCMD("vi " + fileName); err != nil {
						logger.LogError("Failed to execute command: %v", err.Error())
						return
					}

				}
				logger.LogDebug("Existing suspension:")

				go a.GetApp().Draw()
				logger.LogDebug("Redrawing:")
			})

		case tcell.KeyCtrlY:
			row, _ := a.MainView.view.GetSelection()
			if row == 0 {
				// Remove header row
				return event
			}
			resourceType, _ := utils.ParseTableTitle(a.MainView.view.GetTitle())
			res, err := a.sendEvent(model.FrontendEvent{
				EventType:    model.NormalAction,
				ActionName:   string(model.ReadResource),
				ResourceName: a.MainView.view.GetCell(row, 1).Text,
				ResourceType: strings.ToLower(resourceType),
				IsolatorName: a.MainView.view.GetCell(row, 0).Text,
			})
			if err != nil {
				return nil
			}
			go func() {
				dd, _ := yaml.Marshal(res.Result)
				a.SetTextAndSwitchView(string(dd))
				go a.GetApp().Draw()
			}()

			// a.eventChan <- model.Event{
			// 	Type:         string(model.ReadResource),
			// 	RowIndex:     row,
			// 	ResourceName: a.MainView.view.GetCell(row, 1).Text,
			// 	ResourceType: strings.ToLower(resourceType),
			// 	// This will be NA if resource in not isolator specific
			// 	IsolatorName: a.MainView.view.GetCell(row, 0).Text,
			// }

		case tcell.KeyCtrlD:
			row, _ := a.MainView.view.GetSelection()
			if row == 0 {
				return event
			}
			resourceType, _ := utils.ParseTableTitle(a.MainView.view.GetTitle())
			// a.sendEvent(model.FrontendEvent{
			// 	EventType:  model.NormalAction,
			// 	ActionName: string(model.ReadResource),
			// 	RowIndex:   row,
			// })
			// a.eventChan <- model.Event{
			// 	Type:         model.ShowModal,
			// 	RowIndex:     row,
			// 	ResourceName: a.MainView.view.GetCell(row, 1).Text,
			// 	ResourceType: strings.ToLower(resourceType),
			// 	// This will be NA if resource in not isolator specific
			// 	IsolatorName: a.MainView.view.GetCell(row, 0).Text,
			// }
			go func() {
				a.ViewModel(resourceType, a.MainView.view.GetCell(row, 1).Text)
				a.GetApp().Draw()
			}()

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
	return a.application.EnableMouse(false).Run()
}

func (a *Application) SetFlashText(text string) {
	a.rootView.AddItem(a.FlashView.GetView(), 2, 1, true)
	a.FlashView.SetText(text)
	go func() {
		<-time.After(3 * time.Second)
		a.rootView.RemoveItem(a.FlashView.GetView())
		a.application.Draw()
	}()
	go a.application.Draw()
}
