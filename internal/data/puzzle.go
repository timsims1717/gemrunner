package data

import (
	"gemrunner/internal/constants"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
)

var (
	PuzzleView *viewport.ViewPort
	BorderView *viewport.ViewPort
	IMDraw     *imdraw.IMDraw

	CurrLevel  *Level
	CurrPuzzle *Puzzle
	CurrSelect *Selection
	ClipSelect *Selection
	EditorDraw bool

	PuzzleShader string
)

type Level struct {
	Tiles     *Tiles
	Enemies   []*Dynamic
	Players   [constants.MaxPlayers]*Dynamic
	Stats     [constants.MaxPlayers]*PlayerStats
	PControls [constants.MaxPlayers]Controller
	Start     bool
	Failed    bool
	Complete  bool
	DoorsOpen bool

	Puzzle   *Puzzle
	Metadata *PuzzleMetadata
}

type Puzzle struct {
	Tiles *Tiles `json:"tiles"`

	WrenchTiles []*Tile `json:"-"`

	UndoStack  []*Tiles `json:"-"`
	LastChange *Tiles   `json:"-"`
	RedoStack  []*Tiles `json:"-"`

	Click   bool `json:"-"`
	Update  bool `json:"-"`
	Changed bool `json:"-"`

	Metadata *PuzzleMetadata `json:"metadata"`
}

type Tiles struct {
	T [constants.PuzzleHeight][constants.PuzzleWidth]*Tile
}

func (t *Tiles) Get(x, y int) *Tile {
	if x < 0 || y < 0 || x >= constants.PuzzleWidth || y >= constants.PuzzleHeight {
		return nil
	}
	return t.T[y][x]
}

func NewTiles() *Tiles {
	t := [constants.PuzzleHeight][constants.PuzzleWidth]*Tile{}
	for _, row := range t {
		for _, tile := range row {
			tile = &Tile{}
			tile.ToEmpty()
		}
	}
	return &Tiles{
		T: t,
	}
}

type Selection struct {
	Tiles  [][]*Tile
	Offset world.Coords
	Origin world.Coords
	Width  int
	Height int
}

func CreateBlankPuzzle() *Puzzle {
	worldNum := constants.WorldRock
	md := &PuzzleMetadata{
		Name:           "",
		Filename:       "",
		WorldSprite:    constants.WorldSprites[worldNum],
		WorldNumber:    worldNum,
		PrimaryColor:   pixel.ToRGBA(constants.WorldPrimary[worldNum]),
		SecondaryColor: pixel.ToRGBA(constants.WorldSecondary[worldNum]),
		Completed:      false,
	}
	puz := &Puzzle{
		Tiles:    NewTiles(),
		Metadata: md,
	}
	for y := 0; y < constants.PuzzleHeight; y++ {
		for x := 0; x < constants.PuzzleWidth; x++ {
			puz.Tiles.T[y][x] = &Tile{
				Block:  Block(BlockEmpty),
				Coords: world.Coords{X: x, Y: y},
			}
		}
	}
	return puz
}

func (p *Puzzle) CopyTiles() *Tiles {
	tiles := NewTiles()
	for y := 0; y < constants.PuzzleHeight; y++ {
		for x := 0; x < constants.PuzzleWidth; x++ {
			tiles.T[y][x] = p.Tiles.T[y][x].Copy()
		}
	}
	return tiles
}

func (p *Puzzle) HasPlayers() bool {
	hasPlayers := false
playerCheck:
	for _, row := range CurrPuzzle.Tiles.T {
		for _, t := range row {
			if t.Block == BlockPlayer1 ||
				t.Block == BlockPlayer2 ||
				t.Block == BlockPlayer3 ||
				t.Block == BlockPlayer4 {
				hasPlayers = true
				break playerCheck
			}
		}
	}
	return hasPlayers
}
