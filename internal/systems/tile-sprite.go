package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"github.com/gopxl/pixel"
)

type spriteChanger struct {
	SprKey  string
	Batch   string
	Offset  pixel.Vec
	IsAnim  bool
	Reverse bool
	Frames  []int
}

func newSprChanger(key, batch string) *spriteChanger {
	return &spriteChanger{
		SprKey: key,
		Batch:  batch,
	}
}

func (sc *spriteChanger) SetAnim(a bool) *spriteChanger {
	sc.IsAnim = a
	return sc
}

func (sc *spriteChanger) SetReverse(r bool) *spriteChanger {
	sc.Reverse = r
	return sc
}

func (sc *spriteChanger) WithFrames(f []int) *spriteChanger {
	sc.Frames = f
	return sc
}

func (sc *spriteChanger) WithOffset(pos pixel.Vec) *spriteChanger {
	sc.Offset = pos
	return sc
}

func TileSpriteSystemPre() {
	for _, result := range myecs.Manager.Query(myecs.IsTile) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		tile, ok := result.Components[myecs.Tile].(*data.Tile)
		if okO && ok {
			if tile.Update && !data.CurrPuzzleSet.CurrPuzzle.Click {
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
			needsUpdate := NeedsUpdate(tile)
			if needsUpdate {
				drs, hasDraws := result.Entity.GetComponentData(myecs.Drawable)
				sprs := GetTileDrawables(tile)
				if len(sprs) == 0 {
					if hasDraws {
						result.Entity.RemoveComponent(myecs.Drawable)
						result.Entity.RemoveComponent(myecs.Animated)
					}
					continue
				}
				// only change the drawables that have changed
				var changed bool
				currDraws, okD := drs.([]any)
				if !okD {
					currDraws = []any{}
					changed = true
				}
				var hasAnim bool
				for i, spr := range sprs {
					if spr.IsAnim {
						hasAnim = true
					}
					if !okD || i >= len(currDraws) {
						currDraws = append(currDraws, buildDrawable(spr, tile.Live))
						changed = true
						continue
					}
					cd := currDraws[i]
					if cSpr, ok1 := cd.(*img.Sprite); ok1 {
						if !(cSpr.Key == spr.SprKey && cSpr.Batch == spr.Batch &&
							cSpr.Offset == spr.Offset && !spr.IsAnim) {
							currDraws[i] = buildDrawable(spr, tile.Live)
							changed = true
						}
					} else if cAnim, ok2 := cd.(*reanimator.Tree); ok2 {
						if !(cAnim.Default == spr.SprKey && spr.IsAnim) {
							currDraws[i] = buildDrawable(spr, tile.Live)
							changed = true
						}
					} else {
						currDraws[i] = buildDrawable(spr, tile.Live)
						changed = true
					}
				}
				if len(currDraws) > len(sprs) {
					currDraws = currDraws[:len(sprs)]
					changed = true
				}
				if changed {
					result.Entity.AddComponent(myecs.Drawable, currDraws)
					if hasAnim {
						result.Entity.AddComponent(myecs.Animated, currDraws)
					} else {
						result.Entity.RemoveComponent(myecs.Animated)
					}
				}
			}
			tile.LastBlock = tile.Block
		}
	}
}

func buildDrawable(sprCh *spriteChanger, live bool) any {
	if sprCh.IsAnim {
		var a *reanimator.Anim
		if len(sprCh.Frames) > 0 {
			a = reanimator.NewBatchAnimationCustom(sprCh.SprKey, img.Batchers[sprCh.Batch], sprCh.SprKey, sprCh.Frames, reanimator.Loop)
		} else {
			a = reanimator.NewBatchAnimation(sprCh.SprKey, img.Batchers[sprCh.Batch], sprCh.SprKey, reanimator.Loop)
		}
		if sprCh.Reverse {
			a = a.Reverse()
		}
		frames := len(a.S)
		anim := reanimator.NewSimple(a)
		if live {
			anim.SetAnim(sprCh.SprKey, data.CurrLevel.FrameNumber%frames)
		} else {
			anim.SetAnim(sprCh.SprKey, data.Editor.FrameCount%frames)
		}
		return anim
	} else {
		spr := img.NewSprite(sprCh.SprKey, sprCh.Batch).WithOffset(sprCh.Offset)
		return spr
	}
}

func NeedsUpdate(tile *data.Tile) bool {
	if tile.Update || tile.Block != tile.LastBlock {
		return true
	}
	switch tile.Block {
	case data.BlockTurf, data.BlockBedrock,
		data.BlockFall, data.BlockCracked, data.BlockClose, data.BlockPhase,
		data.BlockLadderTurf, data.BlockLadderCrackedTurf, data.BlockLadderExitTurf,
		data.BlockSpike, data.BlockHideout, data.BlockLiquid:
		return true
	default:
		return false
	}
}

// GetTileDrawables returns a list of drawable sprites and/or animations
// returns the list and whether there are any animations
func GetTileDrawables(tile *data.Tile) []*spriteChanger {
	var sprs []*spriteChanger
	switch tile.Block {
	case data.BlockEmpty, data.BlockLadder, data.BlockLadderExit, data.BlockLadderCracked:
	case data.BlockTurf, data.BlockBedrock,
		data.BlockFall, data.BlockCracked, data.BlockClose, data.BlockPhase,
		data.BlockLadderTurf, data.BlockLadderCrackedTurf, data.BlockLadderExitTurf:
		if data.EditorDraw {
			sprs = GetBlockSpritesEditor(tile)
		} else {
			sprs = GetBlockSprites(tile)
		}
	case data.BlockSpike:
		sprs = append(sprs, newSprChanger(GetSpikeSprite(tile), constants.TileBatch))
	case data.BlockLiquid:
		sprs = GetLiquidSprites(tile)
	case data.BlockHideout:
		sprs = append(sprs, newSprChanger(GetHideoutSprite(tile), constants.TileBatch))
		sprs = append(sprs, newSprChanger(constants.TileHideout, constants.TileBatch))
	case data.BlockDemonRegen, data.BlockFlyRegen:
		if !tile.Live {
			sprs = append(sprs, newSprChanger(tile.SpriteString(), constants.TileBatch))
		}
	case data.BlockBomb:
		sprs = append(sprs, newSprChanger(tile.SpriteString(), constants.TileBatch))
		if tile.Metadata.Regenerate && tile.Metadata.BombCross {
			sprs = append(sprs, newSprChanger(constants.ItemBombRegenCross, constants.TileBatch).WithOffset(pixel.V(0, -2)))
		} else if tile.Metadata.BombCross {
			sprs = append(sprs, newSprChanger(constants.ItemBombCross, constants.TileBatch).WithOffset(pixel.V(0, -2)))
		} else if tile.Metadata.Regenerate {
			sprs = append(sprs, newSprChanger(constants.ItemBombRegen, constants.TileBatch).WithOffset(pixel.V(0, -2)))
		}
	case data.BlockBombLit:
		sprs = append(sprs, newSprChanger(constants.ItemBombLit, constants.TileBatch))
		if tile.Metadata.Regenerate && tile.Metadata.BombCross {
			sprs = append(sprs, newSprChanger(constants.ItemBombRegenCross, constants.TileBatch).WithOffset(pixel.V(0, -2)))
		} else if tile.Metadata.BombCross {
			sprs = append(sprs, newSprChanger(constants.ItemBombCross, constants.TileBatch).WithOffset(pixel.V(0, -2)))
		} else if tile.Metadata.Regenerate {
			sprs = append(sprs, newSprChanger(constants.ItemBombRegen, constants.TileBatch).WithOffset(pixel.V(0, -2)))
		}
	case data.BlockGear:
		//if tile.Live {
		//	anim.SetAnim("gear", data.CurrLevel.FrameNumber%4)
		//} else {
		//	anim.SetAnim("gear", data.Editor.FrameCount%4)
		//}
		if (tile.Coords.X+tile.Coords.Y)%2 == 0 {
			sprs = append(sprs, newSprChanger(tile.SpriteString(), constants.TileBatch).SetAnim(true).SetReverse(true).WithFrames([]int{3, 0, 1, 2}))
		} else {
			sprs = append(sprs, newSprChanger(tile.SpriteString(), constants.TileBatch).SetAnim(true))
		}
	default:
		sprs = append(sprs, newSprChanger(tile.SpriteString(), constants.TileBatch))
	}
	var lStr string
	if tile.Live {
		lStr = GetLadderSpriteLive(tile)
	} else {
		lStr = GetLadderSpriteEditor(tile)
	}
	if lStr != "" {
		sprs = append(sprs, newSprChanger(lStr, constants.TileBatch))
	}
	if data.EditorDraw {
		sprs = append(sprs, GetWrenchSprites(tile)...)
	}
	return sprs
}

func GetBlockSpritesEditor(tile *data.Tile) []*spriteChanger {
	var sprs []*spriteChanger
	sprs = append(sprs, newSprChanger(GetBlockSprite(tile), constants.TileBatch))
	if tile.Block == data.BlockFall {
		sprs = append(sprs, newSprChanger(constants.TileFall, constants.TileBatch))
	} else if tile.Block == data.BlockCracked {
		sprs = append(sprs, newSprChanger(constants.TileCracked, constants.TileBatch))
	} else if tile.Block == data.BlockClose {
		sprs = append(sprs, newSprChanger(constants.TileClose, constants.TileBatch))
	} else if tile.Block == data.BlockPhase {
		sprs = append(sprs, newSprChanger(constants.TilePhase, constants.TileBatch))
	}
	return sprs
}

func GetBlockSprites(tile *data.Tile) []*spriteChanger {
	var sprs []*spriteChanger
	if tile.Flags.Collapse &&
		(!tile.Flags.BareFangs && tile.Counter > constants.CollapseCounter) {
		return sprs
	}
	spr := newSprChanger(GetBlockSprite(tile), constants.TileBatch)
	sprs = append(sprs, spr)
	if tile.Block == data.BlockCracked && !tile.Flags.Collapse {
		if tile.Metadata.ShowCrack {
			sprs = append(sprs, newSprChanger(constants.TileCrackedShow, constants.TileBatch))
		}
	}
	return sprs
}

func GetLiquidSprites(tile *data.Tile) []*spriteChanger {
	var sprs []*spriteChanger
	// check position to get correct sprite
	top := true
	bottom := true
	left := true
	right := true
	var a *data.Tile
	if tile.Live {
		a = data.CurrLevel.Get(tile.Coords.X, tile.Coords.Y+1)
	} else {
		a = data.CurrPuzzleSet.CurrPuzzle.Get(tile.Coords.X, tile.Coords.Y+1)
	}
	if a != nil && a.Block == data.BlockLiquid {
		top = false
	}
	var b *data.Tile
	if tile.Live {
		b = data.CurrLevel.Get(tile.Coords.X, tile.Coords.Y-1)
	} else {
		b = data.CurrPuzzleSet.CurrPuzzle.Get(tile.Coords.X, tile.Coords.Y-1)
	}
	if b != nil && b.Block == data.BlockLiquid {
		bottom = false
	}
	var c *data.Tile
	if tile.Live {
		c = data.CurrLevel.Get(tile.Coords.X-1, tile.Coords.Y)
	} else {
		c = data.CurrPuzzleSet.CurrPuzzle.Get(tile.Coords.X-1, tile.Coords.Y)
	}
	if c != nil && c.Block == data.BlockLiquid {
		left = false
	}
	var d *data.Tile
	if tile.Live {
		d = data.CurrLevel.Get(tile.Coords.X+1, tile.Coords.Y)
	} else {
		d = data.CurrPuzzleSet.CurrPuzzle.Get(tile.Coords.X+1, tile.Coords.Y)
	}
	if d != nil && d.Block == data.BlockLiquid {
		right = false
	}
	var sKey string
	if left {
		sKey += "l"
	}
	if bottom {
		sKey += "b"
	}
	if right {
		sKey += "r"
	}
	sprChg := &spriteChanger{
		SprKey: tile.SpriteString(),
		Batch:  constants.TileBatch,
		IsAnim: true,
	}
	frames := make([]int, 8)
	switch data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldLiquid {
	case constants.LiquidBubbles:
		count := 0
		i := tile.Coords.X
		if i%2 == 0 {
			i += 4
		} else {
			i += 2
		}
		if tile.Coords.Y%2 == 0 {
			i += 4
		}
		i %= 8
		for count < 8 {
			frames[count] = i
			count++
			i++
			i %= 8
		}
	case constants.LiquidTiny:
		count := 0
		i := tile.Coords.X
		if i%2 == 0 {
			i += 4
		} else {
			i += 2
		}
		i %= 8
		for count < 8 {
			frames[count] = i
			count++
			i++
			i %= 8
		}
	default:
		frames = []int{0, 1, 2, 3, 4, 5, 6, 7}
	}
	if top {
		sprChg.SprKey = tile.SpriteString()
		sprChg.Frames = frames
	} else {
		switch data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldLiquid {
		case constants.LiquidBubbles, constants.LiquidTiny:
			sprChg.SprKey = fmt.Sprintf("%s_bottom", sprChg.SprKey)
			sprChg.Frames = frames
		case constants.LiquidWaves:
			if bottom {
				sprChg.SprKey = fmt.Sprintf("%s_bottom", sprChg.SprKey)
				sprChg.Frames = frames
			} else {
				sprChg.SprKey = constants.TileLiquid
				sprChg.IsAnim = false
			}
		default:
			sprChg.SprKey = constants.TileLiquid
			sprChg.IsAnim = false
		}
	}
	sprs = append(sprs, sprChg)
	if sKey != "" {
		sprs = append(sprs, newSprChanger(fmt.Sprintf("%s_%s", constants.TileLiquid, sKey), constants.TileBatch))
	}
	return sprs
}

func GetWrenchSprites(tile *data.Tile) []*spriteChanger {
	var sprs []*spriteChanger
	if data.Editor.Mode == data.ModeWrench {
		switch tile.Block {
		case data.BlockPhase:
			sprs = append(sprs, newSprChanger(fmt.Sprintf(constants.UINumber, tile.Metadata.Phase), constants.UIBatch))
		case data.BlockDoorHidden, data.BlockDoorVisible, data.BlockDoorLocked:
			exitIndex := tile.Metadata.ExitIndex + 1
			if tile.Metadata.ExitIndex == -1 {
				exitIndex = data.CurrPuzzleSet.PuzzleIndex + 2
			}
			if exitIndex < 10 {
				sprs = append(sprs, newSprChanger(fmt.Sprintf(constants.UINumber, exitIndex), constants.UIBatch))
			} else if exitIndex < 100 {
				sprs = append(sprs, newSprChanger(fmt.Sprintf(constants.UINumberX, exitIndex/10), constants.UIBatch).WithOffset(pixel.V(-3, 0)))
				sprs = append(sprs, newSprChanger(fmt.Sprintf(constants.UINumberX, exitIndex%10), constants.UIBatch).WithOffset(pixel.V(3, 0)))
			} else if exitIndex < 1000 {
				sprs = append(sprs, newSprChanger(fmt.Sprintf(constants.UINumberX, exitIndex/100), constants.UIBatch).WithOffset(pixel.V(-6, 0)))
				sprs = append(sprs, newSprChanger(fmt.Sprintf(constants.UINumberX, exitIndex%100/10), constants.UIBatch))
				sprs = append(sprs, newSprChanger(fmt.Sprintf(constants.UINumberX, exitIndex%10), constants.UIBatch).WithOffset(pixel.V(6, 0)))
			}
		case data.BlockFly, data.BlockDemon,
			data.BlockCracked,
			data.BlockLadderCracked, data.BlockLadderCrackedTurf:
			for i := 0; i < 4; i++ {
				offset := pixel.V(-4, 4)
				switch i {
				case 0:
					if tile.Metadata.Regenerate &&
						(tile.Block == data.BlockCracked ||
							tile.Block == data.BlockLadderCrackedTurf ||
							tile.Block == data.BlockLadderCracked ||
							tile.Block == data.BlockJetpack ||
							tile.Block == data.BlockDisguise ||
							tile.Block == data.BlockFly ||
							tile.Block == data.BlockDemon) {
						sprs = append(sprs, newSprChanger("tile_ui_regen", constants.UIBatch).WithOffset(offset))
					}
				case 1:
					offset.X = 4
					if tile.Metadata.ShowCrack &&
						(tile.Block == data.BlockCracked ||
							tile.Block == data.BlockLadderCrackedTurf ||
							tile.Block == data.BlockLadderCracked) {
						sprs = append(sprs, newSprChanger("tile_ui_show", constants.UIBatch).WithOffset(offset))
					} else if tile.Metadata.Flipped &&
						tile.Block == data.BlockFly {
						sprs = append(sprs, newSprChanger("tile_ui_flip", constants.UIBatch).WithOffset(offset))
					}
				case 2:
					offset.X = -4
					offset.Y = -4
					if tile.Metadata.EnemyCrack &&
						(tile.Block == data.BlockCracked ||
							tile.Block == data.BlockLadderCrackedTurf ||
							tile.Block == data.BlockLadderCracked) {
						sprs = append(sprs, newSprChanger("tile_ui_enemy", constants.UIBatch).WithOffset(offset))
					}
				case 3:
					offset.X = 4
					offset.Y = -4
				}
			}
		case data.BlockJetpack, data.BlockDisguise, data.BlockFlamethrower:
			yOffset := 0.
			if tile.Metadata.Regenerate {
				sprs = append(sprs, newSprChanger("tile_ui_regen", constants.UIBatch).WithOffset(pixel.V(-4, 4)))
				yOffset = -4
			}
			timer := tile.Metadata.Timer
			if timer == 0 {
				sprs = append(sprs, newSprChanger(constants.UIInfinity, constants.UIBatch).WithOffset(pixel.V(0, yOffset)))
			} else if timer < 10 {
				sprs = append(sprs, newSprChanger(fmt.Sprintf(constants.UINumber, timer), constants.UIBatch).WithOffset(pixel.V(0, yOffset)))
			} else if timer < 100 {
				sprs = append(sprs, newSprChanger(fmt.Sprintf(constants.UINumberX, timer/10), constants.UIBatch).WithOffset(pixel.V(-3, yOffset)))
				sprs = append(sprs, newSprChanger(fmt.Sprintf(constants.UINumberX, timer%10), constants.UIBatch).WithOffset(pixel.V(3, yOffset)))
			} else if timer < 1000 {
				sprs = append(sprs, newSprChanger(fmt.Sprintf(constants.UINumberX, timer/100), constants.UIBatch).WithOffset(pixel.V(-6, yOffset)))
				sprs = append(sprs, newSprChanger(fmt.Sprintf(constants.UINumberX, timer%100/10), constants.UIBatch).WithOffset(pixel.V(0, yOffset)))
				sprs = append(sprs, newSprChanger(fmt.Sprintf(constants.UINumberX, timer%10), constants.UIBatch).WithOffset(pixel.V(6, yOffset)))
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
		a = data.CurrLevel.Get(tile.Coords.X, tile.Coords.Y+1)
	} else {
		a = data.CurrPuzzleSet.CurrPuzzle.Get(tile.Coords.X, tile.Coords.Y+1)
	}
	if a.IsBlock() || a.Block == data.BlockLiquid {
		top = false
	}
	var b *data.Tile
	if tile.Live {
		b = data.CurrLevel.Get(tile.Coords.X, tile.Coords.Y-1)
	} else {
		b = data.CurrPuzzleSet.CurrPuzzle.Get(tile.Coords.X, tile.Coords.Y-1)
	}
	if b.IsBlock() || b.Block == data.BlockHideout {
		bottom = false
	}
	var sKey string
	if top && bottom {
		sKey = fmt.Sprintf("%s%s", tile.SpriteString(), constants.TileBottomTop)
	} else if top {
		sKey = fmt.Sprintf("%s%s", tile.SpriteString(), constants.TileTop)
	} else if bottom {
		sKey = fmt.Sprintf("%s%s", tile.SpriteString(), constants.TileBottom)
	} else if tile.AltBlock == 1 && tile.Block == data.BlockTurf {
		sKey = fmt.Sprintf("%s%s", tile.SpriteString(), constants.TileAlt)
	} else {
		sKey = tile.SpriteString()
	}
	return sKey
}

func GetSpikeSprite(tile *data.Tile) string {
	// check position to get correct sprite
	bottom := true
	var b *data.Tile
	if tile.Live {
		b = data.CurrLevel.Get(tile.Coords.X, tile.Coords.Y-1)
	} else {
		b = data.CurrPuzzleSet.CurrPuzzle.Get(tile.Coords.X, tile.Coords.Y-1)
	}
	if b == nil || (b.IsBlock() && b.Block != data.BlockSpike) || b.Block == data.BlockHideout {
		bottom = false
	}
	var sKey string
	if bottom {
		sKey = fmt.Sprintf("%s%s", tile.SpriteString(), constants.TileBottom)
	} else {
		sKey = tile.SpriteString()
	}
	return sKey
}

func GetHideoutSprite(tile *data.Tile) string {
	// check position to get correct sprite
	top := true
	var a *data.Tile
	if tile.Live {
		a = data.CurrLevel.Get(tile.Coords.X, tile.Coords.Y+1)
	} else {
		a = data.CurrPuzzleSet.CurrPuzzle.Get(tile.Coords.X, tile.Coords.Y+1)
	}
	if a == nil || a.IsBlock() || a.Block == data.BlockHideout {
		top = false
	}
	var sKey string
	if top {
		sKey = fmt.Sprintf("%s%s", tile.SpriteString(), constants.TileTop)
	} else {
		sKey = tile.SpriteString()
	}
	return sKey
}

func GetLadderSpriteLive(tile *data.Tile) string {
	belowTile := data.CurrLevel.Get(tile.Coords.X, tile.Coords.Y-1)
	aboveTile := data.CurrLevel.Get(tile.Coords.X, tile.Coords.Y+1)
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
	belowTile := data.CurrPuzzleSet.CurrPuzzle.Get(tile.Coords.X, tile.Coords.Y-1)
	aboveTile := data.CurrPuzzleSet.CurrPuzzle.Get(tile.Coords.X, tile.Coords.Y+1)
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

func GetTileSpritesSelection(tile *data.Tile) []any {
	var sprs []any
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
			sprs = append(sprs, img.NewSprite(tile.SpriteString(), constants.TileBatch))
		}
		lStr := GetLadderSpriteSelection(tile)
		if lStr != "" {
			sprs = append(sprs, img.NewSprite(lStr, constants.TileBatch))
		}
		sprChng := GetWrenchSprites(tile)
		for _, sc := range sprChng {
			sprs = append(sprs, img.NewSprite(sc.SprKey, sc.Batch))
		}
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
		sKey = fmt.Sprintf("%s%s", tile.SpriteString(), constants.TileBottomTop)
	} else if top {
		sKey = fmt.Sprintf("%s%s", tile.SpriteString(), constants.TileTop)
	} else if bottom {
		sKey = fmt.Sprintf("%s%s", tile.SpriteString(), constants.TileBottom)
	} else {
		sKey = tile.SpriteString()
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
