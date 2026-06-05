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
		case "boss_settings_btn":
			ele.OnClick = OpenBossSettingsDialog
		case "puzzle_settings_btn":
			ele.OnClick = OpenPuzzleSettingsDialog
		case "puzzle_set_settings_btn":
			ele.OnClick = OpenPuzzleSetSettings
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

func OnOpenPuzzleSet() {
	if data.Editor != nil &&
		data.SelectedPuzzleIndex > -1 &&
		data.SelectedPuzzleIndex < len(data.PuzzleSetFileList) {
		err := OpenPuzzleSet(data.PuzzleSetFileList[data.SelectedPuzzleIndex].Filename)
		if err != nil {
			fmt.Println("ERROR:", err)
		}
		if data.CurrPuzzleSet.LastEditedPuzzle > 0 && data.CurrPuzzleSet.LastEditedPuzzle < len(data.CurrPuzzleSet.Puzzles) {
			data.CurrPuzzleSet.SetTo(data.CurrPuzzleSet.LastEditedPuzzle)
		} else {
			data.CurrPuzzleSet.SetToFirst()
		}
		SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
		InitPuzzle(data.CurrentPlayArea)
		SetEditorBoss(data.CurrPuzzleSet.CurrPuzzle.Metadata.Boss)
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
			ele.OnClick = DisposeDialog(constants.DialogPuzzleSettings)
		}
	}
}

func OpenPuzzleSettingsDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		ui.NewDialog(ui.DialogConstructors[constants.DialogPuzzleSettings])
		CustomizeEditorDialog(constants.DialogPuzzleSettings)
		dlg := ui.Dialogs[constants.DialogPuzzleSettings]
		for _, ele := range dlg.Elements {
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
		UpdateDialogView(dlg)
		ui.OpenDialogInStack(constants.DialogPuzzleSettings)
	}
}

func ConfirmPuzzleSettings() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		dlg := ui.Dialogs[constants.DialogPuzzleSettings]
		for _, ele := range dlg.Elements {
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
		ui.Dispose(constants.DialogPuzzleSettings)
		data.CurrPuzzleSet.CurrPuzzle.Changed = true
		UpdateViews()
	}
}

// boss settings dialog

func customizeBossSettings() {
	bossSettingsDlg := ui.Dialogs[constants.DialogBossSettings]
	for _, e := range bossSettingsDlg.Elements {
		ele := e
		switch ele.Key {
		case "enable_boss_check":
			ele.OnClick = ToggleBoss
		case "boss_id":
			ele.Object.Hidden = true
		case "confirm":
			ele.OnClick = ConfirmBossSettings
		case "cancel":
			ele.OnClick = DisposeDialog(constants.DialogBossSettings)
		}
	}
}

func ToggleBoss() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		dlg := ui.Dialogs[constants.DialogBossSettings]
		if dlg.Open {
			bossCheck := dlg.Get("enable_boss_check")
			//currBoss := data.CurrPuzzleSet.CurrPuzzle.Metadata.Boss
			for _, ele := range dlg.Elements {
				switch ele.Key {
				case "boss_name", "select_left", "select_right":
					ele.Object.Hidden = !bossCheck.Checked
				}
			}
		}
	}
}

func OpenBossSettingsDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		ui.NewDialog(ui.DialogConstructors[constants.DialogBossSettings])
		CustomizeEditorDialog(constants.DialogBossSettings)
		dlg := ui.Dialogs[constants.DialogBossSettings]
		currBoss := data.CurrPuzzleSet.CurrPuzzle.Metadata.Boss
		bossEnabled := currBoss != ""
		if currBoss == "" {
			currBoss = constants.BossBlob
		}
		for _, ele := range dlg.Elements {
			switch ele.Key {
			case "boss_name":
				ele.Object.Hidden = !bossEnabled
				ele.Text.SetText(constants.BossNames[currBoss])
			case "boss_id":
				ele.Text.SetText(currBoss)
			case "select_left", "select_right":
				ele.Object.Hidden = !bossEnabled
			case "enable_boss_check":
				ui.SetChecked(ele, bossEnabled)
			}
		}
		UpdateDialogView(dlg)
		ui.OpenDialogInStack(constants.DialogBossSettings)
	}
}

func ConfirmBossSettings() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		dlg := ui.Dialogs[constants.DialogBossSettings]
		if dlg.Get("enable_boss_check").Checked {
			bossId := dlg.Get("boss_id").Text.Raw
			SetEditorBoss(bossId)
		} else {
			RemoveEditorBoss()
			data.CurrPuzzleSet.CurrPuzzle.Metadata.Boss = ""
		}
		ui.Dispose(constants.DialogBossSettings)
		data.CurrPuzzleSet.CurrPuzzle.Changed = true
		UpdateViews()
	}
}

// puzzle set settings

func customizePuzzleSetSettings() {
	dialog := ui.Dialogs[constants.DialogPuzzleSetSettings]
	for _, e := range dialog.Elements {
		ele := e
		switch ele.Key {
		case "confirm":
			ele.OnClick = ConfirmPuzzleSetSettings
		case "sequential_check", "adventure_check":
			ele.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dialog.ViewPort, func(hvc *data.HoverClick) {
				if dialog.Open && dialog.Active && !dialog.Lock && !dialog.Click {
					click := hvc.Input.Get("click")
					if hvc.Hover && click.JustPressed() && !ele.Checked {
						ui.SetChecked(ele, true)
						for _, ele2 := range dialog.Elements {
							if ele2.ElementType == ui.CheckboxElement {
								if (ele2.Key == "sequential_check" || ele2.Key == "adventure_check") &&
									ele2.Key != ele.Key {
									ui.SetChecked(ele2, false)
								}
							}
						}
					}
				}
			}))
		case "cancel":
			ele.OnClick = DisposeDialog(constants.DialogPuzzleSetSettings)
		}
	}
}

func OpenPuzzleSetSettings() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		ui.NewDialog(ui.DialogConstructors[constants.DialogPuzzleSetSettings])
		CustomizeEditorDialog(constants.DialogPuzzleSetSettings)
		dlg := ui.Dialogs[constants.DialogPuzzleSetSettings]
		for _, ele := range dlg.Elements {
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
		UpdateDialogView(dlg)
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
		ui.Dispose(constants.DialogPuzzleSetSettings)
		data.CurrPuzzleSet.CurrPuzzle.Changed = true
		UpdateEditorOptions()
		UpdateViews()
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
		case "puzzle_set_view_title":
			ele.Text.SetText(data.CurrPuzzleSet.CurrPuzzle.Metadata.Name)
			ele.Text.Obj.Pos.X = ele.Text.GetWidth() * -0.5
		case "puzzle_set_view_title_shadow":
			ele.Text.SetText(data.CurrPuzzleSet.CurrPuzzle.Metadata.Name)
			ele.Text.Obj.Pos.X = ele.Text.GetWidth() * -0.5
		case "confirm":
			ele.OnClick = ConfirmRearrangeAdventure
		case "cancel":
			ele.OnClick = DisposeDialog(key)
		//case "zoom_in":
		//	ele.OnClick = ZoomFn(key, true, true)
		//case "zoom_out":
		//	ele.OnClick = ZoomFn(key, true, false)
		case "puzzle_set_view":
			//for _, ele1 := range ele.Elements {
			//	switch ele1.Key {
			//	case "puzzle_center", "puzzle_left", "puzzle_right",
			//		"puzzle_top_center", "puzzle_top_left", "puzzle_top_right",
			//		"puzzle_bot_center", "puzzle_bot_left", "puzzle_bot_right",
			//		"puzzle_float_center", "puzzle_float_left", "puzzle_float_right":
			//		ele1.Border.Style = data.ThinBorderWhite
			//	}
			//}
			SetAdventureViewMovement(ele, rearrangePzl)
		}
	}
	data.AdventureViewGridPos = data.CurrPuzzleSet.CurrPuzzle.Grid
	data.AdventureViewGridMap = make(map[world.Coords]data.AdvViewPzl)
	data.AdventureViewGridArr = make(map[int]world.Coords)
	//switch data.AdventureViewZoomLevel {
	//case 0:
	//	AdventureViewZoomZero(key, true)
	//case 1:
	AdventureViewZoomOne(key)
	pzlView := rearrangePzl.Get("puzzle_set_view")
	for _, e1 := range pzlView.Elements {
		ele1 := e1
		SetPreviewMediumClick(ele1, rearrangePzl, pzlView)
	}
	//case 2:
	//	AdventureViewZoomTwo(key, true)
	//}
	UpdateDialogView(rearrangePzl)
	ui.OpenDialogInStack(key)
}

//func ZoomFn(key string, rearrange, in bool) func() {
//	return func() {
//		if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
//			switch data.AdventureViewZoomLevel {
//			case 0:
//				if in {
//					AdventureViewZoomOne(key, rearrange)
//				}
//			case 1:
//				if in {
//					AdventureViewZoomTwo(key, rearrange)
//				} else {
//					AdventureViewZoomZero(key, rearrange)
//				}
//			case 2:
//				if !in {
//					AdventureViewZoomOne(key, rearrange)
//				}
//			}
//		}
//	}
//}

//func Zoom(key string, rearrange, in bool) {
//	if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
//		switch data.AdventureViewZoomLevel {
//		case 0:
//			if in {
//				AdventureViewZoomOne(key, rearrange)
//			}
//		case 1:
//			if in {
//				AdventureViewZoomTwo(key, rearrange)
//			} else {
//				AdventureViewZoomZero(key, rearrange)
//			}
//		case 2:
//			if !in {
//				AdventureViewZoomOne(key, rearrange)
//			}
//		}
//	}
//}

func ConfirmRearrangeAdventure() {
	if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
		for grid, avp := range data.AdventureViewGridMap {
			if avp.ViewIndex == 0 {
				// set current puzzle to this one
				data.CurrPuzzleSet.SetTo(avp.SetIndex)
			}
			pzl := data.CurrPuzzleSet.Puzzles[avp.SetIndex]
			oldGrid := pzl.Grid
			pzl.Grid = grid
			delete(data.CurrPuzzleSet.PuzzGrid, oldGrid)
			data.CurrPuzzleSet.PuzzGrid[grid] = avp.SetIndex
		}
		SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
		InitPuzzle(data.CurrentPlayArea)
		ui.Dispose(constants.DialogRearrangeAdventureSet)
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
		SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
		InitPuzzle(data.CurrentPlayArea)

		ui.Dispose(constants.DialogRearrangePuzzleSet)
	}
}
