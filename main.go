// Demo code for the Table primitive.
package main

import (
	"github.com/sharadregoti/devops/internal/transformer"
	"github.com/sharadregoti/devops/internal/views"
)

func main() {

	// p := plugin.New()

	p := initPlugin()
	app := views.NewApplication()

	schema := p.GetResourceTypeSchema("pods")

	table := transformer.GetResourceInTableFormat(&schema, p.GetResources("pods"))

	app.MainView.Refresh(table)
	app.GeneralInfoView.Refresh(p.GetGeneralInfo())
	app.IsolatorView.Refresh(map[string]string{"ctrl-d": p.GetDefaultResourceIsolator()})
	app.PluginView.Refresh(map[string]string{"ctrl-a": p.Name(), "ctrl-b": p.Name()})

	app.Start()
}
