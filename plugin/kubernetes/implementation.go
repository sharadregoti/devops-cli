package kubernetes

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/sharadregoti/devops/common"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/shared"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func (d *Kubernetes) Name() string {
	return PluginName
}

// TODO: test & fix this
func (d *Kubernetes) StatusOK() error {
	if d.isOK != nil {
		common.Error(d.logger, fmt.Sprintf("failed to load plugin: %v", errors.Unwrap(d.isOK)))
	}
	return d.isOK
}

func (d *Kubernetes) GetResources(args shared.GetResourcesArgs) ([]interface{}, error) {

	label := ""
	for k, v := range args.Args {
		if k == "labels" {
			selector := labels.NewSelector()
			for labelKey, labelValue := range v.(map[string]interface{}) {
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

	d.logger.Debug(fmt.Sprintf("Found %v %v resources in %v namespace", len(items), args.ResourceType, args.IsolatorID))
	return items, nil
}

// TODO: test & fix this
func (d *Kubernetes) CloseResourceWatcher(resourceType string) error {
	res, ok := d.resourceWatcherChanMap[resourceType]
	if !ok {
		return fmt.Errorf("channel for resource type %s does not exists", resourceType)
	}

	d.logger.Debug("Closing resource watcher", resourceType)
	close(res)
	return nil
}

// TODO: test & fix this
func (d *Kubernetes) WatchResources(resourceType string) (chan shared.WatchResourceResult, error) {
	res, ok := d.resourceWatcherChanMap[resourceType]
	if ok {
		// Send it
		d.logger.Debug("Channel already exists for resource", resourceType)
		return res, nil
	}

	d.logger.Debug("Creating a new channel for resource", resourceType)
	d.lock.Lock()
	c := make(chan shared.WatchResourceResult, 1)
	d.resourceWatcherChanMap[resourceType] = c
	d.lock.Unlock()

	go func() {
		r := d.resourceTypes[resourceType]
		err := getResourcesDynamically(c, d.dynamicClient, context.Background(), r.group, r.version, r.resourceTypeName, "default")
		if err != nil {
			common.Error(d.logger, fmt.Sprintf("failed to watch over resource %s: %w", resourceType, err))
			return
		}
	}()

	return c, nil
}

func (d *Kubernetes) GetResourceTypeSchema(resourceType string) (model.ResourceTransfomer, error) {
	t, ok := d.resourceTypeConfigurations[resourceType]
	if !ok {
		d.logger.Debug(fmt.Sprintf("Schema of resource type %s not found, using the default schema", resourceType))
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
func (d *Kubernetes) GetGeneralInfo() (map[string]string, error) {
	v, err := d.normalClient.ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to load server version: %w", err)
	}

	cc := d.clientConfig.CurrentContext
	user := ""

	context, ok := d.clientConfig.Contexts[cc]
	if ok {
		user = context.AuthInfo
	}

	server := ""
	for _, c := range d.clientConfig.Clusters {
		server = c.Server
	}

	return map[string]string{
		"Context":        cc,
		"Cluster":        server,
		"User":           user,
		"Server Version": v.String(),
	}, nil
}

// TODO: Include plural names as well
func (d *Kubernetes) GetResourceIsolatorType() (string, error) {
	return "namespaces", nil
}

func (d *Kubernetes) GetDefaultResourceIsolator() (string, error) {
	return "default", nil
}

func (d *Kubernetes) GetSupportedActions(resourceType string) (shared.GenericActions, error) {
	return shared.GenericActions{
		IsDelete: true,
		IsUpdate: false,
		IsCreate: false,
	}, nil
}

func (d *Kubernetes) ActionDeleteResource(args shared.ActionDeleteResourceArgs) error {
	return d.deleteResource(context.Background(), args)
}

func (d *Kubernetes) GetSpecficActionList(resourceType string) ([]shared.SpecificAction, error) {
	t, ok := d.resourceTypeConfigurations[resourceType]
	if !ok {
		d.logger.Info(fmt.Sprintf("specific action list schema of resource type %s not found, using the default schema", resourceType))
		t = d.resourceTypeConfigurations["defaults"]
	}

	arr := make([]shared.SpecificAction, 0)
	for _, sa := range t.SpecificActions {
		arr = append(arr, shared.SpecificAction{
			Name:         sa.Name,
			KeyBinding:   sa.KeyBinding,
			ScrrenAction: sa.ScrrenAction,
			OutputType:   sa.OutputType,
			ResourceID:   "",
		})
	}
	return arr, nil
}

func (d *Kubernetes) PerformSpecificAction(args shared.SpecificActionArgs) (shared.SpecificActionResult, error) {

	switch args.ActionName {

	case "describe":
		result, err := d.DescribeResource(args.ResourceType, args.ResourceName, args.IsolatorName)
		if err != nil {
			return shared.SpecificActionResult{}, err
		}

		return shared.SpecificActionResult{
			Result: result,
			// TODO: Output type should come from an core SDK
			OutputType: "string",
		}, nil

	case "decode-secret":
		items, err := d.GetResources(shared.GetResourcesArgs{
			ResourceName: args.ResourceName,
			ResourceType: args.ResourceType,
			IsolatorID:   args.IsolatorName,
		})
		if err != nil {
			return shared.SpecificActionResult{}, err
		}

		result, err := d.decodeSecret(items[0])
		if err != nil {
			return shared.SpecificActionResult{}, err
		}

		return shared.SpecificActionResult{
			Result: result,
			// TODO: Output type should come from an core SDK
			OutputType: "string",
		}, nil

	case "logs":

		containerName := ""
		if args.ResourceType == "containers" {
			parentName := args.Args["parentName"]
			args.ResourceType = "pods"
			containerName = args.ResourceName
			args.ResourceName = parentName.(string)
		}

		res, err := d.getPodLogs(args.ResourceName, args.IsolatorName, containerName)
		if err != nil {
			return shared.SpecificActionResult{}, err
		}

		return shared.SpecificActionResult{
			Result:     res,
			OutputType: "string",
		}, nil

	case "shell":

		containerName := ""
		if args.ResourceType == "containers" {
			parentName := args.Args["parentName"]
			args.ResourceType = "pods"
			containerName = args.ResourceName
			args.ResourceName = parentName.(string)
		}

		res, err := d.execPod(args.ResourceName, args.IsolatorName, containerName)
		if err != nil {
			return shared.SpecificActionResult{}, err
		}

		return shared.SpecificActionResult{
			Result:     res,
			OutputType: "string",
		}, nil

	case "port-forward":

		res, err := d.portForward(args.ResourceName, args.IsolatorName, args.Args)
		if err != nil {
			return shared.SpecificActionResult{}, err
		}

		return shared.SpecificActionResult{
			Result:     res,
			OutputType: "string",
		}, nil

	case "view-pods":
		res, err := d.getPods(context.Background(), args.IsolatorName, args.ResourceName, args.ResourceType)
		if err != nil {
			return shared.SpecificActionResult{}, err
		}
		return shared.SpecificActionResult{
			Result: res,
			// TODO: Output type should come from an core SDK
			OutputType: "invoke-event",
		}, nil

	case "close":
		d.activeChans <- struct{}{}
		return shared.SpecificActionResult{}, nil
	}

	return shared.SpecificActionResult{}, nil
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
			common.Error(d.logger, fmt.Sprintf("failed to decode base64 string: %v", err))
			return "", err
		}
		decodedMap[key] = string(decodedData)
	}

	rawData.(map[string]interface{})["data"] = decodedMap
	finalData, err := yaml.Marshal(rawData)
	if err != nil {
		common.Error(d.logger, fmt.Sprintf("failed to yaml marshal decode secret: %v", err))
		return "", err
	}

	return string(finalData), nil
}
