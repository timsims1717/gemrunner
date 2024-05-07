package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/world"
)

func TileSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsTile) {
		_, okO := result.Components[myecs.Object].(*object.Object)
		tile, ok := result.Components[myecs.Tile].(*data.Tile)
		if okO && ok && data.CurrLevel.Start && tile.Live {
			if tile.Flags.Regen { // if this tile is regenerating
				if reanimator.FrameSwitch {
					tile.Counter++
				}
				if tile.Counter > constants.RegenACounter {
					RemoveMask(tile)
					tile.Flags.Regen = false
					tile.Flags.Collapse = false
					tile.Flags.Cracked = false
					tile.Counter = 0
				}
			} else if tile.Flags.Collapse { // id this tile has collapsed
				if reanimator.FrameSwitch {
					tile.Counter++
				}
				if tile.Counter > constants.CollapseCounter {
					RemoveMask(tile)
					if tile.Metadata.Regenerate {
						if tile.Counter > constants.RegenCounter {
							tile.Flags.Collapse = false
							tile.Flags.Cracked = false
							tile.Flags.Regen = true
							tile.Counter = 0
							AddMask(tile, "regen_mask", false, false)
							// Crush any characters here
							for _, resultC := range myecs.Manager.Query(myecs.IsCharacter) {
								_, okCO := resultC.Components[myecs.Object].(*object.Object)
								ch, okC := resultC.Components[myecs.Dynamic].(*data.Dynamic)
								if okCO && okC && ch.State != data.Dead {
									x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
									chTile := data.CurrLevel.Tiles.Get(x, y)
									if chTile != nil && chTile.Coords.X == tile.Coords.X &&
										(chTile.Coords.Y == tile.Coords.Y) {
										ch.Object.Pos.X = tile.Object.Pos.X
										ch.Object.Pos.Y = tile.Object.Pos.Y
										ch.Flags.Crush = true
										ch.State = data.Hit
									}
								}
							}
						}
					} else {
						tile.Block = data.BlockEmpty
						tile.Flags.Collapse = false
						tile.Flags.Cracked = false
						tile.Counter = 0
					}
				}
			}
			switch tile.Block {
			case data.BlockSpike:
				for _, resultC := range myecs.Manager.Query(myecs.IsCharacter) {
					_, okCO := resultC.Components[myecs.Object].(*object.Object)
					ch, okC := resultC.Components[myecs.Dynamic].(*data.Dynamic)
					if okCO && okC && ch.State == data.Grounded {
						x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
						chTile := data.CurrLevel.Tiles.Get(x, y)
						if chTile != nil && chTile.Coords.X == tile.Coords.X &&
							(chTile.Coords.Y == tile.Coords.Y+1 || chTile.Coords.Y == tile.Coords.Y) {
							ch.Object.Pos.X = tile.Object.Pos.X
							ch.Object.Pos.Y = tile.Object.Pos.Y + world.TileSize*0.75
							ch.Flags.Hit = true
							ch.State = data.Hit
						}
					}
				}
			case data.BlockCracked, data.BlockFall:
				if !tile.Flags.Collapse && !tile.Flags.Regen {
					if tile.Block == data.BlockCracked && tile.Flags.Cracked {
						if reanimator.FrameSwitch {
							tile.Counter++
						}
						if tile.Counter > constants.CrackedCounter {
							tile.Flags.Collapse = true
							tile.Flags.Cracked = false
							tile.Counter = 0
							AddMask(tile, "collapse_mask", false, false)
						}
					} else {
						for _, resultC := range myecs.Manager.Query(myecs.IsCharacter) {
							_, okCO := resultC.Components[myecs.Object].(*object.Object)
							ch, okC := resultC.Components[myecs.Dynamic].(*data.Dynamic)
							if okCO && okC && ch.State == data.Grounded {
								x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
								chTile := data.CurrLevel.Tiles.Get(x, y)
								if chTile != nil && chTile.Coords.X == tile.Coords.X &&
									chTile.Coords.Y == tile.Coords.Y+1 &&
									((ch.Player > -1 && ch.Player < constants.MaxPlayers) ||
										tile.Metadata.EnemyCrack) {
									tile.Counter = 0
									tile.Update = true
									tile.Flags.Regen = false
									if tile.Block == data.BlockCracked {
										tile.Flags.Cracked = true
										tile.Flags.Collapse = false
										AddMask(tile, "cracking_mask", false, false)
									} else if tile.Block == data.BlockFall {
										tile.Flags.Collapse = true
										tile.Flags.Cracked = false
										AddMask(tile, "collapse_mask", false, false)
									}
								}
							}
						}
					}
				}
			case data.BlockPhase:
				if reanimator.FrameSwitch {
					tile.Counter++
				}
				if data.CurrLevel.FrameChange {
					phaseCycle := (data.CurrLevel.FrameCycle % 8) - tile.Metadata.Phase
					if phaseCycle == 0 {
						tile.Flags.PhaseChange = true
					} else if phaseCycle == 4 || phaseCycle == -4 {
						tile.Flags.PhaseChange = true
					}
				}
				if reanimator.FrameSwitch && tile.Flags.PhaseChange {
					if tile.Flags.Collapse {
						tile.Counter = 0
						tile.Flags.Collapse = false
						AddMaskWithTrigger(tile, "phase_1_mask", false, true, func() {
							RemoveMask(tile)
						})
						// Crush any characters here
						for _, resultC := range myecs.Manager.Query(myecs.IsCharacter) {
							_, okCO := resultC.Components[myecs.Object].(*object.Object)
							ch, okC := resultC.Components[myecs.Dynamic].(*data.Dynamic)
							if okCO && okC && ch.State != data.Dead {
								x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
								chTile := data.CurrLevel.Tiles.Get(x, y)
								if chTile != nil && chTile.Coords.X == tile.Coords.X &&
									(chTile.Coords.Y == tile.Coords.Y) {
									ch.Object.Pos.X = tile.Object.Pos.X
									ch.Object.Pos.Y = tile.Object.Pos.Y
									ch.Flags.Crush = true
									ch.State = data.Hit
								}
							}
						}
					} else {
						tile.Counter = 0
						tile.Flags.Collapse = true
						AddMaskWithTrigger(tile, "phase_1_mask", false, false, func() {
							RemoveMask(tile)
						})
					}
					tile.Flags.PhaseChange = false
				}
			case data.BlockLadderCracked, data.BlockLadderCrackedTurf:
				if tile.Flags.LCracked {
					if reanimator.FrameSwitch {
						tile.Counter++
					}
					if tile.Counter > constants.CrackedCounter {
						tile.Flags.LCollapse = true
						tile.Flags.LCracked = false
						tile.Counter = 0
					}
				} else if tile.Flags.LCollapse {
					if tile.Metadata.Regenerate {
						if reanimator.FrameSwitch {
							tile.Counter++
						}
						if tile.Counter > constants.RegenCounter {
							tile.Flags.LCracked = false
							tile.Flags.LCollapse = false
						}
					} else {
						if tile.Block == data.BlockLadderCrackedTurf {
							tile.Block = data.BlockTurf
						} else {
							tile.Block = data.BlockEmpty
						}
						tile.Flags.LCollapse = false
						tile.Flags.LCracked = false
					}
				} else {
					for _, resultC := range myecs.Manager.Query(myecs.IsCharacter) {
						_, okCO := resultC.Components[myecs.Object].(*object.Object)
						ch, okC := resultC.Components[myecs.Dynamic].(*data.Dynamic)
						if okCO && okC {
							x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
							chTile := data.CurrLevel.Tiles.Get(x, y)
							if chTile != nil && chTile.Coords.X == tile.Coords.X &&
								chTile.Coords.Y == tile.Coords.Y &&
								((ch.Player > -1 && ch.Player < constants.MaxPlayers) ||
									tile.Metadata.EnemyCrack) &&
								(ch.State == data.OnLadder || ch.State == data.Leaping) {
								tile.Flags.LCracked = true
								tile.Counter = 0
								tile.Update = true
							}
						}
					}
				}
			}
		}
	}
}

// AddMask creates a new mask animation for the tile and sets the correct layers
// for drawing.
func AddMask(tile *data.Tile, maskKey string, flip, reverse bool) {
	RemoveMask(tile)
	tile.Object.Layer = 31 // Put the tile on top so it's over any characters
	obj := object.New()
	obj.Pos = tile.Object.Pos
	obj.Flip = flip
	obj.Layer = 32 // the mask is one layer higher
	m := myecs.Manager.NewEntity()
	a := reanimator.NewBatchAnimation(maskKey, img.Batchers[constants.TileBatch], maskKey, reanimator.Hold)
	if reverse {
		a = a.Reverse()
	}
	anim := reanimator.NewSimple(a)
	m.AddComponent(myecs.Object, obj)
	m.AddComponent(myecs.Animated, anim)
	m.AddComponent(myecs.Drawable, anim)
	m.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	tile.Mask = m
}

// AddMaskWithTrigger creates a new mask animation for the tile and sets the correct layers
// for drawing. It also sets a trigger to run when the mask is done animating.
func AddMaskWithTrigger(tile *data.Tile, maskKey string, flip, reverse bool, fn func()) {
	RemoveMask(tile)
	tile.Object.Layer = 31 // Put the tile on top so it's over any characters
	obj := object.New()
	obj.Pos = tile.Object.Pos
	obj.Flip = flip
	obj.Layer = 32 // the mask is one layer higher
	m := myecs.Manager.NewEntity()
	a := reanimator.NewBatchAnimation(maskKey, img.Batchers[constants.TileBatch], maskKey, reanimator.Tran)
	if reverse {
		a = a.Reverse()
	}
	a.SetEndTrigger(fn)
	anim := reanimator.NewSimple(a)
	m.AddComponent(myecs.Object, obj)
	m.AddComponent(myecs.Animated, anim)
	m.AddComponent(myecs.Drawable, anim)
	m.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	tile.Mask = m
}

func RemoveMask(tile *data.Tile) {
	if tile.Mask != nil {
		myecs.Manager.DisposeEntity(tile.Mask)
		tile.Mask = nil
		tile.Object.Layer = 10
	}
}
