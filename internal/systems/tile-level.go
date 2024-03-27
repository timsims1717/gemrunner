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
				if tile.Block == data.BlockCracked && tile.Flags.Cracked {
					if reanimator.FrameSwitch {
						tile.Counter++
					}
					if tile.Counter > constants.CrackedCounter {
						tile.Flags.Collapse = true
						tile.Flags.Cracked = false
						tile.Counter = 0
						tile.Object.Layer = 31
						obj := object.New()
						obj.Pos = tile.Object.Pos
						obj.Layer = 32
						m := myecs.Manager.NewEntity()
						anim := reanimator.NewSimple(reanimator.NewBatchAnimation("collapse", img.Batchers[constants.TileBatch], "collapse", reanimator.Hold))
						m.AddComponent(myecs.Object, obj)
						m.AddComponent(myecs.Animated, anim)
						m.AddComponent(myecs.Drawable, anim)
						m.AddComponent(myecs.Temp, myecs.ClearFlag(false))
						tile.Mask = m
					}
				} else if tile.Flags.Regen {
					if reanimator.FrameSwitch {
						tile.Counter++
					}
					if tile.Flags.Collapse {
						tile.Flags.Collapse = false
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
					if tile.Counter > constants.RegenACounter {
						if tile.Mask != nil {
							myecs.Manager.DisposeEntity(tile.Mask)
							tile.Mask = nil
							tile.Object.Layer = 10
						}
						tile.Flags.Regen = false
						tile.Flags.Collapse = false
						tile.Flags.Cracked = false
						tile.Counter = 0
					}
				} else if tile.Flags.Collapse {
					if reanimator.FrameSwitch {
						tile.Counter++
					}
					if tile.Counter > constants.CollapseCounter {
						if tile.Mask != nil {
							myecs.Manager.DisposeEntity(tile.Mask)
							tile.Mask = nil
							tile.Object.Layer = 10
						}
						if tile.Metadata.Regenerate {
							if tile.Counter > constants.RegenCounter {
								tile.Flags.Cracked = false
								tile.Flags.Regen = true
								tile.Counter = 0
								tile.Object.Layer = 31
								obj := object.New()
								obj.Pos = tile.Object.Pos
								obj.Layer = 32
								m := myecs.Manager.NewEntity()
								anim := reanimator.NewSimple(reanimator.NewBatchAnimation("collapse", img.Batchers[constants.TileBatch], "regen", reanimator.Hold))
								m.AddComponent(myecs.Object, obj)
								m.AddComponent(myecs.Animated, anim)
								m.AddComponent(myecs.Drawable, anim)
								m.AddComponent(myecs.Temp, myecs.ClearFlag(false))
								tile.Mask = m
							}
						} else {
							tile.Block = data.BlockEmpty
							tile.Flags.Collapse = false
							tile.Flags.Cracked = false
							tile.Counter = 0
						}
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
								} else if tile.Block == data.BlockFall {
									tile.Flags.Collapse = true
									tile.Flags.Cracked = false
									tile.Object.Layer = 31
									obj := object.New()
									obj.Pos = tile.Object.Pos
									obj.Layer = 32
									m := myecs.Manager.NewEntity()
									anim := reanimator.NewSimple(reanimator.NewBatchAnimation("collapse", img.Batchers[constants.TileBatch], "collapse", reanimator.Hold))
									m.AddComponent(myecs.Object, obj)
									m.AddComponent(myecs.Animated, anim)
									m.AddComponent(myecs.Drawable, anim)
									m.AddComponent(myecs.Temp, myecs.ClearFlag(false))
									tile.Mask = m
								}
							}
						}
					}
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
