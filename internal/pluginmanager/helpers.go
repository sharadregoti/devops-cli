package pluginmanager

import (
	"fmt"

	"github.com/sharadregoti/devops/common"
)

func getPluginPath(name, devopsDir string) string {
	if common.Release {
		return fmt.Sprintf("%s/plugins/%s/%s", devopsDir, name, name)
	}
	return fmt.Sprintf("../../plugins/%s/%s", name, name)
}
	