package main

import (
	hclog "github.com/hashicorp/go-hclog"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/shared"
)

type CurrentPluginContext struct {
	generalInfo         map[string]string
	defaultIsolator     string
	defaultIsolatorType string
	currentIsolator     string

	supportedResourceTypes []string
	currentResourceType    string

	currentResources []interface{}
	currentSchema    model.ResourceTransfomer

	currentSpecficActionList []shared.SpecificAction
}

func InitPluginContext(logger hclog.Logger, p shared.Devops, resourceType string) (*CurrentPluginContext, error) {
	schema, err := p.GetResourceTypeSchema(resourceType)
	if err != nil {
		logger.Error("initial resource fetching failed", err)
		return nil, err
	}
	info, err := p.GetGeneralInfo()
	if err != nil {
		logger.Error("initial resource fetching failed", err)
		return nil, err
	}
	isolator, err := p.GetDefaultResourceIsolator()
	if err != nil {
		logger.Error("initial resource fetching failed", err)
		return nil, err
	}
	resources, err := p.GetResources(shared.GetResourcesArgs{ResourceType: resourceType, IsolatorID: isolator})
	if err != nil {
		logger.Error("initial resource fetching failed", err)
		return nil, err
	}
	resourceTypeList, err := p.GetResourceTypeList()
	if err != nil {
		logger.Error("initial resource fetching failed", err)
		return nil, err
	}

	defaultIsolatorType, err := p.GetResourceIsolatorType()
	if err != nil {
		logger.Error("initial resource fetching failed", err)
		return nil, err
	}

	return &CurrentPluginContext{
		supportedResourceTypes: resourceTypeList,
		generalInfo:            info,
		defaultIsolator:        isolator,
		currentSchema:          schema,
		currentResources:       resources,
		currentResourceType:    resourceType,
		defaultIsolatorType:    defaultIsolatorType,
		currentIsolator:        isolator,
	}, nil
}

func (c *CurrentPluginContext) setCurrentIsolator(isolatorName string) {
	c.currentIsolator = isolatorName
}
