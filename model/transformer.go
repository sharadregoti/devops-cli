package model

type ResourceTransfomer struct {
	Operations []Operations `json:"operations" yaml:"operations"`
}
type JSONPaths struct {
	Path string `json:"path" yaml:"path"`
}
type Operations struct {
	Name         string      `json:"name" yaml:"name"`
	JSONPaths    []JSONPaths `json:"json_paths" yaml:"json_paths"`
	OutputFormat string      `json:"output_format,omitempty" yaml:"output_format,omitempty"`
}
