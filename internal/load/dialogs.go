package load

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems"
	"gemrunner/pkg/img"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

func Dialogs(win *pixelgl.Window) {
	openPuzzleConstructor := &data.DialogConstructor{
		Key:    "open_puzzle",
		Width:  11,
		Height: 10,
		Elements: []data.ElementConstructor{
			{
				Key:      "open_title",
				Text:     "Open Puzzle Group",
				Position: pixel.V(-80, 72),
				Element:  data.TextElement,
			},
			{
				Key:         "cancel_open_puzzle",
				SprKey:      "cancel_btn_big",
				ClickSprKey: "cancel_btn_click_big",
				HelpText:    "Cancel",
				Position:    pixel.V(72, -64),
				Element:     data.ButtonElement,
			},
			{
				Key:         "open_puzzle",
				SprKey:      "open_btn_big",
				ClickSprKey: "open_btn_click_big",
				HelpText:    "Open",
				Position:    pixel.V(52, -64),
				Element:     data.ButtonElement,
			},
			{
				Key:      "puzzle_list",
				HelpText: "The list of puzzles.",
				Element:  data.ScrollElement,
				Position: pixel.V(0, 8),
				Width:    10,
				Height:   7,
			},
		},
	}
	data.NewDialog(openPuzzleConstructor)
	// change_name
	changeNameConstructor := &data.DialogConstructor{
		Key:    "change_name",
		Width:  12,
		Height: 4,
		Elements: []data.ElementConstructor{
			{
				Key:      "change_name_title",
				Text:     "Change Puzzle Name",
				Position: pixel.V(-92, 24),
				Element:  data.TextElement,
			},
			{
				Key:         "cancel_puzzle_name",
				SprKey:      "cancel_btn_big",
				ClickSprKey: "cancel_btn_click_big",
				HelpText:    "Cancel",
				Position:    pixel.V(84, -20),
				Element:     data.ButtonElement,
			},
			{
				Key:         "check_puzzle_name",
				SprKey:      "check_btn_big",
				ClickSprKey: "check_btn_click_big",
				HelpText:    "Confirm",
				Position:    pixel.V(64, -20),
				Element:     data.ButtonElement,
			},
			{
				Key:      "puzzle_name",
				Text:     "Untitled",
				HelpText: "Enter the name of the puzzle here.",
				Element:  data.InputElement,
				Position: pixel.V(0, 4),
				Width:    11,
				Height:   1,
			},
		},
	}
	data.NewDialog(changeNameConstructor)
	editorPanels()
	editorOptPanels()
	customizeDialogs(win)
}

func customizeDialogs(win *pixelgl.Window) {
	for key := range data.Dialogs {
		dialog := data.Dialogs[key]
		b := 0
		for _, e := range dialog.Elements {
			if btn, okB := e.(*data.Button); okB {
				switch btn.Key {
				case "cancel_open_puzzle", "cancel_puzzle_name":
					btn.OnClick = CloseDialog(key)
				case "open_puzzle":
					btn.OnClick = OpenPuzzle
				case "new_btn":
					btn.OnClick = NewPuzzle
				case "open_btn":
					btn.OnClick = OpenDialog("open_puzzle")
				case "quit_btn":
					btn.OnClick = QuitEditor(win)
				case "save_btn":
					btn.OnClick = systems.SavePuzzle
				case "world_btn":
					btn.OnClick = systems.ChangeWorldToNext
				case "name_btn":
					btn.OnClick = OpenDialog("change_name")
				case "test_btn":
					btn.OnClick = TestPuzzle()
				case "check_puzzle_name":
					btn.OnClick = ChangeName
				default:
					switch key {
					case "editor_panel_top", "editor_panel_left":
						btn.OnClick = EditorMode(data.ModeFromSprString(btn.Sprite.Key), btn, dialog)
					default:
						btn.OnClick = Test(fmt.Sprintf("pressed button %s", btn.Key))
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
						//sprB := img.NewSprite("black_square_big", constants.UIBatch)
						sprS := img.NewSprite(bId.String(), constants.TileBatch)
						sprs := []*img.Sprite{sprS}
						if b == data.BlockFall {
							sprs = append(sprs, img.NewSprite(constants.TileFall, constants.TileBatch))
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
				Key:         "brush_btn",
				SprKey:      "brush_btn",
				ClickSprKey: "brush_btn_click",
				HelpText:    "Brush Tool (B)",
				Position:    pixel.V(-7.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				Key:         "line_btn",
				SprKey:      "line_btn",
				ClickSprKey: "line_btn_click",
				HelpText:    "Line Tool (L)",
				Position:    pixel.V(-6.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				Key:         "square_btn",
				SprKey:      "square_btn",
				ClickSprKey: "square_btn_click",
				HelpText:    "Square Tool (H)",
				Position:    pixel.V(-5.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				Key:         "fill_btn",
				SprKey:      "fill_btn",
				ClickSprKey: "fill_btn_click",
				HelpText:    "Fill Tool (G)",
				Position:    pixel.V(-4.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				Key:         "erase_btn",
				SprKey:      "erase_btn",
				ClickSprKey: "erase_btn_click",
				HelpText:    "Erase Tool (E)",
				Position:    pixel.V(-3.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				Key:         "eyedrop_btn",
				SprKey:      "eyedrop_btn",
				ClickSprKey: "eyedrop_btn_click",
				HelpText:    "Eyedrop Tool (Y)",
				Position:    pixel.V(-2.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				Key:         "wrench_btn",
				SprKey:      "wrench_btn",
				ClickSprKey: "wrench_btn_click",
				HelpText:    "Wrench Tool (P)",
				Position:    pixel.V(-1.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				Key:         "select_btn",
				SprKey:      "select_btn",
				ClickSprKey: "select_btn_click",
				HelpText:    "Select Tool (M)",
				Position:    pixel.V(-0.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				Key:         "cut_btn",
				SprKey:      "cut_btn",
				ClickSprKey: "cut_btn_click",
				HelpText:    "Cut Selection (Ctrl+X)",
				Position:    pixel.V(0.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				Key:         "copy_btn",
				SprKey:      "copy_btn",
				ClickSprKey: "copy_btn_click",
				HelpText:    "Copy Selection (Ctrl+C)",
				Position:    pixel.V(1.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				Key:         "paste_btn",
				SprKey:      "paste_btn",
				ClickSprKey: "paste_btn_click",
				HelpText:    "Paste Selection (Ctrl+V)",
				Position:    pixel.V(2.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				Key:         "fliph_btn",
				SprKey:      "fliph_btn",
				ClickSprKey: "fliph_btn_click",
				HelpText:    "Flip Selection Horizontal (U)",
				Position:    pixel.V(3.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				Key:         "flipv_btn",
				SprKey:      "flipv_btn",
				ClickSprKey: "flipv_btn_click",
				HelpText:    "flipv Selection Vertical (K)",
				Position:    pixel.V(4.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				Key:         "undo_btn",
				SprKey:      "undo_btn",
				ClickSprKey: "undo_btn_click",
				HelpText:    "Undo (Ctrl+Z)",
				Position:    pixel.V(5.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				Key:         "redo_btn",
				SprKey:      "redo_btn",
				ClickSprKey: "redo_btn_click",
				HelpText:    "Redo (Ctrl+Shift+Z)",
				Position:    pixel.V(6.5*world.TileSize, 0),
				Element:     data.ButtonElement,
			},
			{
				Key:      "block_select",
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
				Key:         "brush_btn",
				SprKey:      "brush_btn",
				ClickSprKey: "brush_btn_click",
				HelpText:    "Brush Tool (B)",
				Position:    pixel.V(-0.5*world.TileSize, 2.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				Key:         "line_btn",
				SprKey:      "line_btn",
				ClickSprKey: "line_btn_click",
				HelpText:    "Line Tool (L)",
				Position:    pixel.V(0.5*world.TileSize, 2.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				Key:         "square_btn",
				SprKey:      "square_btn",
				ClickSprKey: "square_btn_click",
				HelpText:    "Square Tool (H)",
				Position:    pixel.V(-0.5*world.TileSize, 1.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				Key:         "fill_btn",
				SprKey:      "fill_btn",
				ClickSprKey: "fill_btn_click",
				HelpText:    "Fill Tool (G)",
				Position:    pixel.V(0.5*world.TileSize, 1.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				Key:         "erase_btn",
				SprKey:      "erase_btn",
				ClickSprKey: "erase_btn_click",
				HelpText:    "Erase Tool (E)",
				Position:    pixel.V(-0.5*world.TileSize, 0.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				Key:         "eyedrop_btn",
				SprKey:      "eyedrop_btn",
				ClickSprKey: "eyedrop_btn_click",
				HelpText:    "Eyedrop Tool (Y)",
				Position:    pixel.V(0.5*world.TileSize, 0.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				Key:         "wrench_btn",
				SprKey:      "wrench_btn",
				ClickSprKey: "wrench_btn_click",
				HelpText:    "Wrench Tool (P)",
				Position:    pixel.V(-0.5*world.TileSize, -0.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				Key:         "select_btn",
				SprKey:      "select_btn",
				ClickSprKey: "select_btn_click",
				HelpText:    "Select Tool (M)",
				Position:    pixel.V(-0.5*world.TileSize, -1.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				Key:         "cut_btn",
				SprKey:      "cut_btn",
				ClickSprKey: "cut_btn_click",
				HelpText:    "Cut Selection (Ctrl+X)",
				Position:    pixel.V(0.5*world.TileSize, -1.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				Key:         "copy_btn",
				SprKey:      "copy_btn",
				ClickSprKey: "copy_btn_click",
				HelpText:    "Copy Selection (Ctrl+C)",
				Position:    pixel.V(-0.5*world.TileSize, -2.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				Key:         "paste_btn",
				SprKey:      "paste_btn",
				ClickSprKey: "paste_btn_click",
				HelpText:    "Paste Selection (Ctrl+V)",
				Position:    pixel.V(0.5*world.TileSize, -2.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				Key:         "fliph_btn",
				SprKey:      "fliph_btn",
				ClickSprKey: "fliph_btn_click",
				HelpText:    "Flip Selection Horizontal (U)",
				Position:    pixel.V(-0.5*world.TileSize, -3.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				Key:         "flipv_btn",
				SprKey:      "flipv_btn",
				ClickSprKey: "flipv_btn_click",
				HelpText:    "flipv Selection Vertical (K)",
				Position:    pixel.V(0.5*world.TileSize, -3.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				Key:         "undo_btn",
				SprKey:      "undo_btn",
				ClickSprKey: "undo_btn_click",
				HelpText:    "Undo (Ctrl+Z)",
				Position:    pixel.V(-0.5*world.TileSize, -4.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				Key:         "redo_btn",
				SprKey:      "redo_btn",
				ClickSprKey: "redo_btn_click",
				HelpText:    "Redo (Ctrl+Shift+Z)",
				Position:    pixel.V(0.5*world.TileSize, -4.5*world.TileSize),
				Element:     data.ButtonElement,
			},
			{
				Key:      "block_select",
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

func editorOptPanels() {
	editorOptBottomConstructor := &data.DialogConstructor{
		Key:    "editor_options_bot",
		Width:  9,
		Height: 1,
		Pos:    pixel.V(0, -400),
		Elements: []data.ElementConstructor{
			{
				Key:         "quit_btn",
				SprKey:      "quit_btn",
				ClickSprKey: "quit_btn_click",
				HelpText:    "Quit (Ctrl+Q)",
				Position:    pixel.V(-4*world.TileSize, 0),
			},
			{
				Key:         "new_btn",
				SprKey:      "new_btn",
				ClickSprKey: "new_btn_click",
				HelpText:    "New Puzzle Group (Ctrl+N)",
				Position:    pixel.V(-3*world.TileSize, 0),
			},
			{
				Key:         "save_btn",
				SprKey:      "save_btn",
				ClickSprKey: "save_btn_click",
				HelpText:    "Save Puzzle Group (Ctrl+S)",
				Position:    pixel.V(-2*world.TileSize, 0),
			},
			{
				Key:         "open_btn",
				SprKey:      "open_btn",
				ClickSprKey: "open_btn_click",
				HelpText:    "Open Puzzle Group (Ctrl+O)",
				Position:    pixel.V(-world.TileSize, 0),
			},
			{
				Key:         "prev_btn",
				SprKey:      "prev_btn",
				ClickSprKey: "prev_btn_click",
				HelpText:    "Previous Puzzle",
				Position:    pixel.V(0, 0),
			},
			{
				Key:         "next_btn",
				SprKey:      "next_btn",
				ClickSprKey: "next_btn_click",
				HelpText:    "Next Puzzle",
				Position:    pixel.V(1*world.TileSize, 0),
			},
			{
				Key:         "name_btn",
				SprKey:      "name_btn",
				ClickSprKey: "name_btn_click",
				HelpText:    "Change Name",
				Position:    pixel.V(2*world.TileSize, 0),
			},
			{
				Key:         "world_btn",
				SprKey:      "world_btn",
				ClickSprKey: "world_btn_click",
				HelpText:    "Change World (Tab)",
				Position:    pixel.V(3*world.TileSize, 0),
			},
			{
				Key:         "test_btn",
				SprKey:      "test_btn",
				ClickSprKey: "test_btn_click",
				HelpText:    "Test Puzzle",
				Position:    pixel.V(4*world.TileSize, 0),
			},
		},
	}
	data.NewDialog(editorOptBottomConstructor)
	editorOptRightConstructor := &data.DialogConstructor{
		Key:    "editor_options_right",
		Width:  1,
		Height: 9,
		Pos:    pixel.V(670, 0),
		Elements: []data.ElementConstructor{
			{
				Key:         "quit_btn",
				SprKey:      "quit_btn",
				ClickSprKey: "quit_btn_click",
				HelpText:    "Quit (Ctrl+Q)",
				Position:    pixel.V(0, 4*world.TileSize),
			},
			{
				Key:         "new_btn",
				SprKey:      "new_btn",
				ClickSprKey: "new_btn_click",
				HelpText:    "New Puzzle Group (Ctrl+N)",
				Position:    pixel.V(0, 3*world.TileSize),
			},
			{
				Key:         "save_btn",
				SprKey:      "save_btn",
				ClickSprKey: "save_btn_click",
				HelpText:    "Save Puzzle Group (Ctrl+S)",
				Position:    pixel.V(0, 2*world.TileSize),
			},
			{
				Key:         "open_btn",
				SprKey:      "open_btn",
				ClickSprKey: "open_btn_click",
				HelpText:    "Open Puzzle Group (Ctrl+O)",
				Position:    pixel.V(0, world.TileSize),
			},
			{
				Key:         "prev_btn",
				SprKey:      "prev_btn",
				ClickSprKey: "prev_btn_click",
				HelpText:    "Previous Puzzle",
				Position:    pixel.V(0, 0),
			},
			{
				Key:         "next_btn",
				SprKey:      "next_btn",
				ClickSprKey: "next_btn_click",
				HelpText:    "Next Puzzle",
				Position:    pixel.V(0, -world.TileSize),
			},
			{
				Key:         "name_btn",
				SprKey:      "name_btn",
				ClickSprKey: "name_btn_click",
				HelpText:    "Change Name",
				Position:    pixel.V(0, -2*world.TileSize),
			},
			{
				Key:         "world_btn",
				SprKey:      "world_btn",
				ClickSprKey: "world_btn_click",
				HelpText:    "Change World (Tab)",
				Position:    pixel.V(0, -3*world.TileSize),
			},
			{
				Key:         "test_btn",
				SprKey:      "test_btn",
				ClickSprKey: "test_btn_click",
				HelpText:    "Test Puzzle",
				Position:    pixel.V(0, -4*world.TileSize),
			},
		},
	}
	data.NewDialog(editorOptRightConstructor)
}
