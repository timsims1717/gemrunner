package load

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/content"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
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

func OpenPuzzle() {
	if data.Editor != nil &&
		data.SelectedPuzzleIndex > -1 &&
		data.SelectedPuzzleIndex < len(data.PuzzleInfos) {
		filename := fmt.Sprintf("%s/%s", constants.PuzzlesDir, data.PuzzleInfos[data.SelectedPuzzleIndex].Filename)
		err := systems.OpenPuzzle(filename)
		if err != nil {
			fmt.Println("ERROR:", err)
		}
	}
	data.CloseDialog("open_puzzle")
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
		if data.CurrPuzzle.PuzzleInfo == nil {
			data.CurrPuzzle.PuzzleInfo = &data.PuzzleInfo{}
			data.CurrPuzzle.PuzzleInfo.Name = "test"
			data.CurrPuzzle.PuzzleInfo.Filename = "test.puzzle"
		}
		if data.CurrPuzzle.PuzzleInfo.Name == "" {
			data.CurrPuzzle.PuzzleInfo.Name = "test"
			data.CurrPuzzle.PuzzleInfo.Filename = "test.puzzle"
		}
		//} else {
		data.CurrPuzzle.PuzzleInfo.Filename = fmt.Sprintf("%s.puzzle", data.CurrPuzzle.PuzzleInfo.Name)
		err := systems.SavePuzzle()
		if err != nil {
			fmt.Println("ERROR:", err)
		}
		//}
	}
}

func OnOpenPuzzleDialog() {
	data.SelectedPuzzleIndex = -1
	openPzl := data.Dialogs["open_puzzle"]
	for _, ele := range openPzl.Elements {
		if scroll, ok := ele.(*data.Scroll); ok {
			err := content.LoadPuzzleContent()
			if err != nil {
				fmt.Println("ERROR:", err)
				scroll.Elements = []interface{}{}
				data.UpdateScroll(scroll)
				return
			}
			total := len(data.PuzzleInfos)
			xPos := scroll.ViewPort.CamPos.X - scroll.ViewPort.Rect.W()*0.5 + 4
			//fmt.Println("Puzzle count", total)
			for i := 0; i < total; i++ {
				index := i
				if len(scroll.Elements) <= i {
					ec := data.ElementConstructor{
						Key:         "sub_text",
						Element:     data.TextElement,
						SubElements: nil,
					}
					t := data.CreateTextElement(ec, scroll.ViewPort)
					scroll.Elements = append(scroll.Elements, t)
				}
				t := scroll.Elements[i]
				if txt, okT := t.(*data.Text); okT {
					txt.Key = fmt.Sprintf("open_puzzle_list_%d", i)
					txt.Text.SetPos(pixel.V(xPos, float64(-i)*world.TileSize))
					txt.Text.SetText(data.PuzzleInfos[i].Name)
					if i == 0 {
						txt.Text.SetColor(pixel.ToRGBA(constants.ColorBlue))
						data.SelectedPuzzleIndex = 0
					} else {
						txt.Text.SetColor(pixel.ToRGBA(constants.ColorWhite))
					}
					txt.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, scroll.ViewPort, func(hvc *data.HoverClick) {
						if openPzl.Open && openPzl.Active {
							click := hvc.Input.Get("click")
							if hvc.Hover && click.JustPressed() {
								data.SelectedPuzzleIndex = index
								for _, ie := range scroll.Elements {
									if it, okIT := ie.(*data.Text); okIT {
										it.Text.SetColor(pixel.ToRGBA(constants.ColorWhite))
									}
								}
								txt.Text.SetColor(pixel.ToRGBA(constants.ColorBlue))
								click.Consume()
							}
						}
					}))
				}
			}
			if len(scroll.Elements) > total {
				for i := len(scroll.Elements) - 1; i >= total; i-- {
					t := scroll.Elements[i]
					if txt, okT := t.(*data.Text); okT {
						myecs.Manager.DisposeEntity(txt.Entity)
					}
				}
				if total > 0 {
					scroll.Elements = scroll.Elements[:total]
				} else {
					scroll.Elements = []interface{}{}
				}
			}
			data.UpdateScroll(scroll)
		}
	}
}

func OnChangeNameDialog() {
	if data.CurrPuzzle != nil {
		changeName := data.Dialogs["change_name"]
		for _, ele := range changeName.Elements {
			if in, ok := ele.(*data.Input); ok {
				if data.CurrPuzzle.PuzzleInfo == nil {
					data.CurrPuzzle.PuzzleInfo = &data.PuzzleInfo{}
				}
				if data.CurrPuzzle.PuzzleInfo.Name != "" {
					in.Value = data.CurrPuzzle.PuzzleInfo.Name
				} else {
					in.Value = "Untitled"
				}
				break
			}
		}
	}
}

func ChangeName() {
	if data.CurrPuzzle != nil {
		changeName := data.Dialogs["change_name"]
		for _, ele := range changeName.Elements {
			if in, ok := ele.(*data.Input); ok {
				if data.CurrPuzzle.PuzzleInfo == nil {
					data.CurrPuzzle.PuzzleInfo = &data.PuzzleInfo{}
				}
				if in.Value != "" {
					data.CurrPuzzle.PuzzleInfo.Name = in.Value
					data.CloseDialog("change_name")
				}
				break
			}
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
