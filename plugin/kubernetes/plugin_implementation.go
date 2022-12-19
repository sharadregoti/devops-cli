package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/go-hclog"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/shared"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Kubernetes struct {
	lock sync.RWMutex

	logger hclog.Logger

	isOK error

	// Clients
	normalClient  *kubernetes.Clientset
	dynamicClient dynamic.Interface

	config *rest.Config

	// Kube Config Parser
	clientConfig clientcmdapi.Config

	// Key is resource type, All resources are stored in their plural form
	resourceTypeMap resourceTypeList

	// Stores mapping of resource types corresponding to a schema defined in file
	// Key is resource type
	resourceSchemaTypeMap map[string]model.ResourceTransfomer

	// Key is resource type
	resourceChanMap map[string]chan shared.WatchResourceResult
}

type resourceTypeList map[string]*resourceTypeInfo

type resourceTypeInfo struct {
	group            string
	version          string
	resourceTypeName string
	isNamespaced     bool
}

func New(logger hclog.Logger) (*Kubernetes, error) {

	// Check if the kubectl command exists
	_, err := exec.LookPath("kubectl")
	if err != nil {
		return &Kubernetes{logger: logger, isOK: fmt.Errorf("kubectl command not found: %w", err)}, err
	}

	config := ctrl.GetConfigOrDie()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &Kubernetes{logger: logger, isOK: err}, fmt.Errorf("failed to load kube config: %w", err)
	}

	// Create a dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return &Kubernetes{logger: logger, isOK: err}, fmt.Errorf("failed to load dynamic kube config: %w", err)
	}

	cc, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{},
		&clientcmd.ConfigOverrides{}).RawConfig()
	if err != nil {
		return &Kubernetes{logger: logger, isOK: err}, fmt.Errorf("failed to load kubernetes context: %w", err)
	}

	// List all supported resources
	resources, err := clientset.Discovery().ServerPreferredResources()
	if err != nil {
		return &Kubernetes{logger: logger, isOK: err}, fmt.Errorf("failed to discover kubernetes resource types: %w", err)
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

			resourceTypeMap[r.Name] = &resourceTypeInfo{
				group:            group,
				version:          version,
				resourceTypeName: r.Name,
				isNamespaced:     r.Namespaced,
			}
		}
	}

	resourceSchemaTypeMap := map[string]model.ResourceTransfomer{}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return &Kubernetes{logger: logger, isOK: err}, fmt.Errorf("failed to read directory: %w", err)
	}
	// Create the ".devops" subdirectory if it doesn't exist
	devopsDir := filepath.Join(homeDir, ".devops")

	// Read all resource type scheam
	files, err := ioutil.ReadDir(devopsDir + "/plugins/table_schema")
	if err != nil {
		return &Kubernetes{logger: logger, isOK: err}, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, f := range files {
		data, err := os.ReadFile(devopsDir + "/plugins/table_schema/" + f.Name())
		if err != nil {
			return &Kubernetes{logger: logger, isOK: err}, fmt.Errorf("failed to read table schema file %s: %w", f.Name(), err)
		}

		res := new(model.ResourceTransfomer)
		if err := yaml.Unmarshal(data, res); err != nil {
			return &Kubernetes{logger: logger, isOK: err}, fmt.Errorf("failed to unmarshal table schema data: %w", err)
		}
		resourceSchemaTypeMap[strings.TrimSuffix(f.Name(), ".yaml")] = *res
	}

	return &Kubernetes{
		logger:                logger,
		normalClient:          clientset,
		dynamicClient:         dynamicClient,
		clientConfig:          cc,
		config:                config,
		resourceTypeMap:       resourceTypeMap,
		resourceChanMap:       make(map[string]chan shared.WatchResourceResult),
		resourceSchemaTypeMap: resourceSchemaTypeMap,
	}, nil
}

func (d *Kubernetes) Name() string {
	return "kubernetes"
}

func (d *Kubernetes) StatusOK() error {
	d.logger.Error(fmt.Sprintf("failed to load plugin: %v", errors.Unwrap(d.isOK)))
	return d.isOK
}

func (d *Kubernetes) GetResources(args shared.GetResourcesArgs) ([]interface{}, error) {
	r := d.resourceTypeMap[args.ResourceType]
	if !r.isNamespaced || args.IsolatorID == "all" {
		args.IsolatorID = ""
	}
	items, err := listResourcesDynamically(d.dynamicClient, context.Background(), r.group, r.version, r.resourceTypeName, args.IsolatorID)
	if err != nil {
		return nil, err
	}

	d.logger.Debug(fmt.Sprintf("Found %v %v resources in namespace %v", len(items), args.ResourceType, args.IsolatorID))
	return items, nil
}

func (d *Kubernetes) CloseResourceWatcher(resourceType string) error {
	res, ok := d.resourceChanMap[resourceType]
	if !ok {
		return fmt.Errorf("channel for resource type %s does not exists", resourceType)
	}

	d.logger.Debug("Closing resource watcher", resourceType)
	close(res)
	return nil
}

func (d *Kubernetes) WatchResources(resourceType string) (chan shared.WatchResourceResult, error) {
	res, ok := d.resourceChanMap[resourceType]
	if ok {
		// Send it
		d.logger.Debug("Channel already exists for resource", resourceType)
		return res, nil
	}

	d.logger.Debug("Creating a new channel for resource", resourceType)
	d.lock.Lock()
	c := make(chan shared.WatchResourceResult, 1)
	d.resourceChanMap[resourceType] = c
	d.lock.Unlock()

	go func() {
		r := d.resourceTypeMap[resourceType]
		err := getResourcesDynamically(c, d.dynamicClient, context.Background(), r.group, r.version, r.resourceTypeName, "default")
		if err != nil {
			d.logger.Error("failed to watch over resource %s: %w", resourceType, err)
			return
		}
	}()

	return c, nil
}

func (d *Kubernetes) GetResourceTypeSchema(resourceType string) (model.ResourceTransfomer, error) {
	t, ok := d.resourceSchemaTypeMap[resourceType]
	if !ok {
		d.logger.Info(fmt.Sprintf("schema of resource type %s not found, using the default schema", resourceType))
		return d.resourceSchemaTypeMap["defaults"], nil
	}

	return t, nil
}

func (d *Kubernetes) GetResourceTypeList() ([]string, error) {
	resp := []string{}
	for r := range d.resourceTypeMap {
		resp = append(resp, r)
	}
	return resp, nil
}

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

// Use plural names
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
	r := d.resourceTypeMap[args.ResourceType]
	if !r.isNamespaced || args.IsolatorName == "all" {
		args.IsolatorName = ""
	}

	return deleteResourcesDynamically(d.dynamicClient, context.Background(), r.group, r.version, r.resourceTypeName, args.IsolatorName, args.ResourceName)
}

func (d *Kubernetes) GetSpecficActionList(resourceType string) ([]shared.SpecificAction, error) {
	return []shared.SpecificAction{
		{
			Name:         "describe",
			KeyBinding:   "d",
			ScrrenAction: "view",
			OutputType:   "string",
		},
	}, nil
}

func (d *Kubernetes) PerformSpecificAction(args shared.SpecificActionArgs) (shared.SpecificActionResult, error) {

	switch args.ActionName {

	case "describe":
		result, err := d.DescribeResource(args.ResourceType, args.ResourceName, args.IsolatorName)
		if err != nil {
			return shared.SpecificActionResult{}, err
		}

		return shared.SpecificActionResult{
			Result:     result,
			OutputType: "string",
		}, nil
	}

	return shared.SpecificActionResult{}, nil
}
