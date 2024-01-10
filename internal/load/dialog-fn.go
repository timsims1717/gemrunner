package load

import (
	"fmt"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems"
)

func Test() {
	fmt.Println("test")
}

func CancelDialog(key string) func() {
	return func() {
		systems.CloseDialog(key)
	}
}

func EditorMode(mode data.EditorMode, btn *data.Button, dialog *data.Dialog) func() {
	return func() {
		data.EditorPanel.Mode = mode
		for _, b := range dialog.Buttons {
			b.Entity.AddComponent(myecs.Drawable, b.Sprite)
		}
		btn.Entity.AddComponent(myecs.Drawable, btn.ClickSpr)
	}
}
