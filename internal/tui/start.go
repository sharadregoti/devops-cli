package tui

import (
	"fmt"
	"os"
	"time"

	"github.com/sharadregoti/devops/common"
	"github.com/sharadregoti/devops/utils/logger"
)

func Start(address string) error {
	logger.InitClientLogging()

	app, err := NewApplication(address)
	if err != nil {
		return err
	}

	if err := app.Start(); err != nil {
		common.Error(logger.Loggerf, fmt.Sprintf("failed to start application: %v", err))
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

	return nil
}
