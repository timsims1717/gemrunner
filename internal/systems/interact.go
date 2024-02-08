package systems

import (
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
)

func CollectSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsInteract) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		if okO && !obj.Hidden {
			for _, result1 := range myecs.Manager.Query(myecs.IsInteractor) {
				objC, okOC := result1.Components[myecs.Object].(*object.Object)
				ch, okCh := result1.Components[myecs.Character].(*data.Character)
				if okOC && okCh && !objC.Hidden {
					if obj.Rect.Moved(obj.Pos).Contains(objC.Pos) {
						// it's collected
						if colFn, okC := result.Components[myecs.Interact].(*data.OnTouch); okC {
							colFn.Fn(data.CurrLevel, ch, result.Entity)
						}
					}
				}
			}
		}
	}
}
