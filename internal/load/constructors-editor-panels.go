package load

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/ui"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

var (
	// editor panels
	EditorPanelTopConstructor  *ui.DialogConstructor
	EditorOptBottomConstructor *ui.DialogConstructor
)

func InitEditorPanels() {
	// Editor Panels
	EditorPanelTopConstructor = &ui.DialogConstructor{
		Key:    constants.DialogEditorPanelTop,
		Width:  17,
		Height: 1,
		Pos:    pixel.V(0, 400),
		Elements: []ui.ElementConstructor{
			{
				Key:         "brush_btn",
				SprKey:      "brush_btn",
				SprKey2:     "brush_btn_click",
				HelpText:    "Brush Tool (B)",
				Batch:       constants.UIBatch,
				Position:    pixel.V(-8*world.TileSize, 0),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "line_btn",
				SprKey:      "line_btn",
				SprKey2:     "line_btn_click",
				HelpText:    "Line Tool (L)",
				Batch:       constants.UIBatch,
				Position:    pixel.V(-7*world.TileSize, 0),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "square_btn",
				SprKey:      "square_btn",
				SprKey2:     "square_btn_click",
				HelpText:    "Square Tool (H)",
				Batch:       constants.UIBatch,
				Position:    pixel.V(-6*world.TileSize, 0),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "fill_btn",
				SprKey:      "fill_btn",
				SprKey2:     "fill_btn_click",
				HelpText:    "Fill Tool (G)",
				Batch:       constants.UIBatch,
				Position:    pixel.V(-5*world.TileSize, 0),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "erase_btn",
				SprKey:      "erase_btn",
				SprKey2:     "erase_btn_click",
				HelpText:    "Erase Tool (E)",
				Batch:       constants.UIBatch,
				Position:    pixel.V(-4*world.TileSize, 0),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "eyedrop_btn",
				SprKey:      "eyedrop_btn",
				SprKey2:     "eyedrop_btn_click",
				HelpText:    "Eyedrop Tool (Y)",
				Batch:       constants.UIBatch,
				Position:    pixel.V(-3*world.TileSize, 0),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "wrench_btn",
				SprKey:      "wrench_btn",
				SprKey2:     "wrench_btn_click",
				HelpText:    "Wrench Tool (P)",
				Batch:       constants.UIBatch,
				Position:    pixel.V(-2*world.TileSize, 0),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "wire_btn",
				SprKey:      "wire_btn",
				SprKey2:     "wire_btn_click",
				HelpText:    "Wire Tool (I)",
				Batch:       constants.UIBatch,
				Position:    pixel.V(-world.TileSize, 0),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "select_btn",
				SprKey:      "select_btn",
				SprKey2:     "select_btn_click",
				HelpText:    "Select Tool (M)",
				Batch:       constants.UIBatch,
				Position:    pixel.V(0, 0),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "cut_btn",
				SprKey:      "cut_btn",
				SprKey2:     "cut_btn_click",
				HelpText:    "Cut Selection (Ctrl+X)",
				Batch:       constants.UIBatch,
				Position:    pixel.V(world.TileSize, 0),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "copy_btn",
				SprKey:      "copy_btn",
				SprKey2:     "copy_btn_click",
				HelpText:    "Copy Selection (Ctrl+C)",
				Batch:       constants.UIBatch,
				Position:    pixel.V(2*world.TileSize, 0),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "paste_btn",
				SprKey:      "paste_btn",
				SprKey2:     "paste_btn_click",
				HelpText:    "Paste Selection (Ctrl+V)",
				Batch:       constants.UIBatch,
				Position:    pixel.V(3*world.TileSize, 0),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "fliph_btn",
				SprKey:      "fliph_btn",
				SprKey2:     "fliph_btn_click",
				Batch:       constants.UIBatch,
				HelpText:    "Flip Selection Horizontal (U)",
				Position:    pixel.V(4*world.TileSize, 0),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "flipv_btn",
				SprKey:      "flipv_btn",
				SprKey2:     "flipv_btn_click",
				Batch:       constants.UIBatch,
				HelpText:    "flipv Selection Vertical (K)",
				Position:    pixel.V(5*world.TileSize, 0),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "undo_btn",
				SprKey:      "undo_btn",
				SprKey2:     "undo_btn_click",
				Batch:       constants.UIBatch,
				HelpText:    "Undo (Ctrl+Z)",
				Position:    pixel.V(6*world.TileSize, 0),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "redo_btn",
				SprKey:      "redo_btn",
				SprKey2:     "redo_btn_click",
				Batch:       constants.UIBatch,
				HelpText:    "Redo (Ctrl+Shift+Z)",
				Position:    pixel.V(7*world.TileSize, 0),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "block_select",
				SprKey:      "editor_tile_bg",
				Batch:       constants.UIBatch,
				Position:    pixel.V(8*world.TileSize, 0),
				ElementType: ui.SpriteElement,
			},
		},
	}
	EditorOptBottomConstructor = &ui.DialogConstructor{
		Key:    constants.DialogEditorOptionsBot,
		Width:  9,
		Height: 1,
		Pos:    pixel.V(0, -400),
		Elements: []ui.ElementConstructor{
			{
				Key:      "exit_editor_btn",
				SprKey:   "quit_btn",
				SprKey2:  "quit_btn_click",
				Batch:    constants.UIBatch,
				HelpText: "Quit (Ctrl+Q)",
				Position: pixel.V(-4*world.TileSize, 0),
			},
			{
				Key:      "new_btn",
				SprKey:   "new_btn",
				SprKey2:  "new_btn_click",
				Batch:    constants.UIBatch,
				HelpText: "New Puzzle Group (Ctrl+N)",
				Position: pixel.V(-3*world.TileSize, 0),
			},
			{
				Key:      "save_btn",
				SprKey:   "save_btn",
				SprKey2:  "save_btn_click",
				Batch:    constants.UIBatch,
				HelpText: "Save Puzzle Group (Ctrl+S)",
				Position: pixel.V(-2*world.TileSize, 0),
			},
			{
				Key:      "open_btn",
				SprKey:   "open_btn",
				SprKey2:  "open_btn_click",
				Batch:    constants.UIBatch,
				HelpText: "Open Puzzle Group (Ctrl+O)",
				Position: pixel.V(-world.TileSize, 0),
			},
			{
				Key:      "prev_btn",
				SprKey:   "prev_btn",
				SprKey2:  "prev_btn_click",
				Batch:    constants.UIBatch,
				HelpText: "Previous Puzzle",
				Position: pixel.V(0, 0),
			},
			{
				Key:      "next_btn",
				SprKey:   "next_btn",
				SprKey2:  "next_btn_click",
				Batch:    constants.UIBatch,
				HelpText: "Next Puzzle",
				Position: pixel.V(1*world.TileSize, 0),
			},
			{
				Key:      "name_btn",
				SprKey:   "name_btn",
				SprKey2:  "name_btn_click",
				Batch:    constants.UIBatch,
				HelpText: "Change Name",
				Position: pixel.V(2*world.TileSize, 0),
			},
			{
				Key:      "world_btn",
				SprKey:   "world_btn",
				SprKey2:  "world_btn_click",
				Batch:    constants.UIBatch,
				HelpText: "Change World (Tab)",
				Position: pixel.V(3*world.TileSize, 0),
			},
			{
				Key:      "test_btn",
				SprKey:   "test_btn",
				SprKey2:  "test_btn_click",
				Batch:    constants.UIBatch,
				HelpText: "Test Puzzle",
				Position: pixel.V(4*world.TileSize, 0),
			},
		},
	}
	//f, err := os.Create("assets/ui/editor_panel_left.json")
	//if err != nil {
	//	panic(err)
	//}
	//bts, err := json.Marshal(EditorPanelLeftConstructor)
	//if err != nil {
	//	panic(err)
	//}
	//_, err = f.Write(bts)
	//if err != nil {
	//	panic(err)
	//}
	//f.Close()
}
