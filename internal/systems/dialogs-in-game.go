package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/load"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/state"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"strings"
)

func InGameDialogs(win *pixelgl.Window) {
	ui.NewDialog(load.PauseConstructor)
	ui.NewDialog(load.Player1InvConstructor)
	ui.NewDialog(load.Player2InvConstructor)
	ui.NewDialog(load.Player3InvConstructor)
	ui.NewDialog(load.Player4InvConstructor)
	ui.NewDialog(load.PuzzleTitleConstructor)
	customizeInGameDialogs(win)
}

func DisposeInGameDialogs() {
	for k, d := range ui.Dialogs {
		switch k {
		case constants.DialogPauseMenu,
			constants.DialogPlayer1Inv,
			constants.DialogPlayer2Inv,
			constants.DialogPlayer3Inv,
			constants.DialogPlayer4Inv,
			constants.DialogPuzzleTitle:
			ui.Dispose(d)
		}
	}
}

func customizeInGameDialogs(win *pixelgl.Window) {
	for key := range ui.Dialogs {
		dialog := ui.Dialogs[key]
		for _, e := range dialog.Elements {
			ele := e
			if ele.ElementType == ui.ButtonElement {
				switch ele.Key {
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
			} else if ele.ElementType == ui.ContainerElement {
				switch ele.Key {
				case "pause_resume_ct", "pause_restart_ct", "pause_options_ct", "pause_quit_mm_ct", "pause_quit_full_ct":
					ele.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dialog.ViewPort, func(hvc *data.HoverClick) {
						if dialog.Open && dialog.Active && !dialog.Lock {
							click := hvc.Input.Get("click")
							if hvc.Hover && click.JustPressed() {
								dialog.Click = true
							}
							if hvc.Hover && click.Pressed() && dialog.Click {
								// change border to be reverse, change text to blue
								ele.Border.Style = ui.ThinBorderReverse
								ele.Elements[0].Text.SetColor(pixel.ToRGBA(constants.ColorBlue))
							} else {
								if hvc.Hover && click.JustReleased() && dialog.Click {
									dialog.Click = false
									if ele.OnClick != nil {
										if ele.Delay > 0. {
											dialog.Lock = true
											entity := myecs.Manager.NewEntity()
											entity.AddComponent(myecs.Update, data.NewTimerFunc(func() bool {
												hvc.Input.Get("click").Consume()
												hvc.Input.Get("rClick").Consume()
												ele.OnClick()
												dialog.Lock = false
												myecs.Manager.DisposeEntity(entity)
												return false
											}, ele.Delay))
										} else {
											hvc.Input.Get("click").Consume()
											hvc.Input.Get("rClick").Consume()
											ele.OnClick()
										}
									}
								} else if !click.Pressed() && !click.JustReleased() && dialog.Click {
									dialog.Click = false
									// change border to be normal, change text to white
									ele.Border.Style = ui.ThinBorder
									ele.Elements[0].Text.SetColor(pixel.ToRGBA(constants.ColorWhite))
								} else {
									// change border to be normal, change text to white
									ele.Border.Style = ui.ThinBorder
									ele.Elements[0].Text.SetColor(pixel.ToRGBA(constants.ColorWhite))
								}
							}
						}
					}))
					switch ele.Key {
					case "pause_resume_ct":
						ele.OnClick = func() {
							ui.CloseDialog(constants.DialogPauseMenu)
						}
					case "pause_restart_ct":
						ele.OnClick = func() {
							Restart()
							ui.CloseDialog(constants.DialogPauseMenu)
						}
					case "pause_options_ct":
						ele.OnClick = func() {
							fmt.Println("pause options pressed")
						}
					case "pause_quit_mm_ct":
						ele.OnClick = func() {
							// todo: save
							state.SwitchState(constants.MainMenuKey)
						}
					case "pause_quit_full_ct":
						ele.OnClick = func() {
							// todo: save
							win.SetClosed(true)
						}
					}
				}
			} else if ele.ElementType == ui.SpriteElement {
				if ele.Key == "player_inv_item" {
					ele.Object.Hidden = true
				}
			}
		}
		switch dialog.Key {
		case constants.DialogPlayer1Inv,
			constants.DialogPlayer2Inv,
			constants.DialogPlayer3Inv,
			constants.DialogPlayer4Inv,
			constants.DialogPuzzleTitle:
			dialog.Border.Style = ui.ThinBorder
			dialog.Border.Rect = pixel.R(0, 0, float64(dialog.Border.Width)*world.TileSize, float64(dialog.Border.Height)*world.TileSize)
		}
	}
}

func SetPuzzleTitle() {
	if data.CurrLevel == nil {
		return
	}
	title := data.CurrLevel.Metadata.Name
	dlg := ui.Dialogs[constants.DialogPuzzleTitle]
	txt := dlg.Get("puzzle_title")
	width := txt.Text.Text.BoundsOf(title).W() * txt.Text.RelativeSize
	dlgWidth := world.TileSize * 6
	for width > dlgWidth {
		dlgWidth += world.TileSize
	}
	dlg.Border.Rect = pixel.R(0, 0, dlgWidth, float64(dlg.Border.Height)*world.TileSize)
	dlg.BorderVP.SetRect(pixel.R(0, 0, dlgWidth+1, float64(dlg.Border.Rect.H())+1))
	dlg.ViewPort.SetRect(pixel.R(0, 0, dlgWidth, dlg.Border.Rect.H()))
	for _, ele := range dlg.Elements {
		if ele.Key == "puzzle_title_bg" {
			ele.Text.SetColor(data.CurrLevel.Metadata.PrimaryColor)
		}
		ele.Object.Pos.X = width * -0.5
		ele.Text.SetText(title)
	}
}
