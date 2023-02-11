package core

import (
	"fmt"

	"github.com/sharadregoti/devops/internal/transformer"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/proto"
	"github.com/sharadregoti/devops/utils"
	"google.golang.org/protobuf/types/known/structpb"
)

func (c *CurrentPluginContext) ReadSync(e model.Event) error {
	resources, err := c.plugin.GetResources(&proto.GetResourcesArgs{
		ResourceName: "", //Ingore resource name during sync
		ResourceType: e.ResourceType,
		IsolatorId:   e.IsolatorName,
	})
	if err != nil {
		return err
	}

	schema, err := c.plugin.GetResourceTypeSchema(e.ResourceType)
	if err != nil {
		return err
	}

	tableData, _, err := transformer.GetResourceInTableFormat(c.logger, schema, resources)
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
		SpecificActions: specificActions.Actions,
	})
	return nil
}

func (c *CurrentPluginContext) GetSpecficActionList(e model.Event) (*proto.GetActionListResponse, error) {
	return c.plugin.GetSpecficActionList(e.ResourceType)
}

func (c *CurrentPluginContext) GetLongRunning(e model.Event) map[string]*model.LongRunningInfo {
	tempMap := map[string]*model.LongRunningInfo{}
	for _, v := range c.longRunning {
		if v.GetE().Type == e.Type && v.GetE().ResourceName == e.ResourceName && v.GetE().ResourceType == e.ResourceType && v.GetE().IsolatorName == e.IsolatorName {
			tempMap[v.ID] = v
		}
	}

	return tempMap
}

func (c *CurrentPluginContext) Read(e model.Event) (map[string]interface{}, error) {
	resources, err := c.plugin.GetResources(&proto.GetResourcesArgs{
		ResourceName: e.ResourceName,
		ResourceType: e.ResourceType,
		IsolatorId:   e.IsolatorName,
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
	return c.plugin.ActionDeleteResource(&proto.ActionDeleteResourceArgs{ResourceName: e.ResourceName, ResourceType: e.ResourceType, IsolatorName: e.IsolatorName})
}

func (c *CurrentPluginContext) Create(e model.Event) error {
	m, err := structpb.NewValue(e.Args)
	if err != nil {
		return err
	}
	return c.plugin.ActionCreateResource(&proto.ActionCreateResourceArgs{ResourceName: e.ResourceName, ResourceType: e.ResourceType, IsolatorName: e.IsolatorName, Data: m})
}

func (c *CurrentPluginContext) Update(e model.Event) error {
	m, err := structpb.NewValue(e.Args)
	if err != nil {
		return err
	}
	return c.plugin.ActionUpdateResource(&proto.ActionUpdateResourceArgs{ResourceName: e.ResourceName, ResourceType: e.ResourceType, IsolatorName: e.IsolatorName, Data: m})
}
