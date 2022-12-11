package shared

import "github.com/sharadregoti/devops/model"

type Devops interface {
	Name() string
	MainBox
	SearchBox
	GeneralInfoBox
	ResourceIsolatorBox
}

type MainBox interface {
	GetResources(resourceType string) []interface{}
	GetResourceTypeSchema(resourceType string) model.ResourceTransfomer
}

type SearchBox interface {
	GetResourceTypeList() []string
}

type GeneralInfoBox interface {
	GetGeneralInfo() map[string]string
}

type ResourceIsolatorBox interface {
	GetResourceIsolatorType() string
	GetDefaultResourceIsolator() string
}
