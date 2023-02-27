package tui

import "github.com/sharadregoti/devops/model"

func startInternalEventLoop(eventChan chan model.Event) {
	for event := range eventChan {
		switch event.Type {
		case string(model.ResourceTypeChanged), string(model.RefreshResource):

		}
	}
}
