package systems

import (
	"fmt"
	"gemrunner/internal/ui"
)

func Test(s string) func() {
	return func() {
		fmt.Println(s)
	}
}

func CloseDialog(key string) func() {
	return func() {
		ui.CloseDialog(key)
	}
}

func DisposeDialog(key string) func() {
	return func() {
		ui.Dispose(key)
	}
}

func OpenDialog(key string) func() {
	return func() {
		ui.OpenDialogInStack(key)
	}
}
