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
	ParentView  *ViewPort

	CamSpeed  float64
	CamAccel  float64
	ZoomSpeed float64
	ZoomStep  float64

	interX *gween.Sequence
	interY *gween.Sequence
	interZ *gween.Sequence
	shakeX *gween.Tween
	shakeY *gween.Tween
	shakeZ *gween.Tween
	velX   float64
	velY   float64
	velZ   float64
	limX   *pixel.Vec
	limY   *pixel.Vec
	limZ   *pixel.Vec

	lock  bool
	Mask  color.RGBA
	iLock bool
}

func New(winCan *pixelgl.Canvas) *ViewPort {
	viewPort := &ViewPort{
		CamSpeed:  50.,
		CamAccel:  1000.,
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
	if !v.updateX() {
		fin = false
	}
	if !v.updateY() {
		fin = false
	}
	if !v.updateZ() {
		fin = false
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

func (v *ViewPort) updateX() bool {
	fin := true
	if v.interX != nil {
		x, _, finX := v.interX.Update(timing.DT)
		v.CamPos.X = x
		if finX {
			v.interX = nil
		} else {
			fin = false
		}
	} else if v.velX != 0. {
		v.CamPos.X += v.velX * timing.DT
		fin = false
	}
	if v.velX > 0 {
		v.velX -= v.CamAccel * timing.DT
		if v.velX < 0 {
			v.velX = 0.
		}
	} else if v.velX < 0 {
		v.velX += v.CamAccel * timing.DT
		if v.velX > 0 {
			v.velX = 0.
		}
	}
	if v.limX != nil {
		if v.CamPos.X > v.limX.X {
			v.CamPos.X = v.limX.X
			fin = true
		} else if v.CamPos.X < v.limX.Y {
			v.CamPos.X = v.limX.Y
			fin = true
		}
	}
	return fin
}

func (v *ViewPort) updateY() bool {
	fin := true
	if v.interY != nil {
		y, _, finY := v.interY.Update(timing.DT)
		v.CamPos.Y = y
		if finY {
			v.interY = nil
		} else {
			fin = false
		}
	} else if v.velY != 0. {
		v.CamPos.Y += v.velY * timing.DT
		fin = false
	}
	if v.velY > 0 {
		v.velY -= v.CamAccel * timing.DT
		if v.velY < 0 {
			v.velY = 0.
		}
	} else if v.velY < 0 {
		v.velY += v.CamAccel * timing.DT
		if v.velY > 0 {
			v.velY = 0.
		}
	}
	if v.limY != nil {
		if v.CamPos.Y > v.limY.X {
			v.CamPos.Y = v.limY.X
			fin = true
		} else if v.CamPos.Y < v.limY.Y {
			v.CamPos.Y = v.limY.Y
			fin = true
		}
	}
	return fin
}

func (v *ViewPort) updateZ() bool {
	fin := true
	if v.interZ != nil {
		z, _, finZ := v.interZ.Update(timing.DT)
		v.Zoom = z
		if finZ {
			v.interZ = nil
		} else {
			fin = false
		}
	} else if v.velZ != 0. {
		v.Zoom += v.velZ * timing.DT
		fin = false
	}
	if v.velZ > 0 {
		v.velZ -= v.CamAccel * timing.DT
		if v.velZ < 0 {
			v.velZ = 0.
		}
	} else if v.velZ < 0 {
		v.velZ += v.CamAccel * timing.DT
		if v.velZ > 0 {
			v.velZ = 0.
		}
	}
	if v.limZ != nil {
		if v.Zoom > v.limZ.X {
			v.Zoom = v.limZ.X
			fin = true
		} else if v.Zoom < v.limZ.Y {
			v.Zoom = v.limZ.Y
			fin = true
		}
	}
	return fin
}

func (v *ViewPort) Draw(target pixel.Target) {
	v.Canvas.DrawColorMask(target, v.Mat, v.Mask)
}

func (v *ViewPort) SetRect(r pixel.Rect) *ViewPort {
	v.Rect = r
	//v.Canvas = pixelgl.NewCanvas(r)
	v.Canvas.SetBounds(r)
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
		if dur > 0. {
			v.interX = gween.NewSequence(gween.New(v.CamPos.X, pos.X, dur, ease.InOutQuad))
			v.interY = gween.NewSequence(gween.New(v.CamPos.Y, pos.Y, dur, ease.InOutQuad))
		} else {
			v.interX = nil
			v.interY = nil
			v.CamPos = pos
		}
		v.lock = lock
	}
}

func (v *ViewPort) AddMove(pos pixel.Vec, dur float64) {
	if v.interX == nil {
		v.MoveTo(v.CamPos.Add(pos), dur, false)
	} else {
		if dur > 0. {
			posX := v.CamPos.X
			posY := v.CamPos.Y
			if len(v.interX.Tweens) > 0 {
				posX = v.interX.Tweens[len(v.interX.Tweens)-1].End
			}
			if len(v.interY.Tweens) > 0 {
				posY = v.interY.Tweens[len(v.interY.Tweens)-1].End
			}
			v.interX.Add(gween.New(posX, posX+pos.X, dur, ease.InOutQuad))
			v.interY.Add(gween.New(posY, posY+pos.Y, dur, ease.InOutQuad))
		}
	}
}

func (v *ViewPort) SetVel(vel pixel.Vec) {
	v.velX = vel.X
	v.velY = vel.Y
}

func (v *ViewPort) RemoveXLim() {
	v.limX = nil
}

func (v *ViewPort) SetXLim(min, max float64) {
	lim := pixel.V(max, min)
	v.limX = &lim
}

func (v *ViewPort) RemoveYLim() {
	v.limY = nil
}

func (v *ViewPort) SetYLim(min, max float64) {
	lim := pixel.V(max, min)
	v.limY = &lim
}

func (v *ViewPort) RemoveZLim() {
	v.limY = nil
}

func (v *ViewPort) SetZLim(min, max float64) {
	lim := pixel.V(max, min)
	v.limZ = &lim
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
		v.interZ = gween.NewSequence(gween.New(v.Zoom, v.TargetZoom, v.ZoomSpeed, ease.OutQuad))
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
	return v.Canvas.Bounds().Contains(vec)
	//return v.Rect.Moved(pixel.V(-(v.Rect.W() * 0.5), -(v.Rect.H() * 0.5))).Contains(v.Mat.Unproject(vec))
}

func (v *ViewPort) ProjectWorld(vec pixel.Vec) pixel.Vec {
	//return v.Mat.Unproject(vec).Add(v.PostCamPos)
	if v.ParentView != nil {
		vec = v.ParentView.ProjectWorld(vec)
	}
	return v.Mat.Unproject(vec).Add(v.PostCamPos)
}

func (v *ViewPort) Project(vec pixel.Vec) pixel.Vec {
	return v.Mat.Unproject(vec).Add(v.PostCamPos)
}

func (v *ViewPort) WorldInside(vec pixel.Vec) (bool, pixel.Vec) {
	vec = v.ProjectWorld(vec)
	if v.PointInside(vec) {
		x := vec.X - v.Canvas.Bounds().Max.X
		if math.Abs(x) > v.Canvas.Bounds().W()*0.5 {
			x = vec.X - v.Canvas.Bounds().Min.X
		}
		y := vec.Y - v.Canvas.Bounds().Max.Y
		if math.Abs(y) > v.Canvas.Bounds().H()*0.5 {
			y = vec.Y - v.Canvas.Bounds().Min.Y
		}
		return true, pixel.V(x, y)
	}
	return false, pixel.ZV
}

func (v *ViewPort) ProjectedOut(vec pixel.Vec) pixel.Vec {
	vec = v.Mat.Project(vec.Add(pixel.V(-v.PostCamPos.X, -v.PostCamPos.Y)))
	return vec
}

func (v *ViewPort) Constrain(vec pixel.Vec) pixel.Vec {
	newPos := vec
	if v.CamPos.X+v.Rect.W()*0.5 < vec.X {
		newPos.X = v.CamPos.X + v.Rect.W()*0.5
	} else if v.CamPos.X-v.Rect.W()*0.5 > vec.X {
		newPos.X = v.CamPos.X - v.Rect.W()*0.5
	}
	if v.CamPos.Y+v.Rect.H()*0.5 < vec.Y {
		newPos.Y = v.CamPos.Y + v.Rect.H()*0.5
	} else if v.CamPos.Y-v.Rect.H()*0.5 > vec.Y {
		newPos.Y = v.CamPos.Y - v.Rect.H()*0.5
	}
	return newPos
}

func (v *ViewPort) ConstrainR(vec pixel.Vec, r pixel.Rect) pixel.Vec {
	newPos := vec
	if v.CamPos.X+v.Rect.W()*0.5 < vec.X+r.W()*0.5 {
		newPos.X = v.CamPos.X + v.Rect.W()*0.5 - r.W()*0.5
	} else if v.CamPos.X-v.Rect.W()*0.5 > vec.X-r.W()*0.5 {
		newPos.X = v.CamPos.X - v.Rect.W()*0.5 + r.W()*0.5
	}
	if v.CamPos.Y+v.Rect.H()*0.5 < vec.Y+r.H()*0.5 {
		newPos.Y = v.CamPos.Y + v.Rect.H()*0.5 - r.H()*0.5
	} else if v.CamPos.Y-v.Rect.H()*0.5 > vec.Y-r.H()*0.5 {
		newPos.Y = v.CamPos.Y - v.Rect.H()*0.5 + r.H()*0.5
	}
	return newPos
}

func (v *ViewPort) GetLimX() (float64, float64) {
	return v.limX.X, v.limX.Y
}

func (v *ViewPort) GetLimY() (float64, float64) {
	return v.limY.X, v.limY.Y
}

func (v *ViewPort) GetLimZ() (float64, float64) {
	return v.limZ.X, v.limZ.Y
}
