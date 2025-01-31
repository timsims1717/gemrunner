package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/util"
	"gemrunner/pkg/world"
)

func DynamicSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsDynamic) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		d, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
		isControlled := result.Entity.HasComponent(myecs.Controller)
		if okO && okC && !obj.Hidden && !isControlled && data.CurrLevel.Start {
			if reanimator.FrameSwitch {
				d.ACounter++
				if !result.Entity.HasComponent(myecs.Parent) {
					currPos := d.Object.Pos
					x, y := world.WorldToMap(currPos.X, currPos.Y)
					currTile := data.CurrLevel.Get(x, y)
					if !d.Flags.Floor {
						if d.Flags.Thrown {
							thrown(d, currTile)
						} else {
							falling(d, currTile)
						}
					} else {
						d.Flags.Thrown = false
						d.Flags.JumpL = false
						d.Flags.JumpR = false
					}
				} else {
					currPos := d.Object.Pos.Add(d.Object.Offset)
					x, y := world.WorldToMap(currPos.X, currPos.Y)
					tile := data.CurrLevel.Get(x, y)
					d.LastTile = tile
				}
			}
		}
	}
}

func SmashSystem() {
	if reanimator.FrameSwitch {
		for _, result := range myecs.Manager.Query(myecs.IsSmash) {
			obj, okO := result.Components[myecs.Object].(*object.Object)
			d, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
			s, okS := result.Components[myecs.Smash].(*float64)
			if okO && okC && okS {
				if d.Flags.Floor ||
					obj.Hidden {
					*s = obj.Pos.Y
				}
			}
		}
	}
}

func thrown(ch *data.Dynamic, tile *data.Tile) {
	if (ch.ACounter > constants.ThrownCounter) ||
		(ch.Flags.LeftWall && ch.Flags.JumpL) ||
		(ch.Flags.RightWall && ch.Flags.JumpR) ||
		ch.Flags.Floor {
		ch.Flags.Thrown = false
		ch.Flags.JumpL = false
		ch.Flags.JumpR = false
	} else {
		if ch.Flags.JumpR {
			if !ch.Flags.RightWall {
				ch.Object.Pos.X += constants.ThrownHSpeed
			}
		} else if ch.Flags.JumpL {
			if !ch.Flags.LeftWall {
				ch.Object.Pos.X -= constants.ThrownHSpeed
			}
		}
		if tile.Coords != ch.LastTile.Coords {
			if util.Abs(tile.Coords.X-ch.LastTile.Coords.X) > 1 {
				ch.Object.Pos.Y -= constants.ThrownVSpeed
			}
		} else {
			ch.Object.Pos.Y += constants.ThrownVSpeed
		}
	}
}
