package aws

import (
	"fmt"

	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/shared"
)

func (d *AWS) Name() string {
	return PluginName
}

// TODO: test & fix this
func (d *AWS) StatusOK() error {
	// d.logger.Error(fmt.Sprintf("failed to load plugin: %v", errors.Unwrap(d.isOK)))
	return nil
}

func (d *AWS) GetResources(args shared.GetResourcesArgs) ([]interface{}, error) {
	items, err := d.listResources(args)
	if err != nil {
		return nil, err
	}

	d.logger.Debug(fmt.Sprintf("Found %v %v resources in %v namespace", len(items), args.ResourceType, args.IsolatorID))
	return items, nil
}

// TODO: test & fix this
func (d *AWS) CloseResourceWatcher(resourceType string) error {
	return nil
}

// TODO: test & fix this
func (d *AWS) WatchResources(resourceType string) (chan shared.WatchResourceResult, error) {
	return nil, nil
}

func (d *AWS) GetResourceTypeSchema(resourceType string) (model.ResourceTransfomer, error) {
	t, ok := d.resourceTypeConfigurations[resourceType]
	if !ok {
		d.logger.Debug(fmt.Sprintf("Schema of resource type %s not found, using the default schema", resourceType))
		return d.resourceTypeConfigurations["defaults"], nil
	}

	return t, nil
}

func (d *AWS) GetResourceTypeList() ([]string, error) {
	resp, err := d.getResourceTypes()
	d.logger.Debug(fmt.Sprintf("Total %v resource type exits", len(resp)))
	return resp, err
}

// TODO: test & fix this
func (d *AWS) GetGeneralInfo() (map[string]string, error) {
	return map[string]string{
		"Profile":    "",
		"Account Id": "",
		"Region":     "",
	}, nil
}

// TODO: Include plural names as well
func (d *AWS) GetResourceIsolatorType() (string, error) {
	return "regions", nil
}

func (d *AWS) GetDefaultResourceIsolator() (string, error) {
	return "ap-south-1", nil
}

func (d *AWS) GetSupportedActions(resourceType string) (shared.GenericActions, error) {
	return shared.GenericActions{
		IsDelete: true,
		IsUpdate: false,
		IsCreate: false,
	}, nil
}

func (d *AWS) ActionDeleteResource(args shared.ActionDeleteResourceArgs) error {
	// return d.deleteResource(context.Background(), args)
	return nil
}

func (d *AWS) GetSpecficActionList(resourceType string) ([]shared.SpecificAction, error) {
	return make([]shared.SpecificAction, 0), nil
}

func (d *AWS) PerformSpecificAction(args shared.SpecificActionArgs) (shared.SpecificActionResult, error) {
	return shared.SpecificActionResult{}, nil
}
