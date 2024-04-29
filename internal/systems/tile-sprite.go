package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"github.com/gopxl/pixel"
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
		if okO && ok && !result.Entity.HasComponent(myecs.Animated) {
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
	var sprs []*img.Sprite
	switch tile.Block {
	case data.BlockEmpty, data.BlockLadder, data.BlockLadderExit, data.BlockLadderCracked:
	case data.BlockTurf, data.BlockBedrock,
		data.BlockFall, data.BlockCracked, data.BlockPhase,
		data.BlockLadderTurf, data.BlockLadderCrackedTurf, data.BlockLadderExitTurf:
		if data.EditorDraw {
			sprs = GetBlockSpritesEditor(tile)
		} else {
			sprs = GetBlockSprites(tile)
		}
	case data.BlockSpike:
		sprs = append(sprs, img.NewSprite(GetSpikeSprite(tile), constants.TileBatch))
	case data.BlockDemonRegen, data.BlockFlyRegen:
		if !tile.Live {
			sprs = append(sprs, img.NewSprite(tile.Block.String(), constants.TileBatch))
		}
	default:
		sprs = append(sprs, img.NewSprite(tile.Block.String(), constants.TileBatch))
	}
	var lStr string
	if tile.Live {
		lStr = GetLadderSpriteLive(tile)
	} else {
		lStr = GetLadderSpriteEditor(tile)
	}
	if lStr != "" {
		sprs = append(sprs, img.NewSprite(lStr, constants.TileBatch))
	}
	if data.EditorDraw {
		sprs = append(sprs, GetWrenchSprites(tile)...)
	}
	return sprs
}

func GetBlockSpritesEditor(tile *data.Tile) []*img.Sprite {
	var sprs []*img.Sprite
	sprs = append(sprs, img.NewSprite(GetBlockSprite(tile), constants.TileBatch))
	if tile.Block == data.BlockFall {
		sprs = append(sprs, img.NewSprite(constants.TileFall, constants.TileBatch))
	} else if tile.Block == data.BlockCracked {
		sprs = append(sprs, img.NewSprite(constants.TileCracked, constants.TileBatch))
	} else if tile.Block == data.BlockPhase {
		sprs = append(sprs, img.NewSprite(constants.TilePhase, constants.TileBatch))
	}
	return sprs
}

func GetBlockSprites(tile *data.Tile) []*img.Sprite {
	var sprs []*img.Sprite
	spr := img.NewSprite(GetBlockSprite(tile), constants.TileBatch)
	if tile.Flags.Collapse && tile.Counter > constants.CollapseCounter {
		return sprs
	}
	sprs = append(sprs, spr)
	if tile.Block == data.BlockCracked && !tile.Flags.Collapse {
		if tile.Metadata.ShowCrack {
			sprs = append(sprs, img.NewSprite(constants.TileCrackedShow, constants.TileBatch))
		}
	}
	return sprs
}

func GetWrenchSprites(tile *data.Tile) []*img.Sprite {
	var sprs []*img.Sprite
	if data.Editor.Mode == data.Wrench {
		for i := 0; i < 4; i++ {
			offset := pixel.V(-4, 4)
			switch i {
			case 0:
				if tile.Metadata.Regenerate &&
					(tile.Block == data.BlockCracked ||
						tile.Block == data.BlockLadderCrackedTurf ||
						tile.Block == data.BlockLadderCracked ||
						tile.Block == data.BlockFly ||
						tile.Block == data.BlockDemon) {
					sprs = append(sprs, img.NewSprite("tile_ui_regen", constants.UIBatch).WithOffset(offset))
				}
			case 1:
				offset.X = 4
				if tile.Metadata.ShowCrack &&
					(tile.Block == data.BlockCracked ||
						tile.Block == data.BlockLadderCrackedTurf ||
						tile.Block == data.BlockLadderCracked) {
					sprs = append(sprs, img.NewSprite("tile_ui_show", constants.UIBatch).WithOffset(offset))
				} else if tile.Metadata.Flipped &&
					tile.Block == data.BlockFly {
					sprs = append(sprs, img.NewSprite("tile_ui_flip", constants.UIBatch).WithOffset(offset))
				}
			case 2:
				offset.X = -4
				offset.Y = -4
				if tile.Metadata.EnemyCrack &&
					(tile.Block == data.BlockCracked ||
						tile.Block == data.BlockLadderCrackedTurf ||
						tile.Block == data.BlockLadderCracked) {
					sprs = append(sprs, img.NewSprite("tile_ui_enemy", constants.UIBatch).WithOffset(offset))
				}
			case 3:
				offset.X = 4
				offset.Y = -4
			}
		}
	}
	return sprs
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
	if b.IsBlock() {
		bottom = false
	}
	var sKey string
	if top && bottom {
		sKey = fmt.Sprintf("%s%s", tile.Block.String(), constants.TileBottomTop)
	} else if top {
		sKey = fmt.Sprintf("%s%s", tile.Block.String(), constants.TileTop)
	} else if bottom {
		sKey = fmt.Sprintf("%s%s", tile.Block.String(), constants.TileBottom)
	} else if tile.AltBlock == 1 && tile.Block == data.BlockTurf {
		sKey = fmt.Sprintf("%s%s", tile.Block.String(), constants.TileAlt)
	} else {
		sKey = tile.Block.String()
	}
	return sKey
}

func GetSpikeSprite(tile *data.Tile) string {
	// check position to get correct sprite
	bottom := true
	var b *data.Tile
	if tile.Live {
		b = data.CurrLevel.Tiles.Get(tile.Coords.X, tile.Coords.Y-1)
	} else {
		b = data.CurrPuzzle.Tiles.Get(tile.Coords.X, tile.Coords.Y-1)
	}
	if b == nil || (b.IsBlock() && b.Block != data.BlockSpike) {
		bottom = false
	}
	var sKey string
	if bottom {
		sKey = fmt.Sprintf("%s%s", tile.Block.String(), constants.TileBottom)
	} else {
		sKey = tile.Block.String()
	}
	return sKey
}

func GetLadderSpriteLive(tile *data.Tile) string {
	belowTile := data.CurrLevel.Tiles.Get(tile.Coords.X, tile.Coords.Y-1)
	aboveTile := data.CurrLevel.Tiles.Get(tile.Coords.X, tile.Coords.Y+1)
	var sKey string
	if !tile.Flags.LCollapse &&
		(tile.Block == data.BlockLadder ||
			tile.Block == data.BlockLadderTurf ||
			(tile.Block == data.BlockLadderCracked && !tile.Flags.LCracked) ||
			(tile.Block == data.BlockLadderCrackedTurf && !tile.Flags.LCracked) ||
			(tile.Block == data.BlockLadderExit && data.CurrLevel.DoorsOpen) ||
			(tile.Block == data.BlockLadderExitTurf && data.CurrLevel.DoorsOpen)) {
		if tile.IsLadder() && belowTile != nil && belowTile.IsLadder() {
			if tile.IsBlock() &&
				!aboveTile.IsBlock() {
				sKey = constants.TileLadderLedgeMiddle
			} else {
				sKey = constants.TileLadderMiddle
			}
		} else if tile.IsLadder() {
			if tile.IsBlock() &&
				!aboveTile.IsBlock() {
				sKey = constants.TileLadderLedgeBottom
			} else {
				sKey = constants.TileLadderBottom
			}
		} else if belowTile != nil && belowTile.IsLadder() {
			sKey = constants.TileLadderTop
		} else {
			sKey = ""
		}
	} else if !tile.Flags.LCollapse &&
		(tile.Block == data.BlockLadderCracked || tile.Block == data.BlockLadderCrackedTurf) &&
		tile.Flags.LCracked {
		if tile.IsLadder() && belowTile != nil && belowTile.IsLadder() {
			if tile.IsBlock() &&
				!aboveTile.IsBlock() {
				if tile.Counter > 6 {
					sKey = fmt.Sprintf("%s%d", constants.TileLadderLedgeCrackingM, 3)
				} else {
					sKey = fmt.Sprintf("%s%d", constants.TileLadderLedgeCrackingM, tile.Counter/2)
				}
			} else {
				if tile.Counter > 6 {
					sKey = fmt.Sprintf("%s%d", constants.TileLadderCrackingM, 3)
				} else {
					sKey = fmt.Sprintf("%s%d", constants.TileLadderCrackingM, tile.Counter/2)
				}
			}
		} else if tile.IsLadder() {
			if tile.IsBlock() &&
				!aboveTile.IsBlock() {
				if tile.Counter > 6 {
					sKey = fmt.Sprintf("%s%d", constants.TileLadderLedgeCrackingB, 3)
				} else {
					sKey = fmt.Sprintf("%s%d", constants.TileLadderLedgeCrackingB, tile.Counter/2)
				}
			} else {
				if tile.Counter > 6 {
					sKey = fmt.Sprintf("%s%d", constants.TileLadderCrackingB, 3)
				} else {
					sKey = fmt.Sprintf("%s%d", constants.TileLadderCrackingB, tile.Counter/2)
				}
			}
		}
	} else if belowTile != nil && belowTile.IsLadder() {
		sKey = constants.TileLadderTop
	}
	return sKey
}

func GetLadderSpriteEditor(tile *data.Tile) string {
	belowTile := data.CurrPuzzle.Tiles.Get(tile.Coords.X, tile.Coords.Y-1)
	aboveTile := data.CurrPuzzle.Tiles.Get(tile.Coords.X, tile.Coords.Y+1)
	var sKey string
	if tile.IsLadder() && belowTile != nil && belowTile.IsLadder() {
		if tile.IsBlock() &&
			!aboveTile.IsBlock() {
			switch tile.Block {
			case data.BlockLadder, data.BlockLadderTurf:
				sKey = constants.TileLadderLedgeMiddle
			case data.BlockLadderExit, data.BlockLadderExitTurf:
				sKey = constants.TileExitLadderM
			case data.BlockLadderCracked, data.BlockLadderCrackedTurf:
				sKey = constants.TileLadderLedgeCrackM
			}
		} else {
			switch tile.Block {
			case data.BlockLadder, data.BlockLadderTurf:
				sKey = constants.TileLadderMiddle
			case data.BlockLadderExit, data.BlockLadderExitTurf:
				sKey = constants.TileExitLadderM
			case data.BlockLadderCracked, data.BlockLadderCrackedTurf:
				sKey = constants.TileLadderCrackM
			}
		}
	} else if tile.IsLadder() {
		if tile.IsBlock() &&
			!aboveTile.IsBlock() {
			switch tile.Block {
			case data.BlockLadder, data.BlockLadderTurf:
				sKey = constants.TileLadderLedgeBottom
			case data.BlockLadderExit, data.BlockLadderExitTurf:
				sKey = constants.TileExitLadderB
			case data.BlockLadderCracked, data.BlockLadderCrackedTurf:
				sKey = constants.TileLadderLedgeCrackB
			}
		} else {
			switch tile.Block {
			case data.BlockLadder, data.BlockLadderTurf:
				sKey = constants.TileLadderBottom
			case data.BlockLadderExit, data.BlockLadderExitTurf:
				sKey = constants.TileExitLadderB
			case data.BlockLadderCracked, data.BlockLadderCrackedTurf:
				sKey = constants.TileLadderCrackB
			}
		}
	} else if belowTile != nil && belowTile.IsLadder() {
		switch belowTile.Block {
		case data.BlockLadder, data.BlockLadderTurf, data.BlockLadderCracked, data.BlockLadderCrackedTurf:
			sKey = constants.TileLadderTop
		case data.BlockLadderExit, data.BlockLadderExitTurf:
			sKey = constants.TileExitLadderT
		}
	} else {
		sKey = ""
	}
	return sKey
}

func GetTileSpritesSelection(tile *data.Tile) []*img.Sprite {
	var sprs []*img.Sprite
	if data.EditorDraw {
		switch tile.Block {
		case data.BlockEmpty:
		case data.BlockTurf, data.BlockFall, data.BlockCracked, data.BlockPhase:
			sprs = append(sprs, img.NewSprite(GetSpriteSelection(tile), constants.TileBatch))
			if tile.Block == data.BlockFall {
				sprs = append(sprs, img.NewSprite(constants.TileFall, constants.TileBatch))
			} else if tile.Block == data.BlockCracked {
				sprs = append(sprs, img.NewSprite(constants.TileCracked, constants.TileBatch))
			} else if tile.Block == data.BlockPhase {
				sprs = append(sprs, img.NewSprite(constants.TilePhase, constants.TileBatch))
			}
		default:
			sprs = append(sprs, img.NewSprite(tile.Block.String(), constants.TileBatch))
		}
		lStr := GetLadderSpriteSelection(tile)
		if lStr != "" {
			sprs = append(sprs, img.NewSprite(lStr, constants.TileBatch))
		}
		sprs = append(sprs, GetWrenchSprites(tile)...)
	}
	return sprs
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
		switch tile.Block {
		case data.BlockLadder, data.BlockLadderTurf:
			sKey = constants.TileLadderMiddle
		case data.BlockLadderExit, data.BlockLadderExitTurf:
			sKey = constants.TileExitLadderM
		case data.BlockLadderCracked, data.BlockLadderCrackedTurf:
			sKey = constants.TileLadderCrackM
		}
	} else if tile.IsLadder() {
		switch tile.Block {
		case data.BlockLadder, data.BlockLadderTurf:
			sKey = constants.TileLadderBottom
		case data.BlockLadderExit, data.BlockLadderExitTurf:
			sKey = constants.TileExitLadderB
		case data.BlockLadderCracked, data.BlockLadderCrackedTurf:
			sKey = constants.TileLadderCrackB
		}
	} else if bottom && belowTile != nil {
		switch belowTile.Block {
		case data.BlockLadder, data.BlockLadderTurf, data.BlockLadderCracked, data.BlockLadderCrackedTurf:
			sKey = constants.TileLadderTop
		case data.BlockLadderExit, data.BlockLadderExitTurf:
			sKey = constants.TileExitLadderT
		}
	} else {
		sKey = ""
	}
	return sKey
}
