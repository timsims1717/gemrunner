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

var (
	TileUpdate = false
)

func PuzzleInit() {
	for _, result := range myecs.Manager.Query(myecs.IsTile) {
		myecs.Manager.DisposeEntity(result)
	}
	if data.CurrPuzzle != nil {
		for _, row := range data.CurrPuzzle.Tiles {
			for _, tile := range row {
				obj := object.New()
				obj.Pos = world.MapToWorld(tile.Coords)
				obj.Layer = 2
				myecs.Manager.NewEntity().
					AddComponent(myecs.Object, obj).
					AddComponent(myecs.Tile, tile)
			}
		}
		TileUpdate = true
	}
}

func TileSystem() {
	if TileUpdate {
		for _, result := range myecs.Manager.Query(myecs.IsTile) {
			_, okO := result.Components[myecs.Object].(*object.Object)
			tile, ok := result.Components[myecs.Tile].(*data.Tile)
			if okO && ok {
				switch tile.Block {
				case data.RedRock:
					// main block, check position to get correct sprite
					top := true
					bottom := true
					above := tile.Coords
					above.Y++
					below := tile.Coords
					below.Y--
					if CoordsLegal(above) {
						if data.CurrPuzzle.Tiles[above.Y][above.X].Block != data.Empty {
							top = false
						}
					} else {
						top = false
					}
					if CoordsLegal(below) {
						if data.CurrPuzzle.Tiles[below.Y][below.X].Block != data.Empty {
							bottom = false
						}
					} else {
						bottom = false
					}
					var sKey string
					if top && bottom {
						sKey = fmt.Sprintf("%s_bottom_top", tile.Block.String())
					} else if top {
						sKey = fmt.Sprintf("%s_top", tile.Block.String())
					} else if bottom {
						sKey = fmt.Sprintf("%s_bottom", tile.Block.String())
					} else {
						sKey = tile.Block.String()
					}
					spr := img.NewSprite(sKey, constants.TileFGBatch)
					result.Entity.AddComponent(myecs.Drawable, spr)
				default:
					result.Entity.RemoveComponent(myecs.Drawable)
				}
			}
		}
		TileUpdate = false
	}
}

func ChangeBlock(coords world.Coords, block data.Block) {
	if data.CurrPuzzle != nil {
		if CoordsLegal(coords) {
			// add to puzzle
			tile := data.CurrPuzzle.Tiles[coords.Y][coords.X]
			tile.Block = block
			TileUpdate = true
		} else {
			fmt.Println("error: illegal coordinates for tile")
		}
	} else {
		fmt.Println("error: attempted to add tile when no puzzle is loaded")
	}
}

func DeleteBlock(coords world.Coords) {
	if data.CurrPuzzle != nil {
		if CoordsLegal(coords) {
			tile := data.CurrPuzzle.Tiles[coords.Y][coords.X]
			if tile.Fall {
				tile.Fall = false
			} else if tile.Ladder {
				tile.Ladder = false
			} else if tile.Object != data.Empty {
				tile.Object = data.Empty
			} else if tile.Doodad != data.Empty {
				tile.Doodad = data.Empty
			} else if tile.Block != data.Empty {
				tile.Block = data.Empty
			}
			TileUpdate = true
		} else {
			fmt.Println("error: illegal coordinates for tile")
		}
	} else {
		fmt.Println("error: attempted to add tile when no puzzle is loaded")
	}
}

func CoordsLegal(coords world.Coords) bool {
	return coords.X >= 0 && coords.Y >= 0 && coords.X < constants.PuzzleWidth && coords.Y < constants.PuzzleHeight
}
