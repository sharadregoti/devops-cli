package shared

import "github.com/sharadregoti/devops/proto"

type Devops interface {
	Name() string
	MainBox
	SearchBox
	GeneralInfoBox
	ResourceIsolatorBox
	GenericResourceActions
	ResourceSpecificActions
}

// ############################# MainBox #############################

type GetResourcesArgs struct {
	ResourceName, ResourceType, IsolatorID string
	Args                                   map[string]interface{}
}

type MainBox interface {
	GetResources(args *proto.GetResourcesArgs) ([]interface{}, error)
	WatchResources(resourceType string) (chan WatchResourceResult, chan struct{}, error)
	CloseResourceWatcher(resourceType string) error
	GetResourceTypeSchema(resourceType string) (*proto.ResourceTransformer, error)
}

// ############################# SearchBox #############################

type SearchBox interface {
	GetResourceTypeList() ([]string, error)
}

// type AuthInfo struct {
// 	IdentifyingName  string
// 	Name             string
// 	IsDefault        bool
// 	DefaultIsolators []string
// 	Info             map[string]string
// }

type GeneralInfoBox interface {
	GetAuthInfo() (*proto.AuthInfoResponse, error)
	Connect(authInfo *proto.AuthInfo) error
}

type ResourceIsolatorBox interface {
	GetResourceIsolatorType() (string, error)
	GetDefaultResourceIsolator() (string, error)
}

// type DebuggingBox interface {
// 	GetResourceTypeConditions() error
// }

// type ChatGPTBox interface {
// 	GetResourceTypeConditions() error
// }

// type ActionDeleteResourceArgs struct {
// 	ResourceName, ResourceType, IsolatorName string
// }

// type ActionCreateResourceArgs struct {
// 	ResourceName, ResourceType, IsolatorName string
// 	Data                                     interface{}
// }

// type ActionUpdateResourceArgs struct {
// 	ResourceName, ResourceType, IsolatorName string
// 	Data                                     interface{}
// }

type GenericResourceActions interface {
	GetSupportedActions() (*proto.GetActionListResponse, error)
	ActionDeleteResource(*proto.ActionDeleteResourceArgs) error
	ActionCreateResource(*proto.ActionCreateResourceArgs) error
	ActionUpdateResource(*proto.ActionUpdateResourceArgs) error
}

// type SpecificActionArgs struct {
// 	ActionName string

// 	ResourceName string
// 	ResourceType string

// 	IsolatorName string

// 	Args map[string]interface{}
// }

// type SpecificActionResult struct {
// 	// Temp       string
// 	Result     interface{}
// 	OutputType string
// }

type ResourceSpecificActions interface {
	GetSpecficActionList(resourceType string) (*proto.GetActionListResponse, error)
	PerformSpecificAction(args *proto.SpecificActionArgs) (*proto.SpecificActionResult, error)
}
