package shared

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/sharadregoti/devops/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type GRPCServer struct {
	// This is the real implementation
	Impl Devops
	proto.UnimplementedDevopsServer
}

func (g *GRPCServer) Name(ctx context.Context, e *empty.Empty) (*wrappers.StringValue, error) {
	// Implement the Name method here
	return wrapperspb.String(g.Impl.Name()), nil
}

func (g *GRPCServer) GetResources(ctx context.Context, args *proto.GetResourcesArgs) (*_struct.ListValue, error) {
	// Implement the GetResources method here
	res, err := g.Impl.GetResources(args)
	if err != nil {
		return &_struct.ListValue{}, err
	}

	result := &_struct.ListValue{}
	for _, r := range res {
		m, err := structpb.NewValue(r)
		if err != nil {
			return &_struct.ListValue{}, err
		}
		result.Values = append(result.Values, m)
	}

	return result, nil
}

func (g *GRPCServer) WatchResources(sw *wrappers.StringValue, server proto.Devops_WatchResourcesServer) error {
	// Implement the WatchResources method here
	ch, done, err := g.Impl.WatchResources(sw.Value)
	if err != nil {
		return err
	}

	for {
		select {
		case <-done:
			fmt.Println("Closing server routine")
			return nil
		case v := <-ch:
			m, err := structpb.NewValue(v)
			if err != nil {
				return err
			}
			server.Send(&proto.WatchResourceResult{Type: v.Type, Result: m})
		}
	}

}

func (g *GRPCServer) CloseResourceWatcher(ctx context.Context, w *wrappers.StringValue) (*empty.Empty, error) {
	// Implement the CloseResourceWatcher method here
	return nil, g.Impl.CloseResourceWatcher(w.Value)
}

func (g *GRPCServer) GetResourceTypeSchema(ctx context.Context, t *wrappers.StringValue) (*proto.ResourceTransformer, error) {
	// Implement the GetResourceTypeSchema method here
	return g.Impl.GetResourceTypeSchema(t.Value)
}

func (g *GRPCServer) GetResourceTypeList(ctx context.Context, e *empty.Empty) (*proto.GetResourceTypeListResponse, error) {
	// Implement the GetResourceTypeList method here
	list, err := g.Impl.GetResourceTypeList()
	if err != nil {
		return &proto.GetResourceTypeListResponse{}, err
	}

	return &proto.GetResourceTypeListResponse{ResourceType: list}, nil
}

func (g *GRPCServer) GetAuthInfo(ctx context.Context, e *empty.Empty) (*proto.AuthInfoResponse, error) {
	// Implement the GetAuthInfo method here
	return g.Impl.GetAuthInfo()
}

func (g *GRPCServer) Connect(ctx context.Context, a *proto.AuthInfo) (*empty.Empty, error) {
	// Implement the Connect method here
	return &emptypb.Empty{}, g.Impl.Connect(a)
}

func (g *GRPCServer) GetResourceIsolatorType(ctx context.Context, e *empty.Empty) (*wrappers.StringValue, error) {
	// Implement the GetResourceIsolatorType method here
	str, err := g.Impl.GetResourceIsolatorType()
	return wrapperspb.String(str), err
}

func (g *GRPCServer) GetDefaultResourceIsolator(ctx context.Context, e *empty.Empty) (*wrappers.StringValue, error) {
	// Implement the GetDefaultResourceIsolator method here
	str, err := g.Impl.GetDefaultResourceIsolator()
	return wrapperspb.String(str), err
}

func (g *GRPCServer) GetSupportedActions(ctx context.Context, e *empty.Empty) (*proto.GetActionListResponse, error) {
	// Implement the GetSupportedActions method here
	return g.Impl.GetSupportedActions()
}

func (g *GRPCServer) ActionDeleteResource(ctx context.Context, args *proto.ActionDeleteResourceArgs) (*empty.Empty, error) {
	// Implement the ActionDeleteResource method here
	return &emptypb.Empty{}, g.Impl.ActionDeleteResource(args)
}

func (g *GRPCServer) ActionCreateResource(ctx context.Context, args *proto.ActionCreateResourceArgs) (*empty.Empty, error) {
	// Implement the ActionCreateResource method here

	return &emptypb.Empty{}, g.Impl.ActionCreateResource(args)
}

func (g *GRPCServer) ActionUpdateResource(ctx context.Context, args *proto.ActionUpdateResourceArgs) (*empty.Empty, error) {
	// Implement the ActionUpdateResource method here
	return &emptypb.Empty{}, g.Impl.ActionUpdateResource(args)
}

func (g *GRPCServer) GetSpecficActionList(ctx context.Context, t *wrappers.StringValue) (*proto.GetActionListResponse, error) {
	// Implement the GetSpecficActionList method here
	return g.Impl.GetSpecficActionList(t.Value)
}

func (g *GRPCServer) PerformSpecificAction(ctx context.Context, args *proto.SpecificActionArgs) (*proto.SpecificActionResult, error) {
	return g.Impl.PerformSpecificAction(args)
}
