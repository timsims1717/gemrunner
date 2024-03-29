package systems

import (
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
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
					currTile := data.CurrLevel.Tiles.Get(x, y)
					if !d.Flags.Floor {
						if d.Flags.Throw {
							thrown(d, currTile)
						} else {
							falling(d, currTile)
						}
					} else {
						d.Flags.Throw = false
						d.Flags.JumpL = false
						d.Flags.JumpR = false
					}
				} else {
					currPos := d.Object.Pos.Add(d.Object.Offset)
					x, y := world.WorldToMap(currPos.X, currPos.Y)
					tile := data.CurrLevel.Tiles.Get(x, y)
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
					obj.Hidden ||
					result.Entity.HasComponent(myecs.Parent) {
					*s = obj.Pos.Y
				}
			}
		}
	}
}
