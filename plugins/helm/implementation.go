package helm

import (
	"github.com/sharadregoti/devops-plugin-sdk/proto"
)

func (h *Helm) Name() string {
	return "helm"
}

// TODO: test & fix this
func (h *Helm) Connect(authInfo *proto.AuthInfo) error {

	return nil
}

func (h *Helm) GetResources(args *proto.GetResourcesArgs) ([]interface{}, error) {

	return items, nil
}

// TODO: test & fix this
func (h *Helm) CloseResourceWatcher(resourceType string) error {
	// res, ok := d.resourceWatcherChanMap[resourceType]
	// if !ok {
	// 	return shared.LogError("channel for resource type %s does not exists", resourceType)
	// }

	// d.logger.Debug(fmt.Sprintf("Closing resource watcher %s", resourceType))
	// res.serverDone <- struct{}{}
	// res.watcherDone <- struct{}{}
	return nil
}

// TODO: test & fix this
func (h *Helm) WatchResources(resourceType string) (chan shared.WatchResourceResult, chan struct{}, error) {
	return ch, serverDone, nil
}

func (h *Helm) GetResourceTypeSchema(resourceType string) (*proto.ResourceTransformer, error) {
	return t, nil
}

func (h *Helm) GetResourceTypeList() ([]string, error) {
	return []string{"repos", "charts", "releases"}, nil
}

// TODO: test & fix this
func (h *Helm) GetAuthInfo() (*proto.AuthInfoResponse, error) {

	return authInfo, nil
}

// TODO: Include plural names as well
func (h *Helm) GetResourceIsolatorType() (string, error) {
	return "namespaces", nil
}

func (h *Helm) GetDefaultResourceIsolator() (string, error) {
	return "default", nil
}

func (h *Helm) GetSupportedActions() (*proto.GetActionListResponse, error) {
	genericActions := &proto.GetActionListResponse{
		Actions: []*proto.Action{
			{
				Name:       "read",
				KeyBinding: "ctrl-y",
				OutputType: "string",
			},
			{
				Name:       "create",
				KeyBinding: "ctrl-b",
				OutputType: "string",
			},
			{
				Name:       "edit",
				KeyBinding: "ctrl-e",
				OutputType: "bidirectional",
			},
			{
				Name:       "delete",
				KeyBinding: "ctrl-d",
				OutputType: "nothing",
			},
			{
				Name:       "refresh",
				KeyBinding: "ctrl-r",
				OutputType: "nothing",
			},
		},
	}

	return genericActions, nil
}

func (h *Helm) ActionDeleteResource(args *proto.ActionDeleteResourceArgs) error {
	return nil
}

func (h *Helm) ActionCreateResource(args *proto.ActionCreateResourceArgs) error {
	return nil
}

func (h *Helm) ActionUpdateResource(args *proto.ActionUpdateResourceArgs) error {
	return nil
}

func (h *Helm) GetSpecficActionList(resourceType string) (*proto.GetActionListResponse, error) {
	return &proto.GetActionListResponse{Actions: t.SpecificActions}, nil
}

func (h *Helm) PerformSpecificAction(args *proto.SpecificActionArgs) (*proto.SpecificActionResult, error) {

	switch args.ActionName {

	case "describe":
		// result, err := d.DescribeResource(args.ResourceType, args.ResourceName, args.IsolatorName)
		// if err != nil {
		// 	return nil, err
		// }

		// return &proto.SpecificActionResult{
		// 	Result: result,
		// 	// TODO: Output type should come from an core SDK
		// 	OutputType: "string",
		// }, nil

	case "decode-secret":
		items, err := d.GetResources(&proto.GetResourcesArgs{
			ResourceName: args.ResourceName,
			ResourceType: args.ResourceType,
			IsolatorId:   args.IsolatorName,
		})
		if err != nil {
			return nil, err
		}

		_, err = d.decodeSecret(items[0])
		if err != nil {
			return nil, err
		}

		return &proto.SpecificActionResult{
			Result: nil,
			// TODO: Output type should come from an core SDK
			OutputType: "string",
		}, nil

	case "logs":

		// containerName := ""
		// if args.ResourceType == "containers" {
		// 	parentName := args.Args["parentName"]
		// 	args.ResourceType = "pods"
		// 	containerName = args.ResourceName
		// 	args.ResourceName = parentName.AsInterface().(string)
		// }

		// res, err := d.getPodLogs(args.ResourceName, args.IsolatorName, containerName)
		// if err != nil {
		// 	return nil, err
		// }

		// return shared.SpecificActionResult{
		// 	Result:     res,
		// 	OutputType: "string",
		// }, nil

	case "shell":

		// containerName := ""
		// if args.ResourceType == "containers" {
		// 	parentName := args.Args["parentName"]
		// 	args.ResourceType = "pods"
		// 	containerName = args.ResourceName
		// 	args.ResourceName = parentName.(string)
		// }

		// res, err := d.execPod(args.ResourceName, args.IsolatorName, containerName)
		// if err != nil {
		// 	return shared.SpecificActionResult{}, err
		// }

		// return shared.SpecificActionResult{
		// 	Result:     res,
		// 	OutputType: "string",
		// }, nil

	case "port-forward":

		// res, err := d.portForward(args.ResourceName, args.IsolatorName, args.Args)
		// if err != nil {
		// 	return shared.SpecificActionResult{}, err
		// }

		// return shared.SpecificActionResult{
		// 	Result:     res,
		// 	OutputType: "string",
		// }, nil

	case "view-pods":
		// res, err := d.getPods(context.Background(), args.IsolatorName, args.ResourceName, args.ResourceType)
		// if err != nil {
		// 	return shared.SpecificActionResult{}, err
		// }
		// return shared.SpecificActionResult{
		// 	Result: res,
		// 	// TODO: Output type should come from an core SDK
		// 	OutputType: "invoke-event",
		// }, nil

	case "close":

	}

	return nil, nil
}
