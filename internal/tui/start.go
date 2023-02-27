package tui

import (
	"fmt"

	"github.com/sharadregoti/devops/utils/logger"
)

func Start(address string) error {
	logger.InitClientLogging()

	logger.Loggero.Info("You can find the logs at: ~/.devops/devops-tui.log")

	logger.LogDebug("Starting application...")
	app, err := NewApplication(address)
	if err != nil {
		return err
	}

	logger.LogDebug("Getting app config from server...")
	appConfig, err := app.getAppConfig()
	if err != nil {
		return err
	}
	app.appConfig = appConfig

	// Throw error if no plugins are configured
	if len(appConfig.Plugins) == 0 {
		return fmt.Errorf("no plugins configured")
	}

	logger.LogDebug("Loading data...")
	if err := app.loadPlugin(appConfig.Plugins[0].Name); err != nil {
		return fmt.Errorf("failed to load data: %v", err)
	}

	if err := app.Start(); err != nil {
		return fmt.Errorf("failed to start application: %v", err)
	}

	return nil
}
