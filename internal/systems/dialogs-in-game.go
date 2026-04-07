package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/load"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/state"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

func InGameDialogs(win *pixelgl.Window) {
	ui.NewDialog(load.PauseConstructor)
	ui.NewDialog(ui.DialogConstructors[constants.DialogPlayer1Inv])
	ui.NewDialog(ui.DialogConstructors[constants.DialogPlayer2Inv])
	ui.NewDialog(ui.DialogConstructors[constants.DialogPlayer3Inv])
	ui.NewDialog(ui.DialogConstructors[constants.DialogPlayer4Inv])
	ui.NewDialog(ui.DialogConstructors[constants.DialogPuzzleTitle])
	ui.NewDialog(ui.DialogConstructors[constants.DialogPuzzleTimer])
	CustomizeInGameDialogs(win)
}

func DisposeInGameDialogs() {
	for _, k := range constants.InGameDialogs {
		if dlg, ok := ui.Dialogs[k]; ok {
			ui.DisposeDialog(dlg)
		}
	}
}

func CustomizeInGameDialogs(win *pixelgl.Window) {
	for _, k := range constants.InGameDialogs {
		CustomizeInGameDialog(win, k)
	}
}

func CustomizeInGameDialog(win *pixelgl.Window, key string) {
	if dialog, ok := ui.Dialogs[key]; ok {
		switch key {
		case constants.DialogPlayer1Inv,
			constants.DialogPlayer2Inv,
			constants.DialogPlayer3Inv,
			constants.DialogPlayer4Inv:
			dialog.Border.Style = data.ThinBorder
			dialog.Border.Rect = pixel.R(0, 0, float64(dialog.Border.Width)*world.TileSize, float64(dialog.Border.Height)*world.TileSize)
			if dialog.Key != constants.DialogPuzzleTitle {
				cnt := dialog.Get("player_inv_cnt")
				if cnt != nil {
					cnt.ViewPort.Canvas.SetUniform("uRedPrimary", float32(1))
					cnt.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(1))
					cnt.ViewPort.Canvas.SetUniform("uBluePrimary", float32(1))
					cnt.ViewPort.Canvas.SetUniform("uRedSecondary", float32(1))
					cnt.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(1))
					cnt.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(1))
					cnt.ViewPort.Canvas.SetUniform("uRedDoodad", float32(1))
					cnt.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(1))
					cnt.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(1))
					cnt.ViewPort.Canvas.SetUniform("uRedGoop", float32(1))
					cnt.ViewPort.Canvas.SetUniform("uGreenGoop", float32(1))
					cnt.ViewPort.Canvas.SetUniform("uBlueGoop", float32(1))
					cnt.ViewPort.Canvas.SetUniform("uRedLiquidPrimary", float32(1))
					cnt.ViewPort.Canvas.SetUniform("uGreenLiquidPrimary", float32(1))
					cnt.ViewPort.Canvas.SetUniform("uBlueLiquidPrimary", float32(1))
					cnt.ViewPort.Canvas.SetUniform("uRedLiquidSecondary", float32(1))
					cnt.ViewPort.Canvas.SetUniform("uGreenLiquidSecondary", float32(1))
					cnt.ViewPort.Canvas.SetUniform("uBlueLiquidSecondary", float32(1))
					cnt.ViewPort.Canvas.SetFragmentShader(data.ColorShader)
				}
			}
			for _, e := range dialog.Elements {
				ele := e
				if ele.Key == "player_inv_item" {
					ele.Object.Hidden = true
				}
			}
		case constants.DialogPuzzleTitle, constants.DialogPuzzleTimer:
			dialog.Border.Style = data.ThinBorder
			dialog.Border.Rect = pixel.R(0, 0, float64(dialog.Border.Width)*world.TileSize, float64(dialog.Border.Height)*world.TileSize)
		case constants.DialogPauseMenu:
			for _, e := range dialog.Elements {
				ele := e
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
								ele.Border.Style = data.ThinBorderReverse
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
									ele.Border.Style = data.ThinBorder
									ele.Elements[0].Text.SetColor(pixel.ToRGBA(constants.ColorWhite))
								} else {
									// change border to be normal, change text to white
									ele.Border.Style = data.ThinBorder
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
			}
		}
	}
}

func SetPuzzleTitle(title string, color pixel.RGBA) {
	if title == "" {
		return
	}
	dlg := ui.Dialogs[constants.DialogPuzzleTitle]
	txt := dlg.Get("puzzle_title")
	width := txt.Text.Text.BoundsOf(title).W() * txt.Text.Scalar
	dlgWidth := world.TileSize * 6
	for width > dlgWidth {
		dlgWidth += world.TileSize
	}
	dlg.Border.Rect = pixel.R(0, 0, dlgWidth, float64(dlg.Border.Height)*world.TileSize)
	dlg.BorderVP.SetRect(pixel.R(0, 0, dlgWidth+1, float64(dlg.Border.Rect.H())+1))
	dlg.ViewPort.SetRect(pixel.R(0, 0, dlgWidth, dlg.Border.Rect.H()))
	for _, ele := range dlg.Elements {
		if ele.Key == "puzzle_title_shadow" {
			ele.Text.SetColor(color)
		}
		ele.Object.Pos.X = width * -0.5
		ele.Text.SetText(title)
	}
}

func UpdatePuzzleTimer() {
	timerText := FormatTimePlayed()
	dlg := ui.Dialogs[constants.DialogPuzzleTimer]
	txt := dlg.Get("puzzle_timer")
	txt.Text.SetText(timerText)
	txt.Object.Pos.X = txt.Text.GetWidth() * -0.5
}

// Adventure Transition

func OpenAdventureTransition() {
	key := constants.DialogAdventureTrans
	ui.NewDialog(ui.DialogConstructors[key])
	dlg := ui.Dialogs[key]
	for _, e := range dlg.Elements {
		ele := e
		switch ele.Key {
		case "adventure_title":
			ele.Text.SetText(data.CurrPuzzleSet.Metadata.Name)
		case "puzzle_set_view":
			ele1 := ele.Get("player_symbol")
			if data.CurrLevel != nil {
				grid := data.CurrLevel.Puzzle.Grid
				ele1.Object.SetPos(AdvPuzzleViewPos(grid))
			}
			//blinkTimer := timing.New(0.5)
			startMoveTimer := timing.New(constants.LevelTransSpeed)
			finishTimer := timing.New(constants.LevelTransSpeed * 4)
			nextPos := AdvPuzzleViewPos(data.CurrPuzzleSet.CurrPuzzle.Grid)
			nextIndex := data.CurrPuzzleSet.PuzzleIndex
			ele1.Entity.AddComponent(myecs.Update, data.NewFn(func() {
				//if blinkTimer.UpdateDone() {
				//	ele1.Sprite.Hide = !ele1.Sprite.Hide
				//	blinkTimer.Reset()
				//}
				if startMoveTimer != nil && startMoveTimer.UpdateDone() {
					StartLevelTransitionInDialog(ele1, ele, nextPos)
					startMoveTimer = nil
				}
				if finishTimer.UpdateDone() {
					ele1.Entity.RemoveComponent(myecs.Update)
					ui.CloseDialog(key)
					GoToLevel(nextIndex)
				}
			}))
		}
	}
	data.AdventureViewGridMap = make(map[world.Coords]data.AdvViewPzl)
	data.AdventureViewGridArr = make(map[int]world.Coords)
	AdventureViewZoomOne(key)
	UpdateDialogView(dlg)
	ui.OpenDialogInStack(key)
}
