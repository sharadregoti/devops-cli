package views

import (
	"fmt"
)

func (a *Application) flashLogError(msg string, args ...interface{}) error {
	str := fmt.Sprintf(msg, args...)
	a.SetFlashText(str)
	return fmt.Errorf("%v", str)
}
