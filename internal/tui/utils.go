package tui

import (
	"fmt"

	"github.com/sharadregoti/devops/utils/logger"
)

func (a *Application) flashLogError(msg string, args ...interface{}) error {
	str := fmt.Sprintf(msg, args...)
	a.SetFlashText(str)
	logger.LogError("%v", str)
	return fmt.Errorf("%v", str)
}
