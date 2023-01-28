package model

type Event struct {
	Type string
	// TODO: Remove this as not required
	RowIndex int

	// Resource
	ResourceName string
	ResourceType string

	// Isolator
	IsolatorName string

	// Specific Action
	// TODO: Remove this as not required
	SpecificActionName string

	// Plugin
	PluginName string

	Args map[string]interface{}
}

type EventType string
type NormalEvent string
type InternalEvent string
type SpecficEvent string
type OutputType string

const (
	OutputTypeString       OutputType = "string"
	OutputTypeNothing      OutputType = "nothing"
	OutputTypeStream       OutputType = "stream"
	OutputTypeBidrectional OutputType = "bidirectional"
)

const (
	NormalAction   EventType = "normal-action"
	InternalAction EventType = "internal-action"
	SpecificAction EventType = "specfic-action"
)

const (
	// Generic Actions
	// ReadResource event show entire json/yaml of a resource in full screen view
	// Required fields: RowIndex
	ReadResource NormalEvent = "read"
	// DeleteResource event shows a modal promt for deleting a resource
	// Required fields: ResourceName, ResourceType, IsolatorName
	DeleteResource             = "delete"
	UpdateResource             = "update"
	CreateResource             = "create"
	EditResource   NormalEvent = "edit"

	// ShowModal event shows a modal promt
	// Required fields: ResourceName, ResourceType, IsolatorName
	ShowModal = "show-delete-modal"

	// Resource
	// Required fields
	ResourceTypeChanged InternalEvent = "resource-type-change"
	RefreshResource     InternalEvent = "refresh-resource"
	CloseEventLoop      InternalEvent = "closer-event-loop"

	// Stream
	Close InternalEvent = "close"

	// Isolator
	// AddIsolator     SpecficEvent = "add-isolator"
	IsolatorChanged NormalEvent = "isolator-change"

	// Specific Action
	SpecificActionOccured SpecficEvent = "specific-action-occured"

	ViewNestedResource SpecficEvent = "view-nested-resource"

	// Plugin
	PluginChanged = "plugin-change"

	NestBack = "nest-back"
)
