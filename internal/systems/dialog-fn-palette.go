package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"strings"
)

// palette options

func customizePaletteOptions() {
	paletteOptions := ui.Dialogs[constants.DialogPalette]
	paletteOptions.OnOpen = OnOpenPaletteOptions
	for _, e := range paletteOptions.Elements {
		ele := e
		switch ele.Key {
		case "cancel":
			ele.OnClick = CloseDialog(constants.DialogPalette)
		case "confirm":
			ele.OnClick = OnConfirmPaletteOptions
		default:
			if strings.Contains(ele.Key, "check_color") {
				ele.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, paletteOptions.ViewPort, func(hvc *data.HoverClick) {
					if paletteOptions.Open && paletteOptions.Active &&
						!paletteOptions.Lock && !paletteOptions.Click {
						click := hvc.Input.Get("click")
						if hvc.Hover && click.JustPressed() && !ele.Checked {
							ui.SetChecked(ele, true)
							paletteOptions.Get("palette_current").Text.SetText(ele.HelpText)
							for _, ele2 := range paletteOptions.Elements {
								if ele2.ElementType == ui.CheckboxElement {
									if strings.Contains(ele2.Key, "check_color") &&
										strings.Contains(ele.Key, "check_color") &&
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

func OnOpenPaletteOptions() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		paletteOptions := ui.Dialogs[constants.DialogPalette]
		for _, e := range paletteOptions.Elements {
			if e.ElementType == ui.CheckboxElement {
				ui.SetChecked(e, false)
			}
		}
		var checkbox *ui.Element
		switch data.Editor.PaletteColor {
		case data.ColorDefault:
			checkbox = paletteOptions.Get("white_check_color")
		default:
			checkbox = paletteOptions.Get(fmt.Sprintf("%s_check_color", data.Editor.PaletteColor.String()))
		}
		ui.SetChecked(checkbox, true)
		paletteOptions.Get("palette_current").Text.SetText(checkbox.HelpText)
	}
}

func OnConfirmPaletteOptions() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		paletteOptions := ui.Dialogs[constants.DialogPalette]
		for _, e := range paletteOptions.Elements {
			if e.ElementType == ui.CheckboxElement && e.Checked {
				switch e.Key {
				case "white_check_color":
					data.Editor.PaletteColor = data.ColorDefault
				case "yellow_check_color":
					data.Editor.PaletteColor = data.NonPlayerYellow
				case "brown_check_color":
					data.Editor.PaletteColor = data.NonPlayerBrown
				case "gray_check_color":
					data.Editor.PaletteColor = data.NonPlayerGray
				case "cyan_check_color":
					data.Editor.PaletteColor = data.NonPlayerCyan
				case "blue_check_color":
					data.Editor.PaletteColor = data.PlayerBlue
				case "green_check_color":
					data.Editor.PaletteColor = data.PlayerGreen
				case "purple_check_color":
					data.Editor.PaletteColor = data.PlayerPurple
				case "orange_check_color":
					data.Editor.PaletteColor = data.PlayerOrange
				case "red_check_color":
					data.Editor.PaletteColor = data.NonPlayerRed
				}
			}
		}
		ui.CloseDialog(constants.DialogPalette)
	}
}
