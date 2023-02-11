package shared

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/sharadregoti/devops/proto"
	"google.golang.org/grpc"
)

type DevopsPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl Devops
}

// func (p *DevopsPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
// 	return &DevopsServerRPC{Impl: p.Impl}, nil
// }

// func (DevopsPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
// 	var _ Devops = &DevopsClientRPC{}
// 	return &DevopsClientRPC{client: c}, nil
// }

func (p *DevopsPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterDevopsServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *DevopsPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	var _ Devops = &GRPCClient{}
	return &GRPCClient{client: proto.NewDevopsClient(c)}, nil
}
