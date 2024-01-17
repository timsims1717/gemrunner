package load

import (
	"fmt"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/world"
)

func Test() {
	fmt.Println("test")
}

func CancelDialog(key string) func() {
	return func() {
		data.CloseDialog(key)
	}
}

func EditorMode(mode data.EditorMode, btn *data.Button, dialog *data.Dialog) func() {
	return func() {
		data.Editor.SelectVis = false
		if data.Editor.Mode != mode {
			data.Editor.LastMode = data.Editor.Mode
			data.Editor.LastCoords = world.Coords{X: -1, Y: -1}
		}
		data.Editor.Mode = mode
		for _, b := range dialog.Buttons {
			b.Entity.AddComponent(myecs.Drawable, b.Sprite)
		}
		btn.Entity.AddComponent(myecs.Drawable, btn.ClickSpr)
	}
}
