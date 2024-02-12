package systems

import (
	"fmt"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/typeface"
	"gemrunner/pkg/util"
	"github.com/gopxl/pixel"
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

var currBatches []string

func DrawBatchSystem(target pixel.Target, batchKey string, layers []int) {
	batch := img.Batchers[batchKey]
	if batch == nil {
		fmt.Println("WARNING: Batch with key", batchKey, "does not exist")
		return
	}
	count := 0
	for _, layer := range layers {
		for _, result := range myecs.Manager.Query(myecs.IsDrawable) {
			obj, okO := result.Components[myecs.Object].(*object.Object)
			if okO && obj.Layer == layer && !obj.Hidden {
				draw := result.Components[myecs.Drawable]
				if draw == nil {
					continue
				} else if draws, okD := draw.([]*img.Sprite); okD {
					for _, d := range draws {
						if d == nil {
							continue
						}
						DrawBatchThing(d, obj, batch)
						count++
					}
				} else if anims, okA := draw.([]*reanimator.Tree); okA {
					for _, d := range anims {
						if d == nil {
							continue
						}
						DrawBatchThing(d, obj, batch)
						count++
					}
				} else {
					DrawBatchThing(draw, obj, batch)
					count++
				}
			}
		}
	}
	batch.Draw(target)
	//debug.AddText(fmt.Sprintf("Layer %d: %d entities", layer, count))
}

func DrawBatchThing(draw interface{}, obj *object.Object, batch *img.Batcher) {
	if sprH, ok1 := draw.(*img.Sprite); ok1 {
		if sprH.Batch != "" && sprH.Batch == batch.Key && sprH.Key != "" {
			batch.DrawSpriteColor(sprH.Key, obj.Mat.Moved(sprH.Offset), sprH.Color.Mul(obj.Mask))
		}
	} else if anim, ok2 := draw.(*reanimator.Tree); ok2 {
		res := anim.CurrentSprite()
		if res != nil && res.Batch == batch.Key {
			res.Spr.DrawColorMask(img.Batchers[res.Batch].Batch(), obj.Mat.Moved(res.Off), res.Col.Mul(obj.Mask))
		}
	}
}

func DrawLayerSystem(target pixel.Target, layer int) {
	currBatches = []string{}
	count := 0
	for _, result := range myecs.Manager.Query(myecs.IsDrawable) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		if okO && obj.Layer == layer && !obj.Hidden {
			draw := result.Components[myecs.Drawable]
			if draw == nil {
				continue
			} else if draws, okD := draw.([]*img.Sprite); okD {
				for _, d := range draws {
					if d == nil {
						continue
					}
					DrawThing(d, obj, target)
					count++
				}
			} else if anims, okA := draw.([]*reanimator.Tree); okA {
				for _, d := range anims {
					if d == nil {
						continue
					}
					DrawThing(d, obj, target)
					count++
				}
			} else {
				DrawThing(draw, obj, target)
				count++
			}
		}
	}
	for _, batch := range currBatches {
		img.Batchers[batch].Draw(target)
	}
	//debug.AddText(fmt.Sprintf("Layer %d: %d entities", layer, count))
}

func DrawThing(draw interface{}, obj *object.Object, target pixel.Target) {
	if spr, ok0 := draw.(*pixel.Sprite); ok0 {
		spr.DrawColorMask(target, obj.Mat, obj.Mask)
	} else if sprH, ok1 := draw.(*img.Sprite); ok1 {
		if sprH.Batch != "" && sprH.Key != "" {
			if batch, okB := img.Batchers[sprH.Batch]; okB {
				if !util.ContainsStr(sprH.Batch, currBatches) {
					currBatches = append(currBatches, sprH.Batch)
				}
				batch.DrawSpriteColor(sprH.Key, obj.Mat.Moved(sprH.Offset), sprH.Color.Mul(obj.Mask))
			}
		}
	} else if anim, ok2 := draw.(*reanimator.Tree); ok2 {
		res := anim.CurrentSprite()
		if res != nil {
			if _, okB := img.Batchers[res.Batch]; okB {
				if !util.ContainsStr(res.Batch, currBatches) {
					currBatches = append(currBatches, res.Batch)
				}
				res.Spr.DrawColorMask(img.Batchers[res.Batch].Batch(), obj.Mat.Moved(res.Off), res.Col.Mul(obj.Mask))
			}
		}
	} else if txt, ok3 := draw.(*typeface.Text); ok3 {
		txt.Draw(target)
	}
}
