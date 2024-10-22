package constants

const (
	Title            = "Gem Runner"
	Release          = 0
	Version          = 1
	Build            = 20240322
	ScreenRatioLimit = 0.8

	// Directories
	LinuxDir  = "/.local/share/GemRunnerEditor"
	WinDir    = "/Documents/My Games/GemRunnerEditor"
	MacDir    = "/Library/Application Support/GemRunnerEditor"
	PuzzleDir = "/puzzles"
	SaveDir   = "/saves"
	PuzzleExt = ".puzzle"
	SaveExt   = ".savegame"
	Favorites = ".favorites"

	// World Constants
	TileSize     = 16.
	PuzzleWidth  = 28
	PuzzleHeight = 16

	// Editor Constants
	BlockSelectWidth  = 8.
	BlockSelectHeight = 11.
	RearrangeMoveDur  = 0.2

	// Music Tracks
	TrackBeach   = "beach"
	TrackCrystal = "crystal"
	TrackDark    = "dark"
	TrackDesert  = "desert"
	TrackFungus  = "fungus"
	TrackGilded  = "gilded"
	TrackIce     = "ice"
	TrackJungle  = "jungle"
	TrackLava    = "lava"
	TrackMain    = "main_theme"
	TrackMech    = "mech"
	TrackMenu    = "menu"
	TrackReef    = "reef"
	TrackUrban   = "urban"

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
	UINumber            = "tile_ui_%d"
	UINumberX           = "tile_ui_%dx"

	// UI Constants
	ScrollSpeed = 250.

	// In Game Constants
	FrameRateMax    = 60
	FrameRateMin    = 5
	FrameRateInt    = 5
	FrameCycle      = 8
	NormalGravity   = 2
	MaxPlayers      = 4
	WaitToSwitch    = 3
	ButtonBuffer    = 3
	CrackedCounter  = 16
	RegenCounter    = 128
	RegenACounter   = 7
	CollapseCounter = 7
	ThrownCounter   = 8
	ThrownVSpeed    = 0.9
	ThrownHSpeed    = 3.1
	SmashDistance   = TileSize * 3
	BombFuse        = 24
	ItemRegen       = 16
	SpeedMod        = 0.001
	TextTimer       = 16
	TextProxDist    = 48.
	TextProxBuffer  = 48
	BobInterval     = 16

	// Player Constants
	PlayerWalkSpeed       = 4
	PlayerBarSpeed        = 3.5
	PlayerLeapSpeed       = 4
	PlayerClimbSpeed      = 2
	PlayerSlideSpeed      = 3.5
	PlayerGravity         = 3.25
	PlayerHighJumpSpeed   = 1.5
	PlayerHighJumpHSpeed  = 2.25
	PlayerHighJumpCounter = 12
	PlayerLongJumpVSpeed  = 1
	PlayerLongJumpHSpeed  = 3.75
	PlayerLongJumpCounter = 8
	IdleFrequency         = 10

	// Demon Constants
	DemonWalkSpeed       = 2
	DemonBarSpeed        = 1.5
	DemonLeapSpeed       = 1.5
	DemonClimbSpeed      = 1.25
	DemonSlideSpeed      = 1.5
	DemonGravity         = NormalGravity
	DemonHighJumpSpeed   = 0.9
	DemonHighJumpHSpeed  = 1.75
	DemonHighJumpCounter = 16
	DemonLongJumpVSpeed  = 0.8
	DemonLongJumpHSpeed  = 1.7
	DemonLongJumpCounter = 16

	// Fly Constants
	FlySpeed = 2
)

var (
	// Config
	System     string
	Username   string
	HomeDir    string
	ContentDir string
	PuzzlesDir string
	SavesDir   string
	ConfigFile string

	// Options
	WinWidth    = 1600.
	WinHeight   = 900.
	FrameRate   = 30
	PickedRatio = 1.

	// In Game Vars
	DrawingLayers      = []int{9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30}
	CollapseRegenLayer = []int{31}         // turf that have masks
	CollapseRegenMask  = []int{32}         // turf masks
	EffectsLayer       = []int{33, 34, 35} // digging; explosions; deaths from explosions, text
	TextLayer          = []int{36, 37}     // floating text
)
