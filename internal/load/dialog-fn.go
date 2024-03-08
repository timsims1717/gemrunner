package load

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/content"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/states"
	"gemrunner/internal/systems"
	"gemrunner/pkg/state"
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
		if data.CurrPuzzle == nil {
			data.CurrPuzzle = data.CreateBlankPuzzle()
		} else {
			for _, row := range data.CurrPuzzle.Tiles.T {
				for _, tile := range row {
					tile.ToEmpty()
				}
			}
		}
		systems.PuzzleInit()
		systems.UpdateWorldShaders()
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

func TestPuzzle() func() {
	return func() {
		state.PushState(states.TestStateKey)
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
				data.UpdateScrollBounds(scroll)
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
			data.UpdateScrollBounds(scroll)
		}
	}
}

func OnCrackTileOptions() {
	if data.CurrPuzzle != nil {
		if len(data.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tile selected by wrench")
			data.CloseDialog("cracked_tile_options")
			return
		}
		firstTile := data.CurrPuzzle.WrenchTiles[0]
		crackDialog := data.Dialogs["cracked_tile_options"]
		for _, ele := range crackDialog.Elements {
			if x, ok := ele.(*data.Checkbox); ok {
				switch x.Key {
				case "cracked_tile_regenerate_check":
					data.SetChecked(x, firstTile.Metadata.Regenerate)
				case "cracked_tile_show_check":
					data.SetChecked(x, firstTile.Metadata.ShowCrack)
				case "cracked_tile_enemy_check":
					data.SetChecked(x, firstTile.Metadata.EnemyCrack)
				}
			} else if t, okT := ele.(*data.Text); okT {
				if t.Key == "cracked_tile_title" {
					if firstTile.Block == data.BlockCracked {
						t.Text.SetText("Cracked Turf")
					} else {
						t.Text.SetText("Cracked Ladder")
					}
				}
			}
		}
	}
}

func ChangeCrackTileOptions() {
	if data.CurrPuzzle != nil {
		if len(data.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tile selected by wrench")
			data.CloseDialog("cracked_tile_options")
			return
		}
		crackDialog := data.Dialogs["cracked_tile_options"]
		var regen, show, enemy bool
		for _, ele := range crackDialog.Elements {
			if x, ok := ele.(*data.Checkbox); ok {
				switch x.Key {
				case "cracked_tile_regenerate_check":
					regen = x.Checked
				case "cracked_tile_show_check":
					show = x.Checked
				case "cracked_tile_enemy_check":
					enemy = x.Checked
				}
			}
		}
		for _, tile := range data.CurrPuzzle.WrenchTiles {
			tile.Metadata.Regenerate = regen
			tile.Metadata.ShowCrack = show
			tile.Metadata.EnemyCrack = enemy
		}
		data.CloseDialog("cracked_tile_options")
		data.CurrPuzzle.Update = true
		systems.PushUndoArray(true)
	}
}

func OnChangeNameDialog() {
	if data.CurrPuzzle != nil {
		changeName := data.Dialogs["change_name"]
		for _, ele := range changeName.Elements {
			if in, ok := ele.(*data.Input); ok {
				if data.CurrPuzzle.Metadata == nil {
					data.CurrPuzzle.Metadata = &data.PuzzleMetadata{}
				}
				if data.CurrPuzzle.Metadata.Name != "" {
					if in.Value != data.CurrPuzzle.Metadata.Name {
						in.Value = data.CurrPuzzle.Metadata.Name
						in.Text.SetText(in.Value)
					}
				} else {
					in.Value = "Untitled"
					in.Text.SetText("Untitled")
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
				if data.CurrPuzzle.Metadata == nil {
					data.CurrPuzzle.Metadata = &data.PuzzleMetadata{}
				}
				if in.Value != "" {
					data.CurrPuzzle.Metadata.Name = in.Value
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
		data.CurrPuzzle.Update = true
	}
}
