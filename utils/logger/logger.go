package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	shared "github.com/sharadregoti/devops-plugin-sdk"
	"github.com/sharadregoti/devops/model"
)

var Loggerf hclog.Logger
var Loggero hclog.Logger
var fileWriter io.Writer

func init() {
	// Init plugins
	devopsDir := model.InitCoreDirectory()
	filePath := filepath.Join(devopsDir, "devops.log")
	file, err := getCoreLogFile(filePath)
	if err != nil {
		log.Fatal("Error while creating log file", err)
		return
	}

	fileWriter = file
	Loggero, Loggerf = createLoggers(file)
}

func InitClientLogging() {
	devopsDir := model.InitCoreDirectory()
	filePath := filepath.Join(devopsDir, "devops-tui.log")
	file, err := getCoreLogFile(filePath)
	if err != nil {
		log.Fatal("Error while creating log file", err)
		return
	}

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

func getCoreLogFile(filePath string) (*os.File, error) {
	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// If file does not exist, create it
		file, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		return file, nil
	}

	// If file exists, open it
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
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
