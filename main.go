package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/sharadregoti/devops/internal/transformer"
	"github.com/sharadregoti/devops/internal/views"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/shared"
)

// Updating VS Code Server to version 1ad8d514439d5077d2b0b7ee64d2ce82a9308e5a
// Removing previous installation...
// Installing VS Code Server for x64 (1ad8d514439d5077d2b0b7ee64d2ce82a9308e5a)
// Downloading:  80%

func main() {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Cannot detect home directory", err)
	}

	// Create the ".devops" subdirectory if it doesn't exist
	devopsDir := filepath.Join(homeDir, ".devops")
	if _, err := os.Stat(devopsDir); os.IsNotExist(err) {
		err = os.Mkdir(devopsDir, 0755)
		if err != nil {
			log.Fatal("Cannot create .devops directory", err)
		}
	}

	filePath := filepath.Join(devopsDir, "devops.log")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Cannot create devops.log file", err)
	}
	defer file.Close()

	loggero := hclog.New(&hclog.LoggerOptions{
		Name:   "devops",
		Output: os.Stdout,
		Level:  hclog.Info,
	})

	loggero.Info("Loading config file...")
	c := new(model.Config)
	configBytes, err := os.ReadFile(filepath.Join(devopsDir, "config.yaml"))
	if os.IsNotExist(err) {
		// Load default
		loggero.Info("config.yaml not found, loading default configuration")
		c = &model.Config{
			Plugins: []*model.Plugin{
				{
					Name: "kubernetes",
				},
			},
		}
	} else {
		if err := yaml.Unmarshal(configBytes, c); err != nil {
			loggero.Error("failed to yaml unmarshal config file", err)
			os.Exit(1)
		}
	}

	loggerf := hclog.New(&hclog.LoggerOptions{
		Name:   "devops",
		Output: file,
		Level:  hclog.Debug,
	})

	pc, err := New(loggerf, devopsDir)
	if err != nil {
		os.Exit(1)
	}

	for _, p := range c.Plugins {
		loggero.Info(fmt.Sprintf("Loading plugin: %s", p.Name))
		kp, err := pc.GetPlugin(c.Plugins[0].Name)
		if err != nil {
			os.Exit(1)
		}
		if err := kp.StatusOK(); err != nil {
			loggero.Error("failed to load plugin", err)
			os.Exit(1)
		}
	}

	kp, err := pc.GetPlugin(c.Plugins[0].Name)
	if err != nil {
		os.Exit(1)
	}

	pCtx, err := InitPluginContext(loggerf, kp, "configmaps")
	if err != nil {
		os.Exit(1)
	}

	eventChan := make(chan model.Event, 1)
	defer close(eventChan)

	eventChan <- model.Event{ResourceType: "pods", Type: model.ResourceTypeChanged}
	app := views.NewApplication(loggerf, eventChan)

	app.SearchView.SetResourceTypes(pCtx.supportedResourceTypes)
	app.GeneralInfoView.Refresh(pCtx.generalInfo)
	app.IsolatorView.SetDefault(pCtx.defaultIsolator)

	app.PluginView.Refresh(map[string]string{"ctrl-a": kp.Name()})
	closerChan := make(chan struct{}, 1)
	defer close(closerChan)

	go func() {
		for event := range eventChan {
			loggerf.Debug(fmt.Sprintf("Received new event of type <%s> on resource <%s>, row index <%v>", event.Type, event.ResourceType, event.RowIndex))

			switch event.Type {
			case model.ReadResource:
				data := pCtx.currentResources[event.RowIndex-1]
				dd, _ := yaml.Marshal(data)
				app.SetText(string(dd))
				app.GetApp().Draw()
			case model.DeleteResource:
				event.IsolatorName = pCtx.currentIsolator
				if err := kp.ActionDeleteResource(shared.ActionDeleteResourceArgs{ResourceName: event.ResourceName, ResourceType: event.ResourceType, IsolatorName: event.IsolatorName}); err != nil {
					loggerf.Error("failed to delete resource", err)
					continue
				}
				app.SwitchToMain()

			case model.SpecificActionOccured:

				fnArgs := shared.SpecificActionArgs{
					ActionName:   event.SpecificActionName,
					ResourceName: event.ResourceName,
					ResourceType: pCtx.currentResourceType,
					IsolatorName: pCtx.currentIsolator,
				}

				loggerf.Info("Args", fnArgs)
				res, err := kp.PerformSpecificAction(fnArgs)
				if err != nil {
					loggerf.Error("failed to perform specific action on resource", err)
					continue
				}

				// logger.Info(fmt.Sprintf("Result %s", res.OutputType))
				// continue

				action := shared.SpecificAction{}
				for _, sa := range pCtx.currentSpecficActionList {
					if sa.Name == event.SpecificActionName {
						action = sa
					}
				}
				if action.Name == "" {
					continue
				}

				if action.ScrrenAction == "view" {
					stringData := res.Result.(string)
					app.SetText(stringData)
				}
				app.GetApp().Draw()

			case model.ShowModal:
				app.ViewModel()

			case model.IsolatorChanged:
				event.ResourceType = pCtx.currentResourceType
				pCtx.setCurrentIsolator(event.IsolatorName)
				syncResource(loggerf, event, kp, pCtx, app)

			case model.AddIsolator:
				if event.ResourceType != pCtx.defaultIsolatorType {
					continue
				}
				app.IsolatorView.AddAndRefreshView(event.IsolatorName)

			case model.RefreshResource:
				event.ResourceType = pCtx.currentResourceType
				syncResource(loggerf, event, kp, pCtx, app)

			case model.ResourceTypeChanged:
				// TODO: Handle wrong resource names
				if event.ResourceType == "" {
					loggerf.Debug("False invocation received, resource type is empty")
					continue
				}

				if event.ResourceType == pCtx.currentResourceType {
					loggerf.Debug("Current & new resource type are the same, Ignoring this event")
					continue
				}

				syncResource(loggerf, event, kp, pCtx, app)

				// go func() {
				// 	for {
				// 		select {
				// 		case <-closerChan:
				// 			return
				// 		case <-time.After(5 * time.Second):
				// 			eventChan <- model.Event{
				// 				Type: model.RefreshResource,
				// 			}
				// 		}
				// 	}
				// }()
			default:

			}
		}
	}()

	if err := app.Start(); err != nil {
		loggerf.Error("failed to start application", err)
		os.Exit(1)
	}
}

func syncResource(logger hclog.Logger, event model.Event, kp shared.Devops, pCtx *CurrentPluginContext, app *views.Application) {
	schema, err := kp.GetResourceTypeSchema(event.ResourceType)
	if err != nil {
		logger.Error("failed to fetch resource type schema", err)
		return
	}
	pCtx.currentSchema = schema

	resources, err := kp.GetResources(shared.GetResourcesArgs{ResourceType: event.ResourceType, IsolatorID: pCtx.currentIsolator})
	if err != nil {
		logger.Error("failed to fetch resources", err)
		return
	}
	pCtx.currentResources = resources

	table, err := transformer.GetResourceInTableFormat(&schema, resources)
	if err != nil {
		logger.Error("unable to convert resource data of type into table format", event.ResourceType, err)
		return
	}

	actions, err := kp.GetSupportedActions(event.ResourceType)
	if err != nil {
		logger.Error("unable to get supported actions of resource", event.ResourceType, err)
		return
	}

	app.ActionView.RefreshActions(actions)

	specificActions, err := kp.GetSpecficActionList(event.ResourceType)
	if err != nil {
		logger.Error("unable to get specific actions of resource", event.ResourceType, err)
		return
	}

	if event.ResourceType == pCtx.defaultIsolatorType {
		specificActions = append(specificActions, shared.SpecificAction{Name: "Use", KeyBinding: "u"})
	}
	app.SpecificActionView.RefreshActions(specificActions)
	pCtx.currentSpecficActionList = specificActions

	pCtx.currentResourceType = event.ResourceType

	logger.Debug("Removing search view")
	app.RemoveSearchView()
	logger.Debug("Refreshing table")
	app.MainView.Refresh(table)
	app.MainView.SetTitle(event.ResourceType)
	logger.Debug("Setting focus to main view")
	app.GetApp().SetFocus(app.MainView.GetView())
	app.GetApp().Draw()
	logger.Debug("Activation done")
}
