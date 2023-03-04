package pluginmanager

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	shared "github.com/sharadregoti/devops-plugin-sdk"
	"github.com/sharadregoti/devops/model"
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
	"helm":       &shared.DevopsPlugin{},
	"kubernetes": &shared.DevopsPlugin{},
}

func ValidatePlugins(c *model.Config) error {
	devopsDir := model.InitCoreDirectory()

	for _, p := range c.Plugins {
		fmt.Printf("Checking plugin %s\n", p.Name)
		_, err := os.Stat(getPluginPath(p.Name, devopsDir))
		if os.IsNotExist(err) {
			return fmt.Errorf("Plugin %s does not exists, use devops init command to install the plugin", p.Name)
		}
	}

	return nil
}

func startPlugin(logger hclog.Logger, pluginName, rootDir string) (*PluginClient, error) {
	path := getPluginPath(pluginName, rootDir)
	logger.Debug("Pluging path", path)

	var reader, writer = io.Pipe()

	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  handshakeConfig,
		Plugins:          pluginMap,
		Cmd:              exec.Command(path),
		Logger:           logger,
		SyncStdout:       writer,
		SyncStderr:       writer,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	// TODO: Closer this as well in the Close function
	// defer client.Kill()

	if client.Exited() {
		str := fmt.Sprintf("%s plugin exited", pluginName)
		return nil, fmt.Errorf(str)
	}

	// Connect via GRPC
	rpcClient, err := client.Client()
	if err != nil {
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

func (p *PluginClient) GetStdoutWriter() *io.PipeWriter {
	return p.stdOutWriter
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
		return nil, err
	}

	return raw.(shared.Devops), nil
}
