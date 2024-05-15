package load

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/content"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems"
	"gemrunner/internal/ui"
	"gemrunner/pkg/state"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"strings"
)

// open puzzle dialog

func OnOpenPuzzleDialog() {
	openPzl := ui.Dialogs[constants.DialogOpenPuzzle]
	for _, ele := range openPzl.Elements {
		if ele.ElementType == ui.ScrollElement {
			err := content.LoadPuzzleContent()
			if err != nil {
				fmt.Println("ERROR:", err)
				ele.Elements = []*ui.Element{}
				ui.UpdateScrollBounds(ele)
				return
			}
			total := len(data.PuzzleSetFileList)
			xPos := ele.ViewPort.CamPos.X - ele.ViewPort.Rect.W()*0.5 + 4
			for i := 0; i < total; i++ {
				index := i
				if len(ele.Elements) <= i {
					ec := ui.ElementConstructor{
						Key:         "sub_text",
						ElementType: ui.TextElement,
						SubElements: nil,
					}
					t := ui.CreateTextElement(ec, ele.ViewPort)
					ele.Elements = append(ele.Elements, t)
				}
				txt := ele.Elements[i]
				if txt.ElementType == ui.TextElement {
					txt.Key = fmt.Sprintf("open_puzzle_list_%d", i)
					txt.Text.SetPos(pixel.V(xPos, float64(-i)*world.TileSize))
					txt.Text.SetText(data.PuzzleSetFileList[i].Name)
					if i == 0 {
						txt.Text.SetColor(pixel.ToRGBA(constants.ColorBlue))
						data.SelectedPuzzleIndex = 0
					} else {
						txt.Text.SetColor(pixel.ToRGBA(constants.ColorWhite))
					}
					txt.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, ele.ViewPort, func(hvc *data.HoverClick) {
						if openPzl.Open && openPzl.Active {
							click := hvc.Input.Get("click")
							if hvc.Hover && click.JustPressed() {
								data.SelectedPuzzleIndex = index
								for _, ie := range ele.Elements {
									if ie.ElementType == ui.TextElement {
										ie.Text.SetColor(pixel.ToRGBA(constants.ColorWhite))
									}
								}
								txt.Text.SetColor(pixel.ToRGBA(constants.ColorBlue))
								click.Consume()
							}
						}
					}))
				}
			}
			if len(ele.Elements) > total {
				for i := len(ele.Elements) - 1; i >= total; i-- {
					txt := ele.Elements[i]
					if txt.ElementType == ui.TextElement {
						myecs.Manager.DisposeEntity(txt.Entity)
					}
				}
				if total > 0 {
					ele.Elements = ele.Elements[:total]
				} else {
					ele.Elements = []*ui.Element{}
				}
			}
			ui.UpdateScrollBounds(ele)
		}
	}
}

func OpenOpenPuzzleDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		if data.CurrPuzzleSet.Metadata.Filename == "" && data.CurrPuzzleSet.Changed {
			ui.SetCloseSpcFn(constants.DialogChangeName, func() {
				ui.OpenDialogInStack(constants.DialogOpenPuzzle)
			})
			ui.OpenDialogInStack(constants.DialogChangeName)
		} else if data.CurrPuzzleSet.Changed {
			if systems.SavePuzzleSet() {
				ui.OpenDialogInStack(constants.DialogOpenPuzzle)
			} else {
				ui.SetTempOnClick(constants.DialogUnableToSaveConfirm, "confirm_unable_to_save", func() {
					ui.OpenDialogInStack(constants.DialogOpenPuzzle)
				})
				ui.OpenDialogInStack(constants.DialogUnableToSaveConfirm)
			}
		} else {
			ui.OpenDialogInStack(constants.DialogOpenPuzzle)
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
	ui.CloseDialog(constants.DialogOpenPuzzle)
}

// non dialog puzzle stuff

func NewPuzzle() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		if data.CurrPuzzleSet.Metadata.Filename == "" && data.CurrPuzzleSet.Changed {
			ui.SetCloseSpcFn(constants.DialogChangeName, func() {
				systems.NewPuzzleSet()
			})
			ui.OpenDialogInStack(constants.DialogChangeName)
		} else if data.CurrPuzzleSet.Changed {
			if systems.SavePuzzleSet() {
				systems.NewPuzzleSet()
			} else {
				ui.SetTempOnClick(constants.DialogUnableToSaveConfirm, "confirm_unable_to_save", func() {
					systems.NewPuzzleSet()
				})
				ui.OpenDialogInStack(constants.DialogUnableToSaveConfirm)
			}
		} else {
			systems.NewPuzzleSet()
		}
	}
}

func ExitEditor() func() {
	return func() {
		if data.Editor != nil && data.CurrPuzzleSet != nil {
			if data.CurrPuzzleSet.Metadata.Filename == "" && data.CurrPuzzleSet.Changed {
				ui.SetCloseSpcFn(constants.DialogChangeName, func() {
					state.SwitchState(constants.MainMenuKey)
				})
				ui.OpenDialogInStack(constants.DialogChangeName)
			} else if data.CurrPuzzleSet.Changed {
				if systems.SavePuzzleSet() {
					state.SwitchState(constants.MainMenuKey)
				} else {
					ui.SetTempOnClick(constants.DialogUnableToSaveConfirm, "confirm_unable_to_save", func() {
						state.SwitchState(constants.MainMenuKey)
					})
					ui.OpenDialogInStack(constants.DialogUnableToSaveConfirm)
				}
			} else {
				state.SwitchState(constants.MainMenuKey)
			}
		} else {
			state.SwitchState(constants.MainMenuKey)
		}
	}
}

func TestPuzzle() {
	hasPlayers := data.CurrPuzzleSet.CurrPuzzle.HasPlayers()
	if hasPlayers {
		if data.Editor != nil && data.CurrPuzzleSet != nil {
			if data.CurrPuzzleSet.Metadata.Filename == "" && data.CurrPuzzleSet.Changed {
				ui.SetCloseSpcFn(constants.DialogChangeName, func() {
					state.PushState(constants.TestStateKey)
				})
				ui.OpenDialogInStack(constants.DialogChangeName)
			} else if data.CurrPuzzleSet.Changed {
				if systems.SavePuzzleSet() {
					state.PushState(constants.TestStateKey)
				} else {
					ui.SetTempOnClick(constants.DialogUnableToSaveConfirm, "confirm_unable_to_save", func() {
						state.PushState(constants.TestStateKey)
					})
					ui.OpenDialogInStack(constants.DialogUnableToSaveConfirm)
				}
			} else {
				state.PushState(constants.TestStateKey)
			}
		}
	} else {
		ui.OpenDialogInStack(constants.DialogNoPlayersInPuzzle)
	}
}

func SavePuzzleSet() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		if data.CurrPuzzleSet.Metadata.Filename == "" {
			ui.OpenDialogInStack(constants.DialogChangeName)
		} else if data.CurrPuzzleSet.Changed {
			if !systems.SavePuzzleSet() {
				ui.OpenDialogInStack(constants.DialogUnableToSave)
			}
		}
	}
}

// change name dialog

func OnChangeNameDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		for _, ele := range ui.Dialogs[constants.DialogChangeName].Elements {
			if ele.ElementType == ui.InputElement {
				if data.CurrPuzzleSet.Metadata.Name != "" {
					if ele.Value != data.CurrPuzzleSet.Metadata.Name {
						ele.Value = data.CurrPuzzleSet.Metadata.Name
						ui.ChangeText(ele, ele.Value)
					}
				} else {
					ui.ChangeText(ele, "Untitled")
				}
				break
			}
		}
	}
}

func ChangeName() {
	if data.CurrPuzzleSet != nil {
		changeName := ui.Dialogs[constants.DialogChangeName]
		newName := ""
		for _, ele := range changeName.Elements {
			if ele.ElementType == ui.InputElement {
				if ele.Value != "" {
					newName = ele.Value
				}
				break
			}
		}
		if newName != "" {
			data.CurrPuzzleSet.Metadata.Name = newName
			ui.CloseDialog(constants.DialogChangeName)
			if !systems.SavePuzzleSet() {
				changeName.OnCloseSpc = nil
				ui.OpenDialogInStack(constants.DialogUnableToSave)
			}
		}
	}
}

// change world dialog

func OpenChangeWorldDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		changeWorld := ui.Dialogs[constants.DialogChangeWorld]
		// check if this is a custom world
		data.CustomWorldSelected = data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldNumber == constants.WorldCustom
		data.CustomSelectedBefore = data.CustomWorldSelected
		data.SelectedPrimaryColor = data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor
		data.SelectedSecondaryColor = data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor
		data.SelectedDoodadColor = data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor
		for _, ele := range changeWorld.Elements {
			if strings.Contains(ele.Key, "color_primary") ||
				strings.Contains(ele.Key, "color_secondary") ||
				strings.Contains(ele.Key, "color_doodad") ||
				strings.Contains(ele.Key, "check_primary") ||
				strings.Contains(ele.Key, "check_secondary") ||
				strings.Contains(ele.Key, "check_doodad") ||
				strings.Contains(ele.Key, "primary_text") ||
				strings.Contains(ele.Key, "secondary_text") ||
				strings.Contains(ele.Key, "doodad_text") {
				ele.Object.Hidden = !data.CustomWorldSelected
			}
			switch ele.ElementType {
			case ui.CheckboxElement:
				switch ele.Key {
				case "custom_world_check": // whether Custom World is checked
					ui.SetChecked(ele, data.CustomWorldSelected)
				default:
					updateColorCheckbox(ele)
				}
			case ui.ScrollElement: // world list
				for ctI, ele2 := range ele.Elements {
					if ele2.ElementType == ui.TextElement {
						ele2.Text.Hidden = data.SelectedWorldIndex != ctI/2
					}
				}
			case ui.ContainerElement: // selected world
				if ele.Key == "world_container_selected" {
					for _, ce := range ele.Elements {
						switch ce.Key {
						case "turf_tile":
							ce.Sprite.Key = constants.WorldSprites[data.SelectedWorldIndex]
						case "doodad_tile":
							ce.Sprite.Key = constants.WorldDoodads[data.SelectedWorldIndex]
						case "world_text":
							ce.Text.SetText(constants.WorldNames[data.SelectedWorldIndex])
						}
					}
					pc := pixel.ToRGBA(constants.WorldPrimary[data.SelectedWorldIndex])
					sc := pixel.ToRGBA(constants.WorldSecondary[data.SelectedWorldIndex])
					dc := pixel.ToRGBA(constants.WorldDoodad[data.SelectedWorldIndex])
					ele.ViewPort.Canvas.SetUniform("uRedPrimary", float32(pc.R))
					ele.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(pc.G))
					ele.ViewPort.Canvas.SetUniform("uBluePrimary", float32(pc.B))
					ele.ViewPort.Canvas.SetUniform("uRedSecondary", float32(sc.R))
					ele.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(sc.G))
					ele.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(sc.B))
					ele.ViewPort.Canvas.SetUniform("uRedDoodad", float32(dc.R))
					ele.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(dc.G))
					ele.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(dc.B))
				}
			}
		}
		if data.CustomWorldSelected {
			worldDialogCustomShaders()
		} else {
			worldDialogNormalShaders()
		}
		ui.OpenDialogInStack(constants.DialogChangeWorld)
	}
}

func ConfirmChangeWorld() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		//changeWorld := data.Dialogs[constants.DialogChangeWorld]
		if data.CustomWorldSelected {
			data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldNumber = constants.WorldCustom
		} else {
			data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldNumber = data.SelectedWorldIndex
			data.CurrPuzzleSet.CurrPuzzle.Metadata.MusicTrack = constants.WorldMusic[data.SelectedWorldIndex]
		}
		data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor = data.SelectedPrimaryColor
		data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor = data.SelectedSecondaryColor
		data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor = data.SelectedDoodadColor
		data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldSprite = constants.WorldSprites[data.SelectedWorldIndex]
	}
	systems.UpdateEditorShaders()
	systems.UpdatePuzzleShaders()
	data.CurrPuzzleSet.CurrPuzzle.Update = true
	ui.CloseDialog(constants.DialogChangeWorld)
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
		ui.OpenDialogInStack(constants.DialogAreYouSureDelete)
	}
}

func ConfirmDelete() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		systems.DeletePuzzle()
	}
	ui.CloseDialog(constants.DialogAreYouSureDelete)
}
