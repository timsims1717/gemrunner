package load

import (
	"fmt"
	"gemrunner/internal/systems"
)

func CancelDialog(key string) func() {
	return func() {
		systems.CloseDialogbox(key)
		fmt.Println("test")
	}
}
