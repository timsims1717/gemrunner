package load

import (
	"fmt"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel/pixelgl"
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

func NewPuzzle() {
	if data.Editor != nil {
		// todo: if changes exist, open save/new dialog
		data.CurrPuzzle = data.CreateBlankPuzzle()
		systems.PuzzleInit()
	}
}

func QuitEditor(win *pixelgl.Window) func() {
	//if data.Editor != nil {
	// todo: if in editor, check if any changes, if so, open save/close dialog
	// todo: otherwise, just quit
	//}
	return func() {
		win.SetClosed(true)
	}
}

func SavePuzzle() {
	if data.Editor != nil {
		if data.CurrPuzzle.Title == "" {
			data.CurrPuzzle.Title = "test"
		}
		//} else {
		data.CurrPuzzle.Filename = fmt.Sprintf("%s.puzzle", data.CurrPuzzle.Title)
		err := systems.SavePuzzle()
		if err != nil {
			fmt.Println("ERROR:", err)
		}
		//}
	}
}

func OnOpenPuzzleDialog() {
	openPzl := data.Dialogs["open_puzzle"]
	for _, ele := range openPzl.Elements {
		if scroll, ok := ele.(*data.Scroll); ok {
			scroll.Elements = []interface{}{}
			
		}
	}
}

func WorldDialog() {

}

func EditorMode(mode data.EditorMode, btn *data.Button, dialog *data.Dialog) func() {
	return func() {
		data.Editor.SelectVis = false
		if data.Editor.Mode != mode {
			data.Editor.LastMode = data.Editor.Mode
			data.Editor.LastCoords = world.Coords{X: -1, Y: -1}
		}
		data.Editor.Mode = mode
		for _, e := range dialog.Elements {
			if b, ok := e.(*data.Button); ok {
				b.Entity.AddComponent(myecs.Drawable, b.Sprite)
			}
		}
		btn.Entity.AddComponent(myecs.Drawable, btn.ClickSpr)
	}
}
