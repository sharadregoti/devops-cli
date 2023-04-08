package sdk

import (
	"context"
	"fmt"
	"log"

	"github.com/sharadregoti/devops-plugin-sdk/proto"

	"github.com/golang/protobuf/ptypes/empty"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type GRPCServer struct {
	// This is the real implementation
	Logger hclog.Logger
	Impl   Devops
	proto.UnimplementedDevopsServer
}

func (g *GRPCServer) Name(ctx context.Context, e *empty.Empty) (*wrappers.StringValue, error) {
	return wrapperspb.String(g.Impl.Name()), nil
}

func (g *GRPCServer) GetResources(ctx context.Context, args *proto.GetResourcesArgs) (*_struct.ListValue, error) {
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

func (g *GRPCServer) WatchResources(args *proto.GetResourcesArgs, server proto.Devops_WatchResourcesServer) error {
	// fmt package does not work here. Only log, hclog works
	ch, done, err := g.Impl.WatchResources(args)
	if err != nil {
		// TODO: Errors sent from here does not reach the core client -> core binary
		log.Printf("Error has occured watcher 3")
		return err
	}

	log.Printf("grpc server routine: resource watcher has been started for resource type (%s)", args.ResourceType)
	defer log.Printf("grpc server routine: resource watcher has been stopped for resource type (%s)", args.ResourceType)

	for {
		select {
		case <-done:
			log.Printf("grpc server resource watcher routine: Done received for resource type (%s)", args.ResourceType)
			return nil

		case v := <-ch:
			m, err := structpb.NewValue(v.Result)
			if err != nil {
				fmt.Printf("Error converting the value to structpb.Value: %s\n", err)
				return err
			}

			if err := server.Send(&proto.WatchResourceResult{Type: v.Type, Result: m}); err != nil {
				fmt.Printf("Error sending the result to the client: %s\n", err)
				return err
			}
		}
	}

	return nil
}

func (g *GRPCServer) CloseResourceWatcher(ctx context.Context, w *wrappers.StringValue) (*empty.Empty, error) {
	return &emptypb.Empty{}, g.Impl.CloseResourceWatcher(w.Value)
}

func (g *GRPCServer) GetResourceTypeSchema(ctx context.Context, t *wrappers.StringValue) (*proto.ResourceTransformer, error) {
	return g.Impl.GetResourceTypeSchema(t.Value)
}

func (g *GRPCServer) GetResourceTypeList(ctx context.Context, e *empty.Empty) (*proto.GetResourceTypeListResponse, error) {
	list, err := g.Impl.GetResourceTypeList()
	if err != nil {
		return &proto.GetResourceTypeListResponse{}, err
	}

	return &proto.GetResourceTypeListResponse{ResourceType: list}, nil
}

func (g *GRPCServer) GetAuthInfo(ctx context.Context, e *empty.Empty) (*proto.AuthInfoResponse, error) {
	return g.Impl.GetAuthInfo()
}

func (g *GRPCServer) Connect(ctx context.Context, a *proto.AuthInfo) (*empty.Empty, error) {
	return &emptypb.Empty{}, g.Impl.Connect(a)
}

func (g *GRPCServer) GetResourceIsolatorType(ctx context.Context, e *empty.Empty) (*wrappers.StringValue, error) {
	str, err := g.Impl.GetResourceIsolatorType()
	return wrapperspb.String(str), err
}

func (g *GRPCServer) GetDefaultResourceIsolator(ctx context.Context, e *empty.Empty) (*wrappers.StringValue, error) {
	str, err := g.Impl.GetDefaultResourceIsolator()
	return wrapperspb.String(str), err
}

func (g *GRPCServer) GetSupportedActions(ctx context.Context, e *empty.Empty) (*proto.GetActionListResponse, error) {
	return g.Impl.GetSupportedActions()
}

func (g *GRPCServer) ActionDeleteResource(ctx context.Context, args *proto.ActionDeleteResourceArgs) (*empty.Empty, error) {
	return &emptypb.Empty{}, g.Impl.ActionDeleteResource(args)
}

func (g *GRPCServer) ActionCreateResource(ctx context.Context, args *proto.ActionCreateResourceArgs) (*empty.Empty, error) {

	return &emptypb.Empty{}, g.Impl.ActionCreateResource(args)
}

func (g *GRPCServer) ActionUpdateResource(ctx context.Context, args *proto.ActionUpdateResourceArgs) (*empty.Empty, error) {
	return &emptypb.Empty{}, g.Impl.ActionUpdateResource(args)
}

func (g *GRPCServer) GetSpecficActionList(ctx context.Context, t *wrappers.StringValue) (*proto.GetActionListResponse, error) {
	return g.Impl.GetSpecficActionList(t.Value)
}

func (g *GRPCServer) PerformSpecificAction(ctx context.Context, args *proto.SpecificActionArgs) (*proto.SpecificActionResult, error) {
	return g.Impl.PerformSpecificAction(args)
}
