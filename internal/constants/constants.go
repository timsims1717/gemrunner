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
	BlockSelectWidth  = 8.
	BlockSelectHeight = 7.

	// Batches
	TileBatch = "tiles"
	UIBatch   = "ui"

	// Layers
	BlockSelectLayer = 250

	// UI Sprite Keys
	FancyBorderCorner   = "border_corner"
	FancyBorderStraight = "border_straight"
	ThinBorderWhite     = "white_dot"
	ThinBorderBlue      = "blue_dot"
	TextCaret           = "text_caret"

	// UI Constants
	ScrollSpeed = 150.

	// In Game Constants
	FrameRateMax  = 35
	FrameRateMin  = 5
	FrameRateInt  = 2
	NormalGravity = 5.
	MaxPlayers    = 4
	WaitToSwitch  = 3.
	ButtonBuffer  = 2

	// Player Constants
	PlayerWalkSpeed       = 4.
	PlayerLeapSpeed       = 3.8
	PlayerClimbSpeed      = 2.5
	PlayerSlideSpeed      = 4.5
	PlayerLeapDelay       = 1.5
	PlayerHighJumpSpeed   = 3.0
	PlayerHighJumpHSpeed  = 4.5
	PlayerHighJumpTimer   = 4.8
	PlayerHighJumpCounter = 6
	PlayerLongJumpVSpeed  = 1.8
	PlayerLongJumpHSpeed  = 4.5
	PlayerLongJumpTimer   = 6.2
	PlayerLongJumpCounter = 6
	IdleFrequency         = 10

	// Demon Constants
	DemonWalkSpeed       = 2.7
	DemonLeapSpeed       = 2.2
	DemonClimbSpeed      = 1.6
	DemonSlideSpeed      = 2.2
	DemonLeapDelay       = 1.5
	DemonGravity         = 3.
	DemonHighJumpSpeed   = 1.8
	DemonHighJumpHSpeed  = 3.5
	DemonHighJumpTimer   = 7.2
	DemonHighJumpCounter = 8
	DemonLongJumpVSpeed  = 1.6
	DemonLongJumpHSpeed  = 3.4
	DemonLongJumpTimer   = 7.8
	DemonLongJumpCounter = 8

	// Fly Constants
	FlySpeed = 3.
)

var (
	// Config
	System     string
	HomeDir    string
	ContentDir string
	PuzzlesDir string
	SavesDir   string
	ConfigFile string

	// Options
	FrameRate = 15

	// In Game Vars
	DrawingLayers = []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30}
)
