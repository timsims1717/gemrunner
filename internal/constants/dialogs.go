package constants

// Main Menu Dialogs
const (
	DialogMainMenu   = "main_menu"
	DialogAddPlayers = "add_players"
	DialogPlayLocal  = "play_local"

	DialogPauseMenu   = "pause_menu"
	DialogPlayer1Inv  = "player1_inv"
	DialogPlayer2Inv  = "player2_inv"
	DialogPlayer3Inv  = "player3_inv"
	DialogPlayer4Inv  = "player4_inv"
	DialogPuzzleTitle = "puzzle_title"
)

// Editor Dialogs
const (
	DialogOpenPuzzle          = "open_puzzle"
	DialogChangeName          = "change_name"
	DialogNoPlayersInPuzzle   = "no_players"
	DialogUnableToSave        = "unable_to_save"
	DialogUnableToSaveConfirm = "unable_to_save_confirm"
	DialogChangeWorld         = "change_world"
	DialogAreYouSureDelete    = "are_you_sure_delete"
	DialogPuzzleSettings      = "puzzle_settings"
	DialogCombineSets         = "combine_sets"
	DialogRearrangePuzzleSet  = "rearrange_puzzle_set"

	DialogEditorPanelTop     = "editor_panel_top"
	DialogEditorPanelLeft    = "editor_panel_left"
	DialogEditorOptionsBot   = "editor_options_bot"
	DialogEditorOptionsRight = "editor_options_right"
	DialogEditorBlockSelect  = "block_select"

	DialogCrackedTiles = "cracked_tile_options"
	DialogBomb         = "bomb_options"
	DialogJetpack      = "jetpack_options"
	DialogFloatingText = "floating_text"
	DialogDisguise     = "disguise_options"
)

var DialogKeys = []string{
	DialogJetpack,
	DialogFloatingText,
	DialogDisguise,
}
