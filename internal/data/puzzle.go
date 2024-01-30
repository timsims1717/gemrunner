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

	CurrPuzzle *Puzzle
	CurrSelect *Selection
	ClipSelect *Selection

	PuzzleShader string
)

type Puzzle struct {
	Tiles      *Tiles   `json:"tiles"`
	UndoStack  []*Tiles `json:"-"`
	LastChange *Tiles   `json:"-"`
	RedoStack  []*Tiles `json:"-"`

	Click  bool `json:"-"`
	Update bool `json:"-"`

	WorldSprite    string     `json:"sprite"`
	WorldNumber    int        `json:"world"`
	PrimaryColor   pixel.RGBA `json:"primaryColor"`
	SecondaryColor pixel.RGBA `json:"secondaryColor"`

	Title    string
	Filename string
}

type Tiles struct {
	T [constants.PuzzleHeight][constants.PuzzleWidth]*Tile
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
	puz := &Puzzle{
		Tiles:          NewTiles(),
		WorldNumber:    worldNum,
		WorldSprite:    constants.WorldSprites[worldNum],
		PrimaryColor:   pixel.ToRGBA(constants.WorldPrimary[worldNum]),
		SecondaryColor: pixel.ToRGBA(constants.WorldSecondary[worldNum]),
	}
	for y := 0; y < constants.PuzzleHeight; y++ {
		for x := 0; x < constants.PuzzleWidth; x++ {
			puz.Tiles.T[y][x] = &Tile{
				Block:  Block(Empty),
				Coords: world.Coords{X: x, Y: y},
			}
		}
	}
	return puz
}

func CreateTestPuzzle() *Puzzle {
	puz := &Puzzle{
		Tiles: NewTiles(),
	}
	for y := 0; y < constants.PuzzleHeight; y++ {
		for x := 0; x < constants.PuzzleWidth; x++ {
			block := Empty
			if (x+y)%2 == 0 {
				block = Turf
			}
			puz.Tiles.T[y][x] = &Tile{
				Block:  Block(block),
				Ladder: false,
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
