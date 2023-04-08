package main

import (
	"context"
	"fmt"
	"strings"

	shared "github.com/sharadregoti/devops-plugin-sdk"
	"github.com/sharadregoti/devops-plugin-sdk/proto"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func (h *Helm) Name() string {
	return PluginName
}

// TODO: test & fix this
func (h *Helm) Connect(authInfo *proto.AuthInfo) error {
	if h.isOK != nil {
		return h.isOK
	}

	kubeConfigPath := ""
	for _, kc := range h.kubeCLIconfig.KubeConfigs {
		for _, c := range kc.Contexts {
			if authInfo.IdentifyingName == kc.Name && authInfo.Name == c.Name {
				kubeConfigPath = kc.Path
				break
			}
		}
	}
	shared.LogDebug("Connecting to kubernetes cluster at path %s with context %s", kubeConfigPath, authInfo.Name)

	// contextName := "my-context" // replace with the name of your context
	restConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: authInfo.Name,
		}).ClientConfig()
	// clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(clientcmd.NewDefaultClientConfigLoadingRules(), &clientcmd.ConfigOverrides{CurrentContext: authInfo.Name})
	// restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return shared.LogError("failed to get client config for context %s: %v", authInfo.Name, err)
	}

	// restConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath, "")
	// if err != nil {
	// 	return shared.LogError("failed to build config from flags: %v", err)
	// }

	// Normal client
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return shared.LogError("failed to create normal kube client: %v", err)
	}

	shared.LogDebug("Connecting to kubernetes cluster at path %s with context %s", kubeConfigPath, authInfo.Name)
	h.normalClient = clientset
	h.restConfig = restConfig
	h.currentContext = authInfo.Name
	h.currentKubeConfigPath = kubeConfigPath

	resp := clientset.CoreV1().RESTClient().Get().AbsPath("/").Do(context.Background())
	var statusCode = 0
	if resp.StatusCode(&statusCode); statusCode != 200 {
		return shared.LogError("failed to perform ping test: unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (h *Helm) GetResources(args *proto.GetResourcesArgs) ([]interface{}, error) {
	switch args.ResourceType {
	case "repos":
		return h.listRepos()
	case "namespaces":
		return h.listNamespaces(args)
	case "releases":
		return h.listReleases(args)
	default:
		return nil, shared.LogError("resource type %s not supported", args.ResourceType)
	}
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
func (h *Helm) WatchResources(args *proto.GetResourcesArgs) (chan shared.WatchResourceResult, chan struct{}, error) {
	if h.isOK != nil {
		return nil, nil, h.isOK
	}

	if len(h.resourceWatcherChanMap) > 0 {
		for k, v := range h.resourceWatcherChanMap {
			h.logger.Trace(fmt.Sprintf("Invoking close resource watcher event for resource type %s", k))
			v.serverDone <- struct{}{}
			v.watcherDone <- struct{}{}
		}
		h.resourceWatcherChanMap = make(map[string]channels)
	}

	_, ok := h.resourceWatcherChanMap[args.ResourceType]
	if ok {
		h.logger.Debug(fmt.Sprintf("Channel already exists for resource %s", args.ResourceType))
		return nil, nil, nil
	}

	watcherDone := make(chan struct{}, 1)
	serverDone := make(chan struct{}, 1)
	ch := make(chan shared.WatchResourceResult, 1)

	watcher := h.watchReleases(args, watcherDone)
	// if err != nil {
	// 	return nil, nil, shared.LogError("failed to watch over resource %s: %v", args.ResourceType, err)
	// }

	go func() {
		shared.LogDebug("plugin routine 2: resource watcher has been started for resource type (%s)", args.ResourceType)
		defer shared.LogDebug("plugin routine 2: resource watcher has been stopped for resource type (%s)", args.ResourceType)

		for v := range watcher {
			ch <- shared.WatchResourceResult{
				Type:   strings.ToLower("updated"),
				Result: v,
			}
		}
	}()

	h.resourceWatcherChanMap[args.ResourceType] = channels{watcherDone: watcherDone, serverDone: serverDone}
	return ch, serverDone, nil
}

func (h *Helm) GetResourceTypeSchema(resourceType string) (*proto.ResourceTransformer, error) {
	if h.isOK != nil {
		return nil, h.isOK
	}
	t, ok := h.resourceTypeConfigurations[resourceType]
	if !ok {
		h.logger.Trace(fmt.Sprintf("Schema of resource type %s not found, using the default schema", resourceType))
		return h.resourceTypeConfigurations["defaults"], nil
	}

	return t, nil
}

func (h *Helm) GetResourceTypeList() ([]string, error) {
	return []string{"repos", "releases", "namespaces"}, nil
}

// TODO: test & fix this
func (h *Helm) GetAuthInfo() (*proto.AuthInfoResponse, error) {
	if h.isOK != nil {
		return nil, h.isOK
	}
	authInfo := new(proto.AuthInfoResponse)

	for _, kc := range h.kubeCLIconfig.KubeConfigs {
		for _, c := range kc.Contexts {
			authInfo.AuthInfo = append(authInfo.AuthInfo, &proto.AuthInfo{
				IdentifyingName:  kc.Name,
				Name:             c.Name,
				DefaultIsolators: c.DefaultNamespacesToShow,
				IsDefault:        c.IsDefault,
				Info:             map[string]string{},
				Path:             kc.Path,
			})
		}
	}
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
			// {
			// 	Name:       "create",
			// 	KeyBinding: "ctrl-b",
			// 	OutputType: "string",
			// },
			// {
			// 	Name:       "edit",
			// 	KeyBinding: "ctrl-e",
			// 	OutputType: "bidirectional",
			// },
			// {
			// 	Name:       "delete",
			// 	KeyBinding: "ctrl-d",
			// 	OutputType: "nothing",
			// },
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
	t, ok := h.resourceTypeConfigurations[resourceType]
	if !ok {
		h.logger.Trace(fmt.Sprintf("specific action list schema of resource type %s not found, using the default schema", resourceType))
		t = h.resourceTypeConfigurations["defaults"]
	}

	return &proto.GetActionListResponse{Actions: t.SpecificActions}, nil
}

func (h *Helm) PerformSpecificAction(args *proto.SpecificActionArgs) (*proto.SpecificActionResult, error) {
	return nil, nil
}
