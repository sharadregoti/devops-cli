package views

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/utils/logger"
)

func (a *Application) tableWebsocket(ID string, wdData chan model.WebsocketResponse) error {
	u := url.URL{Scheme: "ws", Host: a.addr, Path: fmt.Sprintf("/v1/ws/%s", ID)}
	logger.LogInfo("Connecting to server... %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	// TODO: Handle websocket closer
	// defer c.Close()

	go func() {
		for {
			msg := new(model.WebsocketResponse)
			err := c.ReadJSON(msg)
			if err != nil {
				a.flashLogError("Failed to read from server", err.Error())
				return
			}
			wdData <- *msg
		}
	}()

	return nil
}

func (a *Application) actionWebsocket(ID string) (*model.WebsocketReadWriter, error) {
	u := url.URL{Scheme: "ws", Host: a.addr, Path: fmt.Sprintf("/v1/ws/action/%s/%s", a.connectionID, ID)}
	logger.LogInfo("Connecting to server... %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	// TODO: Handle websocket closer
	// defer c.Close()

	return model.NewSocketReadWrite(c), nil
}

func (a *Application) getInfo() (*model.Info, error) {
	u := url.URL{Scheme: "http", Host: a.addr, Path: "/v1/info"}
	res, err := http.Get(u.String())
	if err != nil {
		return nil, a.flashLogError("Failed to send get info request: %v", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		return nil, a.flashLogError("Get info request failed: received %v status code from server", res.Status)
	}

	msg := new(model.Info)
	_ = json.NewDecoder(res.Body).Decode(msg)

	return msg, nil

}

func (a *Application) sendEvent(fe model.FrontendEvent) (*model.EventResponse, error) {
	logger.LogDebug("Sending event request... (%s)", fe.ActionName)
	data, _ := json.Marshal(fe)
	u := url.URL{Scheme: "http", Host: a.addr, Path: fmt.Sprintf("/v1/events/%s", a.connectionID)}
	res, err := http.Post(u.String(), "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, a.flashLogError("Failed to send event request: %v", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		msg := new(model.ErrorResponse)
		_ = json.NewDecoder(res.Body).Decode(msg)
		return nil, a.flashLogError("Failed to send event request: received %v status code from server with error %v", res.Status, msg.Message)
	}

	msg := new(model.EventResponse)
	_ = json.NewDecoder(res.Body).Decode(msg)

	return msg, nil
}
