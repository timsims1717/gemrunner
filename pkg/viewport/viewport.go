package viewport

import (
	gween "gemrunner/pkg/gween64"
	"gemrunner/pkg/gween64/ease"
	"gemrunner/pkg/timing"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
	"math/rand"
)

var (
	MainCamera *ViewPort
)

type ViewPort struct {
	Canvas     *pixelgl.Canvas
	Rect       pixel.Rect
	CamPos     pixel.Vec
	PostCamPos pixel.Vec
	Zoom       float64
	TargetZoom float64

	Mat         pixel.Matrix
	PortPos     pixel.Vec
	PostPortPos pixel.Vec
	PostZoom    float64
	PortSize    pixel.Vec

	CamSpeed  float64
	ZoomSpeed float64
	ZoomStep  float64

	interX *gween.Tween
	interY *gween.Tween
	interZ *gween.Tween
	shakeX *gween.Tween
	shakeY *gween.Tween
	shakeZ *gween.Tween

	lock  bool
	Mask  color.RGBA
	iLock bool
}

func New(winCan *pixelgl.Canvas) *ViewPort {
	viewPort := &ViewPort{
		CamSpeed:  50.,
		ZoomSpeed: 1.,
		ZoomStep:  1.2,
		PortSize:  pixel.V(1., 1.),
	}
	viewPort.SetZoom(1.)
	if winCan == nil {
		viewPort.Canvas = pixelgl.NewCanvas(pixel.R(0, 0, 0, 0))
	} else {
		viewPort.Canvas = winCan
	}
	viewPort.Mask = colornames.White
	return viewPort
}

func (v *ViewPort) Update() {
	fin := true
	if v.interX != nil {
		x, finX := v.interX.Update(timing.DT)
		v.CamPos.X = x
		if finX {
			v.interX = nil
		} else {
			fin = false
		}
	}
	if v.interY != nil {
		y, finY := v.interY.Update(timing.DT)
		v.CamPos.Y = y
		if finY {
			v.interY = nil
		} else {
			fin = false
		}
	}
	if v.interZ != nil {
		z, finZ := v.interZ.Update(timing.DT)
		v.Zoom = z
		if finZ {
			v.interZ = nil
		} else {
			fin = false
		}
	}
	if fin && v.lock {
		v.lock = false
	}
	v.PostCamPos = v.CamPos
	if v.shakeX != nil {
		x, finSX := v.shakeX.Update(timing.DT)
		v.PostCamPos.X += x
		if finSX {
			v.shakeX = nil
		}
	}
	if v.shakeY != nil {
		y, finSY := v.shakeY.Update(timing.DT)
		v.PostCamPos.Y += y
		if finSY {
			v.shakeY = nil
		}
	}
	v.PostZoom = v.Zoom
	if v.shakeZ != nil {
		z, finSZ := v.shakeZ.Update(timing.DT)
		v.PostZoom += z
		if finSZ {
			v.shakeZ = nil
		}
	}
	v.PostPortPos = v.PortPos
	if v.iLock {
		v.PostCamPos.X = math.Round(v.PostCamPos.X)
		v.PostCamPos.Y = math.Round(v.PostCamPos.Y)
		v.PostPortPos.X = math.Round(v.PostPortPos.X)
		v.PostPortPos.Y = math.Round(v.PostPortPos.Y)
	}

	hw := v.Rect.W() * 0.5 * (1 / v.Zoom)
	hh := v.Rect.H() * 0.5 * (1 / v.Zoom)
	var r pixel.Rect
	if v.iLock {
		r = pixel.R(math.Round(v.PostCamPos.X-hw), math.Round(v.PostCamPos.Y-hh), math.Round(v.PostCamPos.X+hw), math.Round(v.PostCamPos.Y+hh))
	} else {
		r = pixel.R(v.PostCamPos.X-hw, v.PostCamPos.Y-hh, v.PostCamPos.X+hw, v.PostCamPos.Y+hh)
	}
	v.Canvas.SetBounds(r)
	v.Mat = pixel.IM.ScaledXY(pixel.ZV, v.PortSize).Scaled(pixel.ZV, v.Zoom).Moved(v.PostPortPos)
	v.Canvas.SetColorMask(v.Mask)
}

func (v *ViewPort) SetRect(r pixel.Rect) *ViewPort {
	v.Rect = r
	v.Canvas = pixelgl.NewCanvas(r)
	return v
}

func (v *ViewPort) Stop() {
	v.lock = false
	v.interX = nil
	v.interY = nil
}

func (v *ViewPort) SnapTo(pos pixel.Vec) {
	if !v.lock {
		v.CamPos.X = pos.X
		v.CamPos.Y = pos.Y
	}
}

func (v *ViewPort) MoveTo(pos pixel.Vec, dur float64, lock bool) {
	if !v.lock {
		v.interX = gween.New(v.CamPos.X, pos.X, dur, ease.InOutQuad)
		v.interY = gween.New(v.CamPos.Y, pos.Y, dur, ease.InOutQuad)
		v.lock = lock
	}
}

func (v *ViewPort) Follow(pos pixel.Vec, spd float64) {
	if !v.lock {
		v.CamPos.X += spd * timing.DT * (pos.X - v.CamPos.X)
		v.CamPos.Y += spd * timing.DT * (pos.Y - v.CamPos.Y)
	}
}

func (v *ViewPort) CamLeft() {
	if !v.lock {
		v.CamPos.X -= v.CamSpeed * timing.DT
	}
}

func (v *ViewPort) CamRight() {
	if !v.lock {
		v.CamPos.X += v.CamSpeed * timing.DT
	}
}

func (v *ViewPort) CamDown() {
	if !v.lock {
		v.CamPos.Y -= v.CamSpeed * timing.DT
	}
}

func (v *ViewPort) CamUp() {
	if !v.lock {
		v.CamPos.Y += v.CamSpeed * timing.DT
	}
}

func (v *ViewPort) SetZoom(zoom float64) {
	v.Zoom = zoom
	v.TargetZoom = zoom
}

func (v *ViewPort) ZoomIn(zoom float64) {
	if !v.lock {
		v.TargetZoom *= math.Pow(v.ZoomStep, zoom)
		v.interZ = gween.New(v.Zoom, v.TargetZoom, v.ZoomSpeed, ease.OutQuad)
	}
}

func (v *ViewPort) SetILock(b bool) {
	v.iLock = b
}

func (v *ViewPort) SetColor(col color.RGBA) {
	v.Mask = col
}

func (v *ViewPort) Shake(dur, freq float64) {
	v.shakeX = gween.New((rand.Float64()-0.5)*8., 0., dur, SetSine(freq))
	v.shakeY = gween.New((rand.Float64()-0.5)*8., 0., dur, SetSine(freq))
}

func (v *ViewPort) ZoomShake(dur, freq float64) {
	v.shakeZ = gween.New(0.02, 0., dur, SetSine(freq))
}

func SetSine(freq float64) func(float64, float64, float64, float64) float64 {
	return func(t, b, c, d float64) float64 {
		return b * math.Pow(math.E, -math.Abs(c)*t) * math.Sin(freq*math.Pi*t)
	}
}

func Sine(t, b, c, d float64) float64 {
	return b * math.Pow(math.E, -math.Abs(c)*t) * math.Sin(10.*math.Pi*t)
}

func (v *ViewPort) PointInside(vec pixel.Vec) bool {
	return v.Rect.Moved(pixel.V(-(v.Rect.W() * 0.5), -(v.Rect.H() * 0.5))).Contains(v.Mat.Unproject(vec))
}

func (v *ViewPort) Projected(vec pixel.Vec) pixel.Vec {
	//fmt.Printf("V: (%f,%f)\n", vec.X, vec.Y)
	a := v.Mat.Unproject(vec).Add(v.PostCamPos)
	//fmt.Printf("A: (%f,%f)\n", a.X, a.Y)
	//b := v.Mat.Unproject(vec.Scaled(1 / v.Zoom))
	//fmt.Printf("B: (%f,%f)\n", b.X, b.Y)
	return a
}
