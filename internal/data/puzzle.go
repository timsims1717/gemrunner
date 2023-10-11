package data

import (
	"gemrunner/internal/constants"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/faiface/pixel/imdraw"
)

var (
	PuzzleView *viewport.ViewPort
	BorderView *viewport.ViewPort
	IMDraw     *imdraw.IMDraw

	CurrPuzzle *Puzzle
	CurrSelect *Selection
	ClipSelect *Selection
)

type Puzzle struct {
	Tiles  [constants.PuzzleHeight][constants.PuzzleWidth]*Tile
	Click  bool
	Update bool
	World  string
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
		Tiles: [constants.PuzzleHeight][constants.PuzzleWidth]*Tile{},
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
