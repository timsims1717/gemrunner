package constants

// Main Menu Dialogs
const (
	DialogMainMenu   = "main_menu"
	DialogAddPlayers = "add_players"
	DialogPlayLocal  = "play_local"

	DialogOptions = "options_menu"

	DialogPauseMenu      = "pause_menu"
	DialogPlayer1Inv     = "player1_inv"
	DialogPlayer2Inv     = "player2_inv"
	DialogPlayer3Inv     = "player3_inv"
	DialogPlayer4Inv     = "player4_inv"
	DialogPuzzleTitle    = "puzzle_title"
	DialogPuzzleTimer    = "puzzle_timer"
	DialogAdventureTrans = "adventure_set_transition"
)

// Editor Dialogs
const (
	DialogOpenPuzzle            = "open_puzzle"
	DialogAddPuzzle             = "add_puzzle"
	DialogChangeName            = "change_name"
	DialogNoPlayersInPuzzle     = "no_players"
	DialogUnableToSave          = "unable_to_save"
	DialogUnableToSaveConfirm   = "unable_to_save_confirm"
	DialogChangeWorld           = "change_world"
	DialogAreYouSureDelete      = "are_you_sure_delete"
	DialogPuzzleSettings        = "puzzle_settings"
	DialogPuzzleSetSettings     = "puzzle_set_settings"
	DialogCombineSets           = "combine_sets"
	DialogRearrangePuzzleSet    = "rearrange_puzzle_set"
	DialogRearrangeAdventureSet = "rearrange_adventure_set"

	DialogEditorPanelTop     = "editor_panel_top"
	DialogEditorPanelLeft    = "editor_panel_left"
	DialogEditorOptionsBot   = "editor_options_bot"
	DialogEditorOptionsRight = "editor_options_right"
	DialogEditorBlockSelect  = "block_select"

	DialogBarrier      = "barrier_options"
	DialogCrackedTiles = "cracked_tile_options"
	DialogBomb         = "bomb_options"
	DialogItemOptions  = "item_options"
	DialogFloatingText = "floating_text"
	DialogDoors        = "door_options"
	DialogPalette      = "palette_options"
)

var DialogKeys = []string{
	DialogMainMenu,
	DialogOptions,

	DialogPlayer1Inv,
	DialogPlayer2Inv,
	DialogPlayer3Inv,
	DialogPlayer4Inv,
	DialogPuzzleTitle,
	DialogPuzzleTimer,
	DialogAdventureTrans,

	DialogPlayLocal,

	DialogEditorPanelLeft,
	DialogEditorOptionsRight,
	DialogPuzzleSettings,
	DialogPuzzleSetSettings,
	DialogRearrangePuzzleSet,
	DialogRearrangeAdventureSet,
	DialogNoPlayersInPuzzle,
	DialogAddPuzzle,

	DialogBarrier,
	DialogItemOptions,
	DialogFloatingText,
	DialogDoors,
	DialogPalette,
}

var MainDialogs = []string{
	DialogMainMenu,
	DialogOptions,
	DialogAddPlayers,
	DialogPlayLocal,
}

var InGameDialogs = []string{
	DialogPauseMenu,
	DialogPlayer1Inv,
	DialogPlayer2Inv,
	DialogPlayer3Inv,
	DialogPlayer4Inv,
	DialogPuzzleTitle,
	DialogPuzzleTimer,
}

var EditorDialogs = []string{
	DialogOpenPuzzle,
	DialogAddPuzzle,
	DialogChangeName,
	DialogPuzzleSettings,
	DialogPuzzleSetSettings,
	DialogNoPlayersInPuzzle,
	DialogAreYouSureDelete,
	DialogUnableToSave,
	DialogUnableToSaveConfirm,
	DialogCombineSets,
	DialogBarrier,
	DialogCrackedTiles,
	DialogBomb,
	DialogItemOptions,
	DialogPalette,
	DialogFloatingText,
	DialogEditorPanelLeft,
	DialogEditorPanelTop,
	DialogEditorOptionsRight,
	DialogEditorOptionsBot,
	DialogEditorBlockSelect,
}
