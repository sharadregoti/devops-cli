package sdk

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
)

var Logger hclog.Logger

func LogInfo(msg string, args ...interface{}) {
	if len(args) > 0 {
		Logger.Info(fmt.Sprintf(msg, args...))
	} else {
		Logger.Info(msg)
	}
}

func LogDebug(msg string, args ...interface{}) {
	if len(args) > 0 {
		Logger.Debug(fmt.Sprintf(msg, args...))
	} else {
		Logger.Debug(msg)
	}
}

func LogTrace(msg string, args ...interface{}) {
	if len(args) > 0 {
		Logger.Trace(fmt.Sprintf(msg, args...))
	} else {
		Logger.Trace(msg)
	}
}

func LogError(msg string, args ...interface{}) error {
	if len(args) > 0 {
		Logger.Error(fmt.Sprintf(msg, args...))
		return fmt.Errorf(fmt.Sprintf(msg, args...))
	} else {
		Logger.Error(msg)
		return fmt.Errorf(msg)
	}
}

func GetHCLLogLevel() hclog.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "trace":
		return hclog.Trace
	case "debug":
		return hclog.Debug
	case "info":
		return hclog.Info
	case "warn":
		return hclog.Warn
	case "error":
		return hclog.Error
	}

	return hclog.Info
}
