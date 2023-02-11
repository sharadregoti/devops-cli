package shared

import (
	"net/rpc"
)

type DevopsClientRPC struct{ client *rpc.Client }

// func (d *DevopsClientRPC) Name() string {
// 	var resp string
// 	err := d.client.Call("Plugin.Name", new(interface{}), &resp)
// 	if err != nil {
// 		return ""
// 	}

// 	return resp
// }

// func (d *DevopsClientRPC) Connect(authInfo AuthInfo) error {
// 	var resp error
// 	err := d.client.Call("Plugin.Connect", &authInfo, &resp)
// 	if err != nil {
// 		return err
// 	}

// 	return resp
// }

// func (d *DevopsClientRPC) GetResources(args GetResourcesArgs) ([]interface{}, error) {
// 	var resp []interface{}
// 	err := d.client.Call("Plugin.GetResources", &args, &resp)
// 	if err != nil {
// 		return nil, fmt.Errorf("grpc client function invocation fixed: %w", err)
// 	}

// 	return resp, nil
// }

// func (d *DevopsClientRPC) WatchResources(resourceType string) (WatcheResource, error) {
// 	var resp = make(chan WatchResourceResult, 1)
// 	err := d.client.Call("Plugin.WatchResources", &resourceType, &resp)
// 	if err != nil {
// 		return nil, fmt.Errorf("grpc client function invocation fixed: %w", err)
// 	}

// 	return resp, nil
// }

// func (d *DevopsClientRPC) CloseResourceWatcher(resourceType string) error {
// 	var er error
// 	err := d.client.Call("Plugin.CloseResourceWatcher", &resourceType, &er)
// 	if err != nil {
// 		return fmt.Errorf("grpc client function invocation fixed: %w", err)
// 	}

// 	return er
// }

// func (d *DevopsClientRPC) GetResourceTypeSchema(resourceType string) (ResourceTransformer, error) {
// 	var resp ResourceTransformer
// 	err := d.client.Call("Plugin.GetResourceTypeSchema", &resourceType, &resp)
// 	if err != nil {
// 		return resp, fmt.Errorf("grpc client function invocation fixed: %w", err)
// 	}

// 	return resp, nil
// }

// func (d *DevopsClientRPC) GetResourceTypeList() ([]string, error) {
// 	var resp []string
// 	err := d.client.Call("Plugin.GetResourceTypeList", new(interface{}), &resp)
// 	if err != nil {
// 		return nil, fmt.Errorf("grpc client function invocation fixed: %w", err)
// 	}

// 	return resp, nil
// }

// func (d *DevopsClientRPC) GetAuthInfo() ([]AuthInfo, error) {
// 	var resp []AuthInfo
// 	err := d.client.Call("Plugin.GetAuthInfo", new(interface{}), &resp)
// 	if err != nil {
// 		return nil, fmt.Errorf("grpc client function invocation fixed: %w", err)
// 	}

// 	return resp, nil
// }

// func (d *DevopsClientRPC) GetResourceIsolatorType() (string, error) {
// 	var resp string
// 	err := d.client.Call("Plugin.GetResourceIsolatorType", new(interface{}), &resp)
// 	if err != nil {
// 		return "", fmt.Errorf("grpc client function invocation fixed: %w", err)
// 	}

// 	return resp, nil
// }

// func (d *DevopsClientRPC) GetDefaultResourceIsolator() (string, error) {
// 	var resp string
// 	err := d.client.Call("Plugin.GetDefaultResourceIsolator", new(interface{}), &resp)
// 	if err != nil {
// 		return "", fmt.Errorf("grpc client function invocation fixed: %w", err)
// 	}

// 	return resp, nil
// }

// func (d *DevopsClientRPC) GetSupportedActions() ([]Action, error) {
// 	var resp []Action
// 	var input string
// 	err := d.client.Call("Plugin.GetSupportedActions", &input, &resp)
// 	if err != nil {
// 		return []Action{}, fmt.Errorf("grpc client function invocation fixed: %w", err)
// 	}

// 	return resp, nil
// }

// func (d *DevopsClientRPC) ActionDeleteResource(args ActionDeleteResourceArgs) error {
// 	var er error
// 	err := d.client.Call("Plugin.ActionDeleteResource", &args, &er)
// 	if err != nil {
// 		return fmt.Errorf("grpc client function invocation fixed: %w", err)
// 	}

// 	return err
// }

// func (d *DevopsClientRPC) ActionCreateResource(args ActionCreateResourceArgs) error {
// 	var er error
// 	err := d.client.Call("Plugin.ActionCreateResource", &args, &er)
// 	if err != nil {
// 		return fmt.Errorf("grpc client function invocation fixed: %w", err)
// 	}

// 	return err
// }

// func (d *DevopsClientRPC) ActionUpdateResource(args ActionUpdateResourceArgs) error {
// 	var er error
// 	err := d.client.Call("Plugin.ActionUpdateResource", &args, &er)
// 	if err != nil {
// 		return fmt.Errorf("grpc client function invocation fixed: %w", err)
// 	}

// 	return err
// }

// func (d *DevopsClientRPC) GetSpecficActionList(resourceType string) ([]Action, error) {
// 	var resp []Action
// 	err := d.client.Call("Plugin.GetSpecficActionList", &resourceType, &resp)
// 	if err != nil {
// 		return nil, fmt.Errorf("grpc client function invocation fixed: %w", err)
// 	}

// 	return resp, nil
// }

// func (d *DevopsClientRPC) PerformSpecificAction(args SpecificActionArgs) (SpecificActionResult, error) {
// 	var resp SpecificActionResult
// 	err := d.client.Call("Plugin.PerformSpecificAction", &args, &resp)
// 	if err != nil {
// 		return SpecificActionResult{}, fmt.Errorf("grpc client function invocation fixed: %w", err)
// 	}

// 	return resp, err
// }
