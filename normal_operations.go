package core

import (
	"github.com/sharadregoti/devops/internal/transformer"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/proto"
	"github.com/sharadregoti/devops/utils"
	"github.com/sharadregoti/devops/utils/logger"
	"google.golang.org/protobuf/types/known/structpb"
)

func (c *CurrentPluginContext) ResourceChanged(e model.Event) error {
	// close previous watcher
	// TODO: Writing to c.done does not close this go routine
	// c.done <- struct{}{}

	// TODO: Fix this, CloseResourceWatcher does not do any specific work
	if err := c.plugin.CloseResourceWatcher(""); err != nil {
		return err
	}

	ch, _, err := c.plugin.WatchResources(e.ResourceType)
	if err != nil {
		return logger.LogError("Error while starting watcher", err)
	}

	done := make(chan struct{}, 1)
	// TODO: Writing to c.done does not close this go routine
	c.done = done
	go func() {
		logger.LogDebug("Core binary resource watcher routine has been started for resource type (%s)", e.ResourceType)
		defer logger.LogDebug("Core binary resource watcher routine has been stopped for resource type (%s)", e.ResourceType)
		for {
			select {
			case <-done:
				logger.LogTrace("Core binary resource watcher routine default: Done received for resource type (%s)", e.ResourceType)
				return

			case r := <-ch:
				schema, err := c.plugin.GetResourceTypeSchema(e.ResourceType)
				if err != nil {
					return
				}

				logger.LogTrace("Received new resource from core binary (%v)", r)

				tableData, _, err := transformer.GetResourceInTableFormat(schema, []interface{}{r})
				if err != nil {
					return
				}

				specificActions, err := c.plugin.GetSpecficActionList(e.ResourceType)
				if err != nil {
					return
				}

				c.SendMessage(model.WebsocketResponse{
					// TODO: Fix this, remove 0
					TableName:       utils.GetTableTitle(e.ResourceType, 0),
					Data:            tableData,
					SpecificActions: specificActions.Actions,
				})
			}
		}
	}()

	return nil
}

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

	tableData, _, err := transformer.GetResourceInTableFormat(schema, resources)
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
		logger.LogDebug("Resource is zero")
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
