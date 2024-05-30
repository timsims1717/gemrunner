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
	ui.NewDialog(load.PuzzleSettingsConstructor)
	ui.NewDialog(load.NoPlayersInPuzzleConstructor)
	ui.NewDialog(load.AreYouSureDeleteConstructor)
	ui.NewDialog(load.UnableToSaveConstructor)
	ui.NewDialog(load.UnableToSaveConfirmConstructor)
	ui.NewDialog(load.CombineSetsConstructor)
	ui.NewDialog(load.CrackedTileOptionsConstructor)
	ui.NewDialog(load.BombOptionsConstructor)
	ui.NewDialog(load.JetpackOptionsConstructor)
	editorPanels()
	ui.NewDialog(load.EditorOptBottomConstructor)
	ui.NewDialog(load.EditorOptRightConstructor)
	customizeEditorDialogs(win)
}

func DisposeEditorDialogs() {
	for k, d := range ui.Dialogs {
		switch k {
		case constants.DialogOpenPuzzle,
			constants.DialogChangeName,
			constants.DialogPuzzleSettings,
			constants.DialogNoPlayersInPuzzle,
			constants.DialogAreYouSureDelete,
			constants.DialogUnableToSave,
			constants.DialogUnableToSaveConfirm,
			constants.DialogChangeWorld,
			constants.DialogCombineSets,
			constants.DialogRearrangePuzzleSet,
			constants.DialogCrackedTiles,
			constants.DialogBomb,
			constants.DialogJetpack,
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
		dialog := ui.Dialogs[key]
		b := 0
		for _, e := range dialog.Elements {
			ele := e
			if ele.ElementType == ui.ButtonElement {
				switch ele.Key {
				case "open_puzzle":
					ele.OnClick = OnOpenPuzzle
				case "confirm_combine_puzzle":
					ele.OnClick = OnCombinePuzzleSet
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
				case "check_puzzle_name":
					ele.OnClick = ChangeName
				case "puzzle_settings_btn":
					ele.OnClick = OpenPuzzleSettingsDialog
				case "confirm_puzzle_settings":
					ele.OnClick = ConfirmPuzzleSettings
				case "check_cracked_tile":
					ele.OnClick = ConfirmCrackTileOptions
				case "confirm_bomb_options":
					ele.OnClick = ConfirmBombOptions
				case "bomb_regenerate_delay_minus":
					ele.OnClick = func() {
						IncOrDecBombRegen(false)
					}
				case "bomb_regenerate_delay_plus":
					ele.OnClick = func() {
						IncOrDecBombRegen(true)
					}
				case "confirm_jetpack_options":
					ele.OnClick = ConfirmJetpackOptions
				case "jetpack_regenerate_delay_minus":
					ele.OnClick = func() {
						IncOrDecJetpackRegen(false)
					}
				case "jetpack_regenerate_delay_plus":
					ele.OnClick = func() {
						IncOrDecJetpackRegen(true)
					}
				case "jetpack_timer_minus":
					ele.OnClick = func() {
						IncOrDecJetpackTimer(false)
					}
				case "jetpack_timer_plus":
					ele.OnClick = func() {
						IncOrDecJetpackTimer(true)
					}
				case "check_no_players":
					ele.OnClick = CloseDialog(dialog.Key)
				case "confirm_unable_to_save":
					ele.OnClick = CloseDialog(dialog.Key)
				case "add_btn":
					ele.OnClick = AddPuzzle
				case "prev_btn":
					ele.OnClick = PrevPuzzle
				case "next_btn":
					ele.OnClick = NextPuzzle
				case "delete_btn":
					ele.OnClick = OpenConfirmDelete
				case "confirm_delete":
					ele.OnClick = ConfirmDelete
				default:
					switch dialog.Key {
					case constants.DialogEditorPanelTop, constants.DialogEditorPanelLeft:
						ele.OnClick = EditorMode(data.ModeFromSprString(ele.Sprite.Key), ele, dialog)
					case constants.DialogUnableToSaveConfirm:
						if strings.Contains(ele.Key, "cancel") {
							ele.OnClick = func() {
								ui.SetOnClick(constants.DialogUnableToSaveConfirm, "confirm_unable_to_save", CloseDialog(dialog.Key))
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
			} else if ele.ElementType == ui.SpriteElement {
				switch ele.Key {
				case "block_select":
					beBG := img.NewSprite("editor_tile_bg", constants.UIBatch)
					if dialog.Key == "editor_panel_top" {
						beBG = nil
					}
					beFG := img.NewSprite(data.Block(data.BlockTurf).String(), constants.TileBatch)
					beEx := img.NewSprite("", constants.TileBatch)
					ele.Entity.AddComponent(myecs.Drawable, []*img.Sprite{beBG, beFG, beEx})
					ele.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dialog.ViewPort, func(hvc *data.HoverClick) {
						if data.Editor != nil && dialog.Open && !ui.DialogStackOpen {
							beFG.Key = data.Editor.CurrBlock.String()
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
						sprS := img.NewSprite(bId.String(), constants.TileBatch)
						sprs := []*img.Sprite{sprS}
						switch b {
						case data.BlockFall:
							sprs = append(sprs, img.NewSprite(constants.TileFall, constants.TileBatch))
						case data.BlockPhase:
							sprs = append(sprs, img.NewSprite(constants.TilePhase, constants.TileBatch))
						case data.BlockCracked:
							sprs = append(sprs, img.NewSprite(constants.TileCracked, constants.TileBatch))
						}
						obj := ele.Object
						ele.Entity.AddComponent(myecs.Drawable, sprs)
						ele.Entity.AddComponent(myecs.Block, bId)
						ele.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dialog.ViewPort, func(hvc *data.HoverClick) {
							if data.Editor != nil && dialog.Open && !ui.DialogStackOpen {
								sprS.Key = bId.String()
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
										case data.Brush, data.Line, data.Square, data.Fill:
										default:
											data.Editor.Mode = data.Brush
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
		case constants.DialogJetpack:
			dialog.OnOpen = OnOpenJetpackOptions
		}
	}
}

func editorPanels() {
	ui.NewDialog(load.EditorPanelTopConstructor)
	ui.NewDialog(load.EditorPanelLeftConstructor)
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
			SprKey:      "black_square",
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
	editorPanelLeft.ViewPort.Canvas.SetFragmentShader(data.PuzzleShader)
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
	editorPanelTop.ViewPort.Canvas.SetFragmentShader(data.PuzzleShader)
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
	blockSelect.ViewPort.Canvas.SetFragmentShader(data.PuzzleShader)
}
