package model

type Info struct {
	SessionID       string            `json:"id" yaml:"id"`
	General         map[string]string `json:"general" yaml:"general"`
	Plugins         map[string]string `json:"plugins" yaml:"plugins"`
	Actions         []*Action         `json:"actions" yaml:"actions"`
	ResourceTypes   []string          `json:"resourceTypes" yaml:"resourceTypes"`
	DefaultIsolator string            `json:"defaultIsolator" yaml:"defaultIsolator"`
	IsolatorType    string            `json:"isolatorType" yaml:"isolatorType"`
}

type Action struct {
	// Type can be one of normal, special, internal
	Type EventType `json:"type" yaml:"type"`
	// Name is the name to be shown on UI
	Name       string `json:"name" yaml:"name"`
	KeyBinding string `json:"key_binding" yaml:"key_binding"`
	// Output type can be
	OutputType OutputType             `json:"output_type" yaml:"output_type"`
	Schema     map[string]interface{} `json:"schema" yaml:"schema"`
}

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
	ID              string      `json:"id" yaml:"id"`
	TableName       string      `json:"name" yaml:"name"`
	Data            []*TableRow `json:"data" yaml:"data"`
	SpecificActions []*Action   `json:"specificActions" yaml:"specificActions"`
}

type ErrorResponse struct {
	Message string `json:"message" yaml:"message"`
}

type EventResponse struct {
	ID     string      `json:"id" yaml:"id"`
	Result interface{} `json:"result" yaml:"result"`
}
