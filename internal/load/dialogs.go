package load

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

func Dialogs() {
	openPuzzleConstructor := &data.DialogConstructor{
		Key:    "open_puzzle",
		Width:  7,
		Height: 7,
		Elements: []data.ElementConstructor{
			{
				SprKey:      "cancel_btn_big",
				ClickSprKey: "cancel_btn_click_big",
				HelpText:    "Cancel",
				Position:    pixel.V(32, -32),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "open_btn_big",
				ClickSprKey: "open_btn_click_big",
				HelpText:    "Open",
				Position:    pixel.V(0, -32),
				Element:     data.ButtonElement,
			},
		},
	}
	data.NewDialog(openPuzzleConstructor)
	editorPanels()
	editorOptPanels()
	customizeDialogs()
}

func customizeDialogs() {
	for key := range data.Dialogs {
		dialog := data.Dialogs[key]
		for _, btn := range dialog.Buttons {
			switch btn.Key {
			case "cancel_btn":
				btn.OnClick = CancelDialog(key)
			default:
				switch key {
				case "editor_panel_top", "editor_panel_left":
					btn.OnClick = EditorMode(data.ModeFromSprString(btn.Sprite.Key), btn, dialog)
				case "editor_options_bot", "editor_options_right":
					btn.OnClick = Test
				}
			}
		}
		for i, spr := range dialog.Sprites {
			switch spr.Key {
			case "editor_tile_bg":
				beBG := img.NewSprite("editor_tile_bg", constants.UIBatch)
				if dialog.Key == "editor_panel_top" {
					beBG = nil
				}
				beFG := img.NewSprite(data.Block(data.Turf).String(), constants.BGBatch)
				beEx := img.NewSprite("", constants.BGBatch)
				spr.Entity.AddComponent(myecs.Drawable, []*img.Sprite{beBG, beFG, beEx})
				spr.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.EditorInput, dialog.ViewPort, func(hvc *data.HoverClick) {
					if data.Editor != nil && dialog.Open && !data.DialogStackOpen {
						beFG.Key = data.Editor.CurrBlock.String()
						switch data.Editor.CurrBlock {
						case data.Fall:
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
			case "black_square":
				if i < data.Empty {
					bId := data.Block(i)
					//sprB := img.NewSprite("black_square_big", constants.UIBatch)
					sprS := img.NewSprite(bId.String(), constants.BGBatch)
					sprs := []*img.Sprite{sprS}
					if i == data.Fall {
						sprs = append(sprs, img.NewSprite(constants.TileFall, constants.BGBatch))
					}
					obj := spr.Object
					spr.Entity.AddComponent(myecs.Drawable, sprs)
					spr.Entity.AddComponent(myecs.Block, bId)
					spr.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.EditorInput, dialog.ViewPort, func(hvc *data.HoverClick) {
						if data.Editor != nil && dialog.Open && !data.DialogStackOpen {
							sprS.Key = bId.String()
							click := hvc.Input.Get("click")
							if hvc.Hover && data.Editor.SelectVis {
								dialog.Sprites[len(dialog.Sprites)-1].Object.Pos = obj.Pos
								if click.JustPressed() || click.JustReleased() {
									data.Editor.CurrBlock = bId
									data.Editor.SelectVis = false
									data.Editor.SelectQuick = false
									data.Editor.SelectTimer = nil
									switch data.Editor.Mode {
									case data.Brush, data.Line, data.Square, data.Fill:
									default:
										data.Editor.Mode = data.Brush
									}
									click.Consume()
								}
							}
						}
					}))
				}
			}
		}
	}
}

func editorPanels() {
	editorPanelTopConstructor := &data.DialogConstructor{
		Key:    "editor_panel_top",
		Width:  16,
		Height: 1,
		Pos:    pixel.V(0, 400),
		Elements: []data.ElementConstructor{
			{
				SprKey:      "brush_btn",
				ClickSprKey: "brush_btn_click",
				HelpText:    "Brush Tool (B)",
				Position:    pixel.V(-7.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "line_btn",
				ClickSprKey: "line_btn_click",
				HelpText:    "Line Tool (L)",
				Position:    pixel.V(-6.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "square_btn",
				ClickSprKey: "square_btn_click",
				HelpText:    "Square Tool (H)",
				Position:    pixel.V(-5.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "fill_btn",
				ClickSprKey: "fill_btn_click",
				HelpText:    "Fill Tool (G)",
				Position:    pixel.V(-4.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "erase_btn",
				ClickSprKey: "erase_btn_click",
				HelpText:    "Erase Tool (E)",
				Position:    pixel.V(-3.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "eyedrop_btn",
				ClickSprKey: "eyedrop_btn_click",
				HelpText:    "Eyedrop Tool (Y)",
				Position:    pixel.V(-2.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "wrench_btn",
				ClickSprKey: "wrench_btn_click",
				HelpText:    "Wrench Tool (P)",
				Position:    pixel.V(-1.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "select_btn",
				ClickSprKey: "select_btn_click",
				HelpText:    "Select Tool (M)",
				Position:    pixel.V(-0.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "cut_btn",
				ClickSprKey: "cut_btn_click",
				HelpText:    "Cut Selection (Ctrl+X)",
				Position:    pixel.V(0.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "copy_btn",
				ClickSprKey: "copy_btn_click",
				HelpText:    "Copy Selection (Ctrl+C)",
				Position:    pixel.V(1.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "paste_btn",
				ClickSprKey: "paste_btn_click",
				HelpText:    "Paste Selection (Ctrl+V)",
				Position:    pixel.V(2.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "fliph_btn",
				ClickSprKey: "fliph_btn_click",
				HelpText:    "Flip Selection Horizontal (U)",
				Position:    pixel.V(3.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "flipv_btn",
				ClickSprKey: "flipv_btn_click",
				HelpText:    "flipv Selection Vertical (K)",
				Position:    pixel.V(4.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "undo_btn",
				ClickSprKey: "undo_btn_click",
				HelpText:    "Undo (Ctrl+Z)",
				Position:    pixel.V(5.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "redo_btn",
				ClickSprKey: "redo_btn_click",
				HelpText:    "Redo (Ctrl+Shift+Z)",
				Position:    pixel.V(6.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				SprKey:   "editor_tile_bg",
				Position: pixel.V(7.5*world.TileSize, 0),
				Element:  data.SpriteElement,
			},
		},
	}
	data.NewDialog(editorPanelTopConstructor)
	editorPanelLeftConstructor := &data.DialogConstructor{
		Key:    "editor_panel_left",
		Width:  2,
		Height: 10,
		Pos:    pixel.V(-692, 0),
		Elements: []data.ElementConstructor{
			{
				SprKey:      "brush_btn",
				ClickSprKey: "brush_btn_click",
				HelpText:    "Brush Tool (B)",
				Position:    pixel.V(-0.5*world.TileSize, 2.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "line_btn",
				ClickSprKey: "line_btn_click",
				HelpText:    "Line Tool (L)",
				Position:    pixel.V(0.5*world.TileSize, 2.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "square_btn",
				ClickSprKey: "square_btn_click",
				HelpText:    "Square Tool (H)",
				Position:    pixel.V(-0.5*world.TileSize, 1.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "fill_btn",
				ClickSprKey: "fill_btn_click",
				HelpText:    "Fill Tool (G)",
				Position:    pixel.V(0.5*world.TileSize, 1.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "erase_btn",
				ClickSprKey: "erase_btn_click",
				HelpText:    "Erase Tool (E)",
				Position:    pixel.V(-0.5*world.TileSize, 0.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "eyedrop_btn",
				ClickSprKey: "eyedrop_btn_click",
				HelpText:    "Eyedrop Tool (Y)",
				Position:    pixel.V(0.5*world.TileSize, 0.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "wrench_btn",
				ClickSprKey: "wrench_btn_click",
				HelpText:    "Wrench Tool (P)",
				Position:    pixel.V(-0.5*world.TileSize, -0.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "select_btn",
				ClickSprKey: "select_btn_click",
				HelpText:    "Select Tool (M)",
				Position:    pixel.V(-0.5*world.TileSize, -1.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "cut_btn",
				ClickSprKey: "cut_btn_click",
				HelpText:    "Cut Selection (Ctrl+X)",
				Position:    pixel.V(0.5*world.TileSize, -1.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "copy_btn",
				ClickSprKey: "copy_btn_click",
				HelpText:    "Copy Selection (Ctrl+C)",
				Position:    pixel.V(-0.5*world.TileSize, -2.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "paste_btn",
				ClickSprKey: "paste_btn_click",
				HelpText:    "Paste Selection (Ctrl+V)",
				Position:    pixel.V(0.5*world.TileSize, -2.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "fliph_btn",
				ClickSprKey: "fliph_btn_click",
				HelpText:    "Flip Selection Horizontal (U)",
				Position:    pixel.V(-0.5*world.TileSize, -3.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "flipv_btn",
				ClickSprKey: "flipv_btn_click",
				HelpText:    "flipv Selection Vertical (K)",
				Position:    pixel.V(0.5*world.TileSize, -3.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "undo_btn",
				ClickSprKey: "undo_btn_click",
				HelpText:    "Undo (Ctrl+Z)",
				Position:    pixel.V(-0.5*world.TileSize, -4.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				SprKey:      "redo_btn",
				ClickSprKey: "redo_btn_click",
				HelpText:    "Redo (Ctrl+Shift+Z)",
				Position:    pixel.V(0.5*world.TileSize, -4.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				SprKey:   "editor_tile_bg",
				Position: pixel.V(0, 4*world.TileSize),
				Element:  data.SpriteElement,
			},
		},
	}
	data.NewDialog(editorPanelLeftConstructor)
	blockSelectConstructor := &data.DialogConstructor{
		Key:      "block_select",
		Width:    6,
		Height:   3,
		Pos:      pixel.V(0, 0),
		NoBorder: true,
	}
	w := int(blockSelectConstructor.Width)
	h := int(blockSelectConstructor.Height)
	size := w * h
	b := 0
	for ; b < size; b++ {
		blockSelectConstructor.Elements = append(blockSelectConstructor.Elements, data.ElementConstructor{
			SprKey:   "black_square",
			Position: data.BlockSelectPlacement(b, w, h),
			Element:  data.SpriteElement,
		})
	}
	blockSelectConstructor.Elements = append(blockSelectConstructor.Elements, data.ElementConstructor{
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

func editorOptPanels() {
	editorOptBottomConstructor := &data.DialogConstructor{
		Key:    "editor_options_bot",
		Width:  8,
		Height: 1,
		Pos:    pixel.V(0, -400),
		Elements: []data.ElementConstructor{
			{
				SprKey:      "quit_btn",
				ClickSprKey: "quit_btn_click",
				HelpText:    "Quit (Ctrl+Q)",
				Position:    pixel.V(-3.5*world.TileSize, 0),
			},
			{
				SprKey:      "new_btn",
				ClickSprKey: "new_btn_click",
				HelpText:    "New Puzzle Group (Ctrl+N)",
				Position:    pixel.V(-2.5*world.TileSize, 0),
			},
			{
				SprKey:      "save_btn",
				ClickSprKey: "save_btn_click",
				HelpText:    "Save Puzzle Group (Ctrl+S)",
				Position:    pixel.V(-1.5*world.TileSize, 0),
			},
			{
				SprKey:      "open_btn",
				ClickSprKey: "open_btn_click",
				HelpText:    "Open Puzzle Group (Ctrl+O)",
				Position:    pixel.V(-0.5*world.TileSize, 0),
			},
			{
				SprKey:      "prev_btn",
				ClickSprKey: "prev_btn_click",
				HelpText:    "Previous Puzzle",
				Position:    pixel.V(0.5*world.TileSize, 0),
			},
			{
				SprKey:      "next_btn",
				ClickSprKey: "next_btn_click",
				HelpText:    "Next Puzzle",
				Position:    pixel.V(1.5*world.TileSize, 0),
			},
			{
				SprKey:      "world_btn",
				ClickSprKey: "world_btn_click",
				HelpText:    "Change World (Tab)",
				Position:    pixel.V(2.5*world.TileSize, 0),
			},
			{
				SprKey:      "test_btn",
				ClickSprKey: "test_btn_click",
				HelpText:    "Test Puzzle",
				Position:    pixel.V(3.5*world.TileSize, 0),
			},
		},
	}
	data.NewDialog(editorOptBottomConstructor)
	editorOptRightConstructor := &data.DialogConstructor{
		Key:    "editor_options_right",
		Width:  1,
		Height: 8,
		Pos:    pixel.V(670, 0),
		Elements: []data.ElementConstructor{
			{
				SprKey:      "quit_btn",
				ClickSprKey: "quit_btn_click",
				HelpText:    "Quit (Ctrl+Q)",
				Position:    pixel.V(0, 3.5*world.TileSize),
			},
			{
				SprKey:      "new_btn",
				ClickSprKey: "new_btn_click",
				HelpText:    "New Puzzle Group (Ctrl+N)",
				Position:    pixel.V(0, 2.5*world.TileSize),
			},
			{
				SprKey:      "save_btn",
				ClickSprKey: "save_btn_click",
				HelpText:    "Save Puzzle Group (Ctrl+S)",
				Position:    pixel.V(0, 1.5*world.TileSize),
			},
			{
				SprKey:      "open_btn",
				ClickSprKey: "open_btn_click",
				HelpText:    "Open Puzzle Group (Ctrl+O)",
				Position:    pixel.V(0, 0.5*world.TileSize),
			},
			{
				SprKey:      "prev_btn",
				ClickSprKey: "prev_btn_click",
				HelpText:    "Previous Puzzle",
				Position:    pixel.V(0, -0.5*world.TileSize),
			},
			{
				SprKey:      "next_btn",
				ClickSprKey: "next_btn_click",
				HelpText:    "Next Puzzle",
				Position:    pixel.V(0, -1.5*world.TileSize),
			},
			{
				SprKey:      "world_btn",
				ClickSprKey: "world_btn_click",
				HelpText:    "Change World (Tab)",
				Position:    pixel.V(0, -2.5*world.TileSize),
			},
			{
				SprKey:      "test_btn",
				ClickSprKey: "test_btn_click",
				HelpText:    "Test Puzzle",
				Position:    pixel.V(0, -3.5*world.TileSize),
			},
		},
	}
	data.NewDialog(editorOptRightConstructor)
}
