package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/shared"
)

var Loggerf hclog.Logger
var Loggero hclog.Logger
var fileWriter io.Writer

func init() {
	// Init plugins
	devopsDir := model.InitCoreDirectory()
	file := getCoreLogFile(devopsDir)

	fileWriter = file
	Loggero, Loggerf = createLoggers(file)
	// Loggerf.Info("Yikes 2")
	// Loggerf.Debug("Yikes")
	// Loggerf.Error("Yikes 3")
}

func InitClientLogging() {
	devopsDir := model.InitCoreDirectory()
	file := getClientCoreLogFile(devopsDir)

	fileWriter = file
	Loggero, Loggerf = createLoggers(file)
}

func createLoggers(file *os.File) (loggero, loggerf hclog.Logger) {
	loggero = hclog.New(&hclog.LoggerOptions{
		Name:   "devops",
		Output: os.Stdout,
		Level:  shared.GetHCLLogLevel(),
	})

	loggerf = hclog.New(&hclog.LoggerOptions{
		Name:   "devops",
		Output: file,
		Level:  shared.GetHCLLogLevel(),
	})

	return
}

func GetFileWriter() io.Writer {
	return fileWriter
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

func getClientCoreLogFile(devopsDir string) *os.File {
	// Create the ".devops" subdirectory if it doesn't exist
	filePath := filepath.Join(devopsDir, "devops-client.log")
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

func LogInfo(msg string, args ...interface{}) {
	if len(args) > 0 {
		Loggerf.Info(fmt.Sprintf(msg, args...))
	} else {
		Loggerf.Info(msg)
	}
}

func LogDebug(msg string, args ...interface{}) {
	if len(args) > 0 {
		Loggerf.Debug(fmt.Sprintf(msg, args...))
	} else {
		Loggerf.Debug(msg)
	}
}

func LogTrace(msg string, args ...interface{}) {
	if len(args) > 0 {
		Loggerf.Trace(fmt.Sprintf(msg, args...))
	} else {
		Loggerf.Trace(msg)
	}
}

func LogError(msg string, args ...interface{}) error {
	if len(args) > 0 {
		Loggerf.Error(fmt.Sprintf(msg, args...))
		return fmt.Errorf(fmt.Sprintf(msg, args...))
	} else {
		Loggerf.Error(msg)
		return fmt.Errorf(msg)
	}
}
