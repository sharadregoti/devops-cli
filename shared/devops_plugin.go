package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type DevopsPlugin struct {
	Impl Devops
}

func (p *DevopsPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &DevopsServerRPC{Impl: p.Impl}, nil
}

func (DevopsPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	var _ Devops = &DevopsClientRPC{}
	return &DevopsClientRPC{client: c}, nil
}
