package views

import "fmt"

func (a *Application) flashLogError(msg string, args ...string) error {
	str := fmt.Sprintf(msg, args)
	a.SetFlashText(str)
	return fmt.Errorf("%v", str)
}
