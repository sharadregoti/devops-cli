package sdk

type WatchResourceResult struct {
	Type   string
	Result interface{}
}

type ResourceTransformer struct {
	Operations      []Operations `json:"operations" yaml:"operations"`
	SpecificActions []Action     `json:"specific_actions" yaml:"specific_actions"`
	Styles          []Styles     `json:"styles" yaml:"styles"`
	Nesting         Nesting      `json:"nesting" yaml:"nesting"`
}

type Operations struct {
	Name         string      `json:"name" yaml:"name"`
	JSONPaths    []JSONPaths `json:"json_paths" yaml:"json_paths"`
	OutputFormat string      `json:"output_format,omitempty" yaml:"output_format,omitempty"`
}

type JSONPaths struct {
	Path string `json:"path" yaml:"path"`
}

type Styles struct {
	RowBackgroundColor string   `json:"row_background_color" yaml:"row_background_color"`
	Conditions         []string `json:"conditions" yaml:"conditions"`
}

type Nesting struct {
	IsNested                bool                   `json:"is_nested" yaml:"is_nested"`
	ResourceType            string                 `json:"resource_type" yaml:"resource_type"`
	Args                    map[string]interface{} `json:"args" yaml:"args"`
	IsSelfContainedInParent bool                   `json:"is_self_contained_in_parent" yaml:"is_self_contained_in_parent"`
	ParentDataPaths         []string               `json:"parent_data_paths" yaml:"parent_data_paths"`
}

type Action struct {
	Name       string `json:"name" yaml:"name"`
	KeyBinding string `json:"key_binding" yaml:"key_binding"`
	// ScrrenAction string                 `json:"scrren_action" yaml:"scrren_action"`
	OutputType string                 `json:"output_type" yaml:"output_type"`
	Args       map[string]interface{} `json:"args" yaml:"args"`
	Schema     map[string]interface{} `json:"schema" yaml:"schema"`
	Execution  Execution              `json:"execution" yaml:"execution"`
}

type Execution struct {
	Cmd           string      `json:"cmd" yaml:"cmd"`
	IsLongRunning bool        `json:"is_long_running" yaml:"is_long_running"`
	UserInput     UserInput   `json:"user_input" yaml:"user_input"`
	ServerInput   ServerInput `json:"server_input" yaml:"server_input"`
}

type UserInput struct {
	Required bool                   `json:"required" yaml:"required"`
	Args     map[string]interface{} `json:"args" yaml:"args"`
}

type ServerInput struct {
	Required bool `json:"required" yaml:"required"`
}
