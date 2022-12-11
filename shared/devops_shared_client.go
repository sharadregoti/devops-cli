package shared

import (
	"net/rpc"

	"github.com/sharadregoti/devops/model"
)

type DevopsClientRPC struct{ client *rpc.Client }

func (d *DevopsClientRPC) Name() string {
	var resp string
	err := d.client.Call("Plugin.Name", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}

	return resp
}

func (d *DevopsClientRPC) GetResources(resourceType string) []interface{} {
	var resp []interface{}
	err := d.client.Call("Plugin.GetResources", &resourceType, &resp)
	if err != nil {
		panic(err)
	}

	return resp
}

func (d *DevopsClientRPC) GetResourceTypeSchema(resourceType string) model.ResourceTransfomer {
	var resp model.ResourceTransfomer
	err := d.client.Call("Plugin.GetResourceTypeSchema", &resourceType, &resp)
	if err != nil {
		panic(err)
	}

	return resp
}

func (d *DevopsClientRPC) GetResourceTypeList() []string {
	var resp []string
	err := d.client.Call("Plugin.GetResourceTypeList", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}

	return resp
}

func (d *DevopsClientRPC) GetGeneralInfo() map[string]string {
	var resp map[string]string
	err := d.client.Call("Plugin.GetGeneralInfo", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}

	return resp
}

func (d *DevopsClientRPC) GetResourceIsolatorType() string {
	var resp string
	err := d.client.Call("Plugin.GetResourceIsolatorType", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}

	return resp
}

func (d *DevopsClientRPC) GetDefaultResourceIsolator() string {
	var resp string
	err := d.client.Call("Plugin.GetDefaultResourceIsolator", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}

	return resp
}
