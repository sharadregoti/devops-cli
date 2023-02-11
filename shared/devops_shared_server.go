package shared

// type DevopsServerRPC struct{ Impl Devops }

// func (d *DevopsServerRPC) Name(args interface{}, resp *string) error {
// 	*resp = d.Impl.Name()
// 	return nil
// }

// func (d *DevopsServerRPC) Connect(args *AuthInfo, resp *error) error {
// 	*resp = d.Impl.Connect(*args)
// 	return *resp
// }

// func (d *DevopsServerRPC) GetResources(args *GetResourcesArgs, resp *[]interface{}) error {
// 	var err error
// 	*resp, err = d.Impl.GetResources(*args)
// 	return err
// }

// func (d *DevopsServerRPC) WatchResources(args *string, resp *WatcheResource) error {
// 	var err error
// 	*resp, err = d.Impl.WatchResources(*args)
// 	return err
// }

// func (d *DevopsServerRPC) CloseResourceWatcher(args *string, resp *error) error {
// 	var err error
// 	err = d.Impl.CloseResourceWatcher(*args)
// 	return err
// }

// func (d *DevopsServerRPC) GetResourceTypeSchema(args *string, resp *ResourceTransformer) error {
// 	var err error
// 	*resp, err = d.Impl.GetResourceTypeSchema(*args)
// 	return err
// }

// func (d *DevopsServerRPC) GetResourceTypeList(args interface{}, resp *[]string) error {
// 	var err error
// 	*resp, err = d.Impl.GetResourceTypeList()
// 	return err
// }

// func (d *DevopsServerRPC) GetAuthInfo(args interface{}, resp *[]AuthInfo) error {
// 	var err error
// 	*resp, err = d.Impl.GetAuthInfo()
// 	return err
// }

// func (d *DevopsServerRPC) GetResourceIsolatorType(args interface{}, resp *string) error {
// 	var err error
// 	*resp, err = d.Impl.GetResourceIsolatorType()
// 	return err
// }

// func (d *DevopsServerRPC) GetDefaultResourceIsolator(args interface{}, resp *string) error {
// 	var err error
// 	*resp, err = d.Impl.GetDefaultResourceIsolator()
// 	return err
// }

// func (d *DevopsServerRPC) GetSupportedActions(args *string, resp *[]Action) error {
// 	var err error
// 	*resp, err = d.Impl.GetSupportedActions()
// 	return err
// }

// func (d *DevopsServerRPC) ActionDeleteResource(args *ActionDeleteResourceArgs, resp *error) error {
// 	var err error
// 	err = d.Impl.ActionDeleteResource(*args)
// 	return err
// }

// func (d *DevopsServerRPC) ActionCreateResource(args *ActionCreateResourceArgs, resp *error) error {
// 	var err error
// 	err = d.Impl.ActionCreateResource(*args)
// 	return err
// }

// func (d *DevopsServerRPC) ActionUpdateResource(args *ActionUpdateResourceArgs, resp *error) error {
// 	var err error
// 	err = d.Impl.ActionUpdateResource(*args)
// 	return err
// }

// func (d *DevopsServerRPC) GetSpecficActionList(args *string, resp *[]Action) error {
// 	var err error
// 	*resp, err = d.Impl.GetSpecficActionList(*args)
// 	return err
// }

// func (d *DevopsServerRPC) PerformSpecificAction(args *SpecificActionArgs, resp *SpecificActionResult) error {
// 	var err error
// 	*resp, err = d.Impl.PerformSpecificAction(*args)
// 	return err
// }
