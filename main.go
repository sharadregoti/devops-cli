package core

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ghodss/yaml"
	"github.com/sharadregoti/devops/internal/views"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/shared"
)

// Updating VS Code Server to version 1ad8d514439d5077d2b0b7ee64d2ce82a9308e5a
// Removing previous installation...
// Installing VS Code Server for x64 (1ad8d514439d5077d2b0b7ee64d2ce82a9308e5a)
// Downloading:  80%

var release bool = false

func getPluginPath(name, devopsDir string) string {
	if release {
		return fmt.Sprintf("%s/plugins/%s/%s", devopsDir, name, name)
	}
	return fmt.Sprintf("../../plugin/%s/%s/%s", name, name, name)
}

func Init() {
	devopsDir := initCoreDirectory()
	file := getCoreLogFile(devopsDir)
	defer file.Close()

	c := loadConfig(devopsDir)

	checIfPluginExists(devopsDir, c)

	loggero, loggerf := createLoggers(file)
	if len(c.Plugins) == 0 {
		log.Fatal("No plugins were specified in the configuration, Exitting...")
	}

	// On startup load the first plugin
	initialPlugin := c.Plugins[0]
	fmt.Printf("Loading plugin: %s\n", initialPlugin.Name)

	pc, err := loadPlugin(loggerf, initialPlugin.Name, devopsDir)
	if err != nil {
		os.Exit(1)
	}

	kp, err := pc.GetPlugin(c.Plugins[0].Name)
	if err != nil {
		os.Exit(1)
	}

	if err := kp.StatusOK(); err != nil {
		loggero.Error("failed to load plugin", err)
		os.Exit(1)
	}

	eventChan := make(chan model.Event, 1)
	defer close(eventChan)

	app := views.NewApplication(loggerf, eventChan)

	// Initiate global plugin contexts
	pCtx, err := initPluginContext(loggerf, kp, app, initialPlugin.Name)
	if err != nil {
		os.Exit(1)
	}

	eventChan <- model.Event{
		ResourceType:       pCtx.defaultIsolatorType,
		Type:               model.ResourceTypeChanged,
		RowIndex:           0,
		ResourceName:       "",
		IsolatorName:       "",
		SpecificActionName: "",
	}

	closerChan := make(chan struct{}, 1)
	streamCloserChan := make(chan struct{}, 1)
	defer close(streamCloserChan)
	defer close(closerChan)
	isStreamingOn := false
	// invokingFirstTime := true

	go func() {
		for event := range eventChan {
			loggerf.Debug(fmt.Sprintf("Received new event of type <%s> on resource <%s>, row index <%v>", event.Type, event.ResourceType, event.RowIndex))

			switch event.Type {

			case model.Close:
				if isStreamingOn {
					fnArgs := shared.SpecificActionArgs{
						ActionName: "close",
					}
					_, err := kp.PerformSpecificAction(fnArgs)
					if err != nil {
						loggerf.Error("failed to perform close action on resource", err)
						continue
					}
					streamCloserChan <- struct{}{}
					isStreamingOn = false
				}

			// Specfic action on resource
			case model.SpecificActionOccured:
				fnArgs := shared.SpecificActionArgs{
					ActionName:   event.SpecificActionName,
					ResourceName: event.ResourceName,
					ResourceType: pCtx.currentResourceType,
					IsolatorName: pCtx.currentIsolator,
				}

				loggerf.Debug("Specific action args", fnArgs)
				res, err := kp.PerformSpecificAction(fnArgs)
				if err != nil {
					loggerf.Error("failed to perform specific action on resource", err)
					continue
				}

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
					switch action.OutputType {
					case "string":
						stringData := res.Result.(string)
						app.SetTextAndSwitchView(stringData)
						app.GetApp().Draw()
						continue

					case "stream":
						newReader := bufio.NewReader(reader)
						if newReader == nil {
							loggerf.Debug("Failed to create a reader for stream")
							continue
						}

						app.SetTextAndSwitchView("")
						app.GetApp().Draw()
						go func() {
							isStreamingOn = true
							for {
								select {
								case <-streamCloserChan:
									loggerf.Debug("Streamer go routine closed")
									return
								default:
									data, _, err := newReader.ReadLine()
									// data, err := newReader.ReadString(byte('\n'))
									if err == io.EOF {
										loggerf.Error("EOF received while streaming", err)
										break
									} else if err != nil {
										loggerf.Error("failed to stream data", err)
										break
									}
									w := app.Ta.BatchWriter()
									fmt.Fprintln(w, string(data))
									w.Close()
									app.GetApp().Draw()

								}
							}
						}()
					}

				}

			case model.ShowModal:
				app.ViewModel(event.ResourceType, event.ResourceName)

			// Isolator actions
			case model.IsolatorChanged:
				event.ResourceType = pCtx.currentResourceType
				pCtx.setCurrentIsolator(event.IsolatorName)
				syncResource(loggerf, event, kp, pCtx, app)

			case model.AddIsolator:
				if event.ResourceType != pCtx.defaultIsolatorType {
					continue
				}
				app.IsolatorView.AddAndRefreshView(event.IsolatorName)
				app.GetApp().Draw()

			// Resource Actions
			case model.ReadResource:
				// -1 because, table data index starts with 1 and on
				// The data stored in array starts with 0 index, So 1 table row maps with 0 of array row
				data := pCtx.currentResources[event.RowIndex-1]
				// TODO: Take format type from plugin
				dd, _ := yaml.Marshal(data)
				app.SetTextAndSwitchView(string(dd))
				app.GetApp().Draw()

			case model.DeleteResource:
				if err := kp.ActionDeleteResource(shared.ActionDeleteResourceArgs{ResourceName: event.ResourceName, ResourceType: event.ResourceType, IsolatorName: event.IsolatorName}); err != nil {
					loggerf.Error("failed to delete resource", err)
					continue
				}
				app.SwitchToMain()

			case model.RefreshResource:
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

				// if invokingFirstTime {
				// 	invokingFirstTime = false
				// } else {
				// 	closerChan <- struct{}{}
				// 	invokingFirstTime = true
				// }

				syncResource(loggerf, event, kp, pCtx, app)

				// go func() {
				// 	for {
				// 		select {
				// 		case <-closerChan:
				// 			loggerf.Debug("Closing previous refresh routine")
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
