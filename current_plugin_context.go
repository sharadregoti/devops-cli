package core

import (
	"strings"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/sharadregoti/devops/internal/transformer"
	"github.com/sharadregoti/devops/internal/views"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/shared"
)

type CurrentPluginContext struct {
	currentPluginName string

	generalInfo         map[string]string
	defaultIsolator     string
	defaultIsolatorType string
	currentIsolator     string

	supportedResourceTypes []string
	currentResourceType    string

	currentResources []interface{}
	currentSchema    model.ResourceTransfomer

	currentSpecficActionList []shared.SpecificAction
}

func initPluginContext(logger hclog.Logger, p shared.Devops, app *views.Application, pluginName string) (*CurrentPluginContext, error) {
	// Get known changing infor
	info, err := p.GetGeneralInfo()
	if err != nil {
		logger.Error("initial resource fetching failed", err)
		return nil, err
	}
	isolator, err := p.GetDefaultResourceIsolator()
	if err != nil {
		logger.Error("initial resource fetching failed", err)
		return nil, err
	}
	resourceTypeList, err := p.GetResourceTypeList()
	if err != nil {
		logger.Error("initial resource fetching failed", err)
		return nil, err
	}

	defaultIsolatorType, err := p.GetResourceIsolatorType()
	if err != nil {
		logger.Error("initial resource fetching failed", err)
		return nil, err
	}

	app.SearchView.SetResourceTypes(resourceTypeList)
	app.GeneralInfoView.Refresh(info)
	app.IsolatorView.SetDefault(isolator)
	app.IsolatorView.SetTitle(strings.Title(defaultIsolatorType))

	return &CurrentPluginContext{
		currentPluginName:      pluginName,
		generalInfo:            info,
		defaultIsolatorType:    defaultIsolatorType,
		currentIsolator:        isolator,
		defaultIsolator:        isolator,
		currentResourceType:    "",
		currentResources:       make([]interface{}, 0),
		supportedResourceTypes: resourceTypeList,
		currentSchema:          model.ResourceTransfomer{},
	}, nil
}

func (c *CurrentPluginContext) setCurrentIsolator(isolatorName string) {
	c.currentIsolator = isolatorName
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

	if len(resources) == 0 {
		app.SetFlashText(" !!! No resources exists ")
	}

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
	// if event.Type != model.RefreshResource {
	app.GetApp().Draw()
	// }
	logger.Debug("Activation done")
}
