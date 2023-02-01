package core

import (
	"fmt"

	"github.com/sharadregoti/devops/internal/transformer"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/shared"
	"github.com/sharadregoti/devops/utils"
)

func (c *CurrentPluginContext) ReadSync(e model.Event) error {
	resources, err := c.plugin.GetResources(shared.GetResourcesArgs{
		ResourceName: "", //Ingore resource name during sync
		ResourceType: e.ResourceType,
		IsolatorID:   e.IsolatorName,
	})
	if err != nil {
		return err
	}

	schema, err := c.plugin.GetResourceTypeSchema(e.ResourceType)
	if err != nil {
		return err
	}

	tableData, _, err := transformer.GetResourceInTableFormat(c.logger, &schema, resources)
	if err != nil {
		return err
	}

	specificActions, err := c.plugin.GetSpecficActionList(e.ResourceType)
	if err != nil {
		return err
	}

	c.SendMessage(model.WebsocketResponse{
		TableName:       utils.GetTableTitle(e.ResourceType, len(resources)),
		Data:            tableData,
		SpecificActions: specificActions,
	})
	return nil
}

func (c *CurrentPluginContext) GetSpecficActionList(e model.Event) ([]shared.Action, error) {
	return c.plugin.GetSpecficActionList(e.ResourceType)
}

func (c *CurrentPluginContext) Read(e model.Event) (map[string]interface{}, error) {
	resources, err := c.plugin.GetResources(shared.GetResourcesArgs{
		ResourceName: e.ResourceName,
		ResourceType: e.ResourceType,
		IsolatorID:   e.IsolatorName,
	})
	if err != nil {
		return nil, err
	}
	if len(resources) == 0 {
		return nil, fmt.Errorf("not found")
	}
	return resources[0].(map[string]interface{}), nil
}

func (c *CurrentPluginContext) Delete(e model.Event) error {
	return c.plugin.ActionDeleteResource(shared.ActionDeleteResourceArgs{ResourceName: e.ResourceName, ResourceType: e.ResourceType, IsolatorName: e.IsolatorName})
}

func (c *CurrentPluginContext) Create(e model.Event) error {
	return c.plugin.ActionCreateResource(shared.ActionCreateResourceArgs{ResourceName: e.ResourceName, ResourceType: e.ResourceType, IsolatorName: e.IsolatorName, Data: e.Args})
}

func (c *CurrentPluginContext) Update(e model.Event) error {
	return c.plugin.ActionUpdateResource(shared.ActionUpdateResourceArgs{ResourceName: e.ResourceName, ResourceType: e.ResourceType, IsolatorName: e.IsolatorName, Data: e.Args})
}
