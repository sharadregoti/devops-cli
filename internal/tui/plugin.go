package tui

import (
	"strings"

	"github.com/sharadregoti/devops-plugin-sdk/proto"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/utils"
	"github.com/sharadregoti/devops/utils/logger"
)

func (a *Application) connectAndLoadData(pluginName string, pluginAuth *proto.AuthInfo) error {
	// TODO: Show selection if default not found
	logger.LogDebug("Plugin auth (%s) (%s) (%s)", pluginName, pluginAuth.IdentifyingName, pluginAuth.Name)
	infoRes, err := a.getInfo(pluginName, pluginAuth.IdentifyingName, pluginAuth.Name)
	if err != nil {
		return err
	}

	wsdata := make(chan model.WebsocketResponse, 1)

	if err := a.tableWebsocket(infoRes.SessionID, wsdata); err != nil {
		return err
	}

	// Initialize all boxes
	a.mainPage.normalActionBox.RefreshActions(infoRes.Actions)
	a.mainPage.generalInfoBox.Refresh(infoRes.General)
	a.mainPage.isolatorBox.SetDefault(infoRes.DefaultIsolator)
	a.mainPage.isolatorBox.SetTitle(strings.Title(infoRes.IsolatorType))
	a.mainPage.searchBox.SetResourceTypes(append(infoRes.ResourceTypes, a.settings...))

	a.currentIsolator = infoRes.DefaultIsolator[0]
	a.currentResourceType = infoRes.IsolatorType
	a.sessionID = infoRes.SessionID
	a.wsdata = wsdata

	a.startGoRoutine()

	return nil
}

func (a *Application) startGoRoutine() {
	go func() {
		// a.mainPage.tableBox.view.dra
		logger.LogDebug("Started websocket data go routine for ID (%s)", a.sessionID)
		defer logger.LogDebug("Closing websocket data go routine for ID (%s)", a.sessionID)

		currentTableName := ""
		count := 0
		for {

			select {
			case <-a.closeChan:
				return

			case v := <-a.wsdata:

				if currentTableName != strings.ToLower(v.TableName) {
					count = 0
					a.mainPage.tableBox.view.Clear()
					a.application.Draw()
					currentTableName = strings.ToLower(v.TableName)
				}

				// TODO: Optimize this
				headerRow := v.Data[0]
				a.mainPage.tableBox.SetHeader(headerRow)
				dataRow := v.Data[1]

				switch v.EventType {
				case "added":
					a.mainPage.tableBox.AddRow(dataRow)
					count++
				case "modified", "updated":
					a.mainPage.tableBox.UpdateRow(dataRow)
				case "deleted":
					a.mainPage.tableBox.DeleteRow(dataRow)
					count--
				default:
					logger.LogDebug("Unknown event type recieved from server (%s)", v.EventType)
				}

				a.mainPage.specificActionBox.RefreshActions(v.SpecificActions)
				// // c.appView.ActionView.EnableNesting(rs.currentSchema.Nesting.IsNested)
				a.mainPage.tableBox.SetTitle(utils.GetTableTitle(v.TableName, count))
				a.RemoveSearchView()
				// a.mainPage.tableBox.Refresh(v.Data, 0)
				a.application.SetFocus(a.mainPage.tableBox.view)
				a.application.Draw()
				// a.mainPage.tableBox.view.Draw(tcell.NewSimulationScreen("UTF-8"))
				// a.application.Sync()
				logger.LogDebug("Websocket: received data from server, total length (%v)", len(v.Data))
			}
		}
	}()

}
