package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/load"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/img"
	"gemrunner/pkg/timing"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"strings"
)

func EditorDialogs(win *pixelgl.Window) {
	ui.NewDialog(load.OpenPuzzleConstructor)
	ui.NewDialog(load.ChangeNameConstructor)
	ui.NewDialog(ui.DialogConstructors[constants.DialogPuzzleSettings])
	ui.NewDialog(ui.DialogConstructors[constants.DialogPuzzleSetSettings])
	ui.NewDialog(ui.DialogConstructors[constants.DialogNoPlayersInPuzzle])
	ui.NewDialog(load.AreYouSureDeleteConstructor)
	ui.NewDialog(load.UnableToSaveConstructor)
	ui.NewDialog(load.UnableToSaveConfirmConstructor)
	ui.NewDialog(load.CombineSetsConstructor)
	ui.NewDialog(load.CrackedTileOptionsConstructor)
	ui.NewDialog(load.BombOptionsConstructor)
	ui.NewDialog(ui.DialogConstructors[constants.DialogItemOptions])
	ui.NewDialog(ui.DialogConstructors[constants.DialogFloatingText])
	ui.NewDialog(ui.DialogConstructors[constants.DialogPalette])
	editorPanels()
	ui.NewDialog(load.EditorOptBottomConstructor)
	ui.NewDialog(ui.DialogConstructors[constants.DialogEditorOptionsRight])
	customizeEditorDialogs(win)
}

func DisposeEditorDialogs() {
	for k, d := range ui.Dialogs {
		switch k {
		case constants.DialogOpenPuzzle,
			constants.DialogChangeName,
			constants.DialogPuzzleSettings,
			constants.DialogPuzzleSetSettings,
			constants.DialogNoPlayersInPuzzle,
			constants.DialogAreYouSureDelete,
			constants.DialogUnableToSave,
			constants.DialogUnableToSaveConfirm,
			constants.DialogChangeWorld,
			constants.DialogCombineSets,
			constants.DialogRearrangePuzzleSet,
			constants.DialogCrackedTiles,
			constants.DialogBomb,
			constants.DialogDoors,
			constants.DialogPalette,
			constants.DialogEditorPanelLeft,
			constants.DialogEditorPanelTop,
			constants.DialogEditorOptionsRight,
			constants.DialogEditorOptionsBot,
			constants.DialogEditorBlockSelect:
			ui.DisposeDialog(d)
		}
	}
}

func customizeEditorDialogs(win *pixelgl.Window) {
	for key := range ui.Dialogs {
		CustomizeEditorDialog(key)
	}
}

func CustomizeEditorDialog(key string) {
	switch key {
	case constants.DialogFloatingText:
		customizeFloatingText()
	case constants.DialogPalette:
		customizePaletteOptions()
	case constants.DialogItemOptions:
		customizeItemOptions()
	case constants.DialogPuzzleSettings:
		customizePuzzleSettings()
	default:
		dialog := ui.Dialogs[key]
		b := 0
		for _, e := range dialog.Elements {
			ele := e
			switch ele.Key {
			case "floating_text_value":
				ele.InputType = ui.Special
			}
			if ele.ElementType == ui.ButtonElement {
				switch ele.Key {
				case "confirm":
					switch key {
					case constants.DialogCombineSets:
						ele.OnClick = OnCombinePuzzleSet
					case constants.DialogPuzzleSetSettings:
						ele.OnClick = ConfirmPuzzleSetSettings
					case constants.DialogBomb:
						ele.OnClick = ConfirmBombOptions
					case constants.DialogUnableToSave, constants.DialogUnableToSaveConfirm:
						ele.OnClick = CloseDialog(dialog.Key)
					case constants.DialogAreYouSureDelete:
						ele.OnClick = ConfirmDelete
					case constants.DialogNoPlayersInPuzzle:
						ele.OnClick = CloseDialog(dialog.Key)
					case constants.DialogCrackedTiles:
						ele.OnClick = ConfirmCrackTileOptions
					case constants.DialogChangeName:
						ele.OnClick = ChangeName
					case constants.DialogOpenPuzzle:
						ele.OnClick = OnOpenPuzzle
					}
				case "new_btn":
					ele.OnClick = NewPuzzle
				case "open_btn":
					ele.OnClick = OpenOpenPuzzleDialog
				case "combine_btn":
					ele.OnClick = OpenCombineSetsDialog
				case "rearrange_btn":
					ele.OnClick = OpenRearrangePuzzlesDialog
				case "exit_editor_btn":
					ele.OnClick = ExitEditor
				case "save_btn":
					ele.OnClick = OnSavePuzzleSet
				case "world_btn":
					ele.OnClick = OpenChangeWorldDialog
				case "name_btn":
					ele.OnClick = OpenDialog(constants.DialogChangeName)
				case "test_btn":
					ele.OnClick = TestPuzzle
				case "puzzle_settings_btn":
					ele.OnClick = OpenDialog(constants.DialogPuzzleSettings)
				case "puzzle_set_settings_btn":
					ele.OnClick = OpenPuzzleSetSettingsDialog
				case "bomb_regenerate_delay_minus":
					ele.OnClick = func() {
						ChangeNumberInput(dialog.Get("bomb_regenerate_delay_input"), -1)
					}
				case "bomb_regenerate_delay_plus":
					ele.OnClick = func() {
						ChangeNumberInput(dialog.Get("bomb_regenerate_delay_input"), 1)
					}
				case "add_btn":
					ele.OnClick = AddPuzzle
				case "prev_btn":
					ele.OnClick = PrevPuzzle
				case "next_btn":
					ele.OnClick = NextPuzzle
				case "delete_btn":
					ele.OnClick = OpenConfirmDelete
				default:
					switch dialog.Key {
					case constants.DialogEditorPanelTop, constants.DialogEditorPanelLeft:
						ele.OnClick = EditorMode(data.ModeFromSprString(ele.Sprite.Key), ele, dialog)
					case constants.DialogUnableToSaveConfirm:
						if strings.Contains(ele.Key, "cancel") {
							ele.OnClick = func() {
								ui.SetOnClick(constants.DialogUnableToSaveConfirm, "confirm", CloseDialog(dialog.Key))
								ui.CloseDialog(dialog.Key)
							}
						}
					default:
						if strings.Contains(ele.Key, "cancel") {
							ele.OnClick = CloseDialog(dialog.Key)
						} else if ele.OnClick == nil && ele.OnHold == nil {
							ele.OnClick = Test(fmt.Sprintf("pressed button %s", ele.Key))
						}
					}
				}
			} else if ele.ElementType == ui.CheckboxElement {
				switch ele.Key {
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
				}
			} else if ele.ElementType == ui.SpriteElement {
				switch ele.Key {
				case "block_select":
					beBG := img.NewSprite("editor_tile_bg", constants.UIBatch)
					if dialog.Key == "editor_panel_top" {
						beBG = nil
					}
					beFG := img.NewSprite(data.Block(data.BlockTurf).SpriteString(), constants.TileBatch)
					beEx := img.NewSprite("", constants.TileBatch)
					ele.Entity.AddComponent(myecs.Drawable, []*img.Sprite{beBG, beFG, beEx})
					ele.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dialog.ViewPort, func(hvc *data.HoverClick) {
						if data.Editor != nil && dialog.Open && !ui.DialogStackOpen {
							beFG.Key = data.Editor.CurrBlock.SpriteString()
							switch data.Editor.CurrBlock {
							case data.BlockFall:
								beEx.Key = constants.TileFall
								beEx.Offset.Y = 0
							case data.BlockPhase:
								beEx.Key = constants.TilePhase
								beEx.Offset.Y = 0
							case data.BlockCracked:
								beEx.Key = constants.TileCracked
								beEx.Offset.Y = 0
							case data.BlockClose:
								beEx.Key = constants.TileClose
								beEx.Offset.Y = 0
							case data.BlockHideout:
								beEx.Key = constants.TileHideout
								beEx.Offset.Y = 0
							case data.BlockLiquid:
								beEx.Key = constants.TileLiquidUFG
								beEx.Offset.Y = 0
							default:
								beEx.Key = ""
								beEx.Offset.Y = 0
							}
							data.Editor.Hover = hvc.Hover
							click := hvc.Input.Get("click")
							if hvc.Hover {
								if data.Editor.SelectVis {
									if click.JustPressed() {
										data.Editor.SelectVis = false
										data.Editor.SelectTimer = nil
										data.Editor.SelectQuick = false
										click.Consume()
									} else if click.JustReleased() {
										if data.Editor.SelectTimer != nil && !data.Editor.SelectTimer.Done() {
											data.Editor.SelectQuick = true
										}
									} else if !click.Pressed() && !data.Editor.SelectQuick {
										data.Editor.SelectVis = false
										data.Editor.SelectTimer = nil
									}
								} else {
									data.Editor.SelectQuick = false
									if click.JustPressed() {
										data.Editor.SelectVis = true
										data.Editor.SelectTimer = timing.New(0.2)
									}
								}
							}
						}
					}))
				case "block_select_tile":
					bId := data.BlockList[b]
					if bId != data.BlockEmpty {
						sprS := img.NewSprite(bId.SpriteString(), constants.TileBatch)
						sprs := []*img.Sprite{sprS}
						switch bId {
						case data.BlockFall:
							sprs = append(sprs, img.NewSprite(constants.TileFall, constants.TileBatch))
						case data.BlockPhase:
							sprs = append(sprs, img.NewSprite(constants.TilePhase, constants.TileBatch))
						case data.BlockCracked:
							sprs = append(sprs, img.NewSprite(constants.TileCracked, constants.TileBatch))
						case data.BlockClose:
							sprs = append(sprs, img.NewSprite(constants.TileClose, constants.TileBatch))
						case data.BlockHideout:
							sprs = append(sprs, img.NewSprite(constants.TileHideout, constants.TileBatch))
						case data.BlockLiquid:
							sprs = append(sprs, img.NewSprite(constants.TileLiquidUFG, constants.TileBatch))
						}
						obj := ele.Object
						ele.Entity.AddComponent(myecs.Drawable, sprs)
						ele.Entity.AddComponent(myecs.Block, bId)
						ele.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dialog.ViewPort, func(hvc *data.HoverClick) {
							if data.Editor != nil && dialog.Open && !ui.DialogStackOpen {
								sprS.Key = bId.SpriteString()
								click := hvc.Input.Get("click")
								if hvc.Hover && data.Editor.SelectVis {
									outline := dialog.Elements[len(dialog.Elements)-1]
									if outline.ElementType == ui.SpriteElement {
										outline.Object.Pos = obj.Pos
									}
									if click.JustPressed() || click.JustReleased() {
										data.Editor.CurrBlock = bId
										data.Editor.SelectVis = false
										data.Editor.SelectQuick = false
										data.Editor.SelectTimer = nil
										switch data.Editor.Mode {
										case data.ModeBrush, data.ModeLine, data.ModeSquare, data.ModeFill:
										default:
											data.Editor.Mode = data.ModeBrush
											data.Editor.ModeChanged = true
											data.CurrPuzzleSet.CurrPuzzle.Update = true
										}
										click.Consume()
									}
								}
							}
						}))
					}
					b++
				}
			} else if ele.ElementType == ui.TextElement {
				switch ele.Key {
				case "puzzle_number":
					ele.Text.SetText("0001")
				}
			}
		}
		switch dialog.Key {
		case constants.DialogOpenPuzzle:
			dialog.OnOpen = OnOpenPuzzleDialog
		case constants.DialogChangeName:
			dialog.OnOpen = OnChangeNameDialog
		case constants.DialogCrackedTiles:
			dialog.OnOpen = OnOpenCrackTileOptions
		case constants.DialogBomb:
			dialog.OnOpen = OnOpenBombOptions
		}
	}
}

func editorPanels() {
	ui.NewDialog(load.EditorPanelTopConstructor)
	ui.NewDialog(ui.DialogConstructors[constants.DialogEditorPanelLeft])
	blockSelectConstructor := &ui.DialogConstructor{
		Key:      constants.DialogEditorBlockSelect,
		Width:    constants.BlockSelectWidth,
		Height:   constants.BlockSelectHeight,
		Pos:      pixel.V(0, 0),
		NoBorder: true,
	}
	w := int(blockSelectConstructor.Width)
	h := int(blockSelectConstructor.Height)
	size := w * h
	for i := range data.BlockList {
		if i > size {
			break
		}
		blockSelectConstructor.Elements = append(blockSelectConstructor.Elements, ui.ElementConstructor{
			Key:         "block_select_tile",
			SprKey:      "black_square_16",
			Batch:       constants.UIBatch,
			Position:    data.BlockSelectPlacement(i, w, h),
			ElementType: ui.SpriteElement,
		})
	}
	blockSelectConstructor.Elements = append(blockSelectConstructor.Elements, ui.ElementConstructor{
		Key:         "white_outline",
		SprKey:      "white_outline",
		Batch:       constants.UIBatch,
		Position:    data.BlockSelectPlacement(0, w, h),
		ElementType: ui.SpriteElement,
	})
	ui.NewDialog(blockSelectConstructor)
	editorPanelLeft := ui.Dialogs[constants.DialogEditorPanelLeft]
	editorPanelLeft.ViewPort.Canvas.SetUniform("uRedPrimary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uBluePrimary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uRedSecondary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uRedDoodad", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uRedLiquidPrimary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uGreenLiquidPrimary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uBlueLiquidPrimary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uRedLiquidSecondary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uGreenLiquidSecondary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uBlueLiquidSecondary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetFragmentShader(data.ColorShader)
	editorPanelTop := ui.Dialogs[constants.DialogEditorPanelTop]
	editorPanelTop.ViewPort.Canvas.SetUniform("uRedPrimary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uBluePrimary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uRedSecondary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uRedDoodad", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uRedLiquidPrimary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uGreenLiquidPrimary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uBlueLiquidPrimary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uRedLiquidSecondary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uGreenLiquidSecondary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uBlueLiquidSecondary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetFragmentShader(data.ColorShader)
	blockSelect := ui.Dialogs[constants.DialogEditorBlockSelect]
	blockSelect.ViewPort.Canvas.SetUniform("uRedPrimary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uBluePrimary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uRedSecondary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uRedDoodad", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uRedLiquidPrimary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uGreenLiquidPrimary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uBlueLiquidPrimary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uRedLiquidSecondary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uGreenLiquidSecondary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uBlueLiquidSecondary", float32(0))
	blockSelect.ViewPort.Canvas.SetFragmentShader(data.ColorShader)
}
