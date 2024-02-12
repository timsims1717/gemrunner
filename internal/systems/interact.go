package systems

import (
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
)

func CollectSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsTouch) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		if okO && !obj.Hidden {
			for _, result1 := range myecs.Manager.Query(myecs.IsPlayer) {
				objC, okOC := result1.Components[myecs.Object].(*object.Object)
				ch, okCh := result1.Components[myecs.Dynamic].(*data.Dynamic)
				p, okP := result1.Components[myecs.Player].(data.Player)
				if okOC && okP && okCh && !objC.Hidden {
					// the object's rectangle contains the player's position
					if obj.Rect.Moved(obj.Pos).Contains(objC.Pos) {
						// it's touching
						if colFn, okC := result.Components[myecs.OnTouch].(*data.Interact); okC {
							colFn.Fn(data.CurrLevel, int(p), ch, result.Entity)
						}
					}
				}
			}
		}
	}
}
