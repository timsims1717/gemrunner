package load

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/content"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/states"
	"gemrunner/internal/systems"
	"gemrunner/pkg/object"
	"gemrunner/pkg/state"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"strings"
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

func OpenOpenPuzzleDialog() {
	if data.Editor != nil {
		if data.CurrPuzzle.Metadata.Filename == "" &&
			data.CurrPuzzle.Changed {
			data.SetCloseSpcFn("change_name", func() {
				data.OpenDialogInStack("open_puzzle")
			})
			data.OpenDialogInStack("change_name")
		} else {
			systems.SavePuzzle()
			data.OpenDialogInStack("open_puzzle")
		}
	}
}

func NewPuzzle() {
	if data.Editor != nil {
		if data.CurrPuzzle.Metadata.Filename == "" &&
			data.CurrPuzzle.Changed {
			data.SetCloseSpcFn("change_name", func() {
				systems.NewPuzzle()
			})
			data.OpenDialogInStack("change_name")
		} else {
			systems.SavePuzzle()
			systems.NewPuzzle()
		}
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
	return func() {
		if data.CurrPuzzle.Metadata.Filename == "" &&
			data.CurrPuzzle.Changed {
			data.SetCloseSpcFn("change_name", func() {
				win.SetClosed(true)
			})
			data.OpenDialogInStack("change_name")
		} else {
			systems.SavePuzzle()
			win.SetClosed(true)
		}
	}
}

func OpenChangeWorldDialog() {
	if data.Editor != nil && data.CurrPuzzle != nil {
		changeWorld := data.Dialogs["change_world"]
		// check if this is a custom world
		data.CustomWorldSelected = data.CurrPuzzle.Metadata.WorldNumber == constants.WorldCustom
		data.CustomSelectedBefore = data.CustomWorldSelected
		data.SelectedPrimaryColor = data.CurrPuzzle.Metadata.PrimaryColor
		data.SelectedSecondaryColor = data.CurrPuzzle.Metadata.SecondaryColor
		data.SelectedDoodadColor = data.CurrPuzzle.Metadata.DoodadColor
		for _, ele := range changeWorld.Elements {
			if txt, okT := ele.(*data.Text); okT {
				if o, okO := txt.Entity.GetComponentData(myecs.Object); okO {
					if obj, okO1 := o.(*object.Object); okO1 {
						switch txt.Key {
						case "current_world": // the world text
							obj.Hidden = data.CustomWorldSelected
							txt.Text.SetText(fmt.Sprintf("World - %s", constants.WorldNames[data.CurrPuzzle.Metadata.WorldNumber]))
						case "primary_text", "secondary_text", "doodad_text": // the custom color labels
							obj.Hidden = !data.CustomWorldSelected
						}
					}
				}
			} else if x, ok := ele.(*data.Checkbox); ok {
				switch x.Key {
				case "custom_world_check": // whether Custom World is checked
					data.SetChecked(x, data.CustomWorldSelected)
				default:
					updateColorCheckbox(x)
				}
				if o, okO := x.Entity.GetComponentData(myecs.Object); okO {
					if obj, okO1 := o.(*object.Object); okO1 {
						if strings.Contains(x.Key, "check_primary") ||
							strings.Contains(x.Key, "check_secondary") ||
							strings.Contains(x.Key, "check_doodad") {
							obj.Hidden = !data.CustomWorldSelected
						}
					}
				}
			} else if str1, okS1 := ele.(*data.SprElement); okS1 {
				if o, okO := str1.Entity.GetComponentData(myecs.Object); okO {
					if obj, okO1 := o.(*object.Object); okO1 {
						if strings.Contains(str1.Key, "color_primary") ||
							strings.Contains(str1.Key, "color_secondary") ||
							strings.Contains(str1.Key, "color_doodad") {
							obj.Hidden = !data.CustomWorldSelected
						}
					}
				}
			} else if scr, okScr := ele.(*data.Scroll); okScr { // the list of worlds
				for ctI, ele2 := range scr.Elements {
					if ct, okCt := ele2.(*data.Container); okCt {
						for _, ce := range ct.Elements {
							if it, okIT := ce.(*data.Text); okIT {
								if data.SelectedWorldIndex == ctI {
									it.Text.SetColor(pixel.ToRGBA(constants.ColorBlue))
								} else {
									it.Text.SetColor(pixel.ToRGBA(constants.ColorWhite))
								}
							}
						}
					}
				}
			}
		}
		if data.CustomWorldSelected {
			worldDialogCustomShaders()
		} else {
			worldDialogNormalShaders()
		}
		data.OpenDialogInStack("change_world")
	}
}

func ConfirmChangeWorld() {
	if data.Editor != nil && data.CurrPuzzle != nil {
		//changeWorld := data.Dialogs["change_world"]
		if data.CustomWorldSelected {
			data.CurrPuzzle.Metadata.WorldNumber = constants.WorldCustom
		} else {
			data.CurrPuzzle.Metadata.WorldNumber = data.SelectedWorldIndex
		}
		data.CurrPuzzle.Metadata.PrimaryColor = data.SelectedPrimaryColor
		data.CurrPuzzle.Metadata.SecondaryColor = data.SelectedSecondaryColor
		data.CurrPuzzle.Metadata.DoodadColor = data.SelectedDoodadColor
		data.CurrPuzzle.Metadata.WorldSprite = constants.WorldSprites[data.SelectedWorldIndex]
	}
	systems.UpdateEditorShaders()
	systems.UpdatePuzzleShaders()
	data.CurrPuzzle.Update = true
	data.CloseDialog("change_world")
}

func TestPuzzle() {
	hasPlayers := data.CurrPuzzle.HasPlayers()
	if hasPlayers {
		if data.CurrPuzzle.Metadata.Filename == "" &&
			data.CurrPuzzle.Changed {
			data.SetCloseSpcFn("change_name", func() {
				state.PushState(states.TestStateKey)
			})
			data.OpenDialogInStack("change_name")
		} else {
			systems.SavePuzzle()
			state.PushState(states.TestStateKey)
		}
	} else {
		data.OpenDialogInStack("no_players")
	}
}

func SavePuzzle() {
	if data.Editor != nil {
		if data.CurrPuzzle.Metadata == nil {
			data.CurrPuzzle.Metadata = &data.PuzzleMetadata{}
		}
		if data.CurrPuzzle.Metadata.Filename == "" {
			data.OpenDialogInStack("change_name")
		} else {
			systems.SavePuzzle()
		}
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

func ConfirmCrackTileOptions() {
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
			tile.Metadata.Changed = true
		}
		data.CloseDialog("cracked_tile_options")
		data.CurrPuzzle.Update = true
		data.CurrPuzzle.Changed = true
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
		systems.SavePuzzle()
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
