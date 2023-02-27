package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	shared "github.com/sharadregoti/devops-plugin-sdk"
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
		Name:   "k8s",
		Level:  shared.GetHCLLogLevel(),
		Output: os.Stderr,
	})
	shared.Logger = logger

	shared.LogInfo("Starting kubernetes plugin server")

	pluginK8s, err := New(logger.Named(PluginName))
	if err != nil {
		shared.LogError("failed to initialized kubernetes plugin: %v", err)
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"kubernetes": &shared.DevopsPlugin{Impl: pluginK8s, Logger: logger},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Logger:          logger,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
