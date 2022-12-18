package main

import (
	"encoding/gob"
	"os/exec"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/sharadregoti/devops/shared"
)

type PluginClient struct {
	client plugin.ClientProtocol
	logger hclog.Logger
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"aws":        &shared.DevopsPlugin{},
	"kubernetes": &shared.DevopsPlugin{},
}

func New(logger hclog.Logger) (*PluginClient, error) {
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
	gob.Register(make(chan shared.WatchResourceResult))

	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("./plugin/plugin"),
		Logger:          logger,
	})
	// defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		logger.Error("Failed to initialzed plugin client", err)
		return nil, err
	}

	return &PluginClient{client: rpcClient, logger: logger}, nil
}

func (p *PluginClient) Close() {
	p.client.Close()
}

func (p *PluginClient) GetPlugin(name string) (shared.Devops, error) {
	// Request the plugin
	raw, err := p.client.Dispense(name)
	if err != nil {
		p.logger.Error("failed to get plugin", name, err)
		return nil, err
	}

	return raw.(shared.Devops), nil
}

// func installPlugins(logger hclog.Logger) shared.Devops {
// 	gob.Register(map[string]interface{}{})
// 	gob.Register([]interface{}{})
// 	gob.Register(make(chan shared.WatchResourceResult))

// 	// We're a host! Start by launching the plugin process.
// 	client := plugin.NewClient(&plugin.ClientConfig{
// 		HandshakeConfig: handshakeConfig,
// 		Plugins:         pluginMap,
// 		Cmd:             exec.Command("./plugin/plugin"),
// 		Logger:          logger,
// 	})
// 	// defer client.Kill()

// 	// Connect via RPC
// 	rpcClient, err := client.Client()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Request the plugin
// 	raw, err := rpcClient.Dispense("kubernetes")
// 	if err != nil {
// 		logger.Error("Failed to get plugin", err)
// 		log.Fatal(err)
// 	}

// 	// We should have a Greeter now! This feels like a normal interface
// 	// implementation but is in fact over an RPC connection.
// 	return raw.(shared.Devops)
// }
