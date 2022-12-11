package main

import (
	"encoding/gob"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	aws "github.com/sharadregoti/devops/plugin/aws"
	pnk8s "github.com/sharadregoti/devops/plugin/kubernetes"
	"github.com/sharadregoti/devops/shared"
)

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "plugin-server",
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})

	a := &aws.AWS{}

	pluginK8s := pnk8s.New(logger)

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"aws":        &shared.DevopsPlugin{Impl: a},
		"kubernetes": &shared.DevopsPlugin{Impl: pluginK8s},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
