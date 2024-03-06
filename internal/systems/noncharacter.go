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
			if !result.Entity.HasComponent(myecs.Parent) {
				if reanimator.FrameSwitch {
					currPos := d.Object.Pos
					x, y := world.WorldToMap(currPos.X, currPos.Y)
					currTile := data.CurrLevel.Tiles.Get(x, y)
					if !d.Flags.Floor {
						falling(d, currTile)
					}
				}
			}
		}
	}
}
