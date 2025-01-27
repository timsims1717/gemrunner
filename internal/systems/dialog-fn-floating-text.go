package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"github.com/gopxl/pixel"
	"strconv"
	"strings"
)

// floating text

func customizeFloatingText() {
	floatingTextDlg := ui.Dialogs[constants.DialogFloatingText]
	floatingTextDlg.OnOpen = OnOpenFloatingText
	for _, e := range floatingTextDlg.Elements {
		ele := e
		switch ele.Key {
		case "floating_text_value":
			ele.InputType = ui.Special
		case "confirm":
			ele.OnClick = ConfirmFloatingText
		case "floating_text_time_minus":
			ele.OnClick = func() {
				ChangeNumberInput(floatingTextDlg.Get("floating_text_time_input"), -1)
			}
		case "floating_text_time_plus":
			ele.OnClick = func() {
				ChangeNumberInput(floatingTextDlg.Get("floating_text_time_input"), 1)
			}
		case "floating_text_shadow_check":
			ele.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, floatingTextDlg.ViewPort, func(hvc *data.HoverClick) {
				if floatingTextDlg.Open && floatingTextDlg.Active &&
					!floatingTextDlg.Lock && !floatingTextDlg.Click {
					click := hvc.Input.Get("click")
					if hvc.Hover && click.JustPressed() {
						ui.SetChecked(ele, !ele.Checked)
						if ele.Checked {
							shadowPicked := false
							for _, ele2 := range floatingTextDlg.Elements {
								if strings.Contains(ele2.Key, "check_shadow") && ele2.Checked {
									shadowPicked = true
									break
								}
							}
							if !shadowPicked {
								blueShadow := floatingTextDlg.Get("blue_check_shadow")
								ui.SetChecked(blueShadow, true)
								changeSelectedColor(blueShadow.Key)
							}
						}
					}
				}
			}))
		case "cancel":
			ele.OnClick = CloseDialog(constants.DialogFloatingText)
		default:
			if strings.Contains(ele.Key, "check_color") ||
				strings.Contains(ele.Key, "check_shadow") {
				ele.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, floatingTextDlg.ViewPort, func(hvc *data.HoverClick) {
					if floatingTextDlg.Open && floatingTextDlg.Active &&
						!floatingTextDlg.Lock && !floatingTextDlg.Click {
						click := hvc.Input.Get("click")
						if hvc.Hover && click.JustPressed() && !ele.Checked {
							ui.SetChecked(ele, true)
							changeSelectedColor(ele.Key)
							for _, ele2 := range floatingTextDlg.Elements {
								if ele2.ElementType == ui.CheckboxElement {
									if ((strings.Contains(ele2.Key, "check_color") && strings.Contains(ele.Key, "check_color")) ||
										(strings.Contains(ele2.Key, "check_shadow") && strings.Contains(ele.Key, "check_shadow"))) &&
										ele2.Key != ele.Key {
										ui.SetChecked(ele2, false)
									}
								}
							}
						}
					}
				}))
			} else if ele.OnClick == nil && ele.OnHold == nil {
				ele.OnClick = Test(fmt.Sprintf("pressed button %s", ele.Key))
			}
		}
	}
}

func OnOpenFloatingText() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			ui.CloseDialog(constants.DialogFloatingText)
			return
		}
		theTile := data.CurrPuzzleSet.CurrPuzzle.WrenchTiles[0]
		ftDialog := ui.Dialogs[constants.DialogFloatingText]
		if theTile.TextData == nil {
			data.SelectedTextColor = pixel.ToRGBA(constants.ColorWhite)
			data.SelectedShadowColor = pixel.ToRGBA(constants.ColorBlue)
			ui.ChangeText(ftDialog.Get("floating_text_value"), "")
			for _, ele := range ftDialog.Elements {
				if strings.Contains(ele.Key, "_check_color") {
					ui.SetChecked(ele, false)
				} else if strings.Contains(ele.Key, "_check_shadow") {
					ui.SetChecked(ele, false)
				}
			}
			ui.SetChecked(ftDialog.Get("white_check_color"), true)
			ui.SetChecked(ftDialog.Get("blue_check_shadow"), true)
			ui.SetChecked(ftDialog.Get("floating_text_shadow_check"), true)
			ui.SetChecked(ftDialog.Get("floating_text_show_check"), false)
			ui.SetChecked(ftDialog.Get("floating_text_bob_check"), true)
			ui.ChangeText(ftDialog.Get("floating_text_time_input"), "0")
			// add alignment here
		} else {
			data.SelectedTextColor = theTile.TextData.Color
			data.SelectedShadowColor = theTile.TextData.ShadowCol
			ui.ChangeText(ftDialog.Get("floating_text_value"), theTile.TextData.Raw)
			for _, ele := range ftDialog.Elements {
				if strings.Contains(ele.Key, "_check_color") {
					updateColorCheckbox(ele, theTile.TextData.Color)
				} else if strings.Contains(ele.Key, "_check_shadow") {
					updateColorCheckbox(ele, theTile.TextData.ShadowCol)
				}
			}
			ui.SetChecked(ftDialog.Get("floating_text_shadow_check"), theTile.TextData.HasShadow)
			ui.SetChecked(ftDialog.Get("floating_text_show_check"), theTile.TextData.Prox)
			ui.SetChecked(ftDialog.Get("floating_text_bob_check"), theTile.TextData.Bob)
			timerInput := ftDialog.Get("floating_text_time_input")
			timerInput.InputType = ui.Numeric
			ui.ChangeText(timerInput, strconv.Itoa(theTile.TextData.Timer))
			// add alignment here
		}
	}
}

func ConfirmFloatingText() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			ui.CloseDialog(constants.DialogFloatingText)
			return
		}
		theTile := data.CurrPuzzleSet.CurrPuzzle.WrenchTiles[0]
		ftDialog := ui.Dialogs[constants.DialogFloatingText]

		rawText := ftDialog.Get("floating_text_value").Value

		if rawText == "" {
			if theTile.TextData != nil {
				data.RemoveFloatingText(theTile)
				theTile.TextData = nil
			}
			ui.CloseDialog(constants.DialogFloatingText)
			return
		}

		timer, err := strconv.Atoi(ftDialog.Get("floating_text_time_input").Value)
		if err != nil {
			fmt.Println("WARNING: time input not an int:", err)
			timer = 0
		}

		if theTile.TextData == nil {
			theTile.TextData = &data.FloatingTextData{}
		}
		theTile.TextData.Raw = rawText
		theTile.TextData.Color = data.SelectedTextColor
		theTile.TextData.ShadowCol = data.SelectedShadowColor
		theTile.TextData.Timer = timer
		theTile.TextData.HasShadow = ftDialog.Get("floating_text_shadow_check").Checked
		theTile.TextData.Prox = ftDialog.Get("floating_text_show_check").Checked
		theTile.TextData.Bob = ftDialog.Get("floating_text_bob_check").Checked

		if theTile.FloatingText == nil {
			theTile.FloatingText = data.NewFloatingText().
				WithTile(theTile).
				WithPos(theTile.Object.Pos)
		}

		data.CreateFloatingText(theTile, theTile.TextData)

		ui.CloseDialog(constants.DialogFloatingText)
		data.CurrPuzzleSet.CurrPuzzle.Update = true
		data.CurrPuzzleSet.CurrPuzzle.Changed = true
	}
}

func updateColorCheckbox(x *ui.Element, col pixel.RGBA) {
	key := x.Key[:strings.LastIndex(x.Key, "_")]
	switch key {
	case "white_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorWhite))
	case "red_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorRed))
	case "orange_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorOrange))
	case "green_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorGreen))
	case "cyan_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorCyan))
	case "blue_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorBlue))
	case "purple_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorPurple))
	case "pink_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorPink))
	case "black_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorBlack))
	case "yellow_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorYellow))
	case "gold_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorGold))
	case "brown_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorBrown))
	case "tan_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorTan))
	case "light_gray_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorLightGray))
	case "gray_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorGray))
	case "burnt_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorBurnt))
	}
}
