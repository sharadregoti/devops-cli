package model

type Config struct {
	Plugins []*Plugin `yaml:"plugins"`
}

type Plugin struct {
	Name string `yaml:"name"`
}
