package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
	shared "github.com/sharadregoti/devops-plugin-sdk"
	"github.com/sharadregoti/devops-plugin-sdk/proto"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

func (d *Kubernetes) Name() string {
	return PluginName
}

// TODO: test & fix this
func (d *Kubernetes) Connect(authInfo *proto.AuthInfo) error {
	if d.isOK != nil {
		return d.isOK
	}

	kubeConfigPath := ""
	for _, kc := range d.kubeCLIconfig.KubeConfigs {
		for _, c := range kc.Contexts {
			if authInfo.IdentifyingName == kc.Name && authInfo.Name == c.Name {
				kubeConfigPath = kc.Path
				break
			}
		}
	}

	shared.LogDebug("Connecting to kubernetes cluster at path %s with context %s", kubeConfigPath, authInfo.Name)
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return err
	}

	// Normal client
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("failed to load kube config: %w", err)
	}

	// Dynamic client
	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("failed to load dynamic kube config: %w", err)
	}

	// List all supported resources
	resources, err := clientset.Discovery().ServerPreferredResources()
	if err != nil {
		return fmt.Errorf("failed to discover kubernetes resource types: %w", err)
	}

	// Resource Type List
	resourceTypeMap := make(resourceTypeList)
	for _, resource := range resources {
		for _, r := range resource.APIResources {
			arr := strings.Split(resource.GroupVersion, "/")
			group := ""
			version := ""
			if len(arr) == 1 {
				version = arr[0]
			} else {
				group = arr[0]
				version = arr[1]
			}

			// This is because if you use r.name, it has conflicting values pods for 2 different groups (which overwirte each other)
			name := strings.ToLower(r.Kind)
			if !strings.HasSuffix(name, "s") {
				name = name + "s"
			}

			resourceTypeMap[name] = &resourceTypeInfo{
				group:            group,
				version:          version,
				resourceTypeName: r.Name,
				isNamespaced:     r.Namespaced,
			}
		}
	}

	d.dynamicClient = dynamicClient
	d.normalClient = clientset
	d.resourceTypes = resourceTypeMap

	// authInfo.Info = map[string]string{
	// 	"name":    authInfo.IdentifyingName,
	// 	"context": authInfo.Name,
	// 	"user":    "",
	// 	"server":  "",
	// }

	return nil
}

func (d *Kubernetes) GetResources(args *proto.GetResourcesArgs) ([]interface{}, error) {

	label := ""
	for k, v := range args.Args {
		if k == "labels" {
			selector := labels.NewSelector()
			for labelKey, labelValue := range v.GetStructValue().AsMap() {
				l2, _ := labels.NewRequirement(labelKey, selection.Equals, []string{labelValue.(string)})
				selector = selector.Add(*l2)
			}
			label = selector.String()
		}
	}

	items, err := d.listResources(context.Background(), args, label)
	if err != nil {
		return nil, err
	}

	d.logger.Debug(fmt.Sprintf("Found %v %v resources in %v namespace", len(items), args.ResourceType, args.IsolatorId))
	return items, nil
}

// TODO: test & fix this
func (d *Kubernetes) CloseResourceWatcher(resourceType string) error {
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
func (d *Kubernetes) WatchResources(args *proto.GetResourcesArgs) (chan shared.WatchResourceResult, chan struct{}, error) {
	if len(d.resourceWatcherChanMap) > 0 {
		for k, v := range d.resourceWatcherChanMap {
			d.logger.Trace(fmt.Sprintf("Invoking close resource watcher event for resource type %s", k))
			v.serverDone <- struct{}{}
			v.watcherDone <- struct{}{}
		}
		d.resourceWatcherChanMap = make(map[string]channels)
	}

	_, ok := d.resourceWatcherChanMap[args.ResourceType]
	if ok {
		d.logger.Debug(fmt.Sprintf("Channel already exists for resource %s", args.ResourceType))
		return nil, nil, nil
	}

	watcher, err := d.watchResourceDynamic(context.Background(), args)
	if err != nil {
		return nil, nil, shared.LogError("failed to watch over resource %s: %v", args.ResourceType, err)
	}

	watcherDone := make(chan struct{}, 1)
	serverDone := make(chan struct{}, 1)
	ch := make(chan shared.WatchResourceResult, 1)
	go func() {
		shared.LogDebug("plugin routine: resource watcher has been started for resource type (%s)", args.ResourceType)
		defer shared.LogDebug("plugin routine: resource watcher has been stopped for resource type (%s)", args.ResourceType)

		for {
			select {
			case <-watcherDone:
				shared.LogTrace("plugin routine: Done received for resource type (%s)", args.ResourceType)
				return

			case event, ok := <-watcher.ResultChan():
				if !ok {
					shared.LogDebug("Watcher routine: Watcher channel closed for resource %s", args.ResourceType)
					return
				}

				shared.LogTrace("Watcher routine: got event for resource type (%s), event type (%s)", args.ResourceType, strings.ToLower(string(event.Type)))

				b, err := json.Marshal(event.Object)
				if err != nil {
					shared.LogError("failed to marshal event object: %v", err)
					return
				}

				var rawJson map[string]interface{}
				if err := json.Unmarshal(b, &rawJson); err != nil {
					shared.LogError("Watcher routine: failed to unmarshal event object: %v", err)
					return
				}

				meta := rawJson["metadata"].(map[string]interface{})
				delete(meta, "managedFields")
				rawJson["metadata"] = meta
				if args.ResourceType == "pods" {
					obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(b, nil, nil)
					if err != nil {
						shared.LogError("Watcher routine: failed to unstructure resource: %v", err)
						return
					}
					rawJson["devops"] = map[string]interface{}{
						"customCalculatedStatus": Phase(obj.(*v1.Pod)),
					}
				}

				ch <- shared.WatchResourceResult{
					Type:   strings.ToLower(string(event.Type)),
					Result: rawJson,
				}
			}
		}

	}()

	d.resourceWatcherChanMap[args.ResourceType] = channels{watcherDone: watcherDone, serverDone: serverDone}
	return ch, serverDone, nil
}

func convertToMap(obj runtime.Object) (map[string]interface{}, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func (d *Kubernetes) GetResourceTypeSchema(resourceType string) (*proto.ResourceTransformer, error) {
	t, ok := d.resourceTypeConfigurations[resourceType]
	if !ok {
		d.logger.Trace(fmt.Sprintf("Schema of resource type %s not found, using the default schema", resourceType))
		return d.resourceTypeConfigurations["defaults"], nil
	}

	return t, nil
}

func (d *Kubernetes) GetResourceTypeList() ([]string, error) {
	resp := []string{}
	for r := range d.resourceTypes {
		resp = append(resp, r)
	}

	d.logger.Debug(fmt.Sprintf("Total %v resource type exits", len(resp)))
	return resp, nil
}

// TODO: test & fix this
func (d *Kubernetes) GetAuthInfo() (*proto.AuthInfoResponse, error) {
	authInfo := new(proto.AuthInfoResponse)

	shared.LogDebug("GetAuthInfo Called")

	for _, kc := range d.kubeCLIconfig.KubeConfigs {
		shared.LogDebug("Path Is ==== %s", kc.Path)
		for _, c := range kc.Contexts {
			shared.LogDebug("Path is ==== Love %s %s", kc.Path, c.Name)
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

	// v, err := d.normalClient.ServerVersion()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to load server version: %w", err)
	// }

	// cc := d.clientConfig.CurrentContext
	// user := ""

	// context, ok := d.clientConfig.Contexts[cc]
	// if ok {
	// 	user = context.AuthInfo
	// }

	// server := ""
	// for _, c := range d.clientConfig.Clusters {
	// 	server = c.Server
	// }

	return authInfo, nil
}

// TODO: Include plural names as well
func (d *Kubernetes) GetResourceIsolatorType() (string, error) {
	return "namespaces", nil
}

func (d *Kubernetes) GetDefaultResourceIsolator() (string, error) {
	return "default", nil
}

func (d *Kubernetes) GetSupportedActions() (*proto.GetActionListResponse, error) {
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
			// {
			// 	Name:       "refresh",
			// 	KeyBinding: "ctrl-r",
			// 	OutputType: "nothing",
			// },
		},
	}

	return genericActions, nil
}

func (d *Kubernetes) ActionDeleteResource(args *proto.ActionDeleteResourceArgs) error {
	return d.deleteResource(context.Background(), args)
}

func (d *Kubernetes) ActionCreateResource(args *proto.ActionCreateResourceArgs) error {
	return d.createResource(context.Background(), args)
}

func (d *Kubernetes) ActionUpdateResource(args *proto.ActionUpdateResourceArgs) error {
	return d.updateResource(context.Background(), args)
}

func (d *Kubernetes) GetSpecficActionList(resourceType string) (*proto.GetActionListResponse, error) {
	t, ok := d.resourceTypeConfigurations[resourceType]
	if !ok {
		d.logger.Trace(fmt.Sprintf("specific action list schema of resource type %s not found, using the default schema", resourceType))
		t = d.resourceTypeConfigurations["defaults"]
	}

	return &proto.GetActionListResponse{Actions: t.SpecificActions}, nil
}

func (d *Kubernetes) PerformSpecificAction(args *proto.SpecificActionArgs) (*proto.SpecificActionResult, error) {

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
		// d.activeChans <- struct{}{}
		// return shared.SpecificActionResult{}, nil
	}

	return nil, nil
}

func (d *Kubernetes) decodeSecret(rawData interface{}) (string, error) {

	secretData := rawData.(map[string]interface{})["data"].(map[string]interface{})

	// data, err := json.Marshal(rawData)
	// if err != nil {
	// 	return err.Error()
	// }

	// obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(data, nil, nil)
	// if err != nil {
	// 	log.Fatalf(fmt.Sprintf("Error while decoding YAML object. Err was: %s", err))
	// 	return "", fmt.Errorf("")
	// }

	// switch o := obj.(type) {
	// case *v1.ConfigMap:

	// case *v1.Secret:
	// default:
	// 	fmt.Printf("Type %v is unknown", o)
	// }

	decodedMap := map[string]string{}
	for key, encodedData := range secretData {
		decodedData, err := base64.StdEncoding.DecodeString(string(encodedData.(string)))
		if err != nil {
			return "", shared.LogError("failed to decode base64 string: %v", err)
		}
		decodedMap[key] = string(decodedData)
	}

	rawData.(map[string]interface{})["data"] = decodedMap
	finalData, err := yaml.Marshal(rawData)
	if err != nil {
		return "", shared.LogError("failed to yaml marshal decode secret: %v", err)
	}

	return string(finalData), nil
}
