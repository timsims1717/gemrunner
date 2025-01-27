package load

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/ui"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

var (
	// main menu
	AddPlayersConstructor *ui.DialogConstructor
	PuzzleListEntry       ui.ElementConstructor
	PuzzleTitleItem       ui.ElementConstructor
	PuzzleTitleShadowItem ui.ElementConstructor
	PuzzleAuthorItem      ui.ElementConstructor
	PuzzleNumberSpr       ui.ElementConstructor
	PlayerNumberSpr       ui.ElementConstructor
	NumberItems           ui.ElementConstructor
	FavoriteItem          ui.ElementConstructor
	// options

	// in game
	PauseConstructor       *ui.DialogConstructor
	PuzzleTitleConstructor *ui.DialogConstructor
)

func InitMainMenuConstructors() {
	AddPlayersConstructor = &ui.DialogConstructor{
		Key:    constants.DialogAddPlayers,
		Width:  18,
		Height: 6,
		Elements: []ui.ElementConstructor{
			{
				Key:         "any_button_p2",
				Text:        "Press\na key",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(-48, 16),
				ElementType: ui.TextElement,
				Anchor:      pixel.Right,
			},
			{
				Key:         "any_button_p3",
				Text:        "Press\na key",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(20, 16),
				ElementType: ui.TextElement,
				Anchor:      pixel.Right,
			},
			{
				Key:         "any_button_p4",
				Text:        "Press\na key",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(88, 16),
				ElementType: ui.TextElement,
				Anchor:      pixel.Right,
			},
			{
				Key:         "selected_p1_cnt",
				Position:    pixel.V(-102, 15),
				ElementType: ui.ContainerElement,
				Width:       4 * world.TileSize,
				Height:      4 * world.TileSize,
				SubElements: []ui.ElementConstructor{
					{
						Key:         "selected_p1_spr",
						SprKey:      "player1",
						Batch:       constants.UIBatch,
						Position:    pixel.V(0, 0),
						ElementType: ui.SpriteElement,
					},
				},
			},
			{
				Key:         "selected_p2_cnt",
				Position:    pixel.V(-34, 15),
				ElementType: ui.ContainerElement,
				Width:       4 * world.TileSize,
				Height:      4 * world.TileSize,
				SubElements: []ui.ElementConstructor{
					{
						Key:         "selected_p2_spr",
						SprKey:      "player2",
						Batch:       constants.UIBatch,
						Position:    pixel.V(0, 0),
						ElementType: ui.SpriteElement,
					},
				},
			},
			{
				Key:         "selected_p3_cnt",
				Position:    pixel.V(34, 15),
				ElementType: ui.ContainerElement,
				Width:       4 * world.TileSize,
				Height:      4 * world.TileSize,
				SubElements: []ui.ElementConstructor{
					{
						Key:         "selected_p3_spr",
						SprKey:      "player3",
						Batch:       constants.UIBatch,
						Position:    pixel.V(0, 0),
						ElementType: ui.SpriteElement,
					},
				},
			},
			{
				Key:         "selected_p4_cnt",
				Position:    pixel.V(102, 15),
				ElementType: ui.ContainerElement,
				Width:       4 * world.TileSize,
				Height:      4 * world.TileSize,
				SubElements: []ui.ElementConstructor{
					{
						Key:         "selected_p4_spr",
						SprKey:      "player4",
						Batch:       constants.UIBatch,
						Position:    pixel.V(0, 0),
						ElementType: ui.SpriteElement,
					},
				},
			},
			{
				Key:         "cancel",
				SprKey:      "cancel_btn_big",
				SprKey2:     "cancel_btn_click_big",
				Batch:       constants.UIBatch,
				HelpText:    "Cancel",
				Position:    pixel.V(132, -36),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "confirm",
				SprKey:      "check_btn_big",
				SprKey2:     "check_btn_click_big",
				Batch:       constants.UIBatch,
				HelpText:    "Confirm",
				Position:    pixel.V(112, -36),
				ElementType: ui.ButtonElement,
			},
		},
	}
	PuzzleListEntry = ui.ElementConstructor{
		Key:         "puzzle_container_%d",
		HelpText:    "Select this puzzle.",
		ElementType: ui.ContainerElement,
		Width:       9 * world.TileSize,
		Height:      2 * world.TileSize,
	}
	PuzzleTitleShadowItem = ui.ElementConstructor{
		Key:         "puzzle_title_shadow",
		Position:    pixel.V(-18, 9),
		Color:       pixel.ToRGBA(constants.ColorBlue),
		ElementType: ui.TextElement,
		Anchor:      pixel.Right,
	}
	PuzzleTitleItem = ui.ElementConstructor{
		Key:         "puzzle_title",
		Position:    pixel.V(-18, 8),
		Color:       pixel.ToRGBA(constants.ColorWhite),
		ElementType: ui.TextElement,
		Anchor:      pixel.Right,
	}
	PuzzleAuthorItem = ui.ElementConstructor{
		Key:         "puzzle_author",
		Position:    pixel.V(-12, 0),
		Color:       pixel.ToRGBA(constants.ColorWhite),
		ElementType: ui.TextElement,
		Anchor:      pixel.Right,
	}
	PuzzleNumberSpr = ui.ElementConstructor{
		Key:         "puzzle_num",
		SprKey:      "puzzle_symbol",
		Position:    pixel.V(-32, 8),
		Batch:       constants.UIBatch,
		ElementType: ui.SpriteElement,
	}
	PlayerNumberSpr = ui.ElementConstructor{
		Key:         "player_num",
		SprKey:      "player_symbol",
		Position:    pixel.V(-48, 8),
		Batch:       constants.UIBatch,
		ElementType: ui.SpriteElement,
	}
	NumberItems = ui.ElementConstructor{
		Key:         "number_text",
		Position:    pixel.V(-50, -4),
		Color:       pixel.ToRGBA(constants.ColorWhite),
		ElementType: ui.TextElement,
		Anchor:      pixel.Right,
	}
	FavoriteItem = ui.ElementConstructor{
		Key:         "favorite_symbol_%d",
		SprKey:      "heart_empty",
		SprKey2:     "heart_full",
		Position:    pixel.V(-64, 0),
		Batch:       constants.UIBatch,
		ElementType: ui.CheckboxElement,
	}
	// pause menu
	PauseConstructor = &ui.DialogConstructor{
		Key:    constants.DialogPauseMenu,
		Width:  7,
		Height: 9,
		Elements: []ui.ElementConstructor{
			{
				Key:         "pause_resume_ct",
				Position:    pixel.V(0, 56),
				ElementType: ui.ContainerElement,
				Width:       7*world.TileSize - 8,
				Height:      world.TileSize * 1.5,
				SubElements: []ui.ElementConstructor{
					{
						Key:         "pause_resume_text",
						Text:        "Resume",
						Color:       pixel.ToRGBA(constants.ColorWhite),
						Position:    pixel.V(-16, 0),
						ElementType: ui.TextElement,
						Anchor:      pixel.Right,
					},
				},
			},
			{
				Key:         "pause_restart_ct",
				Position:    pixel.V(0, 28),
				ElementType: ui.ContainerElement,
				Width:       7*world.TileSize - 8,
				Height:      world.TileSize * 1.5,
				SubElements: []ui.ElementConstructor{
					{
						Key:         "pause_restart_text",
						Text:        "Restart",
						Color:       pixel.ToRGBA(constants.ColorWhite),
						Position:    pixel.V(-18, 0),
						ElementType: ui.TextElement,
						Anchor:      pixel.Right,
					},
				},
			},
			{
				Key:         "pause_options_ct",
				Position:    pixel.V(0, 0),
				ElementType: ui.ContainerElement,
				Width:       7*world.TileSize - 8,
				Height:      world.TileSize * 1.5,
				SubElements: []ui.ElementConstructor{
					{
						Key:         "pause_restart_text",
						Text:        "Options",
						Color:       pixel.ToRGBA(constants.ColorWhite),
						Position:    pixel.V(-17, 0),
						ElementType: ui.TextElement,
						Anchor:      pixel.Right,
					},
				},
			},
			{
				Key:         "pause_quit_mm_ct",
				Position:    pixel.V(0, -28),
				ElementType: ui.ContainerElement,
				Width:       7*world.TileSize - 8,
				Height:      world.TileSize * 1.5,
				SubElements: []ui.ElementConstructor{
					{
						Key:         "pause_quit_mm_text",
						Text:        "Quit to Menu",
						Color:       pixel.ToRGBA(constants.ColorWhite),
						Position:    pixel.V(-36, 0),
						ElementType: ui.TextElement,
						Anchor:      pixel.Right,
					},
				},
			},
			{
				Key:         "pause_quit_full_ct",
				Position:    pixel.V(0, -56),
				ElementType: ui.ContainerElement,
				Width:       7*world.TileSize - 8,
				Height:      world.TileSize * 1.5,
				SubElements: []ui.ElementConstructor{
					{
						Key:         "pause_quit_full_text",
						Text:        "Quit to Desktop",
						Color:       pixel.ToRGBA(constants.ColorWhite),
						Position:    pixel.V(-42, 0),
						ElementType: ui.TextElement,
						Anchor:      pixel.Right,
					},
				},
			},
		},
	}
	PuzzleTitleConstructor = &ui.DialogConstructor{
		Key:    constants.DialogPuzzleTitle,
		Width:  6,
		Height: 1,
		Pos:    pixel.V(0, 414),
		Elements: []ui.ElementConstructor{
			{
				Key:         "puzzle_title_bg",
				Text:        "-",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(0, 1),
				ElementType: ui.TextElement,
				Anchor:      pixel.Right,
			},
			{
				Key:         "puzzle_title",
				Text:        "-",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(0, 0),
				ElementType: ui.TextElement,
				Anchor:      pixel.Right,
			},
		},
	}
	//f, err := os.Create("assets/ui/main_menu.json")
	//if err != nil {
	//	panic(err)
	//}
	//bts, err := json.Marshal(MainMenuConstructor)
	//if err != nil {
	//	panic(err)
	//}
	//_, err = f.Write(bts)
	//if err != nil {
	//	panic(err)
	//}
	//f.Close()
}
