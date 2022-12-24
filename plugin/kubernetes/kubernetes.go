package kubernetes

import (
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

var release bool = false

func getSchemaPath(devopsDir string) string {
	if release {
		return devopsDir + "/plugins/kubernetes/resource_config"
	}
	return "../../plugin/kubernetes/resource_config"
}

const PluginName = "kubernetes"

type Kubernetes struct {
	lock sync.RWMutex

	logger hclog.Logger

	isOK error

	// This channel listen for close streaming event from server & closes the corresponding process which is running in background
	activeChans chan struct{}

	// Clients
	normalClient  *kubernetes.Clientset
	dynamicClient dynamic.Interface

	config *rest.Config

	// Kube Config Parser
	clientConfig clientcmdapi.Config

	// Key is resource type, All resources are stored in their plural form
	resourceTypes resourceTypeList

	// Stores mapping of resource types corresponding to a schema defined in file
	// Key is resource type
	resourceTypeConfigurations map[string]model.ResourceTransfomer

	// Key is resource type
	resourceWatcherChanMap map[string]chan shared.WatchResourceResult
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

	// TODO: Use the normal way don't die here
	config := ctrl.GetConfigOrDie()

	// Normal client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &Kubernetes{logger: logger, isOK: err}, fmt.Errorf("failed to load kube config: %w", err)
	}

	// Dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return &Kubernetes{logger: logger, isOK: err}, fmt.Errorf("failed to load dynamic kube config: %w", err)
	}

	// Get details from kube-config
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

	// Read resource configs
	resourceSchemaTypeMap := map[string]model.ResourceTransfomer{}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return &Kubernetes{logger: logger, isOK: err}, fmt.Errorf("failed to read directory: %w", err)
	}
	// Create the ".devops" subdirectory if it doesn't exist
	// TODO: This should be a function in the core binary
	devopsDir := filepath.Join(homeDir, ".devops")

	schemaPath := getSchemaPath(devopsDir)
	// Read all resource type scheam
	files, err := ioutil.ReadDir(schemaPath)
	if err != nil {
		return &Kubernetes{logger: logger, isOK: err}, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, f := range files {
		data, err := os.ReadFile(schemaPath + "/" + f.Name())
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
		logger:                     logger,
		normalClient:               clientset,
		dynamicClient:              dynamicClient,
		clientConfig:               cc,
		config:                     config,
		resourceTypes:              resourceTypeMap,
		resourceTypeConfigurations: resourceSchemaTypeMap,
		resourceWatcherChanMap:     make(map[string]chan shared.WatchResourceResult),
		activeChans:                make(chan struct{}, 1),
	}, nil
}
