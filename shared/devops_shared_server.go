package shared

import "github.com/sharadregoti/devops/model"

type DevopsServerRPC struct{ Impl Devops }

func (d *DevopsServerRPC) Name(args interface{}, resp *string) error {
	*resp = d.Impl.Name()
	return nil
}

func (d *DevopsServerRPC) GetResources(args *string, resp *[]interface{}) error {
	*resp = d.Impl.GetResources(*args)
	return nil
}

func (d *DevopsServerRPC) GetResourceTypeSchema(args *string, resp *model.ResourceTransfomer) error {
	*resp = d.Impl.GetResourceTypeSchema(*args)
	return nil
}

func (d *DevopsServerRPC) GetResourceTypeList(args interface{}, resp *[]string) error {
	*resp = d.Impl.GetResourceTypeList()
	return nil
}

func (d *DevopsServerRPC) GetGeneralInfo(args interface{}, resp *map[string]string) error {
	*resp = d.Impl.GetGeneralInfo()
	return nil
}

func (d *DevopsServerRPC) GetResourceIsolatorType(args interface{}, resp *string) error {
	*resp = d.Impl.GetResourceIsolatorType()
	return nil
}

func (d *DevopsServerRPC) GetDefaultResourceIsolator(args interface{}, resp *string) error {
	*resp = d.Impl.GetDefaultResourceIsolator()
	return nil
}
