package load

import (
	"gemrunner/internal/data"
	"gemrunner/pkg/world"
	"github.com/faiface/pixel"
)

func Dialogs() {
	data.OpenPuzzleConstructor = &data.DialogConstructor{
		Key:    "open_puzzle",
		Width:  7,
		Height: 7,
		Buttons: []data.ButtonConstructor{
			{
				SprKey:      "cancel_btn_big",
				ClickSprKey: "cancel_btn_click_big",
				HelpText:    "Cancel",
				Position:    pixel.V(32, -32),
			},
			{
				SprKey:      "open_btn_big",
				ClickSprKey: "open_btn_click_big",
				HelpText:    "Open",
				Position:    pixel.V(0, -32),
			},
		},
	}
	data.NewDialog(data.OpenPuzzleConstructor)
	data.EditorPanelConstructor = &data.DialogConstructor{
		Key:    "editor_panel",
		Width:  15,
		Height: 1,
		Pos:    pixel.V(0, 400),
		Buttons: []data.ButtonConstructor{
			{
				SprKey:      "brush_btn",
				ClickSprKey: "brush_btn_click",
				HelpText:    "Brush Tool (B)",
				Position:    pixel.V(-7*world.TileSize, 0),
			},
			{
				SprKey:      "line_btn",
				ClickSprKey: "line_btn_click",
				HelpText:    "Line Tool (L)",
				Position:    pixel.V(-6*world.TileSize, 0),
			},
			{
				SprKey:      "square_btn",
				ClickSprKey: "square_btn_click",
				HelpText:    "Square Tool (H)",
				Position:    pixel.V(-5*world.TileSize, 0),
			},
			{
				SprKey:      "fill_btn",
				ClickSprKey: "fill_btn_click",
				HelpText:    "Fill Tool (G)",
				Position:    pixel.V(-4*world.TileSize, 0),
			},
			{
				SprKey:      "erase_btn",
				ClickSprKey: "erase_btn_click",
				HelpText:    "Erase Tool (E)",
				Position:    pixel.V(-3*world.TileSize, 0),
			},
			{
				SprKey:      "eyedrop_btn",
				ClickSprKey: "eyedrop_btn_click",
				HelpText:    "Eyedrop Tool (Y)",
				Position:    pixel.V(-2*world.TileSize, 0),
			},
			{
				SprKey:      "wrench_btn",
				ClickSprKey: "wrench_btn_click",
				HelpText:    "Wrench Tool (P)",
				Position:    pixel.V(-1*world.TileSize, 0),
			},
			{
				SprKey:      "select_btn",
				ClickSprKey: "select_btn_click",
				HelpText:    "Select Tool (M)",
				Position:    pixel.V(0, 0),
			},
			{
				SprKey:      "cut_btn",
				ClickSprKey: "cut_btn_click",
				HelpText:    "Cut Selection (Ctrl+X)",
				Position:    pixel.V(world.TileSize, 0),
			},
			{
				SprKey:      "copy_btn",
				ClickSprKey: "copy_btn_click",
				HelpText:    "Copy Selection (Ctrl+C)",
				Position:    pixel.V(2*world.TileSize, 0),
			},
			{
				SprKey:      "paste_btn",
				ClickSprKey: "paste_btn_click",
				HelpText:    "Paste Selection (Ctrl+V)",
				Position:    pixel.V(3*world.TileSize, 0),
			},
			{
				SprKey:      "fliph_btn",
				ClickSprKey: "fliph_btn_click",
				HelpText:    "Flip Selection Horizontal (U)",
				Position:    pixel.V(4*world.TileSize, 0),
			},
			{
				SprKey:      "flipv_btn",
				ClickSprKey: "flipv_btn_click",
				HelpText:    "flipv Selection Vertical (K)",
				Position:    pixel.V(5*world.TileSize, 0),
			},
			{
				SprKey:      "undo_btn",
				ClickSprKey: "undo_btn_click",
				HelpText:    "Undo (Ctrl+Z)",
				Position:    pixel.V(6*world.TileSize, 0),
			},
			{
				SprKey:      "redo_btn",
				ClickSprKey: "redo_btn_click",
				HelpText:    "Redo (Ctrl+Shift+Z)",
				Position:    pixel.V(7*world.TileSize, 0),
			},
		},
	}
	data.NewDialog(data.EditorPanelConstructor)
	data.EditorOptConstructor = &data.DialogConstructor{
		Key:    "editor_options",
		Width:  8,
		Height: 1,
		Pos:    pixel.V(0, -400),
		Buttons: []data.ButtonConstructor{
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
	data.NewDialog(data.EditorOptConstructor)

	for key, dialog := range data.Dialogs {
		for _, btn := range dialog.Buttons {
			switch btn.Key {
			case "cancel_btn":
				btn.OnClick = CancelDialog(key)
			default:
				switch key {
				case "editor_panel":
					btn.OnClick = EditorMode(data.ModeFromSprString(btn.Sprite.Key), btn, dialog)
				case "editor_options":
					btn.OnClick = Test
				}
			}
		}
	}
}
