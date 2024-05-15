package load

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"strings"
)

func MainDialogs(win *pixelgl.Window) {
	ui.NewDialog(mainMenuConstructor)
	ui.NewDialog(addPlayersConstructor)
	ui.NewDialog(playLocalConstructor)
	customizeMainDialogs(win)
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
				case "play_new_local":
					ele.OnClick = StartPlayNew
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
						ele.Text.Hidden = true
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
								for _, ct1 := range dialog.Elements {
									if ct1.ElementType == ui.ContainerElement {
										if ele.Key == ct1.Key {
											ct1.Border.Style = ui.ThinBorderWhite
										} else {
											ct1.Border.Style = ui.ThinBorderBlue
										}
										for _, txt1 := range ct1.Elements {
											if (ele.Key == "play_main_tab" && txt1.Key == "main_tab_text_shadow") ||
												(ele.Key == "play_custom_tab" && txt1.Key == "custom_tab_text_shadow") {
												txt1.Text.Hidden = false
											} else if txt1.Key == "main_tab_text_shadow" || txt1.Key == "custom_tab_text_shadow" {
												txt1.Text.Hidden = true
											}
										}
									}
								}
								click.Consume()
							}
						}
					}))
				default:
					if strings.Contains(ele.Key, "selected_") {
						ele.Object.Hidden = true
					}
				}
			}
		}
	}
}
