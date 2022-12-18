package shared

import "github.com/sharadregoti/devops/model"

type DevopsServerRPC struct{ Impl Devops }

func (d *DevopsServerRPC) Name(args interface{}, resp *string) error {
	*resp = d.Impl.Name()
	return nil
}

func (d *DevopsServerRPC) GetResources(args *GetResourcesArgs, resp *[]interface{}) error {
	var err error
	*resp, err = d.Impl.GetResources(*args)
	return err
}

func (d *DevopsServerRPC) WatchResources(args *string, resp *chan WatchResourceResult) error {
	var err error
	*resp, err = d.Impl.WatchResources(*args)
	return err
}

func (d *DevopsServerRPC) CloseResourceWatcher(args *string, resp *error) error {
	var err error
	err = d.Impl.CloseResourceWatcher(*args)
	return err
}

func (d *DevopsServerRPC) GetResourceTypeSchema(args *string, resp *model.ResourceTransfomer) error {
	var err error
	*resp, err = d.Impl.GetResourceTypeSchema(*args)
	return err
}

func (d *DevopsServerRPC) GetResourceTypeList(args interface{}, resp *[]string) error {
	var err error
	*resp, err = d.Impl.GetResourceTypeList()
	return err
}

func (d *DevopsServerRPC) GetGeneralInfo(args interface{}, resp *map[string]string) error {
	var err error
	*resp, err = d.Impl.GetGeneralInfo()
	return err
}

func (d *DevopsServerRPC) GetResourceIsolatorType(args interface{}, resp *string) error {
	var err error
	*resp, err = d.Impl.GetResourceIsolatorType()
	return err
}

func (d *DevopsServerRPC) GetDefaultResourceIsolator(args interface{}, resp *string) error {
	var err error
	*resp, err = d.Impl.GetDefaultResourceIsolator()
	return err
}

func (d *DevopsServerRPC) GetSupportedActions(args *string, resp *GenericActions) error {
	var err error
	*resp, err = d.Impl.GetSupportedActions(*args)
	return err
}

func (d *DevopsServerRPC) ActionDeleteResource(args *ActionDeleteResourceArgs, resp *error) error {
	var err error
	err = d.Impl.ActionDeleteResource(*args)
	return err
}

func (d *DevopsServerRPC) GetSpecficActionList(args *string, resp *[]SpecificAction) error {
	var err error
	*resp, err = d.Impl.GetSpecficActionList(*args)
	return err
}

func (d *DevopsServerRPC) PerformSpecificAction(args *SpecificActionArgs, resp *SpecificActionResult) error {
	var err error
	*resp, err = d.Impl.PerformSpecificAction(*args)
	return err
}
