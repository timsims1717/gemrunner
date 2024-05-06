package load

import (
	"fmt"
	"gemrunner/internal/data"
)

func Test(s string) func() {
	return func() {
		fmt.Println(s)
	}
}

func CloseDialog(key string) func() {
	return func() {
		data.CloseDialog(key)
	}
}

func OpenDialog(key string) func() {
	return func() {
		data.OpenDialogInStack(key)
	}
}
