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

// open puzzle dialog

func OnOpenPuzzleDialog() {
	openPzl := data.Dialogs[constants.DialogOpenPuzzle]
	for _, ele := range openPzl.Elements {
		if scroll, ok := ele.(*data.Scroll); ok {
			err := content.LoadPuzzleContent()
			if err != nil {
				fmt.Println("ERROR:", err)
				scroll.Elements = []interface{}{}
				data.UpdateScrollBounds(scroll)
				return
			}
			total := len(data.PuzzleSetFileList)
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
					txt.Text.SetText(data.PuzzleSetFileList[i].Name)
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

func OpenOpenPuzzleDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		if data.CurrPuzzleSet.Metadata.Filename == "" && data.CurrPuzzleSet.Changed {
			data.SetCloseSpcFn(constants.DialogChangeName, func() {
				data.OpenDialogInStack(constants.DialogOpenPuzzle)
			})
			data.OpenDialogInStack(constants.DialogChangeName)
		} else if data.CurrPuzzleSet.Changed {
			if systems.SavePuzzleSet() {
				data.OpenDialogInStack(constants.DialogOpenPuzzle)
			} else {
				data.SetTempOnClick(constants.DialogUnableToSaveConfirm, "confirm_unable_to_save", func() {
					data.OpenDialogInStack(constants.DialogOpenPuzzle)
				})
				data.OpenDialogInStack(constants.DialogUnableToSaveConfirm)
			}
		} else {
			data.OpenDialogInStack(constants.DialogOpenPuzzle)
		}
	}
}

func OpenPuzzle() {
	if data.Editor != nil &&
		data.SelectedPuzzleIndex > -1 &&
		data.SelectedPuzzleIndex < len(data.PuzzleSetFileList) {
		filename := fmt.Sprintf("%s/%s", constants.PuzzlesDir, data.PuzzleSetFileList[data.SelectedPuzzleIndex].Filename)
		err := systems.OpenPuzzle(filename)
		if err != nil {
			err := systems.OpenPuzzleSet(filename)
			if err != nil {
				fmt.Println("ERROR:", err)
			}
		}
	}
	data.CloseDialog(constants.DialogOpenPuzzle)
}

// non dialog puzzle stuff

func NewPuzzle() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		if data.CurrPuzzleSet.Metadata.Filename == "" && data.CurrPuzzleSet.Changed {
			data.SetCloseSpcFn(constants.DialogChangeName, func() {
				systems.NewPuzzleSet()
			})
			data.OpenDialogInStack(constants.DialogChangeName)
		} else if data.CurrPuzzleSet.Changed {
			if systems.SavePuzzleSet() {
				systems.NewPuzzleSet()
			} else {
				data.SetTempOnClick(constants.DialogUnableToSaveConfirm, "confirm_unable_to_save", func() {
					systems.NewPuzzleSet()
				})
				data.OpenDialogInStack(constants.DialogUnableToSaveConfirm)
			}
		} else {
			systems.NewPuzzleSet()
		}
	}
}

func QuitEditor(win *pixelgl.Window) func() {
	return func() {
		if data.Editor != nil && data.CurrPuzzleSet != nil {
			if data.CurrPuzzleSet.Metadata.Filename == "" && data.CurrPuzzleSet.Changed {
				data.SetCloseSpcFn(constants.DialogChangeName, func() {
					win.SetClosed(true)
				})
				data.OpenDialogInStack(constants.DialogChangeName)
			} else if data.CurrPuzzleSet.Changed {
				if systems.SavePuzzleSet() {
					win.SetClosed(true)
				} else {
					data.SetTempOnClick(constants.DialogUnableToSaveConfirm, "confirm_unable_to_save", func() {
						win.SetClosed(true)
					})
					data.OpenDialogInStack(constants.DialogUnableToSaveConfirm)
				}
			} else {
				win.SetClosed(true)
			}
		}
	}
}

func TestPuzzle() {
	hasPlayers := data.CurrPuzzleSet.CurrPuzzle.HasPlayers()
	if hasPlayers {
		if data.Editor != nil && data.CurrPuzzleSet != nil {
			if data.CurrPuzzleSet.Metadata.Filename == "" && data.CurrPuzzleSet.Changed {
				data.SetCloseSpcFn(constants.DialogChangeName, func() {
					state.PushState(states.TestStateKey)
				})
				data.OpenDialogInStack(constants.DialogChangeName)
			} else if data.CurrPuzzleSet.Changed {
				if systems.SavePuzzleSet() {
					state.PushState(states.TestStateKey)
				} else {
					data.SetTempOnClick(constants.DialogUnableToSaveConfirm, "confirm_unable_to_save", func() {
						state.PushState(states.TestStateKey)
					})
					data.OpenDialogInStack(constants.DialogUnableToSaveConfirm)
				}
			} else {
				state.PushState(states.TestStateKey)
			}
		}
	} else {
		data.OpenDialogInStack(constants.DialogNoPlayersInPuzzle)
	}
}

func SavePuzzleSet() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		if data.CurrPuzzleSet.Metadata.Filename == "" {
			data.OpenDialogInStack(constants.DialogChangeName)
		} else if data.CurrPuzzleSet.Changed {
			if !systems.SavePuzzleSet() {
				data.OpenDialogInStack(constants.DialogUnableToSave)
			}
		}
	}
}

// change name dialog

func OnChangeNameDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		changeName := data.Dialogs[constants.DialogChangeName]
		for _, ele := range changeName.Elements {
			if in, ok := ele.(*data.Input); ok {
				if data.CurrPuzzleSet.Metadata.Name != "" {
					if in.Value != data.CurrPuzzleSet.Metadata.Name {
						in.Value = data.CurrPuzzleSet.Metadata.Name
						in.Text.SetText(in.Value)
					}
				} else {
					data.ChangeText(in, "Untitled")
				}
				break
			}
		}
	}
}

func ChangeName() {
	if data.CurrPuzzleSet != nil {
		changeName := data.Dialogs[constants.DialogChangeName]
		newName := ""
		for _, ele := range changeName.Elements {
			if in, ok := ele.(*data.Input); ok {
				if in.Value != "" {
					newName = in.Value
				}
				break
			}
		}
		if newName != "" {
			data.CurrPuzzleSet.Metadata.Name = newName
			data.CloseDialog(constants.DialogChangeName)
			if !systems.SavePuzzleSet() {
				changeName.OnCloseSpc = nil
				data.OpenDialogInStack(constants.DialogUnableToSave)
			}
		}
	}
}

// change world dialog

func OpenChangeWorldDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		changeWorld := data.Dialogs[constants.DialogChangeWorld]
		// check if this is a custom world
		data.CustomWorldSelected = data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldNumber == constants.WorldCustom
		data.CustomSelectedBefore = data.CustomWorldSelected
		data.SelectedPrimaryColor = data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor
		data.SelectedSecondaryColor = data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor
		data.SelectedDoodadColor = data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor
		for _, ele := range changeWorld.Elements {
			if txt, okT := ele.(*data.Text); okT {
				if o, okO := txt.Entity.GetComponentData(myecs.Object); okO {
					if obj, okO1 := o.(*object.Object); okO1 {
						switch txt.Key {
						case "current_world": // the world text
							obj.Hidden = data.CustomWorldSelected
							txt.Text.SetText(fmt.Sprintf("World - %s", constants.WorldNames[data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldNumber]))
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
		data.OpenDialogInStack(constants.DialogChangeWorld)
	}
}

func ConfirmChangeWorld() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		//changeWorld := data.Dialogs[constants.DialogChangeWorld]
		if data.CustomWorldSelected {
			data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldNumber = constants.WorldCustom
		} else {
			data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldNumber = data.SelectedWorldIndex
		}
		data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor = data.SelectedPrimaryColor
		data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor = data.SelectedSecondaryColor
		data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor = data.SelectedDoodadColor
		data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldSprite = constants.WorldSprites[data.SelectedWorldIndex]
	}
	systems.UpdateEditorShaders()
	systems.UpdatePuzzleShaders()
	data.CurrPuzzleSet.CurrPuzzle.Update = true
	data.CloseDialog(constants.DialogChangeWorld)
}

// individual puzzle buttons

func AddPuzzle() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		systems.AddPuzzle()
	}
}

func PrevPuzzle() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		systems.PrevPuzzle()
	}
}

func NextPuzzle() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		systems.NextPuzzle()
	}
}

func DeletePuzzle() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		data.OpenDialogInStack(constants.DialogAreYouSureDelete)
	}
}

func ConfirmDelete() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		systems.DeletePuzzle()
	}
	data.CloseDialog(constants.DialogAreYouSureDelete)
}
