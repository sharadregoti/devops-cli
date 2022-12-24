package model

type ResourceTransfomer struct {
	Operations      []Operations      `json:"operations" yaml:"operations"`
	SpecificActions []SpecificActions `json:"specific_actions" yaml:"specific_actions"`
}
type JSONPaths struct {
	Path string `json:"path" yaml:"path"`
}
type Operations struct {
	Name         string      `json:"name" yaml:"name"`
	JSONPaths    []JSONPaths `json:"json_paths" yaml:"json_paths"`
	OutputFormat string      `json:"output_format,omitempty" yaml:"output_format,omitempty"`
}

type SpecificActions struct {
	Name         string `json:"name" yaml:"name"`
	KeyBinding   string `json:"key_binding" yaml:"key_binding"`
	ScrrenAction string `json:"scrren_action" yaml:"scrren_action"`
	OutputType   string `json:"output_type" yaml:"output_type"`
}
