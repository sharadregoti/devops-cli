package tui

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
	"github.com/sharadregoti/devops-plugin-sdk/proto"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/utils/logger"
)

const (
	cmainPage        = "main-page"
	cdeleteModalPage = "delete-modal-page"
	ctextOnlyPage    = "text-only-page"
	cformPage        = "form-page"
)

type Application struct {
	mainPage        *mainPage
	textOnlyPage    *textOnlyPage
	formPage        *formPage
	deleteModalPage *deleteModalPage
	pages           *tview.Pages
	flashView       *FlashView

	application *tview.Application

	addr      string
	sessionID string

	// Websocket Communication Channel
	wsdata    chan model.WebsocketResponse
	closeChan chan struct{}

	// Application state
	currentIsolator     string
	currentResourceType string
	currentPluginName   string
	settings            []string

	// server
	appConfig          *model.Config
	currentPluginAuths *model.AuthResponse
}

func NewApplication(addr string) (*Application, error) {

	// Initialize pages
	mainPage := newMainPage()
	textOnlyPage := newTextOnlyPage()
	formPage := newFormPage()
	deleteModalPage := newDeleteModalPage()

	// Add pages to the app
	pa := tview.NewPages()
	pa.AddPage(cmainPage, mainPage.flexView, true, true)
	pa.AddPage(cdeleteModalPage, deleteModalPage.view, true, true)
	pa.AddPage(ctextOnlyPage, textOnlyPage.view, true, true)
	pa.AddPage(cformPage, formPage.view, true, true)
	pa.SwitchToPage(cmainPage)

	flash := NewFlashView()

	r := &Application{
		mainPage:        mainPage,
		textOnlyPage:    textOnlyPage,
		formPage:        formPage,
		deleteModalPage: deleteModalPage,
		pages:           pa,
		addr:            addr,
		flashView:       flash,
		application:     tview.NewApplication().SetRoot(pa, true),
		closeChan:       make(chan struct{}, 1),
	}

	return r, nil
}

func (a *Application) loadPlugin(pluginName string) error {
	// Get application session from server
	logger.LogDebug("Getting plugin auths from server...")
	pluginAuths, err := a.getPluginAuths(pluginName)
	if err != nil {
		return err
	}

	if len(pluginAuths.Auths) == 0 {
		// TODO: Show ony search box to enable plugin search or show plugins in UI
		return fmt.Errorf("no authentication found for default plugin %s", pluginName)
	}

	contexts := make([]*model.TableRow, 0)
	contexts = append(contexts, &model.TableRow{
		Data: []string{"ID", "NAME"},
	})

	// Iterate over all auths & add them to settings
	pluginAuth := new(proto.AuthInfo)
	settings := make([]string, 0)
	for _, ai := range pluginAuths.Auths {
		contexts = append(contexts, &model.TableRow{
			Data:  []string{ai.IdentifyingName, ai.Name},
			Color: "lightskyblue",
		})

		if ai.IsDefault {
			pluginAuth = ai
		}

		settings = append(settings, getAuthenticationSetting(ai.IdentifyingName, ai.Name))
	}

	// Iterate over all plugins & add them to settings
	for _, plugin := range a.appConfig.Plugins {
		settings = append(settings, getPluginSetting(plugin.Name))
	}

	a.registerEventHandlers()

	a.currentPluginAuths = pluginAuths
	a.currentPluginName = pluginName
	a.settings = settings
	// sort.Strings(a.settings)

	if pluginAuth.IdentifyingName == "" {
		// Show auth selection
		a.mainPage.tableBox.Refresh(contexts, 0)
		a.mainPage.tableBox.SetTitle("Authentication")
		// TODO: Handle enter on auth table
		return nil
	}

	// Load default plugin auth
	if err := a.connectAndLoadData(a.appConfig.Plugins[0].Name, pluginAuth); err != nil {
		logger.LogError("Failed to load default plugin auth: %v", err)
		// Show auth selection
		a.mainPage.tableBox.Refresh(contexts, 0)
		a.mainPage.tableBox.SetTitle("Authentication")
	}
	return nil
}

func (a *Application) SetTextAndSwitchView(text string) {
	a.textOnlyPage.view.SetText(text)
	a.pages.SwitchToPage(ctextOnlyPage)
}

func (a *Application) ShowForm(formData map[string]interface{}, fe model.FrontendEvent) {
	a.formPage.view.Clear(true)

	for key, value := range formData {
		a.formPage.view.AddInputField(key, value.(string), 0, nil, nil)
	}

	a.formPage.view.AddButton("OK", func() {
		a.pages.SwitchToPage(cmainPage)
		args := map[string]interface{}{}
		for key := range formData {
			fi := a.formPage.view.GetFormItemByLabel(key)
			args[key] = fi.(*tview.InputField).GetText()
		}
		fe.Args = args

		_, err := a.sendEvent(fe)
		if err != nil {
			return
		}
	})
	a.formPage.view.AddButton("Cancel", func() {
		a.pages.SwitchToPage(cmainPage)
	})

	a.pages.SwitchToPage(cformPage)
}

func (a *Application) SwitchToMain() {
	a.pages.SwitchToPage(cmainPage)
}

func (a *Application) ViewModel(rType, rName string) {
	a.deleteModalPage.view.SetText(fmt.Sprintf("Do you want to delete the %s/%s?", rType, rName))
	a.pages.SwitchToPage(cdeleteModalPage)
}

func (a *Application) RemoveSearchView() {
	a.mainPage.flexView.RemoveItem(a.mainPage.searchBox.view)
	a.application.SetFocus(a.mainPage.flexView)
}

func (a *Application) Start() error {
	return a.application.EnableMouse(false).Run()
}

func (a *Application) SetFlashText(text string) {
	a.mainPage.flexView.AddItem(a.flashView.GetView(), 2, 1, true)
	a.flashView.SetText(text)
	go func() {
		<-time.After(3 * time.Second)
		a.mainPage.flexView.RemoveItem(a.flashView.GetView())
		a.application.Draw()
	}()
	go a.application.Draw()
}
