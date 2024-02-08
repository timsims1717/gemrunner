package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/world"
)

func TileSpriteSystemPre() {
	for _, result := range myecs.Manager.Query(myecs.IsTile) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		tile, ok := result.Components[myecs.Tile].(*data.Tile)
		if okO && ok {
			if tile.Update && !data.CurrPuzzle.Click {
				tile.Update = false
			}
			if obj.Hidden {
				obj.Hidden = false
			}
		}
	}
}

func TileSpriteSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsTile) {
		_, okO := result.Components[myecs.Object].(*object.Object)
		tile, ok := result.Components[myecs.Tile].(*data.Tile)
		if okO && ok {
			spr := GetTileSprites(tile)
			if len(spr) == 1 {
				result.Entity.AddComponent(myecs.Drawable, spr[0])
			} else if len(spr) > 0 {
				result.Entity.AddComponent(myecs.Drawable, spr)
			} else {
				result.Entity.RemoveComponent(myecs.Drawable)
			}
		}
	}
}

func GetTileSprites(tile *data.Tile) []*img.Sprite {
	var spr []*img.Sprite
	switch tile.Block {
	case data.Empty, data.Ladder:
	case data.Turf, data.Fall:
		spr = append(spr, img.NewSprite(GetBlockSprite(tile), constants.TileBatch))
		if tile.Block == data.Fall {
			spr = append(spr, img.NewSprite(constants.TileFall, constants.TileBatch))
		}
	default:
		spr = append(spr, img.NewSprite(tile.Block.String(), constants.TileBatch))
	}
	lStr := GetLadderSprite(tile)
	if lStr != "" {
		spr = append(spr, img.NewSprite(lStr, constants.TileBatch))
	}
	return spr
}

func GetBlockSprite(tile *data.Tile) string {
	// check position to get correct sprite
	top := true
	bottom := true
	above := tile.Coords
	above.Y++
	below := tile.Coords
	below.Y--
	if CoordsLegal(above) {
		a := data.CurrPuzzle.Tiles.T[above.Y][above.X].Block
		if a == data.Turf || a == data.Fall {
			top = false
		}
	} else {
		top = false
	}
	if CoordsLegal(below) {
		b := data.CurrPuzzle.Tiles.T[below.Y][below.X].Block
		if b == data.Turf || b == data.Fall {
			bottom = false
		}
	} else {
		bottom = false
	}
	var sKey string
	if top && bottom {
		sKey = fmt.Sprintf("%s%s", tile.Block.String(), constants.TileBottomTop)
	} else if top {
		sKey = fmt.Sprintf("%s%s", tile.Block.String(), constants.TileTop)
	} else if bottom {
		sKey = fmt.Sprintf("%s%s", tile.Block.String(), constants.TileBottom)
	} else {
		sKey = tile.Block.String()
	}
	return sKey
}

func GetLadderSprite(tile *data.Tile) string {
	bottom := false
	below := tile.Coords
	below.Y--
	if CoordsLegal(below) {
		if data.CurrPuzzle.Tiles.T[below.Y][below.X].Ladder {
			bottom = true
		}
	}
	var sKey string
	if tile.Ladder && bottom {
		sKey = constants.TileLadderMiddle
	} else if tile.Ladder {
		sKey = constants.TileLadderBottom
	} else if bottom {
		sKey = constants.TileLadderTop
	} else {
		sKey = ""
	}
	return sKey
}

func GetTileSpritesSelection(tile *data.Tile) []*img.Sprite {
	var spr []*img.Sprite
	switch tile.Block {
	case data.Empty:
	case data.Turf, data.Fall:
		spr = append(spr, img.NewSprite(GetSpriteSelection(tile), constants.TileBatch))
		if tile.Block == data.Fall {
			spr = append(spr, img.NewSprite(constants.TileFall, constants.TileBatch))
		}
	default:
		spr = append(spr, img.NewSprite(tile.Block.String(), constants.TileBatch))
	}
	lStr := GetLadderSpriteSelection(tile)
	if lStr != "" {
		spr = append(spr, img.NewSprite(lStr, constants.TileBatch))
	}
	return spr
}

func GetSpriteSelection(tile *data.Tile) string {
	// check position to get correct sprite
	top := true
	bottom := true
	above := tile.Coords
	above.Y++
	below := tile.Coords
	below.Y--
	if CoordsLegalSelection(above) {
		if data.CurrSelect.Tiles[above.Y][above.X].Block != data.Empty {
			top = false
		}
	}
	if CoordsLegalSelection(below) {
		if data.CurrSelect.Tiles[below.Y][below.X].Block != data.Empty {
			bottom = false
		}
	}
	var sKey string
	if top && bottom {
		sKey = fmt.Sprintf("%s%s", tile.Block.String(), constants.TileBottomTop)
	} else if top {
		sKey = fmt.Sprintf("%s%s", tile.Block.String(), constants.TileTop)
	} else if bottom {
		sKey = fmt.Sprintf("%s%s", tile.Block.String(), constants.TileBottom)
	} else {
		sKey = tile.Block.String()
	}
	return sKey
}

func GetLadderSpriteSelection(tile *data.Tile) string {
	bottom := false
	below := tile.Coords
	below.Y--
	if CoordsLegalSelection(below) {
		if data.CurrSelect.Tiles[below.Y][below.X].Ladder {
			bottom = true
		}
	}
	var sKey string
	if tile.Ladder && bottom {
		sKey = constants.TileLadderMiddle
	} else if tile.Ladder {
		sKey = constants.TileLadderBottom
	} else if bottom {
		sKey = constants.TileLadderTop
	} else {
		sKey = ""
	}
	return sKey
}

func SetBlock(coords world.Coords, block data.Block) {
	if data.CurrPuzzle != nil {
		if CoordsLegal(coords) {
			// add to puzzle
			tile := data.CurrPuzzle.Tiles.T[coords.Y][coords.X]
			if !tile.Update {
				switch block {
				case data.Ladder:
					if tile.Ladder || tile.Block != data.Turf {
						tile.Block = data.Empty
					}
					tile.Ladder = true
				case data.Player1:
					tile.Ladder = false
					// ensure no other player of that type are in puzzle
					for _, row := range data.CurrPuzzle.Tiles.T {
						for _, t := range row {
							if t.Block == block {
								t.Block = data.Empty
							}
						}
					}
					tile.Block = block
				default:
					if tile.Ladder && !(block == data.Turf && tile.Block == data.Empty) {
						tile.Ladder = false
					}
					tile.Block = block
				}
			}
			data.CurrPuzzle.Update = true
			tile.Update = true
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
				if tile.Ladder {
					tile.Ladder = false
				} else if tile.Block != data.Empty {
					tile.Block = data.Empty
				}
				data.CurrPuzzle.Update = true
				tile.Update = true
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
