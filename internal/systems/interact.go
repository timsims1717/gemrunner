package systems

import (
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
)

func TouchSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsTouch) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		if okO && !obj.Hidden {
			for _, result1 := range myecs.Manager.Query(myecs.IsDynamic) {
				objC, okOC := result1.Components[myecs.Object].(*object.Object)
				ch, okCh := result1.Components[myecs.Dynamic].(*data.Dynamic)
				if okOC && okCh && !objC.Hidden && obj.ID != objC.ID {
					// the object's rectangle contains the player's position
					if obj.Rect.Moved(obj.Pos).Contains(objC.Pos) {
						// it's touching
						if colFn, okC := result.Components[myecs.OnTouch].(*data.Interact); okC {
							p := -1
							player, okP := result1.Entity.GetComponentData(myecs.Player)
							if okP {
								p = player.(int)
							}
							colFn.Fn(data.CurrLevel, p, ch, result.Entity)
						}
					}
				}
			}
		}
	}
}
