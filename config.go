package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/sharadregoti/devops/model"
)

// Creates .devops directory if it does not exists
func initCoreDirectory() string {
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

func getCoreLogFile(devopsDir string) *os.File {
	// Create the ".devops" subdirectory if it doesn't exist
	filePath := filepath.Join(devopsDir, "devops.log")
	fmt.Println("Logfile can be found at:", filePath)
	var file *os.File
	var err error
	// if _, err := os.Stat(filePath); os.IsNotExist(err) {
	file, err = os.Create(filePath)
	if err != nil {
		log.Fatal("Cannot create devops.log file", err)
	}
	// TODO: Fix this
	return file
	// } else if !os.IsExist(err) && err != nil {
	// 	log.Fatal("Cannot get stats of devops.log fiel", err)
	// }

	// file, err = os.Open(filePath)
	// if err != nil {
	// 	log.Fatal("Cannot open devops.log file", err)
	// }
	// fmt.Println(filePath)
	// return file
	// defer file.Close()
}

func createLoggers(file *os.File) (loggero, loggerf hclog.Logger) {
	loggero = hclog.New(&hclog.LoggerOptions{
		Name:   "devops",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	loggerf = hclog.New(&hclog.LoggerOptions{
		Name:   "devops",
		Output: file,
		Level:  hclog.Debug,
	})

	return
}

func loadConfig(devopsDir string) *model.Config {
	fmt.Println("Loading config file...")
	c := new(model.Config)
	configBytes, err := os.ReadFile(filepath.Join(devopsDir, "config.yaml"))
	if os.IsNotExist(err) {
		// Load default
		fmt.Println("config.yaml not found, loading default configuration")
		c = &model.Config{
			Plugins: []*model.Plugin{
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
