package core

import (
	"github.com/sharadregoti/devops/proto"
	"github.com/sharadregoti/devops/shared"
)

type tableStack []*resourceStack

func (t *tableStack) length() int {
	return len(*t)
}

func (t *tableStack) upsert(index int, r resourceStack) {
	for i, _ := range *t {
		if i == index {
			(*t)[0] = &r
			return
		}
	}

	// Add if does not exists
	*t = append(*t, &r)
}

func (t *tableStack) resetToParentResource() {
	// Only get the first element
	if len(*t) == 0 {
		return
	}
	*t = (*t)[0:1]
}

type resourceStack struct {
	tableRowNumber           int
	nextResourceArgs         []map[string]interface{}
	currentResourceType      string
	currentResources         []interface{}
	currentSchema            *proto.ResourceTransformer
	currentSpecficActionList *proto.GetActionListResponse
}

type nestedResurce struct {
	nextResourceArgs         []map[string]interface{}
	currentIsolator          string
	currentResourceType      string
	currentResources         []interface{}
	currentSchema            shared.ResourceTransformer
	currentSpecficActionList []shared.Action
}
