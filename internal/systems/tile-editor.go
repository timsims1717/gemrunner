package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/pkg/world"
)

func SetBlock(coords world.Coords, block data.Block) {
	if data.CurrPuzzle != nil {
		if CoordsLegal(coords) {
			// add to puzzle
			tile := data.CurrPuzzle.Tiles.T[coords.Y][coords.X]
			if !tile.Update {
				switch block {
				case data.BlockLadder, data.BlockLadderExit, data.BlockLadderCracked:
					if tile.Ladder == block || tile.Block != data.BlockTurf {
						tile.Block = data.BlockEmpty
					}
				case data.BlockPlayer1, data.BlockPlayer2, data.BlockPlayer3, data.BlockPlayer4:
					tile.Ladder = data.BlockEmpty
					// ensure no other player of that type are in puzzle
					for _, row := range data.CurrPuzzle.Tiles.T {
						for _, t := range row {
							if t.Block == block {
								t.Block = data.BlockEmpty
							}
						}
					}
				case data.BlockDemonRegen:
					for _, row := range data.CurrPuzzle.Tiles.T {
						for _, t := range row {
							if t.Block == data.BlockDemon &&
								t.Metadata.Regenerate &&
								!t.Metadata.Changed {
								t.Metadata.RegenTiles = append(t.Metadata.RegenTiles, coords)
							}
						}
					}
				default:
					if tile.IsLadder() && !(block == data.BlockTurf && tile.Block == data.BlockEmpty) {
						tile.Ladder = data.BlockEmpty
					}
				}
				tile.Block = block
			}
			data.CurrPuzzle.Update = true
			tile.Update = true
			data.CurrPuzzle.Changed = true
		}
	} else {
		fmt.Println("error: attempted to change tile when no puzzle is loaded")
	}
}

func DeleteBlock(coords world.Coords) {
	if data.CurrPuzzle != nil {
		if CoordsLegal(coords) {
			tile := data.CurrPuzzle.Tiles.T[coords.Y][coords.X]
			if !tile.Update {
				if tile.IsLadder() {
					tile.Ladder = data.BlockEmpty
				} else if tile.Block != data.BlockEmpty {
					tile.Block = data.BlockEmpty
				}
				data.CurrPuzzle.Update = true
				tile.Update = true
				data.CurrPuzzle.Changed = true
			}
		}
	} else {
		fmt.Println("error: attempted to delete tile when no puzzle is loaded")
	}
}

func CoordsLegal(coords world.Coords) bool {
	return coords.X >= 0 && coords.Y >= 0 && coords.X < constants.PuzzleWidth && coords.Y < constants.PuzzleHeight
}

func GetClosestLegal(coords world.Coords) world.Coords {
	if CoordsLegal(coords) {
		return coords
	}
	if coords.X < 0 {
		coords.X = 0
	} else if coords.X >= constants.PuzzleWidth {
		coords.X = constants.PuzzleWidth - 1
	}
	if coords.Y < 0 {
		coords.Y = 0
	} else if coords.Y >= constants.PuzzleHeight {
		coords.Y = constants.PuzzleHeight - 1
	}
	return coords
}

func CoordsLegalSelection(coords world.Coords) bool {
	return coords.X >= 0 && coords.Y >= 0 && coords.X < data.CurrSelect.Width && coords.Y < data.CurrSelect.Height
}
