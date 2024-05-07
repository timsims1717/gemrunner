package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
)

func FloatingTextSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsText) {
		_, okO := result.Components[myecs.Object].(*object.Object)
		tf, okTF := result.Components[myecs.Text].(*data.FloatingText)
		if okO && okTF {
			if reanimator.FrameSwitch {
				tf.Counter++
			}
			// temp
			if tf.Temp && tf.Counter > constants.TextTimer {
				myecs.Manager.DisposeEntity(result.Entity)
				continue
			}
			// prox

		}
	}
}
