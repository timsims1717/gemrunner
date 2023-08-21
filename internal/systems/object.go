package systems

import (
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"github.com/faiface/pixel"
	"math"
)

func ObjectSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsObject) {
		if obj, ok := result.Components[myecs.Object].(*object.Object); ok {
			if obj.Kill {
				myecs.Manager.DisposeEntity(result)
			} else {
				obj.PostPos = obj.Pos.Add(obj.Offset)
				if obj.ILock {
					obj.PostPos.X = math.Round(obj.PostPos.X)
					obj.PostPos.Y = math.Round(obj.PostPos.Y)
				}
				obj.Mat = pixel.IM
				if obj.Flip && obj.Flop {
					obj.Mat = obj.Mat.Scaled(pixel.ZV, -1.)
				} else if obj.Flip {
					obj.Mat = obj.Mat.ScaledXY(pixel.ZV, pixel.V(-1., 1.))
				} else if obj.Flop {
					obj.Mat = obj.Mat.ScaledXY(pixel.ZV, pixel.V(1., -1.))
				}
				obj.Mat = obj.Mat.ScaledXY(pixel.ZV, obj.Sca)
				obj.Mat = obj.Mat.Rotated(pixel.ZV, math.Pi*obj.Rot)
				obj.Mat = obj.Mat.Moved(obj.PostPos)
			}
		}
	}
}

func ParentSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasParent) {
		tran, okT := result.Components[myecs.Object].(*object.Object)
		parent, okP := result.Components[myecs.Parent].(*object.Object)
		if okT && okP {
			if parent.Kill {
				myecs.Manager.DisposeEntity(result)
			} else {
				tran.Pos = parent.Pos
				tran.Hide = parent.Hide
			}
		}
	}
}
