package core

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/sharadregoti/devops/common"
	"github.com/sharadregoti/devops/internal/views"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/proto"
	"github.com/sharadregoti/devops/utils/logger"
)

// Updating VS Code Server to version 1ad8d514439d5077d2b0b7ee64d2ce82a9308e5a
// Removing previous installation...
// Installing VS Code Server for x64 (1ad8d514439d5077d2b0b7ee64d2ce82a9308e5a)
// Downloading:  80%

// var release bool = false

type addHelper struct{}

func (*addHelper) SendData(a interface{}) error {
	logger.LogDebug("HI send data from core binary was called")
	return nil
}

func getPluginPath(name, devopsDir string) string {
	if common.Release {
		return fmt.Sprintf("%s/plugins/%s/%s", devopsDir, name, name)
	}
	return fmt.Sprintf("../../plugin/%s/%s/%s", name, name, name)
}

func ListPlugins() ([]string, error) {
	// Init plugins
	devopsDir := model.InitCoreDirectory()

	// Load config
	c := model.LoadConfig(devopsDir)

	checIfPluginExists(devopsDir, c)

	if len(c.Plugins) == 0 {
		return nil, fmt.Errorf("no plugins were specified in the configuration, Exitting...")
	}

	// TODO: Throw error if plugin with same name found...
	plugins := make([]string, 0)
	for _, p := range c.Plugins {
		plugins = append(plugins, p.Name)
	}

	return plugins, nil
}

func InitAndGetAuthInfo(pluginName string) (*proto.AuthInfoResponse, error) {
	devopsDir := model.InitCoreDirectory()

	pc, err := loadPlugin(logger.Loggerf, pluginName, devopsDir)
	if err != nil {
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

	kp, err := pc.GetPlugin(pluginName)
	if err != nil {
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

	auths, err := kp.GetAuthInfo()
	if err != nil {
		return nil, err
	}

	return auths, nil
}

func Start(isTest bool, authInfo *proto.AuthInfo) (*CurrentPluginContext, error) {
	// Init plugins
	devopsDir := model.InitCoreDirectory()

	// Load config
	c := model.LoadConfig(devopsDir)

	checIfPluginExists(devopsDir, c)

	_, loggerf := logger.Loggero, logger.Loggerf
	if len(c.Plugins) == 0 {
		log.Fatal("No plugins were specified in the configuration, Exitting...")
	}

	// On startup load the first plugin
	initialPlugin := c.Plugins[0]
	fmt.Printf("Loading plugin: %s\n", initialPlugin.Name)

	pc, err := loadPlugin(loggerf, initialPlugin.Name, devopsDir)
	if err != nil {
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

	kp, err := pc.GetPlugin(initialPlugin.Name)
	if err != nil {
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

	// authInfos, err := kp.GetAuthInfo()

	// if err := kp.StatusOK(); err != nil {
	// 	common.Error(loggero, fmt.Sprintf("failed to load plugin: %v", err))
	// 	time.Sleep(5 * time.Second)
	// 	os.Exit(1)
	// }

	eventChan := make(chan model.Event, 1)
	// defer close(eventChan)

	// Initiate global plugin contexts
	pCtx, err := initPluginContext(loggerf, kp, initialPlugin.Name, eventChan, pc, authInfo)
	if err != nil {
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

	if isTest {
		return pCtx, err
	}

	// Invoke default
	eventChan <- model.Event{
		ResourceType:       pCtx.defaultIsolatorType,
		Type:               string(model.ResourceTypeChanged),
		RowIndex:           0,
		ResourceName:       "",
		IsolatorName:       "",
		SpecificActionName: "",
	}

	// Start Event Loop
	StartEventLoop(eventChan, pCtx, loggerf)

	// Invoke initial event
	// Now we can invoke frontends
	return pCtx, err
}

func StartTUI() error {
	logger.InitClientLogging()

	app, err := views.NewApplication("localhost:3000")
	if err != nil {
		return err
	}

	if err := app.Start(); err != nil {
		common.Error(logger.Loggerf, fmt.Sprintf("failed to start application: %v", err))
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

	return nil
}

// func Init() {
// 	devopsDir := initCoreDirectory()
// 	file := getCoreLogFile(devopsDir)
// 	// defer file.Close()

// 	c := loadConfig(devopsDir)

// 	checIfPluginExists(devopsDir, c)

// 	loggero, loggerf := createLoggers(file)
// 	if len(c.Plugins) == 0 {
// 		log.Fatal("No plugins were specified in the configuration, Exitting...")
// 	}

// 	// On startup load the first plugin
// 	initialPlugin := c.Plugins[0]
// 	fmt.Printf("Loading plugin: %s\n", initialPlugin.Name)

// 	pc, err := loadPlugin(loggerf, initialPlugin.Name, devopsDir)
// 	if err != nil {
// 		time.Sleep(5 * time.Second)
// 		os.Exit(1)
// 	}

// 	kp, err := pc.GetPlugin(initialPlugin.Name)
// 	if err != nil {
// 		time.Sleep(5 * time.Second)
// 		os.Exit(1)
// 	}

// 	if err := kp.StatusOK(); err != nil {
// 		common.Error(loggero, fmt.Sprintf("failed to load plugin: %v", err))
// 		time.Sleep(5 * time.Second)
// 		os.Exit(1)
// 	}

// 	eventChan := make(chan model.Event, 1)
// 	defer close(eventChan)

// 	app := views.NewApplication(loggerf, eventChan)

// 	// Initiate global plugin contexts
// 	pCtx, err := initPluginContext(loggerf, kp, app, initialPlugin.Name)
// 	if err != nil {
// 		time.Sleep(5 * time.Second)
// 		os.Exit(1)
// 	}

// 	pluginNames := []string{}
// 	for _, p := range c.Plugins {
// 		pluginNames = append(pluginNames, p.Name)
// 	}
// 	app.PluginView.Refresh(pluginNames)

// 	eventChan <- model.Event{
// 		ResourceType:       pCtx.defaultIsolatorType,
// 		Type:               model.ResourceTypeChanged,
// 		RowIndex:           0,
// 		ResourceName:       "",
// 		IsolatorName:       "",
// 		SpecificActionName: "",
// 	}

// 	// closerChan := make(chan struct{}, 1)
// 	streamCloserChan := make(chan struct{}, 1)
// 	defer close(streamCloserChan)
// 	// defer close(closerChan)
// 	isStreamingOn := false
// 	// invokingFirstTime := true

// 	// nestedResourceLevel := -1

// 	// go func() {
// 	// 	http.HandleFunc("/data", pCtx.handle)

// 	// 	// Enable CORS for all routes
// 	// 	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 	// 		w.Header().Set("Access-Control-Allow-Methods", "*")
// 	// 		w.Header().Set("Access-Control-Allow-Headers", "*")
// 	// 	}))

// 	// 	http.ListenAndServe(":9000", nil)
// 	// }()

// 	if err := app.Start(); err != nil {
// 		common.Error(loggerf, fmt.Sprintf("failed to start application: %v", err))
// 		time.Sleep(5 * time.Second)
// 		os.Exit(1)
// 	}
// }

func StartEventLoop(eventChan chan model.Event, pCtx *CurrentPluginContext, loggerf hclog.Logger) {
	go func() {
		defer logger.LogDebug("Exiting start event loop")

		streamCloserChan := make(chan struct{}, 1)
		defer close(streamCloserChan)
		isStreamingOn := false

		for event := range eventChan {
			loggerf.Debug("\n")
			loggerf.Debug("====================================================================================================================")
			data, _ := json.MarshalIndent(event, " ", " ")
			loggerf.Debug(fmt.Sprintf("Received event %v", string(data)))

			switch event.Type {

			case string(model.CloseEventLoop):
				return

			case string(model.ViewNestedResource):
				if !pCtx.getCurrentResource().currentSchema.Nesting.IsNested {
					loggerf.Debug("False alarm for nested resource occured")
				}
				pCtx.syncResource(event)

			case model.NestBack:
				loggerf.Debug("Nest", pCtx.currentNestedResourceLevel)
				pCtx.viewBackwardNestResource(event)

			case model.PluginChanged:
				// pc.Close()
				// pc, err = loadPlugin(loggerf, event.PluginName, devopsDir)
				// if err != nil {
				// 	continue
				// }
				// kp, err = pc.GetPlugin(event.PluginName)
				// if err != nil {
				// 	continue
				// }
				// pCtx, err = initPluginContext(loggerf, kp, app, event.PluginName)
				// if err != nil {
				// 	continue
				// }
				// // Inovked event
				// eventChan <- model.Event{
				// 	ResourceType:       pCtx.defaultIsolatorType,
				// 	Type:               model.ResourceTypeChanged,
				// 	RowIndex:           0,
				// 	ResourceName:       "",
				// 	IsolatorName:       "",
				// 	SpecificActionName: "",
				// }

			case string(model.Close):
				if isStreamingOn {
					// fnArgs := shared.SpecificActionArgs{
					// 	ActionName: "close",
					// }
					// _, err := kp.PerformSpecificAction(fnArgs)
					// if err != nil {
					// 	common.Error(loggerf, fmt.Sprintf("failed to perform close action on resource: %v", err))
					// 	app.SetFlashText(err.Error())
					// 	continue
					// }
					// streamCloserChan <- struct{}{}
					// isStreamingOn = false
				}

			// Specfic action on resource
			// case model.SpecificActionOccured:
			// 	action := model.SpecificActions{}
			// 	rs := pCtx.getCurrentResource()
			// 	for _, sa := range rs.currentSchema.SpecificActions {
			// 		if sa.Name == event.SpecificActionName {
			// 			action = sa
			// 		}
			// 	}
			// 	if action.Name == "" {
			// 		continue
			// 	}

			// 	specArgs := map[string]interface{}{}
			// 	if len(action.Args) > 0 {
			// 		prev := pCtx.getPreviousResource()
			// 		// loggerf.Debug(fmt.Sprintf("DDDD %v, %v", prev.currentResources[prev.tableRowNumber-1], action.Args))
			// 		specArgs = transformer.GetArgs(prev.currentResources[prev.tableRowNumber-1], action.Args)
			// 	}

			// 	fnArgs := shared.SpecificActionArgs{
			// 		ActionName:   event.SpecificActionName,
			// 		ResourceName: event.ResourceName,
			// 		ResourceType: rs.currentResourceType,
			// 		IsolatorName: pCtx.currentIsolator,
			// 		Args:         specArgs,
			// 	}

			// 	loggerf.Debug("Specific action args", fnArgs)
			// 	res, err := kp.PerformSpecificAction(fnArgs)
			// 	if err != nil {
			// 		common.Error(loggerf, fmt.Sprintf("failed to perform specific action on resource: %v", err))
			// 		app.SetFlashText(err.Error())
			// 		continue
			// 	}

			// 	if action.ScrrenAction == "view" {
			// 		switch action.OutputType {
			// 		case "string":
			// 			stringData := res.Result.(string)
			// 			app.SetTextAndSwitchView(stringData)
			// 			app.GetApp().Draw()
			// 			continue

			// 		case "stream":
			// 			newReader := bufio.NewReader(pc.GetStdoutReader())
			// 			if newReader == nil {
			// 				loggerf.Debug("Failed to create a reader for stream")
			// 				continue
			// 			}

			// 			app.SetTextAndSwitchView("")
			// 			app.GetApp().Draw()
			// 			go func() {
			// 				isStreamingOn = true
			// 				for {
			// 					select {
			// 					case <-streamCloserChan:
			// 						loggerf.Debug("Streamer go routine closed")
			// 						return
			// 					default:
			// 						data, _, err := newReader.ReadLine()
			// 						// data, err := newReader.ReadString(byte('\n'))
			// 						if err == io.EOF {
			// 							common.Error(loggerf, fmt.Sprintf("EOF received while streaming: %v", err))
			// 							break
			// 						} else if err != nil {
			// 							common.Error(loggerf, fmt.Sprintf("failed to stream data: %v", err))
			// 							break
			// 						}
			// 						w := app.Ta.BatchWriter()
			// 						fmt.Fprintln(w, string(data))
			// 						w.Close()
			// 						app.GetApp().Draw()

			// 					}
			// 				}
			// 			}()
			// 		}

			// 	}

			case model.ShowModal:
				// app.ViewModel(event.ResourceType, event.ResourceName)
				// app.GetApp().Draw()

			// Isolator actions
			case string(model.IsolatorChanged):
				event.ResourceType = pCtx.getCurrentResource().currentResourceType
				// pCtx.setCurrentIsolator(event.IsolatorName)
				// pCtx.resetToParentResource()
				pCtx.syncResource(event)

			// case string(model.AddIsolator):
			// 	if event.ResourceType != pCtx.defaultIsolatorType {
			// 		continue
			// 	}
			// pCtx.clearNestedResource()
			// app.IsolatorView.AddAndRefreshView(event.IsolatorName)
			// app.GetApp().Draw()

			// Resource Actions
			case string(model.ReadResource):
				// -1 because, table data index starts with 1 and on
				// The data stored in array starts with 0 index, So 1 table row maps with 0 of array row
				// data := pCtx.getCurrentResource().currentResources[event.RowIndex-1]
				// // TODO: Take format type from plugin
				// dd, _ := yaml.Marshal(data)
				// app.SetTextAndSwitchView(string(dd))
				// app.GetApp().Draw()

			case model.DeleteResource:
				// if err := kp.ActionDeleteResource(shared.ActionDeleteResourceArgs{ResourceName: event.ResourceName, ResourceType: pCtx.getCurrentResource().currentResourceType, IsolatorName: pCtx.currentIsolator}); err != nil {
				// 	common.Error(loggerf, fmt.Sprintf("failed to delete resource: %v", err))
				// 	continue
				// }
				// app.SwitchToMain()

			case string(model.RefreshResource):
				pCtx.syncResource(event)

			case string(model.ResourceTypeChanged):
				// TODO: Handle wrong resource names
				if event.ResourceType == "" {
					loggerf.Debug("False invocation received, resource type is empty")
					continue
				}

				// if pCtx.areWeViewingNestedResource() {
				// 	loggerf.Debug("Nested resource activated")
				// 	pCtx.syncNestResource(event)
				// 	continue
				// }

				if rs := pCtx.getCurrentResource(); rs != nil && event.ResourceType == rs.currentResourceType {
					loggerf.Debug("Current & new resource type are the same, Ignoring this event")
					continue
				}

				// if invokingFirstTime {
				// 	invokingFirstTime = false
				// } else {
				// 	closerChan <- struct{}{}
				// 	invokingFirstTime = true
				// }

				pCtx.syncResource(event)

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
}
