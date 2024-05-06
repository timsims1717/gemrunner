package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/pkg/world"
)

func SetBlock(coords world.Coords, block data.Block) {
	if data.CurrPuzzleSet.CurrPuzzle != nil {
		if CoordsLegal(coords) {
			// add to puzzle
			tile := data.CurrPuzzleSet.CurrPuzzle.Tiles.Get(coords.X, coords.Y)
			if !tile.Update {
				if tile.Block != block {
					tile.Metadata = data.DefaultMetadata()
					RemoveLinkedTiles(tile)
				}
				switch block {
				case data.BlockTurf:
					if tile.IsLadder() {
						switch tile.Block {
						case data.BlockLadder:
							tile.Block = data.BlockLadderTurf
						case data.BlockLadderExit:
							tile.Block = data.BlockLadderExitTurf
						case data.BlockLadderCracked:
							tile.Block = data.BlockLadderCrackedTurf
						default:
							tile.Block = block
						}
					} else {
						tile.Block = block
					}
				case data.BlockPhase:
					rt := fmt.Sprintf("%d", tile.Metadata.Phase)
					tile.WrenchTxt.SetText(rt)
					tile.Block = block
				case data.BlockLadder, data.BlockLadderExit, data.BlockLadderCracked:
					if tile.Block == data.BlockTurf ||
						(tile.Block == data.BlockLadderTurf && block != data.BlockLadder) ||
						(tile.Block == data.BlockLadderCrackedTurf && block != data.BlockLadderCracked) ||
						(tile.Block == data.BlockLadderExitTurf && block != data.BlockLadderExit) {
						switch block {
						case data.BlockLadder:
							tile.Block = data.BlockLadderTurf
						case data.BlockLadderExit:
							tile.Block = data.BlockLadderExitTurf
						case data.BlockLadderCracked:
							tile.Block = data.BlockLadderCrackedTurf
						}
					} else {
						tile.Block = block
					}
				case data.BlockPlayer1, data.BlockPlayer2, data.BlockPlayer3, data.BlockPlayer4:
					// ensure no other player of that type are in puzzle
					for _, row := range data.CurrPuzzleSet.CurrPuzzle.Tiles.T {
						for _, t := range row {
							if t.Block == block {
								t.Block = data.BlockEmpty
							}
						}
					}
					tile.Block = block
				case data.BlockDemonRegen:
					for _, row := range data.CurrPuzzleSet.CurrPuzzle.Tiles.T {
						for _, t := range row {
							if t.Block == data.BlockDemon &&
								t.Metadata.Regenerate &&
								!t.Metadata.Changed {
								LinkTiles(tile, t)
							}
						}
					}
					tile.Block = block
				case data.BlockDemon:
					for _, row := range data.CurrPuzzleSet.CurrPuzzle.Tiles.T {
						for _, t := range row {
							if t.Block == data.BlockDemonRegen &&
								!t.Metadata.Changed {
								LinkTiles(tile, t)
							}
						}
					}
					tile.Block = block
				case data.BlockFlyRegen:
					for _, row := range data.CurrPuzzleSet.CurrPuzzle.Tiles.T {
						for _, t := range row {
							if t.Block == data.BlockFly &&
								t.Metadata.Regenerate &&
								!t.Metadata.Changed {
								LinkTiles(tile, t)
							}
						}
					}
					tile.Block = block
				case data.BlockFly:
					for _, row := range data.CurrPuzzleSet.CurrPuzzle.Tiles.T {
						for _, t := range row {
							if t.Block == data.BlockFlyRegen &&
								!t.Metadata.Changed {
								LinkTiles(tile, t)
							}
						}
					}
					tile.Block = block
				default:
					tile.Block = block
				}
			}
			data.CurrPuzzleSet.CurrPuzzle.Update = true
			tile.Update = true
			data.CurrPuzzleSet.CurrPuzzle.Changed = true
		}
	} else {
		fmt.Println("error: attempted to change tile when no puzzle is loaded")
	}
}

func DeleteBlock(coords world.Coords) {
	if data.CurrPuzzleSet.CurrPuzzle != nil {
		if CoordsLegal(coords) {
			tile := data.CurrPuzzleSet.CurrPuzzle.Tiles.Get(coords.X, coords.Y)
			if !tile.Update {
				if tile.IsLadder() {
					switch tile.Block {
					case data.BlockLadderTurf, data.BlockLadderCrackedTurf, data.BlockLadderExitTurf:
						tile.Block = data.BlockTurf
					default:
						tile.Block = data.BlockEmpty
					}
				} else if tile.Block != data.BlockEmpty {
					tile.Block = data.BlockEmpty
				}
				data.CurrPuzzleSet.CurrPuzzle.Update = true
				tile.Update = true
				data.CurrPuzzleSet.CurrPuzzle.Changed = true
				RemoveLinkedTiles(tile)
				tile.Metadata = data.DefaultMetadata()
			}
		}
	} else {
		fmt.Println("error: attempted to delete tile when no puzzle is loaded")
	}
}

func LinkTiles(tileA, tileB *data.Tile) {
	if !world.CoordsIn(tileB.Coords, tileA.Metadata.LinkedTiles) {
		tileA.Metadata.LinkedTiles = append(tileA.Metadata.LinkedTiles, tileB.Coords)
	}
	if !world.CoordsIn(tileA.Coords, tileB.Metadata.LinkedTiles) {
		tileB.Metadata.LinkedTiles = append(tileB.Metadata.LinkedTiles, tileA.Coords)
	}
}

func RemoveLinkedTiles(tile *data.Tile) {
	for _, ltc := range tile.Metadata.LinkedTiles {
		lt := data.CurrPuzzleSet.CurrPuzzle.Tiles.Get(ltc.X, ltc.Y)
		for i, t := range lt.Metadata.LinkedTiles {
			if t == tile.Coords {
				if len(lt.Metadata.LinkedTiles) < 2 {
					lt.Metadata.LinkedTiles = []world.Coords{}
				} else {
					lt.Metadata.LinkedTiles = append(lt.Metadata.LinkedTiles[:i], lt.Metadata.LinkedTiles[i+1:]...)
				}
				break
			}
		}
	}
	tile.Metadata.LinkedTiles = []world.Coords{}
}

// UpdateLinkedTiles should be called if a selection gets placed.
func UpdateLinkedTiles(tile *data.Tile) {
	for _, ltc := range tile.Metadata.LinkedTiles {
		lt := data.CurrPuzzleSet.CurrPuzzle.Tiles.Get(ltc.X, ltc.Y)
		if !world.CoordsIn(tile.Coords, lt.Metadata.LinkedTiles) {
			lt.Metadata.LinkedTiles = append(lt.Metadata.LinkedTiles, tile.Coords)
		}
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
