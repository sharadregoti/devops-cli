package core

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/sharadregoti/devops/common"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/shared"
)

type PluginClient struct {
	client plugin.ClientProtocol
	logger hclog.Logger

	name string

	stdErrReader *io.PipeReader
	stdErrWriter *io.PipeWriter

	stdOutReader *io.PipeReader
	stdOutWriter *io.PipeWriter
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

func checIfPluginExists(rooDir string, c *model.Config) {
	for _, p := range c.Plugins {
		fmt.Printf("Checking plugin %s\n", p.Name)
		_, err := os.Stat(getPluginPath(p.Name, rooDir))
		if os.IsNotExist(err) {
			log.Fatalf("Plugin %s does not exists, use devops init command to install the plugin", p.Name)
		}
	}
}

func loadPlugin(logger hclog.Logger, pluginName, rootDir string) (*PluginClient, error) {
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
	gob.Register(make(chan shared.WatchResourceResult))
	gob.Register(time.Time{})

	path := getPluginPath(pluginName, rootDir)
	logger.Debug("Pluging path", path)

	var reader, writer = io.Pipe()

	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(path),
		Logger:          logger,
		SyncStdout:      writer,
		SyncStderr:      writer,
	})
	// defer client.Kill()

	if client.Exited() {
		str := fmt.Sprintf("%s plugin exited", pluginName)
		common.Error(logger, str)
		return nil, fmt.Errorf(str)
	}

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		common.Error(logger, fmt.Sprintf("Failed to initialzed plugin client: %v", err))
		return nil, err
	}

	return &PluginClient{
		name: pluginName,
		// TODO: Set error writer as well
		stdErrReader: nil,
		stdErrWriter: nil,
		stdOutReader: reader,
		stdOutWriter: writer,
		client:       rpcClient,
		logger:       logger,
	}, nil
}

func (p *PluginClient) GetStdoutReader() *io.PipeReader {
	return p.stdOutReader
}

func (p *PluginClient) Close() {
	// TODO: close error writer as well
	p.stdOutReader.Close()
	p.stdOutWriter.Close()
	p.client.Close()
}

func (p *PluginClient) GetPlugin(name string) (shared.Devops, error) {
	// Request the plugin
	raw, err := p.client.Dispense(name)
	if err != nil {
		common.Error(p.logger, fmt.Sprintf("failed to get plugin: %v, %v", name, err))
		return nil, err
	}

	return raw.(shared.Devops), nil
}
