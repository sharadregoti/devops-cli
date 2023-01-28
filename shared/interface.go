package shared

import "github.com/sharadregoti/devops/model"

type DevopsHelper interface {
	SendData(data interface{}) error
}

type Devops interface {
	Name() string
	StatusOK() error
	MainBox
	SearchBox
	GeneralInfoBox
	ResourceIsolatorBox
	GenericResourceActions
	ResourceSpecificActions
}

type GetResourcesArgs struct {
	ResourceName, ResourceType, IsolatorID string
	Args                                   map[string]interface{}
}

type MainBox interface {
	GetResources(args GetResourcesArgs) ([]interface{}, error)
	WatchResources(resourceType string) (chan WatchResourceResult, error)
	CloseResourceWatcher(resourceType string) error
	GetResourceTypeSchema(resourceType string) (model.ResourceTransfomer, error)
}

type SearchBox interface {
	GetResourceTypeList() ([]string, error)
}

type GeneralInfoBox interface {
	GetGeneralInfo() (map[string]string, error)
}

type ResourceIsolatorBox interface {
	GetResourceIsolatorType() (string, error)
	GetDefaultResourceIsolator() (string, error)
}

type DebuggingBox interface {
	GetResourceTypeConditions() error
}

type ChatGPTBox interface {
	GetResourceTypeConditions() error
}

type GenericActions struct {
	// Read is supported by default
	IsDelete bool
	IsUpdate bool
	IsCreate bool
}

type ActionDeleteResourceArgs struct {
	ResourceName, ResourceType, IsolatorName string
}

type GenericResourceActions interface {
	GetSupportedActions(resourceType string) (GenericActions, error)
	ActionDeleteResource(args ActionDeleteResourceArgs) error
	// ActionApplyResource(data interface{}, resourceType string) error
}

type SpecificAction struct {
	Name       string
	KeyBinding string

	// View, Edit, Confirmation
	ScrrenAction string

	// String, Object
	OutputType string
	ResourceID string
}

type SpecificActionArgs struct {
	ActionName string

	ResourceName string
	ResourceType string

	IsolatorName string

	Args map[string]interface{}
}

type SpecificActionResult struct {
	// Temp       string
	Result     interface{}
	OutputType string
}

type ResourceSpecificActions interface {
	GetSpecficActionList(resourceType string) ([]SpecificAction, error)
	PerformSpecificAction(args SpecificActionArgs) (SpecificActionResult, error)
}
