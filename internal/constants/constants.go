package constants

import (
	"gemrunner/internal/data/config"
	"github.com/gopxl/pixel"
)

const (
	Title            = "Gem Runner"
	Release          = 0
	Version          = 4
	Build            = 20250930
	ScreenRatioLimit = 0.8

	// Directories
	ContentDirName = "GemRunner/"
	LinuxDir       = "/.local/share/" + ContentDirName
	WinDir         = "/Documents/My Games/" + ContentDirName
	MacDir         = "/Library/Application Support/" + ContentDirName
	PuzzleDir      = "puzzles"
	SaveDir        = "saves"
	ReplayDir      = "replays"
	PuzzleExt      = ".puzzle"
	SaveExt        = ".savegame"
	ReplayPath     = "%s_%d_%s.replay"
	Favorites      = ".favorites"
	ConfigFilename = "config.toml"

	// World Constants
	TileSize        = 16.
	PuzzleWidth     = 28
	PuzzleHeight    = 16
	PuzzleMaxWidth  = 28
	PuzzleMaxHeight = 16
	PuzzleMinWidth  = 6
	PuzzleMinHeight = 6

	// Editor Constants
	BlockSelectWidth         = 8
	BlockSelectHeight        = 11
	RearrangeMoveDur         = 0.2
	UndoStackSize            = 50
	AdventureViewScrollSpeed = 60

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
	UIInfinity          = "tile_ui_infinity"
	UINumber            = "tile_ui_%d"
	UINumberX           = "tile_ui_%dx"
	UIRegenerate        = "tile_ui_regen"
	UIEnemy             = "tile_ui_enemy"
	UIShow              = "tile_ui_show"
	UIFlip              = "tile_ui_flip"
	UIDoor              = "tile_ui_door"
	UILock              = "tile_ui_lock"
	UIUnlock            = "tile_ui_unlock"

	// UI Constants
	ScrollSpeed = 500

	// In Game Constants
	FrameRateMax    = 90
	FrameRateMin    = 10
	FrameRateInt    = 5
	FrameCycle      = 16
	NormalGravity   = 2.5
	MaxPlayers      = 4
	WaitToSwitch    = 3
	ButtonBuffer    = 3
	CrackedCounter  = 16
	FangsCounter    = 40
	DrillCounter    = -8
	FlamethrowerCnt = 8
	GoopBucketCnt   = 7
	RegenCounter    = 176
	RegenACounter   = 7
	CollapseCounter = 7
	ThrownCounter   = 8
	ThrownVSpeed    = 0.9
	ThrownHSpeed    = 3.1
	SmashDistance   = TileSize * 3
	BombFuse        = 32
	ItemRegen       = 24
	SpeedMod        = 0.001
	TextTimer       = 16
	TextProxDist    = 48
	TextProxBuffer  = 48
	TextBobInterval = 16
	FakePlayerR     = 16
	PickUpGemChance = 50

	// Player Constants
	PlayerWalkSpeed       = 3
	PlayerGoopSpeed       = 1
	PlayerBarSpeed        = 2.75
	PlayerLeapSpeed       = 3
	PlayerClimbSpeed      = 1.75
	PlayerSlideSpeed      = 2.75
	PlayerGravity         = 4
	PlayerHighJumpSpeed   = 1.5
	PlayerHighJumpHSpeed  = 2.25
	PlayerHighJumpCounter = 12
	PlayerLongJumpVSpeed  = 1
	PlayerLongJumpHSpeed  = 3.75
	PlayerLongJumpCounter = 8
	IdleFrequency         = 10
	SmallBombInv          = 4

	// Demon Constants
	DemonWalkSpeed       = 1.4
	DemonGoopSpeed       = 0.5
	DemonBarSpeed        = 1.45
	DemonLeapSpeed       = 1.4
	DemonClimbSpeed      = 0.9
	DemonSlideSpeed      = 1.8
	DemonGravity         = NormalGravity
	DemonHighJumpSpeed   = 0.9
	DemonHighJumpHSpeed  = 1.75
	DemonHighJumpCounter = 16
	DemonLongJumpVSpeed  = 0.8
	DemonLongJumpHSpeed  = 1.7
	DemonLongJumpCounter = 16
	DemonInHoleCounter   = 56

	// Other Character Constants
	FlySpeed  = 1.4
	SlugSpeed = 1
)

var (
	// Config
	System        string
	Username      string
	HomeDir       string
	ContentDir    string
	PuzzlesDir    string
	SavesDir      string
	ReplaysDir    string
	ConfigFile    string
	Configuration config.Configuration

	// Options
	PickedRatio = 1.
	Resolutions = []pixel.Vec{
		//pixel.V(640, 480),
		pixel.V(1280, 720),
		pixel.V(1366, 768),
		pixel.V(1600, 900),
		pixel.V(1920, 1080),
		pixel.V(2560, 1440),
		pixel.V(3840, 2160),
	}
	Scanlines int32

	// In Game Vars
	DrawingLayers      = []int{9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30}
	CollapseRegenLayer = []int{31}         // turf that have masks
	CollapseRegenMask  = []int{32}         // turf masks
	EffectsLayer       = []int{33, 34, 35} // digging; explosions; deaths from explosions, text
	TextLayer          = []int{36, 37}     // floating text

	LevelTransSpeed = 0.25
)
