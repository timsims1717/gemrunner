package data

import (
	"gemrunner/internal/constants"
	"github.com/google/uuid"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
	"time"
)

var (
	PuzzleSetFileList      []PuzzleSetMetadata
	PuzzleSetSortedList    []PuzzleSetMetadata
	CustomPuzzleListLoaded bool
	FavoritesList          []string
	SelectedPuzzleIndex    int

	CustomWorldSelected    bool
	CustomSelectedBefore   bool
	SelectedWorldIndex     int
	SelectedPrimaryColor   pixel.RGBA
	SelectedSecondaryColor pixel.RGBA
	SelectedDoodadColor    pixel.RGBA
	SelectedTextColor      pixel.RGBA
	SelectedShadowColor    pixel.RGBA

	PuzzleSetViewIsMoving bool
	PuzzleSetViewIndex    int
	PuzzleSetViewPuzzles  []int
	PuzzleSetViewAllowEnd bool
	RearrangeLeftX        = constants.TileSize * -8
	RearrangeRightX       = constants.TileSize * 8
	RearrangeFloatX       = constants.TileSize * -16

	Players       []Player
	MenuInputUsed = pxginput.KeyboardMouse
	MainJoystick  = -1
)

type PuzzleMetadata struct {
	Name                 string     `json:"title"`
	Author               string     `json:"author"`
	Filename             string     `json:"filename"`
	Width                int        `json:"width"`
	Height               int        `json:"height"`
	WorldSprite          string     `json:"sprite"`
	WorldNumber          int        `json:"world"`
	WorldLiquid          string     `json:"liquid"`
	PrimaryColor         pixel.RGBA `json:"primaryColor"`
	SecondaryColor       pixel.RGBA `json:"secondaryColor"`
	DoodadColor          pixel.RGBA `json:"doodadColor"`
	LiquidPrimaryColor   pixel.RGBA `json:"liquidColor1"`
	LiquidSecondaryColor pixel.RGBA `json:"liquidColor2"`
	ShaderMode           int        `json:"shader"`
	ShaderSpeed          float32    `json:"-"`
	ShaderY              float32    `json:"-"`
	ShaderX              float32    `json:"-"`
	ShaderCustom         float32    `json:"-"`
	ParticleMode         int        `json:"particles"`
	MusicTrack           string     `json:"musicTrack"`
	HubLevel             bool       `json:"hubLevel"`
	SecretLevel          bool       `json:"secretLevel"`
	Darkness             bool       `json:"darkness"`
	Completed            bool       `json:"completed"`
}

type PuzzleSetMetadata struct {
	UUID       *uuid.UUID `json:"uuid"`
	Name       string     `json:"title"`
	Filename   string     `json:"filename"`
	Author     string     `json:"author"`
	Adventure  bool       `json:"adventure"`
	NumPlayers int        `json:"numPlayers"`
	NumPuzzles int        `json:"numPuzzles"`
	Favorite   bool       `json:"favorite"`
	Publish    bool       `json:"publish"`
	Local      bool       `json:"local"`
	Online     bool       `json:"online"`
	RecentPlay *time.Time `json:"recentPlay"`
	RecordDate *time.Time `json:"recordDate"`
	Downloads  int        `json:"downloads"`
	Favorites  int        `json:"favorites"`
	Desc       string     `json:"desc"`
}

type Player struct {
	PlayerNum int
	Keyboard  bool
	Gamepad   pixelgl.Joystick
}
