package model

type Event struct {
	Type     string
	RowIndex int

	// Resource
	ResourceName string
	ResourceType string

	// Isolator
	IsolatorName string

	// Specific Action
	SpecificActionName string
}

const (
	// Generic Actions
	// ReadResource event show entire json/yaml of a resource in full screen view
	// Required fields: RowIndex
	ReadResource = "read"
	// DeleteResource event shows a modal promt for deleting a resource
	// Required fields: ResourceName, ResourceType, IsolatorName
	DeleteResource = "delete"
	UpdateResource = "update"
	CreateResource = "create"

	// ShowModal event shows a modal promt
	// Required fields: ResourceName, ResourceType, IsolatorName
	ShowModal = "show-delete-modal"

	// Resource
	// Required fields
	ResourceTypeChanged = "resource-type-change"
	RefreshResource     = "refresh-resource"

	// Stream
	Close = "close"

	// Isolator
	AddIsolator     = "add-isolator"
	IsolatorChanged = "isolator-change"

	// Specific Action
	SpecificActionOccured = "specific-action-occured"
)
