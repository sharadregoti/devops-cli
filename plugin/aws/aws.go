package aws

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/sharadregoti/devops/model"
	"gopkg.in/yaml.v2"
)

const PluginName = "aws"

var release bool = false

func getSchemaPath(devopsDir string) string {
	if release {
		return devopsDir + "/plugins/aws/resource_config"
	}
	return "../../plugin/aws/resource_config"
}

type AWS struct {
	logger                     hclog.Logger
	isOK                       error
	resourceTypeConfigurations map[string]model.ResourceTransfomer
}

func New(logger hclog.Logger) (*AWS, error) {
	// Read resource configs
	resourceSchemaTypeMap := map[string]model.ResourceTransfomer{}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return &AWS{logger: logger, isOK: err}, fmt.Errorf("failed to read directory: %w", err)
	}
	// Create the ".devops" subdirectory if it doesn't exist
	// TODO: This should be a function in the core binary
	devopsDir := filepath.Join(homeDir, ".devops")

	schemaPath := getSchemaPath(devopsDir)
	// Read all resource type scheam
	files, err := ioutil.ReadDir(schemaPath)
	if err != nil {
		return &AWS{logger: logger, isOK: err}, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, f := range files {
		if f.Name() != "defaults.yaml" {
			continue
		}

		data, err := os.ReadFile(schemaPath + "/" + f.Name())
		if err != nil {
			return &AWS{logger: logger, isOK: err}, fmt.Errorf("failed to read table schema file %s: %w", f.Name(), err)
		}

		res := new(model.ResourceTransfomer)
		if err := yaml.Unmarshal(data, res); err != nil {
			return &AWS{logger: logger, isOK: err}, fmt.Errorf("failed to unmarshal table schema data: %w", err)
		}
		resourceSchemaTypeMap[strings.TrimSuffix(f.Name(), ".yaml")] = *res
	}

	return &AWS{
		logger:                     logger,
		isOK:                       nil,
		resourceTypeConfigurations: resourceSchemaTypeMap,
	}, nil
}
