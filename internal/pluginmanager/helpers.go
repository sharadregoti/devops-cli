package pluginmanager

import (
	"fmt"
	"os"

	"github.com/sharadregoti/devops/common"
	"github.com/sharadregoti/devops/model"
)

func getPluginPath(name, devopsDir string) string {
	if common.Release {
		return fmt.Sprintf("%s/plugins/%s/%s", devopsDir, name, name)
	}
	return fmt.Sprintf("../../plugins/%s/%s", name, name)
}

func ListPlugins() ([]*model.Plugin, error) {
	devopsDir := model.InitCoreDirectory()

	var path string = "../../plugins"
	if common.Release {
		path = fmt.Sprintf("%s/plugins", devopsDir)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var plugins []*model.Plugin
	for _, entry := range entries {
		if entry.IsDir() {
			plugins = append(plugins, &model.Plugin{Name: entry.Name(), IsDefault: entry.Name() == "kubernetes"})
		}
	}
	return plugins, nil
}
