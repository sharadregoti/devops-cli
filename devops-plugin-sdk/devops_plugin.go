package sdk

import (
	"context"

	"github.com/sharadregoti/devops-plugin-sdk/proto"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type DevopsPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Logger hclog.Logger
	Impl   Devops
}

func (p *DevopsPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterDevopsServer(s, &GRPCServer{Impl: p.Impl, Logger: p.Logger})
	return nil
}

func (p *DevopsPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	var _ Devops = &GRPCClient{}
	return &GRPCClient{client: proto.NewDevopsClient(c)}, nil
}
