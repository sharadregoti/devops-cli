package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/go-hclog"
	shared "github.com/sharadregoti/devops-plugin-sdk"
	"github.com/sharadregoti/devops-plugin-sdk/proto"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var release bool = false

func getConfigPath(devopsDir string) string {
	if release {
		return devopsDir + "/plugins/helm"
	}
	return "../../plugins/helm"
}

func getSchemaPath(devopsDir string) string {
	if release {
		return devopsDir + "/plugins/helm/resource_config"
	}
	return "../../plugins/helm/resource_config"
}

const PluginName = "helm"

type Helm struct {
	logger hclog.Logger
	isOK   error
	// Stores mapping of resource types corresponding to a schema defined in file
	// Key is resource type
	resourceTypeConfigurations map[string]*proto.ResourceTransformer

	// Key is resource type
	// resourceWatcherChanMap map[string]channels
	kubeCLIconfig *Config

	// Key is resource type
	resourceWatcherChanMap map[string]channels

	currentKubeConfigPath string

	restConfig     *rest.Config
	normalClient   *kubernetes.Clientset
	currentContext string
}

type channels struct {
	watcherDone chan struct{}
	serverDone  chan struct{}
}

type Config struct {
	KubeConfigs []*KubeConfigs `json:"kube_configs" yaml:"kube_configs"`
}
type Contexts struct {
	Name                    string   `json:"name" yaml:"name"`
	DefaultNamespacesToShow []string `json:"default_namespaces_to_show" yaml:"default_namespaces_to_show"`
	ReadOnly                bool     `json:"read_only" yaml:"read_only"`
	IsDefault               bool     `json:"is_default" yaml:"is_default"`
}
type KubeConfigs struct {
	Name     string      `json:"name" yaml:"name"`
	Path     string      `json:"path" yaml:"path"`
	Contexts []*Contexts `json:"contexts" yaml:"contexts"`
}

func New(logger hclog.Logger) (*Helm, error) {

	// Check if the kubectl command exists
	_, err := exec.LookPath("helm")
	if err != nil {
		return &Helm{logger: logger, isOK: shared.LogError("helm command not found: %v", err)}, err
	}

	// Read resource configs
	resourceSchemaTypeMap := map[string]*proto.ResourceTransformer{}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return &Helm{logger: logger, isOK: err}, shared.LogError("failed to read directory: %v", err)
	}
	// Create the ".devops" subdirectory if it doesn't exist
	// TODO: This should be a function in the core binary
	devopsDir := filepath.Join(homeDir, ".devops")

	schemaPath := getSchemaPath(devopsDir)
	// Read all resource type scheam
	files, err := ioutil.ReadDir(schemaPath)
	if err != nil {
		return &Helm{logger: logger, isOK: err}, shared.LogError("failed to read directory: %v", err)
	}

	for _, f := range files {
		if f.Name() == "config.yaml" {
			continue
		}

		data, err := os.ReadFile(schemaPath + "/" + f.Name())
		if err != nil {
			return &Helm{logger: logger, isOK: err}, shared.LogError("failed to read table schema file %s: %v", f.Name(), err)
		}

		res := new(proto.ResourceTransformer)
		if err := yaml.Unmarshal(data, res); err != nil {
			return &Helm{logger: logger, isOK: err}, shared.LogError("failed to unmarshal table schema data: %v", err)
		}
		resourceSchemaTypeMap[strings.TrimSuffix(f.Name(), ".yaml")] = res
	}

	data, err := os.ReadFile(getConfigPath(devopsDir) + "/" + "config.yaml")
	if err != nil {
		return &Helm{logger: logger, isOK: err}, shared.LogError("failed to read config.yaml file %s: %v", "config.yaml", err)
	}

	res := new(Config)
	if err := yaml.Unmarshal(data, res); err != nil {
		return &Helm{logger: logger, isOK: err}, shared.LogError("failed to unmarshal config.yaml data: %v", err)
	}

	defaultKubeConfigLocation := filepath.Join(homeDir, ".kube", "config")

	defaultFound := false
	for _, kubeConfig := range res.KubeConfigs {
		if kubeConfig.Path == defaultKubeConfigLocation {
			if kubeConfig.Name == "" {
				kubeConfig.Name = "default"
			}
			defaultFound = true
		}
	}

	if !defaultFound {
		// Add default config
		res.KubeConfigs = append(res.KubeConfigs, &KubeConfigs{
			Name: "default",
			Path: defaultKubeConfigLocation},
		)
	}

	// Expand all the contexts
	for i, kc := range res.KubeConfigs {
		kubeconfig, err := clientcmd.LoadFromFile(kc.Path)
		if err != nil {
			return &Helm{logger: logger, isOK: err}, shared.LogError("failed to build config from flags: %v", err)
		}

		for k, ctx := range kubeconfig.Contexts {

			isContextFound := false
			for _, c := range kc.Contexts {
				if c.Name == k {
					isContextFound = true
				}
			}

			namespaces := []string{"all", "default"}
			if ctx.Namespace != "" && ctx.Namespace != "default" {
				namespaces = append(namespaces, ctx.Namespace)
			}

			for i, localConfigCtx := range kc.Contexts {
				isAllFound := false
				for _, n := range localConfigCtx.DefaultNamespacesToShow {
					if n == "all" {
						isAllFound = true
					}
				}

				if !isAllFound {
					kc.Contexts[i].DefaultNamespacesToShow = append(kc.Contexts[i].DefaultNamespacesToShow, "all")
				}
			}

			if !isContextFound {
				res.KubeConfigs[i].Contexts = append(res.KubeConfigs[i].Contexts, &Contexts{
					Name:                    k,
					IsDefault:               kubeconfig.CurrentContext == k,
					DefaultNamespacesToShow: namespaces,
				})
			}

		}

		for _, c := range res.KubeConfigs[i].Contexts {
			c.IsDefault = kubeconfig.CurrentContext == c.Name
		}
	}

	return &Helm{
		logger:                     logger,
		resourceTypeConfigurations: resourceSchemaTypeMap,
		resourceWatcherChanMap:     make(map[string]channels),
		kubeCLIconfig:              res,
	}, nil
}
