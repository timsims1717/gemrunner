package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/content"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/state"
	"gemrunner/pkg/util"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

// open puzzle dialog

func OnOpenPuzzleDialog() {
	openPzl := ui.Dialogs[constants.DialogOpenPuzzle]
	for _, ele := range openPzl.Elements {
		if ele.ElementType == ui.ScrollElement {
			content.LoadLocalPuzzleList()
			//if err != nil {
			//	fmt.Println("ERROR:", err)
			//	ele.Elements = []*ui.Element{}
			//	ui.UpdateScrollBounds(ele)
			//	return
			//}
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
			if SavePuzzleSet() {
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

func OnOpenPuzzle() {
	if data.Editor != nil &&
		data.SelectedPuzzleIndex > -1 &&
		data.SelectedPuzzleIndex < len(data.PuzzleSetFileList) {
		err := OpenPuzzleSet(data.PuzzleSetFileList[data.SelectedPuzzleIndex].Filename)
		if err != nil {
			fmt.Println("ERROR:", err)
		}
	}
	ui.CloseDialog(constants.DialogOpenPuzzle)
}

// non dialog puzzle stuff

func NewPuzzle() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		if data.CurrPuzzleSet.Metadata.Filename == "" && data.CurrPuzzleSet.Changed {
			ui.SetCloseSpcFn(constants.DialogChangeName, func() {
				NewPuzzleSet()
			})
			ui.OpenDialogInStack(constants.DialogChangeName)
		} else if data.CurrPuzzleSet.Changed {
			if SavePuzzleSet() {
				NewPuzzleSet()
			} else {
				ui.SetTempOnClick(constants.DialogUnableToSaveConfirm, "confirm_unable_to_save", func() {
					NewPuzzleSet()
				})
				ui.OpenDialogInStack(constants.DialogUnableToSaveConfirm)
			}
		} else {
			NewPuzzleSet()
		}
	}
}

func ExitEditor() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		if data.CurrPuzzleSet.Metadata.Filename == "" && data.CurrPuzzleSet.Changed {
			ui.SetCloseSpcFn(constants.DialogChangeName, func() {
				state.SwitchState(constants.MainMenuKey)
			})
			ui.OpenDialogInStack(constants.DialogChangeName)
		} else if data.CurrPuzzleSet.Changed {
			if SavePuzzleSet() {
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
				if SavePuzzleSet() {
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

func OnSavePuzzleSet() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		if data.CurrPuzzleSet.Metadata.Filename == "" {
			ui.OpenDialogInStack(constants.DialogChangeName)
		} else if data.CurrPuzzleSet.Changed {
			if !SavePuzzleSet() {
				ui.OpenDialogInStack(constants.DialogUnableToSave)
			}
		}
	}
}

// change name dialog

func OnChangeNameDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		changeName := ui.Dialogs[constants.DialogChangeName]
		inEle := changeName.Get("puzzle_name")
		if data.CurrPuzzleSet.Metadata.Name != "" {
			if inEle.Value != data.CurrPuzzleSet.Metadata.Name {
				ui.ChangeText(inEle, data.CurrPuzzleSet.Metadata.Name)
			}
		} else {
			ui.ChangeText(inEle, "Untitled")
		}
		changeName.Get("change_name_error").Object.Hidden = true
	}
}

func ChangeName() {
	if data.CurrPuzzleSet != nil {
		changeName := ui.Dialogs[constants.DialogChangeName]
		newName := changeName.Get("puzzle_name").Value
		if newName == "" {
			errorTxt := changeName.Get("change_name_error")
			errorTxt.Text.SetText("Name can't be empty.")
			errorTxt.Object.Hidden = false
		} else if util.ProfanityDetector.IsProfane(newName) {
			errorTxt := changeName.Get("change_name_error")
			errorTxt.Text.SetText("Name can't contain\nprofanity.")
			errorTxt.Object.Hidden = false
		} else {
			data.CurrPuzzleSet.Metadata.Name = newName
			ui.CloseDialog(constants.DialogChangeName)
			if !SavePuzzleSet() {
				changeName.OnCloseSpc = nil
				ui.OpenDialogInStack(constants.DialogUnableToSave)
			}
		}
	}
}

// puzzle settings dialog

func OpenPuzzleSettingsDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		for _, ele := range ui.Dialogs[constants.DialogPuzzleSettings].Elements {
			switch ele.Key {
			case "puzzle_name":
				if data.CurrPuzzleSet.CurrPuzzle.Metadata.Name != "" {
					if ele.Value != data.CurrPuzzleSet.CurrPuzzle.Metadata.Name {
						ui.ChangeText(ele, data.CurrPuzzleSet.CurrPuzzle.Metadata.Name)
					}
				} else {
					ui.ChangeText(ele, "Untitled")
				}
			case "puzzle_author":
				if data.CurrPuzzleSet.CurrPuzzle.Metadata.Author != "" {
					if ele.Value != data.CurrPuzzleSet.CurrPuzzle.Metadata.Author {
						ui.ChangeText(ele, data.CurrPuzzleSet.CurrPuzzle.Metadata.Author)
					}
				} else {
					ui.ChangeText(ele, constants.Username)
				}
			case "puzzle_hub_check":
				ui.SetChecked(ele, data.CurrPuzzleSet.CurrPuzzle.Metadata.HubLevel)
				ele.Object.Hidden = !data.CurrPuzzleSet.Metadata.Adventure
			case "puzzle_hub_label":
				ele.Object.Hidden = !data.CurrPuzzleSet.Metadata.Adventure
			case "puzzle_secret_check":
				ui.SetChecked(ele, data.CurrPuzzleSet.CurrPuzzle.Metadata.SecretLevel)
				ele.Object.Hidden = data.CurrPuzzleSet.Metadata.Adventure
			case "puzzle_secret_label":
				ele.Object.Hidden = data.CurrPuzzleSet.Metadata.Adventure
			case "puzzle_darkness_check":
				ui.SetChecked(ele, data.CurrPuzzleSet.CurrPuzzle.Metadata.Darkness)
			}
		}
		ui.OpenDialogInStack(constants.DialogPuzzleSettings)
	}
}

func ConfirmPuzzleSettings() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		for _, ele := range ui.Dialogs[constants.DialogPuzzleSettings].Elements {
			switch ele.Key {
			case "puzzle_name":
				data.CurrPuzzleSet.CurrPuzzle.Metadata.Name = ele.Value
			case "puzzle_author":
				data.CurrPuzzleSet.CurrPuzzle.Metadata.Author = ele.Value
			case "puzzle_hub_check":
				data.CurrPuzzleSet.CurrPuzzle.Metadata.HubLevel = ele.Checked
			case "puzzle_secret_check":
				data.CurrPuzzleSet.CurrPuzzle.Metadata.SecretLevel = ele.Checked
			case "puzzle_darkness_check":
				data.CurrPuzzleSet.CurrPuzzle.Metadata.Darkness = ele.Checked
			}
		}
		ui.CloseDialog(constants.DialogPuzzleSettings)
		data.CurrPuzzleSet.CurrPuzzle.Changed = true
		//if !SavePuzzleSet() {
		//	ui.OpenDialogInStack(constants.DialogUnableToSave)
		//}
	}
}

// combine puzzles

func OpenCombineSetsDialog() {
	combinePzl := ui.Dialogs[constants.DialogCombineSets]
	for _, ele := range combinePzl.Elements {
		if ele.ElementType == ui.ScrollElement {
			content.LoadLocalPuzzleList()
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
					txt.Key = fmt.Sprintf("combine_puzzle_list_%d", i)
					txt.Text.SetPos(pixel.V(xPos, float64(-i)*world.TileSize))
					txt.Text.SetText(data.PuzzleSetFileList[i].Name)
					if i == 0 {
						txt.Text.SetColor(pixel.ToRGBA(constants.ColorBlue))
						data.SelectedPuzzleIndex = 0
					} else {
						txt.Text.SetColor(pixel.ToRGBA(constants.ColorWhite))
					}
					txt.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, ele.ViewPort, func(hvc *data.HoverClick) {
						if combinePzl.Open && combinePzl.Active {
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
	ui.OpenDialogInStack(constants.DialogCombineSets)
}

func OnCombinePuzzleSet() {
	if data.Editor != nil &&
		data.SelectedPuzzleIndex > -1 &&
		data.SelectedPuzzleIndex < len(data.PuzzleSetFileList) {
		err := CombinePuzzleSet(data.PuzzleSetFileList[data.SelectedPuzzleIndex].Filename)
		if err != nil {
			fmt.Println("ERROR:", err)
		}
	}
	ui.CloseDialog(constants.DialogCombineSets)
}

// individual puzzle buttons

func OpenConfirmDelete() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		ui.OpenDialogInStack(constants.DialogAreYouSureDelete)
	}
}

func ConfirmDelete() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		DeletePuzzle()
	}
	ui.CloseDialog(constants.DialogAreYouSureDelete)
}
