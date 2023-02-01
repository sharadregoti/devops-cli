package shared

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
	WatchResources(resourceType string) (WatcheResource, error)
	CloseResourceWatcher(resourceType string) error
	GetResourceTypeSchema(resourceType string) (ResourceTransfomer, error)
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

type ActionDeleteResourceArgs struct {
	ResourceName, ResourceType, IsolatorName string
}

type ActionCreateResourceArgs struct {
	ResourceName, ResourceType, IsolatorName string
	Data                                     interface{}
}

type ActionUpdateResourceArgs struct {
	ResourceName, ResourceType, IsolatorName string
	Data                                     interface{}
}

type GenericResourceActions interface {
	GetSupportedActions() ([]Action, error)
	ActionDeleteResource(args ActionDeleteResourceArgs) error
	ActionCreateResource(ActionCreateResourceArgs) error
	ActionUpdateResource(ActionUpdateResourceArgs) error
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
	GetSpecficActionList(resourceType string) ([]Action, error)
	PerformSpecificAction(args SpecificActionArgs) (SpecificActionResult, error)
}
