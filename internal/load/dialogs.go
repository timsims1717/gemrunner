package load

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems"
	"gemrunner/pkg/img"
	"gemrunner/pkg/timing"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"strings"
)

func Dialogs(win *pixelgl.Window) {
	data.NewDialog(openPuzzleConstructor)
	data.NewDialog(changeNameConstructor)
	data.NewDialog(noPlayersInPuzzle)
	data.NewDialog(crackedTileOptionsConstructor)
	editorPanels()
	data.NewDialog(editorOptBottomConstructor)
	data.NewDialog(editorOptRightConstructor)
	customizeDialogs(win)
}

func customizeDialogs(win *pixelgl.Window) {
	for key := range data.Dialogs {
		dialog := data.Dialogs[key]
		b := 0
		for _, e := range dialog.Elements {
			if btn, okB := e.(*data.Button); okB {
				switch btn.Key {
				case "open_puzzle":
					btn.OnClick = OpenPuzzle
				case "new_btn":
					btn.OnClick = NewPuzzle
				case "open_btn":
					btn.OnClick = OpenOpenPuzzleDialog
				case "quit_btn":
					btn.OnClick = QuitEditor(win)
				case "save_btn":
					btn.OnClick = SavePuzzle
				case "world_btn":
					btn.OnClick = systems.ChangeWorldToNext
				case "name_btn":
					btn.OnClick = OpenDialog("change_name")
				case "test_btn":
					btn.OnClick = TestPuzzle
				case "check_puzzle_name":
					btn.OnClick = ChangeName
				case "check_cracked_tile":
					btn.OnClick = ChangeCrackTileOptions
				case "check_no_players":
					btn.OnClick = CloseDialog(key)
				default:
					switch key {
					case "editor_panel_top", "editor_panel_left":
						btn.OnClick = EditorMode(data.ModeFromSprString(btn.Sprite.Key), btn, dialog)
					default:
						if strings.Contains(btn.Key, "cancel") {
							btn.OnClick = CloseDialog(key)
						} else if btn.OnClick == nil && btn.OnHeld == nil {
							btn.OnClick = Test(fmt.Sprintf("pressed button %s", btn.Key))
						}
					}
				}
			} else if spr, okS := e.(*data.SprElement); okS {
				switch spr.Key {
				case "block_select":
					beBG := img.NewSprite("editor_tile_bg", constants.UIBatch)
					if dialog.Key == "editor_panel_top" {
						beBG = nil
					}
					beFG := img.NewSprite(data.Block(data.BlockTurf).String(), constants.TileBatch)
					beEx := img.NewSprite("", constants.TileBatch)
					spr.Entity.AddComponent(myecs.Drawable, []*img.Sprite{beBG, beFG, beEx})
					spr.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dialog.ViewPort, func(hvc *data.HoverClick) {
						if data.Editor != nil && dialog.Open && !data.DialogStackOpen {
							beFG.Key = data.Editor.CurrBlock.String()
							switch data.Editor.CurrBlock {
							case data.BlockFall:
								beEx.Key = constants.TileFall
							default:
								beEx.Key = ""
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
					if b < data.BlockEmpty {
						bId := data.Block(b)
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
						obj := spr.Object
						spr.Entity.AddComponent(myecs.Drawable, sprs)
						spr.Entity.AddComponent(myecs.Block, bId)
						spr.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dialog.ViewPort, func(hvc *data.HoverClick) {
							if data.Editor != nil && dialog.Open && !data.DialogStackOpen {
								sprS.Key = bId.String()
								click := hvc.Input.Get("click")
								if hvc.Hover && data.Editor.SelectVis {
									wo := dialog.Elements[len(dialog.Elements)-1]
									if outline, ok := wo.(*data.SprElement); ok {
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
											data.CurrPuzzle.Update = true
										}
										click.Consume()
									}
								}
							}
						}))
					}
					b++
				}
			}
		}
		switch key {
		case "open_puzzle":
			dialog.OnOpen = OnOpenPuzzleDialog
		case "change_name":
			dialog.OnOpen = OnChangeNameDialog
		case "cracked_tile_options":
			dialog.OnOpen = OnCrackTileOptions
		}
	}
}

func editorPanels() {
	data.NewDialog(editorPanelTopConstructor)
	data.NewDialog(editorPanelLeftConstructor)
	blockSelectConstructor := &data.DialogConstructor{
		Key:      "block_select",
		Width:    constants.BlockSelectWidth,
		Height:   constants.BlockSelectHeight,
		Pos:      pixel.V(0, 0),
		NoBorder: true,
	}
	w := int(blockSelectConstructor.Width)
	h := int(blockSelectConstructor.Height)
	size := w * h
	b := 0
	for ; b < size; b++ {
		blockSelectConstructor.Elements = append(blockSelectConstructor.Elements, data.ElementConstructor{
			Key:      "block_select_tile",
			SprKey:   "black_square",
			Position: data.BlockSelectPlacement(b, w, h),
			Element:  data.SpriteElement,
		})
	}
	blockSelectConstructor.Elements = append(blockSelectConstructor.Elements, data.ElementConstructor{
		Key:      "white_outline",
		SprKey:   "white_outline",
		Position: data.BlockSelectPlacement(0, w, h),
		Element:  data.SpriteElement,
	})
	data.NewDialog(blockSelectConstructor)
	editorPanelLeft := data.Dialogs["editor_panel_left"]
	editorPanelLeft.ViewPort.Canvas.SetUniform("uRedPrimary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uBluePrimary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uRedSecondary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetFragmentShader(data.PuzzleShader)
	editorPanelTop := data.Dialogs["editor_panel_top"]
	editorPanelTop.ViewPort.Canvas.SetUniform("uRedPrimary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uBluePrimary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uRedSecondary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetFragmentShader(data.PuzzleShader)
	blockSelect := data.Dialogs["block_select"]
	blockSelect.ViewPort.Canvas.SetUniform("uRedPrimary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uBluePrimary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uRedSecondary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(0))
	blockSelect.ViewPort.Canvas.SetFragmentShader(data.PuzzleShader)
}
