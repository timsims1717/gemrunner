package data

import (
	"gemrunner/internal/constants"
	"gemrunner/pkg/world"
)

type Block int

const (
	RedRock = iota
	RedStone
	RedBrick
	RedDirt
	Empty
)

func (b Block) String() string {
	switch b {
	case RedRock:
		return "red_rock"
	case RedStone:
		return "red_stone"
	case RedBrick:
		return "red_brick"
	case RedDirt:
		return "red_dirt"
	}
	return "empty"
}

type Tile struct {
	Block  Block
	Doodad Block
	Object Block
	Ladder bool
	Fall   bool
	Coords world.Coords
}

var (
	CurrPuzzle *Puzzle
)

type Puzzle struct {
	Tiles [constants.PuzzleHeight][constants.PuzzleWidth]*Tile
}

func CreateBlankPuzzle() *Puzzle {
	puz := &Puzzle{
		Tiles: [constants.PuzzleHeight][constants.PuzzleWidth]*Tile{},
	}
	for y := 0; y < constants.PuzzleHeight; y++ {
		for x := 0; x < constants.PuzzleWidth; x++ {
			puz.Tiles[y][x] = &Tile{
				Ladder: false,
				Fall:   false,
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
				block = RedRock
			}
			puz.Tiles[y][x] = &Tile{
				Block:  Block(block),
				Ladder: false,
				Fall:   false,
				Coords: world.Coords{X: x, Y: y},
			}
		}
	}
	return puz
}
