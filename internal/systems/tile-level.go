package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/world"
)

func TileSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsTile) {
		_, okO := result.Components[myecs.Object].(*object.Object)
		tile, ok := result.Components[myecs.Tile].(*data.Tile)
		if okO && ok && tile.Live {
			switch tile.Block {
			case data.BlockCracked:
				if tile.Flags.Cracked {
					if reanimator.FrameSwitch {
						tile.Counter++
					}
					if tile.Counter > constants.CrackedCollapse {
						tile.Flags.Collapse = true
						tile.Flags.Cracked = false
						tile.Counter = 0
					}
				} else if tile.Flags.Collapse {
					if tile.Metadata.Regenerate {
						if reanimator.FrameSwitch {
							tile.Counter++
						}
						if tile.Counter > constants.CollapseCounter {
							tile.Flags.Cracked = false
							tile.Flags.Collapse = false
						}
					} else {
						tile.Block = data.BlockEmpty
						tile.Flags.Collapse = false
						tile.Flags.Cracked = false
					}
				} else {
					for _, resultC := range myecs.Manager.Query(myecs.IsCharacter) {
						_, okCO := resultC.Components[myecs.Object].(*object.Object)
						ch, okC := resultC.Components[myecs.Dynamic].(*data.Dynamic)
						if okCO && okC {
							x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
							chTile := data.CurrLevel.Tiles.Get(x, y)
							if chTile != nil && chTile.Coords.X == tile.Coords.X &&
								chTile.Coords.Y == tile.Coords.Y+1 &&
								((ch.Player > -1 && ch.Player < constants.MaxPlayers) ||
									tile.Metadata.EnemyCrack) && ch.State == data.Grounded {
								tile.Flags.Cracked = true
								tile.Counter = 0
								tile.Update = true
							}
						}
					}
				}
			}
			switch tile.Ladder {
			case data.BlockLadderCracked:
				if tile.Flags.LCracked {
					if reanimator.FrameSwitch {
						tile.Counter++
					}
					if tile.Counter > constants.CrackedCollapse {
						tile.Flags.LCollapse = true
						tile.Flags.LCracked = false
						tile.Counter = 0
					}
				} else if tile.Flags.LCollapse {
					if tile.Metadata.Regenerate {
						if reanimator.FrameSwitch {
							tile.Counter++
						}
						if tile.Counter > constants.CollapseCounter {
							tile.Flags.LCracked = false
							tile.Flags.LCollapse = false
						}
					} else {
						tile.Ladder = data.BlockEmpty
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
