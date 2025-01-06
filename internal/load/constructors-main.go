package load

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/ui"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

var (
	// main menu
	MainMenuConstructor   *ui.DialogConstructor
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
	Player1InvConstructor  *ui.DialogConstructor
	Player2InvConstructor  *ui.DialogConstructor
	Player3InvConstructor  *ui.DialogConstructor
	Player4InvConstructor  *ui.DialogConstructor
	PuzzleTitleConstructor *ui.DialogConstructor
)

func InitMainMenuConstructors() {
	MainMenuConstructor = &ui.DialogConstructor{
		Key:    constants.DialogMainMenu,
		Width:  28,
		Height: 16,
		Elements: []ui.ElementConstructor{
			{
				Key:         "play_local_game_btn",
				SprKey:      "play_local",
				SprKey2:     "play_local_click",
				Batch:       constants.UIBatch,
				HelpText:    "Play a local game with up to four players.",
				Position:    pixel.V(-144, -64),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "play_online_game_btn",
				SprKey:      "play_online",
				SprKey2:     "play_online_click",
				Batch:       constants.UIBatch,
				HelpText:    "Play an online game with up to four players.",
				Position:    pixel.V(-72, -64),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "start_editor_btn",
				SprKey:      "start_editor",
				SprKey2:     "start_editor_click",
				Batch:       constants.UIBatch,
				HelpText:    "Create your own puzzle sets.",
				Position:    pixel.V(0, -64),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "start_options_btn",
				SprKey:      "start_options",
				SprKey2:     "start_options_click",
				Batch:       constants.UIBatch,
				HelpText:    "Change the game options.",
				Position:    pixel.V(72, -64),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "quit_btn",
				SprKey:      "main_quit",
				SprKey2:     "main_quit_click",
				Batch:       constants.UIBatch,
				HelpText:    "Quit the game.",
				Position:    pixel.V(144, -64),
				ElementType: ui.ButtonElement,
			},
		},
	}
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
				Key:         "cancel_add_players",
				SprKey:      "cancel_btn_big",
				SprKey2:     "cancel_btn_click_big",
				Batch:       constants.UIBatch,
				HelpText:    "Cancel",
				Position:    pixel.V(132, -36),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "confirm_add_players",
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
		Key:         "favorite_symbol",
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
	Player1InvConstructor = &ui.DialogConstructor{
		Key:    constants.DialogPlayer1Inv,
		Width:  6,
		Height: 1,
		Pos:    pixel.V(-528, -414),
		Elements: []ui.ElementConstructor{
			{
				Key:         "player_score",
				Text:        "0000000",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(-12, 5),
				ElementType: ui.TextElement,
				Anchor:      pixel.Right,
			},
			{
				Key:         "player_skull",
				SprKey:      "player1_skull",
				Batch:       constants.UIBatch,
				HelpText:    "Player 1's total deaths",
				Position:    pixel.V(-42, -4),
				ElementType: ui.SpriteElement,
			},
			{
				Key:         "player_deaths",
				Text:        "x000",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(-36, -3),
				ElementType: ui.TextElement,
				Anchor:      pixel.Right,
			},
			{
				Key:         "player_gem",
				SprKey:      "player1_gem",
				Batch:       constants.UIBatch,
				HelpText:    "Player 1's total gems",
				Position:    pixel.V(-6, -3),
				ElementType: ui.SpriteElement,
			},
			{
				Key:         "player_gems",
				Text:        "x0000",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(0, -3),
				ElementType: ui.TextElement,
				Anchor:      pixel.Right,
			},
			{
				Key:         "player_inv_cnt",
				Position:    pixel.V(40, 0),
				ElementType: ui.ContainerElement,
				Width:       world.TileSize,
				Height:      world.TileSize,
			},
			{
				Key:         "player_inv_item",
				SprKey:      "rock",
				Batch:       constants.TileBatch,
				Position:    pixel.V(40, 0),
				ElementType: ui.SpriteElement,
			},
			{
				Key:         "player_inv_item_2",
				SprKey:      "rock",
				Batch:       constants.TileBatch,
				Position:    pixel.V(40, 0),
				ElementType: ui.SpriteElement,
			},
		},
	}
	Player2InvConstructor = &ui.DialogConstructor{
		Key:    constants.DialogPlayer2Inv,
		Width:  6,
		Height: 1,
		Pos:    pixel.V(-264, -414),
		Elements: []ui.ElementConstructor{
			{
				Key:         "player_score",
				Text:        "0000000",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(-12, 5),
				ElementType: ui.TextElement,
				Anchor:      pixel.Right,
			},
			{
				Key:         "player_skull",
				SprKey:      "player2_skull",
				Batch:       constants.UIBatch,
				HelpText:    "Player 2's total deaths",
				Position:    pixel.V(-42, -4),
				ElementType: ui.SpriteElement,
			},
			{
				Key:         "player_deaths",
				Text:        "x000",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(-36, -3),
				ElementType: ui.TextElement,
				Anchor:      pixel.Right,
			},
			{
				Key:         "player_gem",
				SprKey:      "player2_gem",
				Batch:       constants.UIBatch,
				HelpText:    "Player 2's total gems",
				Position:    pixel.V(-6, -3),
				ElementType: ui.SpriteElement,
			},
			{
				Key:         "player_gems",
				Text:        "x0000",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(0, -3),
				ElementType: ui.TextElement,
				Anchor:      pixel.Right,
			},
			{
				Key:         "player_inv_cnt",
				Position:    pixel.V(40, 0),
				ElementType: ui.ContainerElement,
				Width:       world.TileSize,
				Height:      world.TileSize,
			},
			{
				Key:         "player_inv_item",
				SprKey:      "rock",
				Batch:       constants.TileBatch,
				Position:    pixel.V(40, 0),
				ElementType: ui.SpriteElement,
			},
			{
				Key:         "player_inv_item_2",
				SprKey:      "rock",
				Batch:       constants.TileBatch,
				Position:    pixel.V(40, 0),
				ElementType: ui.SpriteElement,
			},
		},
	}
	Player3InvConstructor = &ui.DialogConstructor{
		Key:    constants.DialogPlayer3Inv,
		Width:  6,
		Height: 1,
		Pos:    pixel.V(264, -414),
		Elements: []ui.ElementConstructor{
			{
				Key:         "player_score",
				Text:        "0000000",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(-12, 5),
				ElementType: ui.TextElement,
				Anchor:      pixel.Right,
			},
			{
				Key:         "player_skull",
				SprKey:      "player3_skull",
				Batch:       constants.UIBatch,
				HelpText:    "Player 3's total deaths",
				Position:    pixel.V(-42, -4),
				ElementType: ui.SpriteElement,
			},
			{
				Key:         "player_deaths",
				Text:        "x000",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(-36, -3),
				ElementType: ui.TextElement,
				Anchor:      pixel.Right,
			},
			{
				Key:         "player_gem",
				SprKey:      "player3_gem",
				Batch:       constants.UIBatch,
				HelpText:    "Player 3's total gems",
				Position:    pixel.V(-6, -3),
				ElementType: ui.SpriteElement,
			},
			{
				Key:         "player_gems",
				Text:        "x0000",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(0, -3),
				ElementType: ui.TextElement,
				Anchor:      pixel.Right,
			},
			{
				Key:         "player_inv_cnt",
				Position:    pixel.V(40, 0),
				ElementType: ui.ContainerElement,
				Width:       world.TileSize,
				Height:      world.TileSize,
			},
			{
				Key:         "player_inv_item",
				SprKey:      "rock",
				Batch:       constants.TileBatch,
				Position:    pixel.V(40, 0),
				ElementType: ui.SpriteElement,
			},
			{
				Key:         "player_inv_item_2",
				SprKey:      "rock",
				Batch:       constants.TileBatch,
				Position:    pixel.V(40, 0),
				ElementType: ui.SpriteElement,
			},
		},
	}
	Player4InvConstructor = &ui.DialogConstructor{
		Key:    constants.DialogPlayer4Inv,
		Width:  6,
		Height: 1,
		Pos:    pixel.V(528, -414),
		Elements: []ui.ElementConstructor{
			{
				Key:         "player_score",
				Text:        "0000000",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(-12, 5),
				ElementType: ui.TextElement,
				Anchor:      pixel.Right,
			},
			{
				Key:         "player_skull",
				SprKey:      "player4_skull",
				Batch:       constants.UIBatch,
				HelpText:    "Player 4's total deaths",
				Position:    pixel.V(-42, -4),
				ElementType: ui.SpriteElement,
			},
			{
				Key:         "player_deaths",
				Text:        "x000",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(-36, -3),
				ElementType: ui.TextElement,
				Anchor:      pixel.Right,
			},
			{
				Key:         "player_gem",
				SprKey:      "player4_gem",
				Batch:       constants.UIBatch,
				HelpText:    "Player 4's total gems",
				Position:    pixel.V(-6, -3),
				ElementType: ui.SpriteElement,
			},
			{
				Key:         "player_gems",
				Text:        "x0000",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(0, -3),
				ElementType: ui.TextElement,
				Anchor:      pixel.Right,
			},
			{
				Key:         "player_inv_cnt",
				Position:    pixel.V(40, 0),
				ElementType: ui.ContainerElement,
				Width:       world.TileSize,
				Height:      world.TileSize,
			},
			{
				Key:         "player_inv_item",
				SprKey:      "rock",
				Batch:       constants.TileBatch,
				Position:    pixel.V(40, 0),
				ElementType: ui.SpriteElement,
			},
			{
				Key:         "player_inv_item_2",
				SprKey:      "rock",
				Batch:       constants.TileBatch,
				Position:    pixel.V(40, 0),
				ElementType: ui.SpriteElement,
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
	//f, err := os.Create("assets/ui/player4.json")
	//if err != nil {
	//	panic(err)
	//}
	//bts, err := json.Marshal(Player4InvConstructor)
	//if err != nil {
	//	panic(err)
	//}
	//_, err = f.Write(bts)
	//if err != nil {
	//	panic(err)
	//}
	//f.Close()
}
