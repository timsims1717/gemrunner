package systems

import (
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/typeface"
	"gemrunner/pkg/viewport"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func AnimationSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasAnimation) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		theAnim := result.Components[myecs.Animated]
		if okO && !obj.Hidden {
			if theAnim == nil {
				continue
			} else if anims, okS := theAnim.([]*reanimator.Tree); okS {
				for _, anim := range anims {
					anim.Update()
				}
			} else if anim, okA := theAnim.(*reanimator.Tree); okA {
				anim.Update()
			}
		}
	}
}

func DrawSystem(win *pixelgl.Window, layer int) {
	count := 0
	for _, result := range myecs.Manager.Query(myecs.IsDrawable) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		if okO && obj.Layer == layer && !obj.Hidden {
			draw := result.Components[myecs.Drawable]
			var target pixel.Target
			target = win
			if ok := result.Entity.HasComponent(myecs.DrawTarget); ok {
				if tar, okT := result.Entity.GetComponentData(myecs.DrawTarget); okT {
					if vp, okV := tar.(*viewport.ViewPort); okV {
						target = vp.Canvas
					} else {
						target = tar.(pixel.Target)
					}
				}
			}
			if draw == nil {
				continue
			} else if draws, okD := draw.([]*img.Sprite); okD {
				for _, d := range draws {
					DrawThing(d, obj, target)
					count++
				}
			} else if anims, okA := draw.([]*reanimator.Tree); okA {
				for _, d := range anims {
					DrawThing(d, obj, target)
					count++
				}
			} else {
				DrawThing(draw, obj, target)
				count++
			}
		}
	}
	//debug.AddText(fmt.Sprintf("Layer %d: %d entities", layer, count))
}

func DrawThing(draw interface{}, obj *object.Object, target pixel.Target) {
	if spr, ok0 := draw.(*pixel.Sprite); ok0 {
		spr.DrawColorMask(target, obj.Mat, obj.Mask)
	} else if sprH, ok1 := draw.(*img.Sprite); ok1 {
		if sprH.Batch != "" && sprH.Key != "" {
			if batch, okB := img.Batchers[sprH.Batch]; okB {
				batch.DrawSpriteColor(sprH.Key, obj.Mat.Moved(sprH.Offset), sprH.Color.Mul(obj.Mask))
			}
		}
	} else if anim, ok2 := draw.(*reanimator.Tree); ok2 {
		res := anim.CurrentSprite()
		if res != nil {
			if _, okB := img.Batchers[res.Batch]; okB {
				res.Spr.DrawColorMask(img.Batchers[res.Batch].Batch(), obj.Mat.Moved(res.Off), res.Col.Mul(obj.Mask))
			}
		}
	} else if txt, ok3 := draw.(*typeface.Text); ok3 {
		txt.Draw(target)
	}
}