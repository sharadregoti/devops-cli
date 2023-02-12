package shared

import (
	"context"
	"fmt"

	empty "github.com/golang/protobuf/ptypes/empty"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/sharadregoti/devops/proto"
)

type GRPCClient struct {
	client proto.DevopsClient
}

func (g *GRPCClient) Name() string {
	res, err := g.client.Name(context.Background(), &empty.Empty{})
	if err != nil {
		return ""
	}
	return res.Value
}

func (g *GRPCClient) GetResources(args *proto.GetResourcesArgs) ([]interface{}, error) {
	res, err := g.client.GetResources(context.Background(), args)
	if err != nil {
		return nil, err
	}

	resp := make([]interface{}, len(res.Values))
	for i, v := range res.Values {
		resp[i] = v.GetStructValue().AsMap()
	}
	return resp, nil
}

func (g *GRPCClient) WatchResources(resourceType string) (chan WatchResourceResult, chan struct{}, error) {
	resp, err := g.client.WatchResources(context.Background(), &wrappers.StringValue{Value: resourceType})
	if err != nil {
		return nil, nil, err
	}

	ch := make(chan WatchResourceResult, 1)
	done := make(chan struct{}, 1)
	// TODO: We cannot close the go routine. As resp.Recv() is blocking call & break only when the plugin exits
	go func() {
		fmt.Printf("grpc client routine: resource watcher has been started for resource type (%s)\n", resourceType)
		defer fmt.Printf("grpc client routine: resource watcher has been stopped for resource type (%s)\n", resourceType)

		for {
			select {
			case <-done:
				fmt.Printf("grpc client routine: resource watcher Done received for resource type (%s)", resourceType)
				return

			default:
				// This closes when the plugin exists
				res, err := resp.Recv()
				if err != nil {
					// done <- struct{}{}
					fmt.Printf("grpc client routine: Error while watching resource for type (%s) got error: %s\n", resourceType, err)
					return
				}

				ch <- WatchResourceResult{
					Type:   res.Type,
					Result: res.Result.GetStructValue().AsMap(),
				}
			}
		}
	}()

	return ch, done, err
}

func (g *GRPCClient) CloseResourceWatcher(resourceType string) error {
	_, err := g.client.CloseResourceWatcher(context.Background(), &wrappers.StringValue{Value: resourceType})
	return err
}

// GetResourceTypeSchema function retrieves the schema for a given resource type.
func (c *GRPCClient) GetResourceTypeSchema(resourceType string) (*proto.ResourceTransformer, error) {
	res, err := c.client.GetResourceTypeSchema(context.Background(), &wrappers.StringValue{Value: resourceType})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetResourceTypeList function retrieves a list of resource types.
func (c *GRPCClient) GetResourceTypeList() ([]string, error) {
	response, err := c.client.GetResourceTypeList(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return response.ResourceType, nil
}

// GetAuthInfo function retrieves authentication information.
func (c *GRPCClient) GetAuthInfo() (*proto.AuthInfoResponse, error) {
	response, err := c.client.GetAuthInfo(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (g *GRPCClient) Connect(args *proto.AuthInfo) error {
	_, err := g.client.Connect(context.Background(), args)
	if err != nil {
		return err
	}
	return nil
}

// GetResourceIsolatorType function retrieves the resource isolator type.
func (c *GRPCClient) GetResourceIsolatorType() (string, error) {
	response, err := c.client.GetResourceIsolatorType(context.Background(), &empty.Empty{})
	if err != nil {
		return "", err
	}
	return response.Value, nil
}

// GetDefaultResourceIsolator function retrieves the default resource isolator.
func (c *GRPCClient) GetDefaultResourceIsolator() (string, error) {
	response, err := c.client.GetDefaultResourceIsolator(context.Background(), &empty.Empty{})
	if err != nil {
		return "", err
	}
	return response.Value, nil
}

func (c *GRPCClient) GetSupportedActions() (*proto.GetActionListResponse, error) {
	response, err := c.client.GetSupportedActions(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *GRPCClient) ActionDeleteResource(args *proto.ActionDeleteResourceArgs) error {
	_, err := c.client.ActionDeleteResource(context.Background(), args)
	return err
}

func (c *GRPCClient) ActionCreateResource(args *proto.ActionCreateResourceArgs) error {
	_, err := c.client.ActionCreateResource(context.Background(), args)
	return err
}

func (c *GRPCClient) ActionUpdateResource(args *proto.ActionUpdateResourceArgs) error {
	_, err := c.client.ActionUpdateResource(context.Background(), args)
	return err
}

func (c *GRPCClient) GetSpecficActionList(resourceType string) (*proto.GetActionListResponse, error) {
	resourceTypeValue := &wrappers.StringValue{Value: resourceType}
	response, err := c.client.GetSpecficActionList(context.Background(), resourceTypeValue)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *GRPCClient) PerformSpecificAction(specificActionArgs *proto.SpecificActionArgs) (*proto.SpecificActionResult, error) {
	response, err := c.client.PerformSpecificAction(context.Background(), specificActionArgs)
	if err != nil {
		return nil, err
	}
	return response, nil
}
