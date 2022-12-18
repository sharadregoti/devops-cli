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
	ReadResource   = "read"
	UpdateResource = "update"
	DeleteResource = "delete"
	CreateResource = "create"

	ShowModal = "show-delete-modal"

	// Resource
	ResourceTypeChanged = "resource-type-change"
	RefreshResource     = "refresh-resource"

	// Isolator
	AddIsolator     = "add-isolator"
	IsolatorChanged = "isolator-change"

	// Specific Action
	SpecificActionOccured = "specific-action-occured"
)
