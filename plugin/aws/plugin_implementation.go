package plugin

import (
	"github.com/hashicorp/go-hclog"
	"github.com/sharadregoti/devops/model"
)

type AWS struct {
	logger hclog.Logger
}

func (d *AWS) Name() string {
	return "aws"
}

func (d *AWS) GetResources(resourceType string) []interface{} {
	return make([]interface{}, 0)
}

func (d *AWS) GetResourceTypeSchema(resourceType string) model.ResourceTransfomer {
	return model.ResourceTransfomer{}
}

func (d *AWS) GetResourceTypeList() []string {
	return []string{"ec2"}
}

func (d *AWS) GetGeneralInfo() map[string]string {
	return map[string]string{
		"account": "1234",
		"cluster": "1234",
		"type":    "1234",
		"name":    "1234",
	}
}

func (d *AWS) GetResourceIsolatorType() string {
	return "region"
}

func (d *AWS) GetDefaultResourceIsolator() string {
	return "ap-south-1"
}
