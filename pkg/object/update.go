package object

import (
	"github.com/gopxl/pixel"
	"math"
)

func (obj *Object) Update() {
	obj.LastPos = obj.PostPos
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
	obj.Mat = obj.Mat.Rotated(pixel.ZV, obj.Rot)
	obj.Mat = obj.Mat.Moved(obj.PostPos)
}
