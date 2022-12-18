package shared

import (
	"fmt"
	"net/rpc"

	"github.com/sharadregoti/devops/model"
)

type DevopsClientRPC struct{ client *rpc.Client }

func (d *DevopsClientRPC) Name() string {
	var resp string
	err := d.client.Call("Plugin.Name", new(interface{}), &resp)
	if err != nil {
		return ""
	}

	return resp
}

func (d *DevopsClientRPC) GetResources(args GetResourcesArgs) ([]interface{}, error) {
	var resp []interface{}
	err := d.client.Call("Plugin.GetResources", &args, &resp)
	if err != nil {
		return nil, fmt.Errorf("grpc client function invocation fixed: %w", err)
	}

	return resp, nil
}

func (d *DevopsClientRPC) WatchResources(resourceType string) (chan WatchResourceResult, error) {
	var resp = make(chan WatchResourceResult, 1)
	err := d.client.Call("Plugin.WatchResources", &resourceType, &resp)
	if err != nil {
		return nil, fmt.Errorf("grpc client function invocation fixed: %w", err)
	}

	return resp, nil
}

func (d *DevopsClientRPC) CloseResourceWatcher(resourceType string) error {
	var er error
	err := d.client.Call("Plugin.CloseResourceWatcher", &resourceType, &er)
	if err != nil {
		return fmt.Errorf("grpc client function invocation fixed: %w", err)
	}

	return er
}

func (d *DevopsClientRPC) GetResourceTypeSchema(resourceType string) (model.ResourceTransfomer, error) {
	var resp model.ResourceTransfomer
	err := d.client.Call("Plugin.GetResourceTypeSchema", &resourceType, &resp)
	if err != nil {
		return resp, fmt.Errorf("grpc client function invocation fixed: %w", err)
	}

	return resp, nil
}

func (d *DevopsClientRPC) GetResourceTypeList() ([]string, error) {
	var resp []string
	err := d.client.Call("Plugin.GetResourceTypeList", new(interface{}), &resp)
	if err != nil {
		return nil, fmt.Errorf("grpc client function invocation fixed: %w", err)
	}

	return resp, nil
}

func (d *DevopsClientRPC) GetGeneralInfo() (map[string]string, error) {
	var resp map[string]string
	err := d.client.Call("Plugin.GetGeneralInfo", new(interface{}), &resp)
	if err != nil {
		return nil, fmt.Errorf("grpc client function invocation fixed: %w", err)
	}

	return resp, nil
}

func (d *DevopsClientRPC) GetResourceIsolatorType() (string, error) {
	var resp string
	err := d.client.Call("Plugin.GetResourceIsolatorType", new(interface{}), &resp)
	if err != nil {
		return "", fmt.Errorf("grpc client function invocation fixed: %w", err)
	}

	return resp, nil
}

func (d *DevopsClientRPC) GetDefaultResourceIsolator() (string, error) {
	var resp string
	err := d.client.Call("Plugin.GetDefaultResourceIsolator", new(interface{}), &resp)
	if err != nil {
		return "", fmt.Errorf("grpc client function invocation fixed: %w", err)
	}

	return resp, nil
}

func (d *DevopsClientRPC) GetSupportedActions(resourceType string) (GenericActions, error) {
	var resp GenericActions
	err := d.client.Call("Plugin.GetSupportedActions", &resourceType, &resp)
	if err != nil {
		return GenericActions{}, fmt.Errorf("grpc client function invocation fixed: %w", err)
	}

	return resp, nil
}

func (d *DevopsClientRPC) ActionDeleteResource(args ActionDeleteResourceArgs) error {
	var er error
	err := d.client.Call("Plugin.ActionDeleteResource", &args, &er)
	if err != nil {
		return fmt.Errorf("grpc client function invocation fixed: %w", err)
	}

	return err
}

func (d *DevopsClientRPC) GetSpecficActionList(resourceType string) ([]SpecificAction, error) {
	var resp []SpecificAction
	err := d.client.Call("Plugin.GetSpecficActionList", &resourceType, &resp)
	if err != nil {
		return nil, fmt.Errorf("grpc client function invocation fixed: %w", err)
	}

	return resp, nil
}

func (d *DevopsClientRPC) PerformSpecificAction(args SpecificActionArgs) (SpecificActionResult, error) {
	var resp SpecificActionResult
	err := d.client.Call("Plugin.PerformSpecificAction", &args, &resp)
	if err != nil {
		return SpecificActionResult{}, fmt.Errorf("grpc client function invocation fixed: %w", err)
	}

	return resp, err
}
