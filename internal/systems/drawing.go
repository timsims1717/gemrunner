package systems

import (
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func AnimationSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasAnimation) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		anim, ok := result.Components[myecs.Animation].(*reanimator.Tree)
		if okO && ok && !obj.Hidden {
			anim.Update()
		}
	}
}

func DrawSystem(win *pixelgl.Window, layer int) {
	for _, result := range myecs.Manager.Query(myecs.IsDrawable) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		if okO && obj.Layer == layer {
			draw := result.Components[myecs.Drawable]
			if draw == nil {
				continue
			} else if draws, okD := draw.([]interface{}); okD {
				for _, d := range draws {
					DrawThing(d, obj, win)
				}
			} else {
				DrawThing(draw, obj, win)
			}
		}
	}
}

func DrawThing(draw interface{}, obj *object.Object, win *pixelgl.Window) {
	if spr, ok0 := draw.(*pixel.Sprite); ok0 {
		spr.Draw(win, obj.Mat)
	} else if sprH, ok1 := draw.(*img.Sprite); ok1 {
		if batch, okB := img.Batchers[sprH.Batch]; okB {
			batch.DrawSpriteColor(sprH.Key, obj.Mat.Moved(sprH.Offset), sprH.Color)
		}
	} else if anim, ok2 := draw.(*reanimator.Tree); ok2 {
		sprA := anim.CurrentSprite()
		if batch, okB := img.Batchers[sprA.Batch]; okB {
			batch.DrawSpriteColor(sprA.Key, obj.Mat.Moved(sprA.Offset), sprA.Color)
		}
	}
}
