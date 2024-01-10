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
	Tiles [constants.PuzzleHeight][constants.PuzzleWidth]*Tile `json:"tiles"`

	Click  bool   `json:"-"`
	Update bool   `json:"-"`
	World  string `json:"world"`

	PrimaryColor   pixel.RGBA `json:"primaryColor"`
	SecondaryColor pixel.RGBA `json:"secondaryColor"`

	Filename string
}

type Selection struct {
	Tiles  [][]*Tile
	Offset world.Coords
	Origin world.Coords
	Width  int
	Height int
}

func CreateBlankPuzzle() *Puzzle {
	puz := &Puzzle{
		Tiles:          [constants.PuzzleHeight][constants.PuzzleWidth]*Tile{},
		World:          constants.WorldRock,
		PrimaryColor:   pixel.ToRGBA(constants.ColorGray),
		SecondaryColor: pixel.ToRGBA(constants.ColorGreen),
	}
	for y := 0; y < constants.PuzzleHeight; y++ {
		for x := 0; x < constants.PuzzleWidth; x++ {
			puz.Tiles[y][x] = &Tile{
				Block:  Block(Empty),
				Coords: world.Coords{X: x, Y: y},
			}
		}
	}
	return puz
}

func CreateTestPuzzle() *Puzzle {
	puz := &Puzzle{
		Tiles: [constants.PuzzleHeight][constants.PuzzleWidth]*Tile{},
	}
	for y := 0; y < constants.PuzzleHeight; y++ {
		for x := 0; x < constants.PuzzleWidth; x++ {
			block := Empty
			if (x+y)%2 == 0 {
				block = Turf
			}
			puz.Tiles[y][x] = &Tile{
				Block:  Block(block),
				Ladder: false,
				Coords: world.Coords{X: x, Y: y},
			}
		}
	}
	return puz
}

func (p *Puzzle) Copy() *Puzzle {
	puz := &Puzzle{
		Tiles: [constants.PuzzleHeight][constants.PuzzleWidth]*Tile{},
	}
	for y := 0; y < constants.PuzzleHeight; y++ {
		for x := 0; x < constants.PuzzleWidth; x++ {
			puz.Tiles[y][x] = p.Tiles[y][x].Copy()
		}
	}
	puz.Update = true
	return puz
}
