package systems

import (
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
)

func EnemyCollisionSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsDynamic) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		ch, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
		if okO && okC && !obj.Hidden && ch.Enemy > -1 &&
			ch.State != data.Hit && ch.State != data.Dead {
			setEnemyCollisionFlags(ch)
			chPos := ch.Object.Pos
			pPos := ch.Object.LastPos
			for _, result1 := range myecs.Manager.Query(myecs.IsDynamic) {
				obj1, okO1 := result1.Components[myecs.Object].(*object.Object)
				ch1, okC1 := result1.Components[myecs.Dynamic].(*data.Dynamic)
				if okO1 && okC1 && !obj1.Hidden && ch1.Enemy > -1 &&
					ch1.State != data.Hit && ch1.State != data.Dead {
					if obj.Rect.Intersects(obj1.Rect) {

					}
				}
			}
		}
	}
}

func setEnemyCollisionFlags(ch *data.Dynamic) {
	ch.Flags.EnemyL = false
	ch.Flags.EnemyR = false
	ch.Flags.EnemyU = false
	ch.Flags.EnemyD = false
}
