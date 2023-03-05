package tui

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
		logger.LogDebug("Started websocket reader routine for ID (%s)", ID)
		defer logger.LogDebug("Closing websocket reader routine for ID (%s)", ID)

		for {
			// TODO: This should be closed as per server close repsonse
			msg := new(model.WebsocketResponse)
			err := c.ReadJSON(msg)
			if err != nil {
				logger.LogDebug("Closing websocket connection for ID %s", ID)
				a.flashLogError("Failed to read from server: %s", err.Error())
				return
			}
			wdData <- *msg
		}
	}()

	return nil
}

func (a *Application) actionWebsocket(ID string) (*model.WebsocketReadWriter, error) {
	u := url.URL{Scheme: "ws", Host: a.addr, Path: fmt.Sprintf("/v1/ws/action/%s/%s", a.sessionID, ID)}
	logger.LogInfo("Connecting to server... %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	// TODO: Handle websocket closer
	// defer c.Close()

	return model.NewSocketReadWrite(c), nil
}

func (a *Application) getAppConfig() (*model.Config, error) {
	u := url.URL{Scheme: "http", Host: a.addr, Path: "/v1/config"}
	res, err := http.Get(u.String())
	if err != nil {
		return nil, a.flashLogError("Failed to send get app config request: %v", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		msg := new(model.ErrorResponse)
		_ = json.NewDecoder(res.Body).Decode(msg)
		return nil, a.flashLogError("Get app config request failed: received %v status code from server with error %v", res.Status, msg.Message)
	}

	msg := new(model.Config)
	_ = json.NewDecoder(res.Body).Decode(msg)

	return msg, nil
}

func (a *Application) getPluginAuths(pluginName string) (*model.AuthResponse, error) {
	u := url.URL{Scheme: "http", Host: a.addr, Path: fmt.Sprintf("/v1/auth/%s", pluginName)}
	res, err := http.Get(u.String())
	if err != nil {
		return nil, a.flashLogError("Failed to send get plugin auths request: %v", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		msg := new(model.ErrorResponse)
		_ = json.NewDecoder(res.Body).Decode(msg)
		return nil, a.flashLogError("Get plugin auth request failed: received %v status code from server with error %v", res.Status, msg.Message)
	}

	msg := new(model.AuthResponse)
	_ = json.NewDecoder(res.Body).Decode(msg)

	return msg, nil
}

func (a *Application) getInfo(pluginName, authId, contextId string) (*model.InfoResponse, error) {
	u := url.URL{Scheme: "http", Host: a.addr, Path: fmt.Sprintf("/v1/connect/%s/%s/%s", pluginName, authId, contextId)}
	res, err := http.Get(u.String())
	if err != nil {
		return nil, a.flashLogError("Failed to send get info request: %v", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		msg := new(model.ErrorResponse)
		_ = json.NewDecoder(res.Body).Decode(msg)
		return nil, a.flashLogError("Get info request failed: received %v status code from server with error %v", res.Status, msg.Message)
	}

	msg := new(model.InfoResponse)
	_ = json.NewDecoder(res.Body).Decode(msg)

	return msg, nil
}

func (a *Application) sendEvent(fe model.FrontendEvent) (*model.EventResponse, error) {
	logger.LogDebug("Sending event request... (%s)", fe.ActionName)
	data, err := json.Marshal(fe)
	if err != nil {
		return nil, a.flashLogError("Failed to send event request, json marshal error : %v", err.Error())
	}
	u := url.URL{Scheme: "http", Host: a.addr, Path: fmt.Sprintf("/v1/events/%s", a.sessionID)}
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
