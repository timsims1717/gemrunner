package object

import (
	"fmt"
	"gemrunner/pkg/util"
	"github.com/gopxl/pixel"
	"golang.org/x/image/colornames"
)

var objIndex = uint32(0)

var ILock bool

type Object struct {
	ID     string
	Hidden bool
	Loaded bool
	Killed bool

	Pos  pixel.Vec
	Mat  pixel.Matrix
	Rot  float64
	Sca  pixel.Vec
	Flip bool
	Flop bool
	Rect pixel.Rect

	HalfWidth  float64
	HalfHeight float64

	PostPos pixel.Vec
	LastPos pixel.Vec
	Offset  pixel.Vec

	Mask  pixel.RGBA
	Layer int
	IntA  int

	ILock        bool
	HideChildren bool
}

func New() *Object {
	return &Object{
		Sca: pixel.Vec{
			X: 1.,
			Y: 1.,
		},
		Mask:  pixel.ToRGBA(colornames.White),
		ILock: ILock,
	}
}

func (obj *Object) WithID(code string) *Object {
	obj.ID = fmt.Sprintf("%s-%d", code, objIndex)
	objIndex++
	return obj
}

func (obj *Object) PointInside(vec pixel.Vec) bool {
	return obj.Rect.Moved(obj.PostPos).Contains(vec)
}

func (obj *Object) SetRect(r pixel.Rect) {
	obj.Rect = util.RectToOrigin(r).Moved(pixel.V(r.W()*-0.5, r.H()*-0.5))
	obj.HalfWidth = obj.Rect.W() * 0.5
	obj.HalfHeight = obj.Rect.H() * 0.5
}
