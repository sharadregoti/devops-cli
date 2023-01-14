package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/sharadregoti/devops/common"
	"github.com/sharadregoti/devops/internal/transformer"
	"github.com/sharadregoti/devops/internal/views"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/shared"
	"github.com/sharadregoti/devops/utils"
)

// Node represents a node of linked list
type Node struct {
	value int
	next  *Node
}

// LinkedList represents a linked list
type LinkedList struct {
	head *Node
	len  int
}

type CurrentPluginContext struct {
	logger hclog.Logger

	currentPluginName string

	// This field indicates current nest level
	// The value corresponds to the
	currentNestedResourceLevel int
	// This field holds the nested resources of the parent resource

	plugin shared.Devops

	appView *views.Application

	generalInfo         map[string]string
	defaultIsolator     string
	defaultIsolatorType string
	currentIsolator     string

	supportedResourceTypes []string
	// currentResourceType    string

	// currentResources []interface{}
	// currentSchema    model.ResourceTransfomer

	currentGenericActions shared.GenericActions
	// currentSpecficActionList []shared.SpecificAction

	tableStack tableStack
}

type tableStack []*resourceStack

func (t *tableStack) length() int {
	return len(*t)
}

func (t *tableStack) upsert(index int, r resourceStack) {
	for i, _ := range *t {
		if i == index {
			(*t)[0] = &r
			return
		}
	}

	// Add if does not exists
	*t = append(*t, &r)
}

func (t *tableStack) resetToParentResource() {
	// Only get the first element
	if len(*t) == 0 {
		return
	}
	*t = (*t)[0:1]
}

type resourceStack struct {
	tableRowNumber           int
	nextResourceArgs         []map[string]interface{}
	currentResourceType      string
	currentResources         []interface{}
	currentSchema            model.ResourceTransfomer
	currentSpecficActionList []shared.SpecificAction
}

type nestedResurce struct {
	nextResourceArgs         []map[string]interface{}
	currentIsolator          string
	currentResourceType      string
	currentResources         []interface{}
	currentSchema            model.ResourceTransfomer
	currentSpecficActionList []shared.SpecificAction
}

func initPluginContext(logger hclog.Logger, p shared.Devops, app *views.Application, pluginName string) (*CurrentPluginContext, error) {
	// Get known changing infor
	info, err := p.GetGeneralInfo()
	if err != nil {
		common.Error(logger, fmt.Sprintf("initial resource fetching failed: %v", err))
		return nil, err
	}
	isolator, err := p.GetDefaultResourceIsolator()
	if err != nil {
		common.Error(logger, fmt.Sprintf("initial resource fetching failed: %v", err))
		return nil, err
	}
	resourceTypeList, err := p.GetResourceTypeList()
	if err != nil {
		common.Error(logger, fmt.Sprintf("initial resource fetching failed: %v", err))
		return nil, err
	}

	defaultIsolatorType, err := p.GetResourceIsolatorType()
	if err != nil {
		common.Error(logger, fmt.Sprintf("initial resource fetching failed: %v", err))
		return nil, err
	}

	app.SearchView.SetResourceTypes(resourceTypeList)
	app.GeneralInfoView.Refresh(info)
	app.IsolatorView.SetDefault(isolator)
	app.IsolatorView.SetTitle(strings.Title(defaultIsolatorType))

	return &CurrentPluginContext{
		logger:              logger,
		currentPluginName:   pluginName,
		plugin:              p,
		appView:             app,
		generalInfo:         info,
		defaultIsolatorType: defaultIsolatorType,
		currentIsolator:     isolator,
		defaultIsolator:     isolator,
		// currentResourceType:        "",
		// currentResources:           make([]interface{}, 0),
		supportedResourceTypes: resourceTypeList,
		// currentSchema:              model.ResourceTransfomer{},
		currentNestedResourceLevel: 0,
		tableStack:                 make([]*resourceStack, 0),
	}, nil
}

func (c *CurrentPluginContext) resetToParentResource() {
	c.tableStack.resetToParentResource()
	c.currentNestedResourceLevel = 0
}

func (c *CurrentPluginContext) setCurrentIsolator(isolatorName string) {
	c.currentIsolator = isolatorName
}

// func (c *CurrentPluginContext) getSpecificActionList() []shared.SpecificAction {
// 	if c.areWeViewingNestedResource() {
// 		nest := c.nestedResources[c.currentNestedResourceLevel-1]
// 		return nest.currentSpecficActionList
// 	}

// 	return c.currentSpecficActionList
// }

func (c *CurrentPluginContext) areWeViewingNestedResource() bool {
	return c.currentNestedResourceLevel > 0
}

func (c *CurrentPluginContext) viewBackwardNestResource(event model.Event) {
	if c.currentNestedResourceLevel-1 == -1 {
		return
	}
	c.currentNestedResourceLevel--
	if c.currentNestedResourceLevel == 0 {
		c.resetToParentResource()
	}
	c.setAppView()
}

func (c *CurrentPluginContext) getCurrentResource() *resourceStack {
	if c.tableStack.length() == 0 {
		return nil
	}
	return c.tableStack[c.currentNestedResourceLevel]
}

func (c *CurrentPluginContext) getPreviousResource() *resourceStack {
	if c.tableStack.length() == 0 {
		return nil
	}
	return c.tableStack[c.currentNestedResourceLevel-1]
}

func (c *CurrentPluginContext) syncResource(event model.Event) {
	var rs *resourceStack
	var resourceLevel int
	fnArgs := map[string]interface{}{}
	if event.Type == model.ViewNestedResource {
		rs = &resourceStack{
			currentResourceType: c.getCurrentResource().currentSchema.Nesting.ResourceType,
		}
		resourceLevel = c.currentNestedResourceLevel + 1
		fnArgs = c.getCurrentResource().nextResourceArgs[event.RowIndex-1]
		c.getCurrentResource().tableRowNumber = event.RowIndex
	} else if event.Type == model.ReadResource {
		rs = c.getCurrentResource()
	} else if event.Type == model.ResourceTypeChanged || event.Type == model.RefreshResource {
		c.resetToParentResource()
		rs = &resourceStack{
			currentResourceType: event.ResourceType,
		}
	} else if event.Type == model.IsolatorChanged {
		c.setCurrentIsolator(event.IsolatorName)
		c.resetToParentResource()
		rs = &resourceStack{
			currentResourceType: event.ResourceType,
		}
	} else {
		return
	}

	schema, err := c.plugin.GetResourceTypeSchema(rs.currentResourceType)
	if err != nil {
		common.Error(c.logger, fmt.Sprintf("failed to fetch resource type schema: %v", err))
		return
	}

	var resources []interface{}
	// TODO: Remove enent type condition from here
	// if parent := c.getCurrentResource(); event.RowIndex > 0 && parent != nil && parent.currentSchema.Nesting.IsSelfContainedInParent {
	if parent := c.getCurrentResource(); event.RowIndex > 0 && parent != nil && schema.Nesting.IsSelfContainedInParent {
		c.logger.Debug("Getting nested resource from parent")
		resources, err = transformer.GetSelfContainedResource(schema.Nesting.ParentDataPaths, parent.currentResources[event.RowIndex-1])
		if err != nil {
			common.Error(c.logger, err.Error())
			c.appView.SetFlashText(err.Error())
			return
		}
	} else {
		resources, err = c.plugin.GetResources(shared.GetResourcesArgs{ResourceType: rs.currentResourceType, IsolatorID: c.currentIsolator, Args: fnArgs})
		if err != nil {
			common.Error(c.logger, fmt.Sprintf("failed to fetch resources: %v", err))
			c.appView.SetFlashText(err.Error())
			return
		}
	}

	if len(resources) == 0 {
		c.appView.SetFlashText("!!! No resources exists ")
	}

	// table, err := transformer.GetResourceInTableFormat(&schema, resources)
	// if err != nil {
	// 	common.Error(logger, "unable to convert resource data of type into table format", event.ResourceType, err)
	// 	return
	// }

	actions, err := c.plugin.GetSupportedActions(rs.currentResourceType)
	if err != nil {
		common.Error(c.logger, fmt.Sprintf("unable to get supported actions of resource: %v, %v", rs.currentResourceType, err))
		c.appView.SetFlashText(err.Error())
		return
	}
	c.currentGenericActions = actions

	specificActions, err := c.plugin.GetSpecficActionList(rs.currentResourceType)
	if err != nil {
		common.Error(c.logger, fmt.Sprintf("unable to get specific actions of resource: %v, %v", rs.currentResourceType, err))
		c.appView.SetFlashText(err.Error())
		return
	}

	if rs.currentResourceType == c.defaultIsolatorType {
		specificActions = append(specificActions, shared.SpecificAction{Name: "Use", KeyBinding: "u"})
	}
	// c.currentSpecficActionList = specificActions

	// c.currentResourceType = event.ResourceType
	c.currentNestedResourceLevel = resourceLevel

	c.tableStack.upsert(resourceLevel, resourceStack{
		currentResourceType:      rs.currentResourceType,
		currentResources:         resources,
		currentSchema:            schema,
		currentSpecficActionList: specificActions,
		nextResourceArgs:         []map[string]interface{}{},
	})

	c.setAppView()

	// c.appView.SpecificActionView.RefreshActions(specificActions)
	// c.appView.ActionView.RefreshActions(actions)
	// c.appView.RemoveSearchView()
	// c.appView.MainView.Refresh(table)
	// c.appView.MainView.SetTitle(event.ResourceType)
	// c.appView.GetApp().SetFocus(c.appView.MainView.GetView())
	// c.appView.GetApp().Draw()
}

func (c *CurrentPluginContext) handle(w http.ResponseWriter, req *http.Request) {
	rs := c.tableStack[c.currentNestedResourceLevel]
	table, _, err := transformer.GetResourceInTableFormat(c.logger, &rs.currentSchema, rs.currentResources)
	if err != nil {
		common.Error(c.logger, fmt.Sprintf("unable to convert resource data of type into table format: %v, %v", rs.currentResourceType, err))
		c.appView.SetFlashText(err.Error())
		return
	}
	temp := map[string]interface{}{
		"specificActions": c.getCurrentResource().currentSpecficActionList,
		"resources": map[string]interface{}{
			"headers": table[0],
			"columns": table[1:],
		},
		"schema":       c.getCurrentResource().currentSchema,
		"resourceType": c.getCurrentResource().currentResourceType,
	}
	// currentResourceType      string
	// currentResources         []interface{}
	// currentSchema            model.ResourceTransfomer

	SendResponse(context.Background(), w, 200, temp)
}

// SendResponse sends an http response
func SendResponse(ctx context.Context, w http.ResponseWriter, statusCode int, body interface{}) error {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(body)
}

func (c *CurrentPluginContext) setAppView() {

	rs := c.tableStack[c.currentNestedResourceLevel]
	table, nestArgs, err := transformer.GetResourceInTableFormat(c.logger, &rs.currentSchema, rs.currentResources)
	if err != nil {
		common.Error(c.logger, fmt.Sprintf("unable to convert resource data of type into table format: %v, %v", rs.currentResourceType, err))
		c.appView.SetFlashText(err.Error())
		return
	}

	if rs.currentSchema.Nesting.IsNested {
		rs.nextResourceArgs = nestArgs
	}

	c.appView.SpecificActionView.RefreshActions(rs.currentSpecficActionList)
	c.appView.ActionView.RefreshActions(c.currentGenericActions)
	c.appView.ActionView.EnableNesting(rs.currentSchema.Nesting.IsNested)
	c.appView.RemoveSearchView()
	c.appView.MainView.Refresh(table, rs.tableRowNumber)
	c.appView.MainView.SetTitle(utils.GetTableTitle(rs.currentResourceType, len(rs.currentResources)))
	c.appView.GetApp().SetFocus(c.appView.MainView.GetView())
	c.appView.GetApp().Draw()
}

// func (c *CurrentPluginContext) syncNestResource(row int, eventType string) {
// 	nest := c.nestedResources[c.currentNestedResourceLevel]
// 	schema, err := c.plugin.GetResourceTypeSchema(nest.currentResourceType)
// 	if err != nil {
// 		common.Error(c.logger, fmt.Sprintf("failed to fetch resource type schema: %v", err))
// 		return
// 	}
// 	nest.currentSchema = schema

// 	var resources []interface{}
// 	if eventType == model.NestBack {
// 		resources = nest.currentResources
// 	} else {
// 		if c.currentNestedResourceLevel > 0 {
// 			parentNest := c.nestedResources[c.currentNestedResourceLevel-1]
// 			if parentNest.currentSchema.Nesting.IsSelfContainedInParent {
// 				parentNest := c.nestedResources[c.currentNestedResourceLevel-1]
// 				resources, err = transformer.GetSelfContainedResource(&parentNest.currentSchema, parentNest.currentResources[row-1])
// 				if err != nil {
// 					common.Error(c.logger, err.Error())
// 					return
// 				}
// 			}

// 		} else {
// 			// Row will alwyas be greater than 0
// 			// But array index start 0, And a row == 1 indicates, 0 in the index
// 			// So we are doing -1
// 			fnArgs := nest.nextResourceArgs[row-1]

// 			resources, err = c.plugin.GetResources(shared.GetResourcesArgs{ResourceType: nest.currentResourceType, IsolatorID: nest.currentIsolator, Args: fnArgs})
// 			if err != nil {
// 				common.Error(c.logger, fmt.Sprintf("failed to fetch resources: %v", err))
// 				return
// 			}
// 		}
// 		nest.currentResources = resources
// 	}

// 	if len(resources) == 0 {
// 		c.appView.SetFlashText("!!! No resources exists ")
// 	}

// 	c.logger.Debug(fmt.Sprintf("Length of nest resource %v", len(nest.currentResources)))
// 	table, nestArgs, err := transformer.GetResourceInTableFormat(&nest.currentSchema, nest.currentResources)
// 	if err != nil {
// 		common.Error(c.logger, fmt.Sprintf("unable to convert resource data of type into table format: %v, %v", nest.currentResourceType, err))
// 		return
// 	}

// 	specificActions, err := c.plugin.GetSpecficActionList(nest.currentResourceType)
// 	if err != nil {
// 		common.Error(c.logger, fmt.Sprintf("unable to get specific actions of resource: %v, %v", nest.currentResourceType, err))
// 		return
// 	}

// 	if nest.currentResourceType == c.defaultIsolatorType {
// 		specificActions = append(specificActions, shared.SpecificAction{Name: "Use", KeyBinding: "u"})
// 	}
// 	nest.currentSpecficActionList = specificActions

// 	// nest.currentResourceType = nest.currentResourceType

// 	if nest.currentSchema.Nesting.IsNested {
// 		c.logger.Debug("Super nesting is enabled", nest.currentSchema.Nesting.ResourceType)
// 		c.logger.Debug("Data", len(c.nestedResources) >= c.currentNestedResourceLevel, len(c.nestedResources), c.currentNestedResourceLevel)
// 		// TODO: Fix remvoe entType from this function

// 		if eventType != model.NestBack && len(c.nestedResources) >= c.currentNestedResourceLevel {
// 			c.logger.Debug("Incremening")
// 			c.currentNestedResourceLevel++
// 			c.nestedResources = append(c.nestedResources, &nestedResurce{
// 				nextResourceArgs:    nestArgs,
// 				currentResourceType: nest.currentSchema.Nesting.ResourceType,
// 				currentIsolator:     c.currentIsolator,
// 			})
// 		}
// 	}

// 	c.appView.SpecificActionView.RefreshActions(nest.currentSpecficActionList)
// 	c.appView.ActionView.RefreshActions(c.currentGenericActions)
// 	c.appView.ActionView.EnableNesting(nest.currentSchema.Nesting.IsNested)
// 	c.appView.RemoveSearchView()
// 	c.appView.MainView.Refresh(table)
// 	c.appView.MainView.SetTitle(nest.currentResourceType)
// 	c.appView.GetApp().SetFocus(c.appView.MainView.GetView())
// 	c.appView.GetApp().Draw()
// }
