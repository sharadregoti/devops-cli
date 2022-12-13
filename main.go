// Demo code for the Table primitive.
package main

import (
	"os"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/sharadregoti/devops/internal/transformer"
	"github.com/sharadregoti/devops/internal/views"
	"github.com/sharadregoti/devops/model"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "devops",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	logger.Info("Starting the process...")

	p := installPlugins()
	c := make(chan model.Event, 1)
	defer close(c)
	app := views.NewApplication(c)

	schema := p.GetResourceTypeSchema("pods")
	table := transformer.GetResourceInTableFormat(&schema, p.GetResources("pods"))
	app.MainView.Refresh(table)
	app.SearchView.SetResourceTypes(p.GetResourceTypeList())
	app.GeneralInfoView.Refresh(p.GetGeneralInfo())
	app.IsolatorView.Refresh(map[string]string{"ctrl-d": p.GetDefaultResourceIsolator()})
	app.PluginView.Refresh(map[string]string{"ctrl-a": p.Name(), "ctrl-b": p.Name()})

	go func() {
		for e := range c {
			schema := p.GetResourceTypeSchema(e.ResourceType)
			table := transformer.GetResourceInTableFormat(&schema, p.GetResources(e.ResourceType))
			app.MainView.Refresh(table)
			app.SearchView.SetResourceTypes(p.GetResourceTypeList())
			app.GeneralInfoView.Refresh(p.GetGeneralInfo())
			app.IsolatorView.Refresh(map[string]string{"ctrl-d": p.GetDefaultResourceIsolator()})
			app.PluginView.Refresh(map[string]string{"ctrl-a": p.Name(), "ctrl-b": p.Name()})
			app.GetApp().SetFocus(app.MainView.GetView())
			// app.
		}
	}()

	// for r := range p.WatchResources("pods") {
	// 	table := transformer.GetResourceInTableFormat(&schema, []interface{}{r.Result})
	// 	app.MainView.Refresh(table)
	// }

	app.Start()
}

// Event listener chaiye, that take appropiate actions
// What are the events
// 1. Search action completed
// 	-> Rerender table

// Event invoker
// 1. Done Complete Function of Search
