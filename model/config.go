package model

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
)

// Config stores app config
type Config struct {
	Server  *Server   `json:"server" yaml:"server" binding:"required"`
	Plugins []*Plugin `json:"plugins" yaml:"plugins" binding:"required"`
}

type Server struct {
	Address string `json:"address" yaml:"address" binding:"required"`
}

type Plugin struct {
	Name      string `json:"name" yaml:"name" binding:"required"`
	IsDefault bool   `json:"isDefault" yaml:"isDefault" binding:"required"`
}

// Creates .devops directory if it does not exists
func InitCoreDirectory() string {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Cannot detect home directory", err)
	}

	// Create the ".devops" subdirectory if it doesn't exist
	devopsDir := filepath.Join(homeDir, ".devops")
	if _, err := os.Stat(devopsDir); os.IsNotExist(err) {
		err = os.Mkdir(devopsDir, 0755)
		if err != nil {
			log.Fatal("Cannot create .devops directory", err)
		}
	} else if !os.IsExist(err) && err != nil {
		log.Fatal("Cannot get stats of directory", err)
	}

	return devopsDir
}

func LoadConfig(devopsDir string) *Config {
	fmt.Println("Loading config file...")
	c := new(Config)
	configBytes, err := os.ReadFile(filepath.Join(devopsDir, "config.yaml"))
	if os.IsNotExist(err) {
		// Load default
		fmt.Println("config.yaml not found, loading default configuration")
		c = &Config{
			Plugins: []*Plugin{
				{
					Name: "kubernetes",
				},
				// {
				// 	Name: "aws",
				// },
			},
		}
	} else if !os.IsExist(err) && err != nil {
		log.Fatal("failed to read config.yaml file", err)
	}

	if err := yaml.Unmarshal(configBytes, c); err != nil {
		log.Fatal("failed to yaml unmarshal config file", err)
	}
	return c
}
