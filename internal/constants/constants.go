package constants

const (
	Title   = "Gem Runner"
	Release = 0
	Version = 1
	Build   = 20230125

	// Directories
	LinuxDir  = "/.local/share/GemRunnerEditor"
	WinDir    = "/Documents/My Games/GemRunnerEditor"
	MacDir    = "/Library/Application Support/GemRunnerEditor"
	PuzzleDir = "/puzzles"
	SaveDir   = "/saves"
	PuzzleExt = ".puzzle"
	SaveExt   = ".savegame"

	// World Constants
	TileSize     = 16.
	PuzzleWidth  = 28
	PuzzleHeight = 16

	// Editor Constants
	BlockSelectWidth = 6.

	// Batches
	BGBatch = "tile_bg"
	FGBatch = "tile_fg"
	UIBatch = "ui"

	// Layers
	BlockSelectLayer = 250

	// UI Sprite Keys
	FancyBorderCorner   = "border_corner"
	FancyBorderStraight = "border_straight"
	ThinBorderWhite     = "white_dot"
	ThinBorderBlue      = "blue_dot"
	TextCaret           = "text_caret"
)

var (
	// Config
	System     string
	HomeDir    string
	ContentDir string
	PuzzlesDir string
	SavesDir   string
	ConfigFile string
)
