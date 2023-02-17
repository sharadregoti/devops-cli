package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/sharadregoti/devops-plugin-sdk/proto"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/utils"
	"github.com/sharadregoti/devops/utils/logger"
	"gopkg.in/yaml.v2"
)

func (a *Application) registerEventHandlers() {
	a.deleteModalPageEventHandler()
	a.textOnlyPageEventHandler()
	a.mainPageEventHandler()
}

func (a *Application) mainPageEventHandler() error {
	a.mainPage.searchBox.view.Autocomplete().SetDoneFunc(func(key tcell.Key) {
		// TODO: when invalid resource type is send, the server may/may not return an errror. Error is not shown to the user
		searchResult := a.mainPage.searchBox.view.GetText()

		if strings.HasPrefix(searchResult, settingPlugin) {
			pluginName := parsePluginSetting(searchResult)
			if pluginName == a.currentPluginName {
				a.mainPage.searchBox.view.SetText("")
				a.SetFlashText("Plugin already set to " + pluginName)
				return
			}

			if err := a.loadPlugin(pluginName); err != nil {
				a.flashLogError(err.Error())
			}
			return
		}

		if strings.HasPrefix(searchResult, settingAuthentication) {
			identifyingName, name := parseAuthenticationSetting(searchResult)
			// TODO: Need to close connection from server
			a.closeChan <- struct{}{}
			go func() {
				a.mainPage.searchBox.view.SetText("")
				if err := a.connectAndLoadData(a.currentPluginName, &proto.AuthInfo{IdentifyingName: identifyingName, Name: name}); err != nil {
					a.flashLogError(err.Error())
				}
			}()
			// go a.application.Draw()
			return
		}

		if a.currentResourceType == searchResult {
			a.mainPage.searchBox.view.SetText("")
			go a.application.Draw()
			return
		}

		_, err := a.sendEvent(model.FrontendEvent{
			EventType:    model.InternalAction,
			ActionName:   string(model.ResourceTypeChanged),
			ResourceType: a.mainPage.searchBox.view.GetText(),
			ResourceName: "",
			IsolatorName: a.currentIsolator,
			PluginName:   a.currentPluginName})
		if err != nil {
			a.flashLogError(err.Error())
		}
		a.currentResourceType = searchResult
		a.mainPage.searchBox.view.SetText("")
	})

	// This is here because during the search we if we press any rune key that corresponding function get triggered
	a.mainPage.tableBox.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// if event.Key() == tcell.KeyEnter && normalActionBox.nestingEnabled {
		// 	row, _ := a.mainPage.tableBox.view.GetSelection()
		// 	if row == 0 {
		// 		return nil
		// 	}
		// 	// c <- model.Event{
		// 	// 	Type:         string(model.ViewNestedResource),
		// 	// 	RowIndex:     row,
		// 	// 	ResourceName: m.view.GetCell(row, 1).Text,
		// 	// }
		// 	return event
		// }

		if event.Key() == tcell.KeyEnter {
			// If title is "Authentication" then only move forward
			if a.mainPage.tableBox.view.GetTitle() != "Authentication" {
				return nil
			}

			row, _ := a.mainPage.tableBox.view.GetSelection()
			if row == 0 {
				return nil
			}

			go func() {
				id := a.mainPage.tableBox.view.GetCell(row, 0).Text
				name := a.mainPage.tableBox.view.GetCell(row, 1).Text
				if err := a.connectAndLoadData(a.currentPluginName, &proto.AuthInfo{IdentifyingName: id, Name: name}); err != nil {
					a.flashLogError(err.Error())
				}
			}()
			// c <- model.Event{
			// 	Type:         string(model.ViewNestedResource),
			// 	RowIndex:     row,
			// 	ResourceName: m.view.GetCell(row, 1).Text,
			// }
			return event
		}

		if event.Key() == tcell.KeyRune {
			// Specific Action Checks
			row, _ := a.mainPage.tableBox.view.GetSelection()
			if row == 0 {
				return nil
			}

			for _, action := range a.mainPage.specificActionBox.actions {

				// TODO: Validate key bindings
				stringToRune := action.KeyBinding[0]
				if event.Rune() == rune(stringToRune) {
					resourceType, _ := utils.ParseTableTitle(a.mainPage.tableBox.view.GetTitle())

					fe := model.FrontendEvent{
						EventType:    model.SpecificAction,
						ActionName:   action.Name,
						ResourceType: strings.ToLower(resourceType),
						ResourceName: a.mainPage.tableBox.view.GetCell(row, 1).Text,
						IsolatorName: a.currentIsolator,
						PluginName:   a.currentPluginName,
					}
					if action.Execution != nil && action.Execution.UserInput != nil && action.Execution.UserInput.Required {
						go func() {
							fe.ActionName = string(model.SpecificActionResolveArgs)
							fe.Args = utils.GetMapInterface(action.Execution.UserInput.Args)
							eRes, err := a.sendEvent(fe)
							if err != nil {
								return
							}
							// rewrite the action name
							fe.ActionName = action.Name
							a.ShowForm(eRes.Result.(map[string]interface{}), fe)
							a.application.Draw()
						}()
						return event
					}

					eRes, err := a.sendEvent(fe)
					if err != nil {
						return event
					}

					if action.OutputType == model.OutputTypeString {
						stringData := eRes.Result.(string)
						go func() {
							a.SetTextAndSwitchView(stringData)
							a.application.Draw()
						}()
					}

					switch action.OutputType {
					case model.OutputTypeBidrectional, model.OutputTypeStream:
						logger.LogDebug("Executing action command: %v", eRes.Result.(string))
						a.application.Suspend(func() {
							fmt.Fprintf(os.Stdout, "\033[2J\033[H")
							logger.LogDebug("Starting suspension:")
							if err := utils.ExecuteCMD(eRes.Result.(string)); err != nil {
								logger.LogError("Failed to execute command: %v", err.Error())
								return
							}
							logger.LogDebug("Existing suspension:")

							go a.application.Draw()
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

	a.mainPage.flexView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

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
			row, _ := a.mainPage.tableBox.view.GetSelection()
			if event.Rune() == 'u' {
				if row == 0 {
					return event
				}
				go func() {
					// pCtx.clearNestedResource()
					a.mainPage.isolatorBox.AddAndRefreshView(a.mainPage.tableBox.view.GetCell(row, 1).Text)
					a.application.Draw()
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

			for i := range a.mainPage.isolatorBox.currentKeyMap {
				// This check ensure title count does not get dublicated when same namespace events is triggered
				numToRune := fmt.Sprintf("%s", i)[0]
				if event.Rune() == rune(numToRune) {
					if a.currentIsolator == a.mainPage.isolatorBox.currentKeyMap[i] {
						return nil
					}
					resourceType, _ := utils.ParseTableTitle(a.mainPage.tableBox.view.GetTitle())

					_, err := a.sendEvent(model.FrontendEvent{
						EventType:    model.NormalAction,
						ActionName:   string(model.IsolatorChanged),
						ResourceName: "",
						ResourceType: strings.ToLower(resourceType),
						IsolatorName: a.mainPage.isolatorBox.currentKeyMap[i],
						PluginName:   a.currentPluginName})
					if err != nil {
						return event
					}
					a.mainPage.tableBox.view.Clear()
					a.currentIsolator = a.mainPage.isolatorBox.currentKeyMap[i]
					a.closeChan <- struct{}{}
					a.startGoRoutine()
					// a.eventChan <- model.Event{
					// 	Type:         string(model.IsolatorChanged),
					// 	IsolatorName: a.IsolatorView.currentKeyMap[i],
					// }
				}
			}

		case tcell.KeyCtrlR:
			row, _ := a.mainPage.tableBox.view.GetSelection()
			resourceType, _ := utils.ParseTableTitle(a.mainPage.tableBox.view.GetTitle())
			_, err := a.sendEvent(model.FrontendEvent{
				EventType:    model.InternalAction,
				ActionName:   string(model.RefreshResource),
				ResourceName: a.mainPage.tableBox.view.GetCell(row, 1).Text,
				ResourceType: strings.ToLower(resourceType),
				IsolatorName: a.mainPage.tableBox.view.GetCell(row, 0).Text,
			})
			if err != nil {
				return nil
			}

		case tcell.KeyCtrlL:
			row, _ := a.mainPage.tableBox.view.GetSelection()
			if row == 0 {
				// Remove header row
				return event
			}
			resourceType, _ := utils.ParseTableTitle(a.mainPage.tableBox.view.GetTitle())
			_, err := a.sendEvent(model.FrontendEvent{
				EventType:    model.NormalAction,
				ActionName:   string(model.ViewLongRunning),
				ResourceName: a.mainPage.tableBox.view.GetCell(row, 1).Text,
				ResourceType: strings.ToLower(resourceType),
				IsolatorName: a.mainPage.tableBox.view.GetCell(row, 0).Text,
			})
			if err != nil {
				return event
			}

		case tcell.KeyCtrlE:
			row, _ := a.mainPage.tableBox.view.GetSelection()
			if row == 0 {
				// Remove header row
				return event
			}
			resourceType, _ := utils.ParseTableTitle(a.mainPage.tableBox.view.GetTitle())

			a.application.Suspend(func() {
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
					ResourceName: a.mainPage.tableBox.view.GetCell(row, 1).Text,
					ResourceType: strings.ToLower(resourceType),
					IsolatorName: a.mainPage.tableBox.view.GetCell(row, 0).Text,
				})
				if err != nil {
					return
				}
				// delete(res.Result.(map[string]interface{}), "devops")
				dd, _ := yaml.Marshal(res.Result)

				_, err = f.WriteString(string(dd))
				if err != nil {
					logger.LogError("Failed to write to temp file: %v", err.Error())
					return
				}
				f.Close()

				// TODO: Edit will fail
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
						ResourceName: a.mainPage.tableBox.view.GetCell(row, 1).Text,
						ResourceType: strings.ToLower(resourceType),
						IsolatorName: a.mainPage.tableBox.view.GetCell(row, 0).Text,
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

				go a.application.Draw()
				logger.LogDebug("Redrawing:")
			})

		case tcell.KeyCtrlY:
			row, _ := a.mainPage.tableBox.view.GetSelection()
			if row == 0 {
				// Remove header row
				return event
			}
			resourceType, _ := utils.ParseTableTitle(a.mainPage.tableBox.view.GetTitle())
			res, err := a.sendEvent(model.FrontendEvent{
				EventType:    model.NormalAction,
				ActionName:   string(model.ReadResource),
				ResourceName: a.mainPage.tableBox.view.GetCell(row, 1).Text,
				ResourceType: strings.ToLower(resourceType),
				IsolatorName: a.mainPage.tableBox.view.GetCell(row, 0).Text,
			})
			if err != nil {
				return nil
			}
			go func() {
				// delete(res.Result.(map[string]interface{}), "devops")
				dd, _ := yaml.Marshal(res.Result)
				a.SetTextAndSwitchView(string(dd))
				go a.application.Draw()
			}()

			// a.eventChan <- model.Event{
			// 	Type:         string(model.ReadResource),
			// 	RowIndex:     row,
			// 	ResourceName: a.mainPage.tableBox.view.GetCell(row, 1).Text,
			// 	ResourceType: strings.ToLower(resourceType),
			// 	// This will be NA if resource in not isolator specific
			// 	IsolatorName: a.mainPage.tableBox.view.GetCell(row, 0).Text,
			// }

		case tcell.KeyCtrlD:
			row, _ := a.mainPage.tableBox.view.GetSelection()
			if row == 0 {
				return event
			}
			resourceType, _ := utils.ParseTableTitle(a.mainPage.tableBox.view.GetTitle())
			// a.sendEvent(model.FrontendEvent{
			// 	EventType:  model.NormalAction,
			// 	ActionName: string(model.ReadResource),
			// 	RowIndex:   row,
			// })
			// a.eventChan <- model.Event{
			// 	Type:         model.ShowModal,
			// 	RowIndex:     row,
			// 	ResourceName: a.mainPage.tableBox.view.GetCell(row, 1).Text,
			// 	ResourceType: strings.ToLower(resourceType),
			// 	// This will be NA if resource in not isolator specific
			// 	IsolatorName: a.mainPage.tableBox.view.GetCell(row, 0).Text,
			// }
			go func() {
				a.ViewModel(resourceType, a.mainPage.tableBox.view.GetCell(row, 1).Text)
				a.application.Draw()
			}()

		case tcell.KeyCtrlA:
			if a.mainPage.flexView.GetItemCount() == 2 {
				a.mainPage.flexView.AddItem(a.mainPage.searchBox.view, 2, 1, true)
				a.application.SetFocus(a.mainPage.searchBox.view)
			} else {
				a.mainPage.searchBox.view.SetText("")
				a.mainPage.flexView.RemoveItem(a.mainPage.searchBox.view)
				a.application.SetFocus(a.mainPage.tableBox.view)
			}

		}

		return event
	})

	return nil
}

func (a *Application) textOnlyPageEventHandler() error {
	a.textOnlyPage.view.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			a.pages.SwitchToPage(cmainPage)
		}
	})

	return nil
}

func (a *Application) deleteModalPageEventHandler() error {
	a.deleteModalPage.view.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Yes" {
			row, _ := a.mainPage.tableBox.view.GetSelection()
			resourceType, _ := utils.ParseTableTitle(a.mainPage.tableBox.view.GetTitle())

			_, err := a.sendEvent(model.FrontendEvent{
				EventType:    model.NormalAction,
				ActionName:   string(model.DeleteResource),
				ResourceName: a.mainPage.tableBox.view.GetCell(row, 1).Text,
				ResourceType: strings.ToLower(resourceType),
				IsolatorName: a.mainPage.tableBox.view.GetCell(row, 0).Text,
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
			a.pages.SwitchToPage(cmainPage)
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
			a.pages.SwitchToPage(cmainPage)
		}
	})

	return nil
}
