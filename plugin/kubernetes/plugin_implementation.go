package kubernetes

import (
	"context"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/go-hclog"
	"github.com/sharadregoti/devops/model"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Kubernetes struct {
	logger          hclog.Logger
	normalClient    *kubernetes.Clientset
	dynamicClient   dynamic.Interface
	clientConfig    clientcmdapi.Config
	resourceTypeMap resourceTypeList
}

type resourceTypeList map[string]*resourceTypeInfo

type resourceTypeInfo struct {
	group            string
	version          string
	resourceTypeName string
}

func New(logger hclog.Logger) *Kubernetes {
	logger = logger.ResetNamed("kubernetes")
	config := ctrl.GetConfigOrDie()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Create a dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	cc, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{},
		&clientcmd.ConfigOverrides{}).RawConfig()
	if err != nil {
		panic(err)
	}

	// List all supported resources
	resources, err := clientset.Discovery().ServerPreferredResources()
	if err != nil {
		panic(err)
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
			}
		}
	}

	return &Kubernetes{
		logger:          logger,
		normalClient:    clientset,
		dynamicClient:   dynamicClient,
		clientConfig:    cc,
		resourceTypeMap: resourceTypeMap,
	}
}

func (d *Kubernetes) Name() string {
	return "kubernetes"
}

func (d *Kubernetes) GetResources(resourceType string) []interface{} {
	r := d.resourceTypeMap[resourceType]
	items, err := getResourcesDynamically(d.dynamicClient, context.Background(), r.group, r.version, r.resourceTypeName, "default")
	if err != nil {
		panic(err)
	}
	return items
}

func (d *Kubernetes) GetResourceTypeSchema(resourceType string) model.ResourceTransfomer {
	data, err := os.ReadFile("plugin/kubernetes/table_schema/pods.yaml")
	if err != nil {
		panic(err)
	}

	res := new(model.ResourceTransfomer)
	if err := yaml.Unmarshal(data, res); err != nil {
		panic(err)
	}

	return *res
}

func (d *Kubernetes) GetResourceTypeList() []string {
	resp := []string{}
	for r := range d.resourceTypeMap {
		resp = append(resp, r)
	}
	return resp
}

func (d *Kubernetes) GetGeneralInfo() map[string]string {
	v, err := d.normalClient.ServerVersion()
	if err != nil {
		panic(err)
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
	}
}

func (d *Kubernetes) GetResourceIsolatorType() string {
	return "namespace"
}

func (d *Kubernetes) GetDefaultResourceIsolator() string {
	return "all"
}
