package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
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
	case data.BlockEmpty, data.BlockLadder, data.BlockLadderExit, data.BlockLadderCracked:
	case data.BlockTurf, data.BlockFall, data.BlockCracked, data.BlockPhase:
		if data.EditorDraw || !tile.Flags.Collapse {
			spr = append(spr, img.NewSprite(GetBlockSprite(tile), constants.TileBatch))
			if tile.Block == data.BlockFall && data.EditorDraw {
				spr = append(spr, img.NewSprite(constants.TileFall, constants.TileBatch))
			} else if tile.Block == data.BlockCracked {
				if data.EditorDraw || (tile.Flags.Cracked && tile.Counter > 2) {
					spr = append(spr, img.NewSprite(constants.TileCracked, constants.TileBatch))
				} else if tile.Flags.Cracked {
					spr = append(spr, img.NewSprite(fmt.Sprintf("%s%d", constants.TileCracking, tile.Counter), constants.TileBatch))
				}
			} else if tile.Block == data.BlockPhase && data.EditorDraw {
				spr = append(spr, img.NewSprite(constants.TilePhase, constants.TileBatch))
			}
		}
	default:
		spr = append(spr, img.NewSprite(tile.Block.String(), constants.TileBatch))
	}
	var lStr string
	if tile.Live {
		lStr = GetLadderSpriteLive(tile)
	} else {
		lStr = GetLadderSpriteEditor(tile)
	}
	if lStr != "" {
		spr = append(spr, img.NewSprite(lStr, constants.TileBatch))
	}
	return spr
}

func GetBlockSprite(tile *data.Tile) string {
	// check position to get correct sprite
	top := true
	bottom := true
	var a *data.Tile
	if tile.Live {
		a = data.CurrLevel.Tiles.Get(tile.Coords.X, tile.Coords.Y+1)
	} else {
		a = data.CurrPuzzle.Tiles.Get(tile.Coords.X, tile.Coords.Y+1)
	}
	if a.IsBlock() {
		top = false
	}
	var b *data.Tile
	if tile.Live {
		b = data.CurrLevel.Tiles.Get(tile.Coords.X, tile.Coords.Y-1)
	} else {
		b = data.CurrPuzzle.Tiles.Get(tile.Coords.X, tile.Coords.Y-1)
	}
	if b == nil || b.IsBlock() {
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

func GetLadderSpriteLive(tile *data.Tile) string {
	belowTile := data.CurrLevel.Tiles.Get(tile.Coords.X, tile.Coords.Y-1)
	var sKey string
	if !tile.Flags.LCollapse &&
		(tile.Ladder == data.BlockLadder ||
			(tile.Ladder == data.BlockLadderCracked && !tile.Flags.LCracked) ||
			(tile.Ladder == data.BlockLadderExit && data.CurrLevel.DoorsOpen)) {
		if tile.IsLadder() && belowTile != nil && belowTile.IsLadder() {
			sKey = constants.TileLadderMiddle
		} else if tile.IsLadder() {
			sKey = constants.TileLadderBottom
		} else if belowTile != nil && belowTile.IsLadder() {
			sKey = constants.TileLadderTop
		} else {
			sKey = ""
		}
	} else if !tile.Flags.LCollapse &&
		tile.Ladder == data.BlockLadderCracked &&
		tile.Flags.LCracked {
		if tile.IsLadder() && belowTile != nil && belowTile.IsLadder() {
			if tile.Counter > 6 {
				sKey = fmt.Sprintf("%s%d", constants.TileLadderCrackingM, 3)
			} else {
				sKey = fmt.Sprintf("%s%d", constants.TileLadderCrackingM, tile.Counter/2)
			}
		} else if tile.IsLadder() {
			if tile.Counter > 6 {
				sKey = fmt.Sprintf("%s%d", constants.TileLadderCrackingB, 3)
			} else {
				sKey = fmt.Sprintf("%s%d", constants.TileLadderCrackingB, tile.Counter/2)
			}
		}
	} else if belowTile != nil && belowTile.IsLadder() {
		sKey = constants.TileLadderTop
	}
	return sKey
}

func GetLadderSpriteEditor(tile *data.Tile) string {
	belowTile := data.CurrPuzzle.Tiles.Get(tile.Coords.X, tile.Coords.Y-1)
	var sKey string
	if tile.IsLadder() && belowTile != nil && belowTile.IsLadder() {
		switch tile.Ladder {
		case data.BlockLadder:
			sKey = constants.TileLadderMiddle
		case data.BlockLadderExit:
			sKey = constants.TileExitLadderM
		case data.BlockLadderCracked:
			sKey = constants.TileLadderCrackM
		}
	} else if tile.IsLadder() {
		switch tile.Ladder {
		case data.BlockLadder:
			sKey = constants.TileLadderBottom
		case data.BlockLadderExit:
			sKey = constants.TileExitLadderB
		case data.BlockLadderCracked:
			sKey = constants.TileLadderCrackB
		}
	} else if belowTile != nil && belowTile.IsLadder() {
		switch belowTile.Ladder {
		case data.BlockLadder, data.BlockLadderCracked:
			sKey = constants.TileLadderTop
		case data.BlockLadderExit:
			sKey = constants.TileExitLadderT
		}
	} else {
		sKey = ""
	}
	return sKey
}

func GetTileSpritesSelection(tile *data.Tile) []*img.Sprite {
	var spr []*img.Sprite
	switch tile.Block {
	case data.BlockEmpty:
	case data.BlockTurf, data.BlockFall, data.BlockCracked, data.BlockPhase:
		spr = append(spr, img.NewSprite(GetSpriteSelection(tile), constants.TileBatch))
		if tile.Block == data.BlockFall && data.EditorDraw {
			spr = append(spr, img.NewSprite(constants.TileFall, constants.TileBatch))
		} else if tile.Block == data.BlockCracked && data.EditorDraw {
			spr = append(spr, img.NewSprite(constants.TileCracked, constants.TileBatch))
		} else if tile.Block == data.BlockPhase && data.EditorDraw {
			spr = append(spr, img.NewSprite(constants.TilePhase, constants.TileBatch))
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
		if data.CurrSelect.Tiles[above.Y][above.X].IsBlock() {
			top = false
		}
	}
	if CoordsLegalSelection(below) {
		if data.CurrSelect.Tiles[below.Y][below.X].IsBlock() {
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
	var belowTile *data.Tile
	if CoordsLegalSelection(below) {
		belowTile = data.CurrSelect.Tiles[below.Y][below.X]
		bottom = belowTile.IsLadder()
	}
	var sKey string
	if tile.IsLadder() && bottom {
		switch tile.Ladder {
		case data.BlockLadder:
			sKey = constants.TileLadderMiddle
		case data.BlockLadderExit:
			sKey = constants.TileExitLadderM
		case data.BlockLadderCracked:
			sKey = constants.TileLadderCrackM
		}
	} else if tile.IsLadder() {
		switch tile.Ladder {
		case data.BlockLadder:
			sKey = constants.TileLadderBottom
		case data.BlockLadderExit:
			sKey = constants.TileExitLadderB
		case data.BlockLadderCracked:
			sKey = constants.TileLadderCrackB
		}
	} else if bottom && belowTile != nil {
		switch belowTile.Ladder {
		case data.BlockLadder, data.BlockLadderCracked:
			sKey = constants.TileLadderTop
		case data.BlockLadderExit:
			sKey = constants.TileExitLadderT
		}
	} else {
		sKey = ""
	}
	return sKey
}
