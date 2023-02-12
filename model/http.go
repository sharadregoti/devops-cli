package model

import (
	"github.com/gorilla/websocket"
	"github.com/sharadregoti/devops-plugin-sdk/proto"
)

type WebsocketReadWriter struct {
	Socket *websocket.Conn
}

func NewSocketReadWrite(conn *websocket.Conn) *WebsocketReadWriter {
	return &WebsocketReadWriter{
		Socket: conn,
	}
}

func (srw WebsocketReadWriter) Write(p []byte) (n int, err error) {
	err = srw.Socket.WriteMessage(websocket.TextMessage, p)
	return len(p), err
}

func (srw WebsocketReadWriter) Read(p []byte) (n int, err error) {
	_, b, err := srw.Socket.ReadMessage()
	for i, d := range b {
		p[i] = d
	}
	return len(b), err
}

type InfoResponse struct {
	SessionID       string            `json:"id" yaml:"id"`
	General         map[string]string `json:"general" yaml:"general"`
	Actions         []*proto.Action   `json:"actions" yaml:"actions"`
	ResourceTypes   []string          `json:"resourceTypes" yaml:"resourceTypes"`
	DefaultIsolator []string          `json:"defaultIsolator" yaml:"defaultIsolator"`
	IsolatorType    string            `json:"isolatorType" yaml:"isolatorType"`
}

// type Action struct {
// 	// Type can be one of normal, special, internal
// 	Type EventType `json:"type" yaml:"type"`
// 	// Name is the name to be shown on UI
// 	Name       string `json:"name" yaml:"name"`
// 	KeyBinding string `json:"key_binding" yaml:"key_binding"`
// 	// Output type can be
// 	OutputType            OutputType             `json:"output_type" yaml:"output_type"`
// }

type FrontendEvent struct {
	EventType    EventType              `json:"eventType" yaml:"eventType"`
	ActionName   string                 `json:"name" yaml:"name"`
	ResourceType string                 `json:"resourceType" yaml:"resourceType"`
	ResourceName string                 `json:"resourceName" yaml:"resourceName"`
	IsolatorName string                 `json:"isolatorName" yaml:"isolatorName"`
	PluginName   string                 `json:"pluginName" yaml:"pluginName"`
	Args         map[string]interface{} `json:"args" yaml:"args"`
}

type WebsocketResponse struct {
	ID              string          `json:"id" yaml:"id"`
	TableName       string          `json:"name" yaml:"name"`
	EventType       string          `json:"eventType" yaml:"eventType"`
	Data            []*TableRow     `json:"data" yaml:"data"`
	SpecificActions []*proto.Action `json:"specificActions" yaml:"specificActions"`
}

type ErrorResponse struct {
	Message string `json:"message" yaml:"message"`
}

type EventResponse struct {
	ID     string      `json:"id" yaml:"id"`
	Result interface{} `json:"result" yaml:"result"`
}

type LongRunningInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
	e       *Event
}

func (lri *LongRunningInfo) SetE(e *Event) {
	lri.e = e
}

func (lri *LongRunningInfo) GetE() *Event {
	return lri.e
}

type AuthResponse struct {
	Auths []*proto.AuthInfo `json:"auths" yaml:"auths"`
}
