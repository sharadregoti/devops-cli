package main

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/sharadregoti/devops/common"
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
	go common.ConnLoggingInit()

	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "kubernetes-plugin-server",
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
	gob.Register(make(chan shared.WatchResourceResult))

	logger.Info("Starting kubernetes plugin server")

	pluginK8s, err := pnk8s.New(logger.Named(pnk8s.PluginName))
	if err != nil {
		common.Error(logger, fmt.Sprintf("failed to initialized kubernetes plugin: %v", err))
		// os.Exit(1)
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"kubernetes": &shared.DevopsPlugin{Impl: pluginK8s},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Logger:          logger,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
