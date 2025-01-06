package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/load"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
	"strings"
)

func MainDialogs(win *pixelgl.Window) {
	ui.NewDialog(load.MainMenuConstructor)
	ui.NewDialog(load.AddPlayersConstructor)
	ui.NewDialog(ui.DialogConstructors[constants.DialogPlayLocal])
	customizeMainDialogs(win)
}

func DisposeMainDialogs() {
	for k, d := range ui.Dialogs {
		switch k {
		case constants.DialogMainMenu,
			constants.DialogAddPlayers,
			constants.DialogPlayLocal:
			ui.DisposeDialog(d)
		}
	}
}

func customizeMainDialogs(win *pixelgl.Window) {
	for key := range ui.Dialogs {
		dialog := ui.Dialogs[key]
		for _, e := range dialog.Elements {
			ele := e
			if ele.ElementType == ui.ButtonElement {
				switch ele.Key {
				case "play_local_game_btn":
					ele.OnClick = OpenAddPlayers
				case "start_editor_btn":
					ele.OnClick = StartEditor
				case "quit_btn":
					ele.OnClick = QuitGame(win)
				case "confirm_add_players":
					ele.OnClick = ConfirmAddPlayers
				default:
					switch dialog.Key {
					default:
						if strings.Contains(ele.Key, "cancel") {
							ele.OnClick = CloseDialog(dialog.Key)
						} else if ele.OnClick == nil && ele.OnHold == nil {
							ele.OnClick = Test(fmt.Sprintf("pressed button %s", ele.Key))
						}
					}
				}
			} else if ele.ElementType == ui.TextElement {
				switch ele.Key {
				default:
					if strings.Contains(ele.Key, "any_button") {
						ele.Text.Hide()
					}
				}
			} else if ele.ElementType == ui.ContainerElement {
				if ele.Key == "play_main_tab" {
					ele.Border.Style = ui.ThinBorderWhite
				} else if ele.Key == "play_custom_tab" {
					ele.Border.Style = ui.ThinBorderBlue
				}
				switch ele.Key {
				case "play_main_tab", "play_custom_tab":
					for _, ce := range ele.Elements {
						switch ce.Key {
						case "main_tab_text_shadow", "custom_tab_text_shadow":
							ce.Text.SetColor(pixel.ToRGBA(constants.ColorBlue))
						}
					}
					ele.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dialog.ViewPort, func(hvc *data.HoverClick) {
						if dialog.Open && dialog.Active {
							click := hvc.Input.Get("click")
							if hvc.Hover && click.JustPressed() {
								for _, dle := range dialog.Elements {
									switch dle.Key {
									case "play_main_tab", "play_custom_tab":
										if dle.ElementType == ui.ContainerElement {
											if ele.Key == dle.Key {
												dle.Border.Style = ui.ThinBorderWhite
											} else {
												dle.Border.Style = ui.ThinBorderBlue
											}
											for _, txt1 := range dle.Elements {
												if (ele.Key == "play_main_tab" && txt1.Key == "main_tab_text_shadow") ||
													(ele.Key == "play_custom_tab" && txt1.Key == "custom_tab_text_shadow") {
													txt1.Text.Show()
												} else if txt1.Key == "main_tab_text_shadow" || txt1.Key == "custom_tab_text_shadow" {
													txt1.Text.Hide()
												}
											}
										}
									case "main_tab_display":
										dle.Object.Hidden = ele.Key == "play_custom_tab"
									case "custom_tab_display":
										dle.Object.Hidden = ele.Key == "play_main_tab"
									}
								}
								click.Consume()
							}
						}
					}))
				case "main_tab_display":
					ele.Object.Hidden = true
				case "custom_tab_display":
					for _, ctEle := range ele.Elements {
						switch ctEle.Key {
						case "custom_new_game":
							ctEle.OnClick = StartCustomNew
						case "custom_continue_game":
							ctEle.OnClick = StartCustomContinue
						}
					}
				default:
					if strings.Contains(ele.Key, "selected_") {
						ele.Object.Hidden = true
					}
				}
			}
		}
	}
}

func AddPlayersDialog(win *pixelgl.Window) bool {
	if data.MenuInput.Get("escape").JustPressed() {
		data.MenuInput.Get("escape").Consume()
		if len(data.Players) < 2 {
			return true
		} else {
			var tCntKey, tTxtKey, nTxtKey string
			switch len(data.Players) {
			case 2:
				tCntKey = "selected_p2_cnt"
				tTxtKey = "any_button_p2"
				nTxtKey = "any_button_p3"
			case 3:
				tCntKey = "selected_p3_cnt"
				tTxtKey = "any_button_p3"
				nTxtKey = "any_button_p4"
			case 4:
				tCntKey = "selected_p4_cnt"
				tTxtKey = "any_button_p4"
			}
			data.Players = data.Players[:len(data.Players)-1]
			for _, ele := range ui.Dialogs[constants.DialogAddPlayers].Elements {
				if ele.Key == tCntKey {
					ele.Object.Hidden = true
				} else if ele.Key == nTxtKey {
					ele.Text.Hide()
				} else if ele.Key == tTxtKey {
					ele.Text.Show()
				}
			}
		}
	}
	if len(data.Players) < constants.MaxPlayers {
		var playerFound bool
		nextPlayer := data.Player{
			PlayerNum: len(data.Players),
		}
		var nCntKey, nTxtKey, pTxtKey string
		switch len(data.Players) {
		case 1:
			pTxtKey = "any_button_p2"
			nCntKey = "selected_p2_cnt"
			nTxtKey = "any_button_p3"
		case 2:
			pTxtKey = "any_button_p3"
			nCntKey = "selected_p3_cnt"
			nTxtKey = "any_button_p4"
		case 3:
			pTxtKey = "any_button_p4"
			nCntKey = "selected_p4_cnt"
		}
		joysticks := pxginput.GetAllGamepads(win)
		for _, js := range joysticks {
			pressed := pxginput.GetAllJustPressedGamepad(win, js)
			if len(pressed) > 0 {
				playerFound = true
				nextPlayer.Gamepad = js
				break
			}
		}
		if !playerFound && data.MenuInput.Get("space").JustPressed() {
			playerFound = true
			nextPlayer.Keyboard = true
		}
		if playerFound {
			data.Players = append(data.Players, nextPlayer)
			for _, ele := range ui.Dialogs[constants.DialogAddPlayers].Elements {
				if ele.Key == nCntKey {
					ele.Object.Hidden = false
				} else if ele.Key == nTxtKey {
					ele.Text.Show()
				} else if ele.Key == pTxtKey {
					ele.Text.Hide()
				}
			}
		}
	}
	return false
}
