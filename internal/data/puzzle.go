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

	PuzzleShader string
)

type Level struct {
	Tiles    *Tiles
	Chars    []*Dynamic
	Players  [constants.MaxPlayers]*Dynamic
	Stats    [constants.MaxPlayers]*PlayerStats
	Start    bool
	Failed   bool
	Complete bool

	Puzzle   *Puzzle
	Metadata *PuzzleMetadata
}

type Puzzle struct {
	Tiles *Tiles `json:"tiles"`

	UndoStack  []*Tiles `json:"-"`
	LastChange *Tiles   `json:"-"`
	RedoStack  []*Tiles `json:"-"`

	Click  bool `json:"-"`
	Update bool `json:"-"`

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
	return &Tiles{
		T: [constants.PuzzleHeight][constants.PuzzleWidth]*Tile{},
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
		WorldSprite:    constants.WorldSprites[worldNum],
		WorldNumber:    worldNum,
		PrimaryColor:   pixel.ToRGBA(constants.WorldPrimary[worldNum]),
		SecondaryColor: pixel.ToRGBA(constants.WorldSecondary[worldNum]),
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
