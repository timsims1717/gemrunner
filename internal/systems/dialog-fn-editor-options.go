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
	"golang.org/x/image/colornames"
	"strconv"
)

// options panel

func customizeEditorOptions(key string) {
	editorOptions := ui.Dialogs[key]
	for _, e := range editorOptions.Elements {
		ele := e
		switch ele.Key {
		case "puzzle_number":
			ele.Text.SetText("0001")
		case "test_btn":
			ele.OnClick = TestPuzzle
		case "save_btn":
			ele.OnClick = OnSavePuzzleSet
		case "open_btn":
			ele.OnClick = OpenOpenPuzzleDialog
		case "new_btn":
			ele.OnClick = NewPuzzle
		case "combine_btn":
			ele.OnClick = OpenCombineSetsDialog
		case "add_btn":
			ele.OnClick = InsertPuzzle
		case "delete_btn":
			ele.OnClick = OpenConfirmDelete
		case "left_puz_btn":
			ele.OnClick = PrevPuzzle
		case "right_puz_btn":
			ele.OnClick = NextPuzzle
		case "up_puz_btn":
			ele.OnClick = UpPuzzle
		case "down_puz_btn":
			ele.OnClick = DownPuzzle
		case "rearrange_btn":
			ele.OnClick = OpenRearrangeDialog
		case "world_btn":
			ele.OnClick = OpenChangeWorldDialog
		case "name_btn":
			ele.OnClick = OpenDialog(constants.DialogChangeName)
		case "puzzle_settings_btn":
			ele.OnClick = OpenDialog(constants.DialogPuzzleSettings)
		case "puzzle_set_settings_btn":
			ele.OnClick = OpenPuzzleSetSettingsDialog
		case "exit_editor_btn":
			ele.OnClick = ExitEditor
		}
	}
}

func UpdateEditorOptions() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		optionsKeys := []string{constants.DialogEditorOptionsRight, constants.DialogEditorOptionsBot}
		for _, key := range optionsKeys {
			editorOptions := ui.Dialogs[key]
			for _, e := range editorOptions.Elements {
				ele := e
				switch ele.Key {
				case "add_btn":
					ele.OnClick = OpenDialog(constants.DialogAddPuzzle)
				case "puzzle_number":
					if data.CurrPuzzleSet.Metadata.Adventure {
						ele.Text.SetText(fmt.Sprintf("%02d,%02d", data.CurrPuzzleSet.CurrPuzzle.Grid.X, data.CurrPuzzleSet.CurrPuzzle.Grid.Y))
					} else {
						ele.Text.SetText(fmt.Sprintf("%04d", data.CurrPuzzleSet.PuzzleIndex+1))
					}
				case "left_puz_btn":
					if data.CurrPuzzleSet.Metadata.Adventure {
						ele.OnClick = LeftPuzzle
					} else {
						ele.OnClick = PrevPuzzle
					}
				case "right_puz_btn":
					if data.CurrPuzzleSet.Metadata.Adventure {
						ele.OnClick = RightPuzzle
					} else {
						ele.OnClick = NextPuzzle
					}
				case "up_puz_btn":
					ele.Object.Hidden = !data.CurrPuzzleSet.Metadata.Adventure
				case "down_puz_btn":
					ele.Object.Hidden = !data.CurrPuzzleSet.Metadata.Adventure
				}
			}
		}
	}
}

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
			if total < data.SelectedPuzzleIndex {
				data.SelectedPuzzleIndex = 0
			}
			xPos := ele.ViewPort.CamPos.X - ele.ViewPort.Rect.W()*0.5 + 4
			for i := 0; i < total; i++ {
				yPos := float64(-i) * world.TileSize
				index := i
				if len(ele.Elements) <= i {
					ec := ui.ElementConstructor{
						Key:         "sub_text",
						ElementType: ui.TextElement,
						SubElements: nil,
						Anchor:      pixel.Right,
					}
					t := ui.CreateTextElement(ec, ele.ViewPort)
					ele.Elements = append(ele.Elements, t)
				}
				txt := ele.Elements[i]
				if txt.ElementType == ui.TextElement {
					txt.Key = fmt.Sprintf("open_puzzle_list_%d", i)
					txt.Text.SetPos(pixel.V(xPos, yPos))
					txt.Text.SetText(data.PuzzleSetFileList[i].Name)
					if i == data.SelectedPuzzleIndex {
						txt.Text.SetColor(pixel.ToRGBA(constants.ColorBlue))
						ui.MoveScrollToInclude(ele, yPos+world.HalfSize+2, yPos-world.HalfSize)
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
			ui.MoveToScrollTop(ele)
		}
	}
}

func OpenOpenPuzzleDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		if data.CurrPuzzleSet.Metadata.Filename == "" && data.CurrPuzzleSet.NeedToSave {
			ui.SetCloseSpcFn(constants.DialogChangeName, func() {
				ui.OpenDialogInStack(constants.DialogOpenPuzzle)
			})
			ui.OpenDialogInStack(constants.DialogChangeName)
		} else if data.CurrPuzzleSet.NeedToSave {
			if SavePuzzleSet() {
				ui.OpenDialogInStack(constants.DialogOpenPuzzle)
			} else {
				ui.SetTempOnClick(constants.DialogUnableToSaveConfirm, "confirm", func() {
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
		if data.CurrPuzzleSet.Metadata.Filename == "" && data.CurrPuzzleSet.NeedToSave {
			ui.SetCloseSpcFn(constants.DialogChangeName, func() {
				NewPuzzleSet()
			})
			ui.OpenDialogInStack(constants.DialogChangeName)
		} else if data.CurrPuzzleSet.NeedToSave {
			if SavePuzzleSet() {
				NewPuzzleSet()
			} else {
				ui.SetTempOnClick(constants.DialogUnableToSaveConfirm, "confirm", func() {
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
		if data.CurrPuzzleSet.Metadata.Filename == "" && data.CurrPuzzleSet.NeedToSave {
			ui.SetCloseSpcFn(constants.DialogChangeName, func() {
				state.SwitchState(constants.MainMenuKey)
			})
			ui.OpenDialogInStack(constants.DialogChangeName)
		} else if data.CurrPuzzleSet.NeedToSave {
			if SavePuzzleSet() {
				state.SwitchState(constants.MainMenuKey)
			} else {
				ui.SetTempOnClick(constants.DialogUnableToSaveConfirm, "confirm", func() {
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
			if data.CurrPuzzleSet.Metadata.Filename == "" && data.CurrPuzzleSet.NeedToSave {
				ui.SetCloseSpcFn(constants.DialogChangeName, func() {
					state.PushState(constants.TestStateKey)
				})
				ui.OpenDialogInStack(constants.DialogChangeName)
			} else if data.CurrPuzzleSet.NeedToSave {
				if SavePuzzleSet() {
					state.PushState(constants.TestStateKey)
				} else {
					ui.SetTempOnClick(constants.DialogUnableToSaveConfirm, "confirm", func() {
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
		} else if data.CurrPuzzleSet.NeedToSave {
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
				ui.SetText(inEle, data.CurrPuzzleSet.Metadata.Name)
			}
		} else {
			ui.SetText(inEle, "Untitled")
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

// add puzzle dialog

func customizeAddPuzzle() {
	addPuzzleDlg := ui.Dialogs[constants.DialogAddPuzzle]
	addPuzzleDlg.OnOpen = OpenAddPuzzleDialog
	for _, e := range addPuzzleDlg.Elements {
		ele := e
		switch ele.Key {
		case "cancel":
			ele.OnClick = CloseDialog(constants.DialogAddPuzzle)
		case "add_up":
			ele.OnClick = func() {
				AddPuzzleUp()
				ui.CloseDialog(constants.DialogAddPuzzle)
			}
		case "add_down":
			ele.OnClick = func() {
				AddPuzzleDown()
				ui.CloseDialog(constants.DialogAddPuzzle)
			}
		case "add_right":
			ele.OnClick = func() {
				AddPuzzleRight()
				ui.CloseDialog(constants.DialogAddPuzzle)
			}
		case "add_left":
			ele.OnClick = func() {
				AddPuzzleLeft()
				ui.CloseDialog(constants.DialogAddPuzzle)
			}
		}
	}
}

func OpenAddPuzzleDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		currPuzzle := data.CurrPuzzleSet.CurrPuzzle
		currGrid := currPuzzle.Grid
		for _, e := range ui.Dialogs[constants.DialogAddPuzzle].Elements {
			ele := e
			switch ele.Key {
			case "add_up":
				grid := currGrid
				grid.Y++
				if _, ok := data.CurrPuzzleSet.PuzzGrid[grid]; ok {
					ele.Object.Mask = pixel.ToRGBA(constants.ColorDisable)
					ele.Disable(true)
				} else {
					ele.Object.Mask = pixel.ToRGBA(colornames.White)
					ele.Disable(false)
				}
			case "add_down":
				grid := currGrid
				grid.Y--
				if _, ok := data.CurrPuzzleSet.PuzzGrid[grid]; ok {
					ele.Object.Mask = pixel.ToRGBA(constants.ColorDisable)
					ele.Disable(true)
				} else {
					ele.Object.Mask = pixel.ToRGBA(colornames.White)
					ele.Disable(false)
				}
			case "add_right":
				grid := currGrid
				grid.X++
				if _, ok := data.CurrPuzzleSet.PuzzGrid[grid]; ok {
					ele.Object.Mask = pixel.ToRGBA(constants.ColorDisable)
					ele.Disable(true)
				} else {
					ele.Object.Mask = pixel.ToRGBA(colornames.White)
					ele.Disable(false)
				}
			case "add_left":
				grid := currGrid
				grid.X--
				if _, ok := data.CurrPuzzleSet.PuzzGrid[grid]; ok {
					ele.Object.Mask = pixel.ToRGBA(constants.ColorDisable)
					ele.Disable(true)
				} else {
					ele.Object.Mask = pixel.ToRGBA(colornames.White)
					ele.Disable(false)
				}
			}
		}
	}
}

// puzzle settings dialog

func customizePuzzleSettings() {
	puzzleSettingsDlg := ui.Dialogs[constants.DialogPuzzleSettings]
	puzzleSettingsDlg.OnOpen = OpenPuzzleSettingsDialog
	for _, e := range puzzleSettingsDlg.Elements {
		ele := e
		switch ele.Key {
		case "confirm":
			ele.OnClick = ConfirmPuzzleSettings
		case "puzzle_width_minus":
			ele.OnClick = func() {
				ChangeNumberInputWithLimits(puzzleSettingsDlg.Get("puzzle_width_input"), -1, constants.PuzzleMinWidth, constants.PuzzleMaxWidth)
			}
		case "puzzle_width_plus":
			ele.OnClick = func() {
				ChangeNumberInputWithLimits(puzzleSettingsDlg.Get("puzzle_width_input"), 1, constants.PuzzleMinWidth, constants.PuzzleMaxWidth)
			}
		case "puzzle_height_minus":
			ele.OnClick = func() {
				ChangeNumberInputWithLimits(puzzleSettingsDlg.Get("puzzle_height_input"), -1, constants.PuzzleMinHeight, constants.PuzzleMaxHeight)
			}
		case "puzzle_height_plus":
			ele.OnClick = func() {
				ChangeNumberInputWithLimits(puzzleSettingsDlg.Get("puzzle_height_input"), 1, constants.PuzzleMinHeight, constants.PuzzleMaxHeight)
			}
		case "cancel":
			ele.OnClick = CloseDialog(constants.DialogPuzzleSettings)
		}
	}
}

func OpenPuzzleSettingsDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		for _, ele := range ui.Dialogs[constants.DialogPuzzleSettings].Elements {
			switch ele.Key {
			case "puzzle_name":
				if data.CurrPuzzleSet.CurrPuzzle.Metadata.Name != "" {
					if ele.Value != data.CurrPuzzleSet.CurrPuzzle.Metadata.Name {
						ui.SetText(ele, data.CurrPuzzleSet.CurrPuzzle.Metadata.Name)
					}
				} else {
					ui.SetText(ele, "Untitled")
				}
			case "puzzle_author":
				if data.CurrPuzzleSet.CurrPuzzle.Metadata.Author != "" {
					if ele.Value != data.CurrPuzzleSet.CurrPuzzle.Metadata.Author {
						ui.SetText(ele, data.CurrPuzzleSet.CurrPuzzle.Metadata.Author)
					}
				} else {
					ui.SetText(ele, constants.Username)
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
			case "puzzle_width_input":
				ele.InputType = ui.Numeric
				ui.SetText(ele, fmt.Sprintf("%d", data.CurrPuzzleSet.CurrPuzzle.Metadata.Width))
			case "puzzle_height_input":
				ele.InputType = ui.Numeric
				ui.SetText(ele, fmt.Sprintf("%d", data.CurrPuzzleSet.CurrPuzzle.Metadata.Height))
			}
		}
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
			case "puzzle_width_input":
				wi, err := strconv.Atoi(ele.Text.Raw)
				if err != nil {
					fmt.Println("WARNING: width is not an int:", err)
					wi = constants.PuzzleWidth
				}
				data.CurrPuzzleSet.CurrPuzzle.SetWidth(wi)
			case "puzzle_height_input":
				hi, err := strconv.Atoi(ele.Text.Raw)
				if err != nil {
					fmt.Println("WARNING: height is not an int:", err)
					hi = constants.PuzzleHeight
				}
				data.CurrPuzzleSet.CurrPuzzle.SetHeight(hi)
			}
		}
		ui.CloseDialog(constants.DialogPuzzleSettings)
		data.CurrPuzzleSet.CurrPuzzle.Changed = true
		UpdateViews()
		//if !SavePuzzleSet() {
		//	ui.OpenDialogInStack(constants.DialogUnableToSave)
		//}
	}
}

// puzzle set settings

func OpenPuzzleSetSettingsDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		for _, ele := range ui.Dialogs[constants.DialogPuzzleSetSettings].Elements {
			switch ele.Key {
			case "puzzle_set_name":
				if data.CurrPuzzleSet.Metadata.Name != "" {
					if ele.Value != data.CurrPuzzleSet.Metadata.Name {
						ui.SetText(ele, data.CurrPuzzleSet.Metadata.Name)
					}
				} else {
					ui.SetText(ele, "Untitled")
				}
			case "puzzle_set_author":
				if data.CurrPuzzleSet.Metadata.Author != "" {
					if ele.Value != data.CurrPuzzleSet.Metadata.Author {
						ui.SetText(ele, data.CurrPuzzleSet.Metadata.Author)
					}
				} else {
					ui.SetText(ele, constants.Username)
				}
			case "sequential_check":
				ui.SetChecked(ele, !data.CurrPuzzleSet.Metadata.Adventure)
			case "adventure_check":
				ui.SetChecked(ele, data.CurrPuzzleSet.Metadata.Adventure)
			}
		}
		ui.OpenDialogInStack(constants.DialogPuzzleSetSettings)
	}
}

func ConfirmPuzzleSetSettings() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		for _, ele := range ui.Dialogs[constants.DialogPuzzleSetSettings].Elements {
			switch ele.Key {
			case "puzzle_set_name":
				data.CurrPuzzleSet.Metadata.Name = ele.Value
			case "puzzle_set_author":
				data.CurrPuzzleSet.Metadata.Author = ele.Value
			case "adventure_check":
				data.CurrPuzzleSet.Metadata.Adventure = ele.Checked
			}
		}
		ui.CloseDialog(constants.DialogPuzzleSetSettings)
		data.CurrPuzzleSet.CurrPuzzle.Changed = true
		UpdateEditorOptions()
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
			ui.MoveToScrollTop(ele)
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

func OpenRearrangeDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		if data.CurrPuzzleSet.Metadata.Adventure {
			OpenRearrangeAdventureDialog()
		} else {
			OpenRearrangePuzzlesDialog()
		}
	}
}

// rearrange puzzles (adventure)

func OpenRearrangeAdventureDialog() {
	key := constants.DialogRearrangeAdventureSet
	ui.NewDialog(ui.DialogConstructors[key])
	rearrangePzl := ui.Dialogs[key]
	for _, e := range rearrangePzl.Elements {
		ele := e
		switch ele.Key {
		case "confirm":
			ele.OnClick = ConfirmRearrangeAdventure
		case "cancel":
			ele.OnClick = DisposeDialog(key)
		case "zoom_in":
			ele.OnClick = Zoom(key, true, true)
		case "zoom_out":
			ele.OnClick = Zoom(key, true, false)
		case "puzzle_set_view":
			for _, ele1 := range ele.Elements {
				switch ele1.Key {
				case "puzzle_center", "puzzle_left", "puzzle_right",
					"puzzle_top_center", "puzzle_top_left", "puzzle_top_right",
					"puzzle_bot_center", "puzzle_bot_left", "puzzle_bot_right",
					"puzzle_float_center", "puzzle_float_left", "puzzle_float_right":
					ele1.Border.Style = ui.ThinBorderWhite
				}
			}
		}
	}
	switch data.AdventureViewZoomLevel {
	case 0:
		AdventureViewZoomZero(key, true)
	case 1:
		AdventureViewZoomOne(key, true)
	case 2:
		AdventureViewZoomTwo(key, true)
	}
	UpdateDialogView(rearrangePzl)
	ui.OpenDialogInStack(key)
}

func Zoom(key string, rearrange, in bool) func() {
	return func() {
		if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
			switch data.AdventureViewZoomLevel {
			case 0:
				if in {
					AdventureViewZoomOne(key, rearrange)
				}
			case 1:
				if in {
					AdventureViewZoomTwo(key, rearrange)
				} else {
					AdventureViewZoomZero(key, rearrange)
				}
			case 2:
				if !in {
					AdventureViewZoomOne(key, rearrange)
				}
			}
		}
	}
}

func ConfirmRearrangeAdventure() {
	if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {

		ui.Dispose(constants.DialogRearrangeAdventureSet)
	}
}

var keepTheseKeys1 = []string{
	"puzzle_center", "puzzle_left", "puzzle_right",
	"puzzle_top_center", "puzzle_top_left", "puzzle_top_right",
	"puzzle_bot_center", "puzzle_bot_left", "puzzle_bot_right",
	"puzzle_float_center", "puzzle_float_left", "puzzle_float_right",
}

func AdventureViewZoomZero(key string, rearrange bool) {
	if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
		dialog := ui.Dialogs[key]
		pzlView := dialog.Get("puzzle_set_view")
		var keepTheseKeys []string
		data.AdventureViewZoomLevel = 0
		data.PuzzleSetViewIsMoving = false
		data.PuzzleSetViewIndex = data.CurrPuzzleSet.PuzzleIndex
		pzlCoords := data.CurrPuzzleSet.CurrPuzzle.Grid
		for grid, i := range data.CurrPuzzleSet.PuzzGrid {
			pzl := data.CurrPuzzleSet.Puzzles[i]
			col := pzl.Metadata.PrimaryColor

			// positions
			pos0 := pixel.V(float64(grid.X*8), float64(grid.Y*5))
			if grid.X < 0 {
				pos0.X += 1.
			}
			if grid.Y < 0 {
				pos0.Y += 1.
			}
			if pzlCoords == grid { // move camera here
				pzlView.ViewPort.CamPos = pos0
			}

			// small, zoom level 0
			z0Key := fmt.Sprintf("z0_%d_%d", grid.X, grid.Y)
			z0ec := ui.ElementConstructor{
				Key:         z0Key,
				SprKey:      "white_level_preview",
				Batch:       constants.UIBatch,
				Position:    pos0,
				ElementType: ui.SpriteElement,
			}
			z0Ele := ui.CreateSpriteElement(z0ec)
			z0Ele.Object.Mask = col
			pzlView.Elements = append(pzlView.Elements, z0Ele)
			keepTheseKeys = append(keepTheseKeys, z0Key)
			ResetAdventureView(key, keepTheseKeys, 0)
		}
	}
}

func AdventureViewZoomOne(key string, rearrange bool) {
	if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
		dialog := ui.Dialogs[key]
		pzlView := dialog.Get("puzzle_set_view")
		var keepTheseKeys []string
		data.AdventureViewZoomLevel = 1
		data.PuzzleSetViewIsMoving = false
		data.PuzzleSetViewIndex = data.CurrPuzzleSet.PuzzleIndex
		pzlCoords := data.CurrPuzzleSet.CurrPuzzle.Grid
		for grid, i := range data.CurrPuzzleSet.PuzzGrid {
			pzl := data.CurrPuzzleSet.Puzzles[i]
			col := pzl.Metadata.PrimaryColor

			// positions
			pos1 := pixel.V(float64(grid.X*29)+0.5, float64(grid.Y*17)+0.5)
			if grid.X < 0 {
				pos1.X += 1.
			}
			if grid.Y < 0 {
				pos1.Y += 1.
			}
			if pzlCoords == grid { // move camera here
				pzlView.ViewPort.CamPos = pos1
			}

			// medium, zoom level 1
			z1Key := fmt.Sprintf("z1_%d_%d", grid.X, grid.Y)
			pic := CreatePuzzlePreviewMedium(pzl)
			pSpr := pixel.NewSprite(pic, pic.Bounds())
			z1ec := ui.ElementConstructor{
				Key:         z1Key,
				Batch:       constants.UIBatch,
				Position:    pos1,
				ElementType: ui.SpriteElement,
			}
			z1Ele := ui.CreatePixelSpriteElement(z1ec, pSpr)
			z1Ele.Object.Mask = col
			z1Ele.Object.Hidden = data.AdventureViewZoomLevel != 1
			pzlView.Elements = append(pzlView.Elements, z1Ele)
			z1BKey := fmt.Sprintf("z1b_%d_%d", grid.X, grid.Y)
			z1Bec := ui.ElementConstructor{
				Key:         z1BKey,
				SprKey:      "white_level_preview_border",
				Batch:       constants.UIBatch,
				Position:    pos1,
				ElementType: ui.SpriteElement,
			}
			z1BEle := ui.CreateSpriteElement(z1Bec)
			z1BEle.Object.Hidden = data.AdventureViewZoomLevel != 1
			pzlView.Elements = append(pzlView.Elements, z1BEle)
			keepTheseKeys = append(keepTheseKeys, z1Key)
			keepTheseKeys = append(keepTheseKeys, z1BKey)
			ResetAdventureView(key, keepTheseKeys, 1)
		}
	}
}

func AdventureViewZoomTwo(key string, rearrange bool) {
	if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
		dialog := ui.Dialogs[key]
		pzlView := dialog.Get("puzzle_set_view")
		data.AdventureViewZoomLevel = 2
		data.PuzzleSetViewIsMoving = false
		data.PuzzleSetViewIndex = data.CurrPuzzleSet.PuzzleIndex
		pzlCoords := data.CurrPuzzleSet.CurrPuzzle.Grid

		CreatePuzzlePreview(pzlView.Get("puzzle_center"), data.PuzzleSetViewIndex)
		CreatePuzzlePreview(pzlView.Get("puzzle_left"), data.CurrPuzzleSet.GetGrid(world.Coords{X: pzlCoords.X - 1, Y: pzlCoords.Y}))
		CreatePuzzlePreview(pzlView.Get("puzzle_right"), data.CurrPuzzleSet.GetGrid(world.Coords{X: pzlCoords.X + 1, Y: pzlCoords.Y}))
		CreatePuzzlePreview(pzlView.Get("puzzle_top_center"), data.CurrPuzzleSet.GetGrid(world.Coords{X: pzlCoords.X, Y: pzlCoords.Y + 1}))
		CreatePuzzlePreview(pzlView.Get("puzzle_top_left"), data.CurrPuzzleSet.GetGrid(world.Coords{X: pzlCoords.X - 1, Y: pzlCoords.Y + 1}))
		CreatePuzzlePreview(pzlView.Get("puzzle_top_right"), data.CurrPuzzleSet.GetGrid(world.Coords{X: pzlCoords.X + 1, Y: pzlCoords.Y + 1}))
		CreatePuzzlePreview(pzlView.Get("puzzle_bot_center"), data.CurrPuzzleSet.GetGrid(world.Coords{X: pzlCoords.X, Y: pzlCoords.Y - 1}))
		CreatePuzzlePreview(pzlView.Get("puzzle_bot_left"), data.CurrPuzzleSet.GetGrid(world.Coords{X: pzlCoords.X - 1, Y: pzlCoords.Y - 1}))
		CreatePuzzlePreview(pzlView.Get("puzzle_bot_right"), data.CurrPuzzleSet.GetGrid(world.Coords{X: pzlCoords.X + 1, Y: pzlCoords.Y - 1}))

		pzlView.ViewPort.CamPos = pzlView.Get("puzzle_center").Object.Pos

		ResetAdventureView(key, nil, 2)
	}
}

func ResetAdventureView(key string, keepTheseKeys []string, zoom int) {
	if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
		dialog := ui.Dialogs[key]
		pzlView := dialog.Get("puzzle_set_view")
		for i := len(pzlView.Elements) - 1; i >= 0; i-- {
			ele := pzlView.Elements[i]
			// remove extra elements
			if util.ContainsStr(ele.Key, keepTheseKeys1) {
				if zoom != 2 { // dispose of puzzle view
					ui.DisposeSubElements(ele.Elements)
					ele.Object.Hidden = true
				}
			} else if !util.ContainsStr(ele.Key, keepTheseKeys) {
				if i < len(pzlView.Elements)+1 {
					pzlView.Elements = append(pzlView.Elements[:i], pzlView.Elements[i+1:]...)
				} else if i == 0 {
					pzlView.Elements = []*ui.Element{}
				} else {
					pzlView.Elements = pzlView.Elements[:i]
				}
				ui.DisposeSubElements(ele.Elements)
				myecs.Manager.DisposeEntity(ele.Entity)
			}
		}
	}
}

// rearrange puzzles (sequence)

func OpenRearrangePuzzlesDialog() {
	ui.NewDialog(ui.DialogConstructors[constants.DialogRearrangePuzzleSet])
	rearrangePzl := ui.Dialogs[constants.DialogRearrangePuzzleSet]
	for _, e := range rearrangePzl.Elements {
		ele := e
		switch ele.Key {
		case "confirm":
			ele.OnClick = ConfirmRearrangedPuzzles
		case "rearrange_next":
			ele.OnHold = PuzzleSetViewNextPuzzle(rearrangePzl)
			ele.OnClick = PuzzleSetViewNextPuzzle(rearrangePzl)
		case "rearrange_prev":
			ele.OnHold = PuzzleSetViewPrevPuzzle(rearrangePzl)
			ele.OnClick = PuzzleSetViewPrevPuzzle(rearrangePzl)
		case "rearrange_swap_next":
			ele.OnHold = PuzzleSetViewSwapNext(rearrangePzl)
			ele.OnClick = PuzzleSetViewSwapNext(rearrangePzl)
		case "rearrange_swap_prev":
			ele.OnHold = PuzzleSetViewSwapPrev(rearrangePzl)
			ele.OnClick = PuzzleSetViewSwapPrev(rearrangePzl)
		case "rearrange_end":
			ele.OnHold = PuzzleSetViewSwapEnd(rearrangePzl)
			ele.OnClick = PuzzleSetViewSwapEnd(rearrangePzl)
		case "rearrange_begin":
			ele.OnHold = PuzzleSetViewSwapToBegin(rearrangePzl)
			ele.OnClick = PuzzleSetViewSwapToBegin(rearrangePzl)
		case "cancel":
			ele.OnClick = DisposeDialog(constants.DialogRearrangePuzzleSet)
		}
	}
	UpdateDialogView(rearrangePzl)
	data.PuzzleSetViewAllowEnd = false
	data.PuzzleSetViewIsMoving = false
	data.PuzzleSetViewIndex = data.CurrPuzzleSet.PuzzleIndex
	data.PuzzleSetViewPuzzles = make([]int, len(data.CurrPuzzleSet.Puzzles))
	for i := range data.PuzzleSetViewPuzzles {
		data.PuzzleSetViewPuzzles[i] = i
	}
	pzlView := rearrangePzl.Get("puzzle_set_view")
	CreatePuzzlePreview(pzlView.Get("puzzle_center"), data.PuzzleSetViewIndex)
	CreatePuzzlePreview(pzlView.Get("puzzle_left"), data.PuzzleSetViewIndex-1)
	CreatePuzzlePreview(pzlView.Get("puzzle_right"), data.PuzzleSetViewIndex+1)
	ResetPuzzleSetView(rearrangePzl)
	PuzzleSetViewNameAndNum(rearrangePzl, data.PuzzleSetViewIndex)
	ui.OpenDialogInStack(constants.DialogRearrangePuzzleSet)
}

func ConfirmRearrangedPuzzles() {
	if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
		// rearrange puzzle set
		var newPuzzles []*data.Puzzle
		for _, i := range data.PuzzleSetViewPuzzles {
			newPuzzles = append(newPuzzles, data.CurrPuzzleSet.Puzzles[i])
		}
		data.CurrPuzzleSet.Puzzles = newPuzzles
		// go to currently selected puzzle
		data.CurrPuzzleSet.SetTo(data.PuzzleSetViewIndex)
		PuzzleInit()

		ui.Dispose(constants.DialogRearrangePuzzleSet)
	}
}
