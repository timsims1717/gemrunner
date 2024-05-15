package load

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/ui"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

var (
	mainMenuConstructor   *ui.DialogConstructor
	addPlayersConstructor *ui.DialogConstructor
	playLocalConstructor  *ui.DialogConstructor
	puzzleListEntry       ui.ElementConstructor
	puzzleTitleItem       ui.ElementConstructor
	puzzleTitleShadowItem ui.ElementConstructor
	puzzleNumberSpr       ui.ElementConstructor
	playerNumberSpr       ui.ElementConstructor
	numberItems           ui.ElementConstructor
	favoriteItem          ui.ElementConstructor
)

func InitMainMenuConstructors() {
	mainMenuConstructor = &ui.DialogConstructor{
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
	addPlayersConstructor = &ui.DialogConstructor{
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
			},
			{
				Key:         "any_button_p3",
				Text:        "Press\na key",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(20, 16),
				ElementType: ui.TextElement,
			},
			{
				Key:         "any_button_p4",
				Text:        "Press\na key",
				Color:       pixel.ToRGBA(constants.ColorWhite),
				Position:    pixel.V(88, 16),
				ElementType: ui.TextElement,
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
	playLocalConstructor = &ui.DialogConstructor{
		Key:    constants.DialogPlayLocal,
		Width:  11,
		Height: 12,
		Elements: []ui.ElementConstructor{
			{
				Key:         "play_main_tab",
				Position:    pixel.V(-40, 79),
				ElementType: ui.ContainerElement,
				Width:       4 * world.TileSize,
				Height:      2 * world.TileSize,
				SubElements: []ui.ElementConstructor{
					{
						Key:         "main_tab_text_shadow",
						Text:        "Main Story",
						Color:       pixel.ToRGBA(constants.ColorBlue),
						Position:    pixel.V(-31, 1),
						ElementType: ui.TextElement,
					},
					{
						Key:         "main_tab_text",
						Text:        "Main Story",
						Color:       pixel.ToRGBA(constants.ColorWhite),
						Position:    pixel.V(-31, 0),
						ElementType: ui.TextElement,
					},
				},
			},
			{
				Key:         "play_custom_tab",
				Position:    pixel.V(40, 79),
				ElementType: ui.ContainerElement,
				Width:       4 * world.TileSize,
				Height:      2 * world.TileSize,
				SubElements: []ui.ElementConstructor{
					{
						Key:         "custom_tab_text_shadow",
						Text:        "Custom",
						Color:       pixel.ToRGBA(constants.ColorBlue),
						Position:    pixel.V(-18, 1),
						ElementType: ui.TextElement,
					},
					{
						Key:         "custom_tab_text",
						Text:        "Custom",
						Color:       pixel.ToRGBA(constants.ColorWhite),
						Position:    pixel.V(-18, 0),
						ElementType: ui.TextElement,
					},
				},
			},
			//{
			//	Key:         "play_continue_local",
			//	SprKey:      "play_continue_btn_big",
			//	SprKey2:     "play_continue_btn_click_big",
			//	Batch:       constants.UIBatch,
			//	HelpText:    "Start where you left off",
			//	Position:    pixel.V(36, -84),
			//	ElementType: ui.ButtonElement,
			//},
			{
				Key:         "play_new_local",
				SprKey:      "play_btn_big",
				SprKey2:     "play_btn_click_big",
				Batch:       constants.UIBatch,
				HelpText:    "Start at the Beginning",
				Position:    pixel.V(56, -84),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "cancel_play_local",
				SprKey:      "cancel_btn_big",
				SprKey2:     "cancel_btn_click_big",
				Batch:       constants.UIBatch,
				HelpText:    "Cancel",
				Position:    pixel.V(76, -84),
				ElementType: ui.ButtonElement,
			},
			{
				Key:         "custom_puzzle_list",
				Batch:       constants.UIBatch,
				ElementType: ui.ScrollElement,
				Position:    pixel.V(0, -7),
				Width:       10 * world.TileSize,
				Height:      8 * world.TileSize,
			},
		},
	}
	puzzleListEntry = ui.ElementConstructor{
		Key:         "puzzle_container_%d",
		HelpText:    "Select this puzzle.",
		ElementType: ui.ContainerElement,
		Width:       9 * world.TileSize,
		Height:      2 * world.TileSize,
	}
	puzzleTitleShadowItem = ui.ElementConstructor{
		Key:         "puzzle_title_shadow",
		Position:    pixel.V(-18, 9),
		Color:       pixel.ToRGBA(constants.ColorBlue),
		ElementType: ui.TextElement,
	}
	puzzleTitleItem = ui.ElementConstructor{
		Key:         "puzzle_title",
		Position:    pixel.V(-18, 8),
		Color:       pixel.ToRGBA(constants.ColorWhite),
		ElementType: ui.TextElement,
	}
	puzzleNumberSpr = ui.ElementConstructor{
		Key:         "puzzle_num",
		SprKey:      "puzzle_symbol",
		Position:    pixel.V(-32, 8),
		Batch:       constants.UIBatch,
		ElementType: ui.SpriteElement,
	}
	playerNumberSpr = ui.ElementConstructor{
		Key:         "player_num",
		SprKey:      "player_symbol",
		Position:    pixel.V(-48, 8),
		Batch:       constants.UIBatch,
		ElementType: ui.SpriteElement,
	}
	numberItems = ui.ElementConstructor{
		Key:         "number_text",
		Position:    pixel.V(-50, -4),
		Color:       pixel.ToRGBA(constants.ColorWhite),
		ElementType: ui.TextElement,
	}
	favoriteItem = ui.ElementConstructor{
		Key:         "favorite_symbol",
		SprKey:      "heart_empty",
		SprKey2:     "heart_full",
		Position:    pixel.V(-64, 0),
		Batch:       constants.UIBatch,
		ElementType: ui.CheckboxElement,
	}
}
