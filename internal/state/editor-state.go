package state

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/systems"
	"gemrunner/pkg/img"
	"gemrunner/pkg/state"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/util"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
	"image/color"
)

var (
	editorInput = &pxginput.Input{
		Buttons: map[string]*pxginput.ButtonSet{
			"click":      pxginput.NewJoyless(pixelgl.MouseButtonLeft),
			"rightClick": pxginput.NewJoyless(pixelgl.MouseButtonRight),
		},
		Mode: pxginput.KeyboardMouse,
	}
)

type editorState struct {
	*state.AbstractState
	PuzzleView *viewport.ViewPort
	BorderView *viewport.ViewPort
}

func (s *editorState) Unload() {

}

func (s *editorState) Load(done chan struct{}) {
	data.CurrPuzzle = data.CreateTestPuzzle()
	systems.PuzzleInit()
	systems.EditorInit()
	s.PuzzleView = viewport.New(nil)
	s.PuzzleView.SetRect(pixel.R(0, 0, world.TileSize*constants.PuzzleWidth, world.TileSize*constants.PuzzleHeight))
	s.PuzzleView.CamPos = pixel.V(world.TileSize*0.5*(constants.PuzzleWidth-1), world.TileSize*0.5*(constants.PuzzleHeight-1))
	s.BorderView = viewport.New(nil)
	s.BorderView.SetRect(pixel.R(0, 0, world.TileSize*(constants.PuzzleWidth+1), world.TileSize*(constants.PuzzleHeight+1)))
	s.BorderView.CamPos = pixel.V(world.TileSize*0.5*(constants.PuzzleWidth), world.TileSize*0.5*(constants.PuzzleHeight))
	s.UpdateViews()
	done <- struct{}{}
}

func (s *editorState) Update(win *pixelgl.Window) {
	if debugInput.Get("camUp").Pressed() {
		s.PuzzleView.PortSize.Y += 1. * timing.DT
	} else if debugInput.Get("camDown").Pressed() {
		s.PuzzleView.PortSize.Y -= 1. * timing.DT
	}
	if debugInput.Get("camRight").Pressed() {
		s.PuzzleView.PortSize.X += 1. * timing.DT
	} else if debugInput.Get("camLeft").Pressed() {
		s.PuzzleView.PortSize.X -= 1. * timing.DT
	}
	if debugInput.Get("debugSP").JustPressed() {
		s.PuzzleView.ZoomIn(1.)
	} else if debugInput.Get("debugSM").JustPressed() {
		s.PuzzleView.ZoomIn(-1.)
	}

	editorInput.Update(win, viewport.MainCamera.Mat)

	data.EditorPane.Hover = data.EditorPane.ViewPort.PointInside(editorInput.World) || (data.EditorPane.SelectVis && data.EditorPane.BlockSelect.PointInside(editorInput.World))
	if data.EditorPane.Hover {
		pos := data.EditorPane.ViewPort.Projected(editorInput.World)
		bsObj := data.EditorPane.BlockView.Object
		if bsObj.Rect.Moved(bsObj.Pos).Contains(pos) {
			if editorInput.Get("click").Pressed() {
				data.EditorPane.SelectVis = true
			}
		}
	} else {
		data.EditorPane.SelectVis = false
		if editorInput.Get("rightClick").Pressed() {
			if s.PuzzleView.PointInside(editorInput.World) {
				projPos := util.ProjectedPoint(editorInput.World, s.PuzzleView.Rect, s.PuzzleView.Mat)
				x, y := world.WorldToMap(projPos.X, projPos.Y)
				coords := world.Coords{X: x, Y: y}
				systems.DeleteBlock(coords)
			}
		} else if editorInput.Get("click").Pressed() {
			if s.PuzzleView.PointInside(editorInput.World) {
				projPos := util.ProjectedPoint(editorInput.World, s.PuzzleView.Rect, s.PuzzleView.Mat)
				x, y := world.WorldToMap(projPos.X, projPos.Y)
				coords := world.Coords{X: x, Y: y}
				systems.ChangeBlock(coords, data.RedRock)
			}
		}
	}

	systems.HoverSystem(editorInput)
	systems.TileSystem()
	systems.ParentSystem()
	systems.ObjectSystem()
	systems.AnimationSystem()

	//s.UpdateViews()

	s.BorderView.Update()
	s.PuzzleView.Update()
	data.EditorPane.ViewPort.Update()
	bSelect := data.EditorPane.BlockSelect
	bSelect.PortSize = data.EditorPane.ViewPort.PortSize
	bSelect.PortPos = data.EditorPane.ViewPort.PortPos
	bSelect.PortPos.X += (bSelect.Canvas.Bounds().W() + world.TileSize + 8.) * bSelect.PortSize.X * 0.5
	bSelect.PortPos.Y -= (data.EditorPane.ViewPort.Canvas.Bounds().H() - world.TileSize + 6.) * bSelect.PortSize.Y * 0.5
	data.EditorPane.BlockSelect.Update()
}

func (s *editorState) Draw(win *pixelgl.Window) {
	// draw border
	s.BorderView.Canvas.Clear(constants.BlackColor)
	systems.BorderSystem(1)
	img.Batchers[constants.UIBatch].Draw(s.BorderView.Canvas)
	img.Clear()
	s.BorderView.Canvas.Draw(win, s.BorderView.Mat)
	// draw puzzle
	s.PuzzleView.Canvas.Clear(constants.BlackColor)
	systems.DrawSystem(win, 2)
	img.Batchers[constants.TileBGBatch].Draw(s.PuzzleView.Canvas)
	img.Batchers[constants.TileFGBatch].Draw(s.PuzzleView.Canvas)
	img.Clear()
	s.PuzzleView.Canvas.Draw(win, s.PuzzleView.Mat)
	// draw editor panel
	data.EditorPane.ViewPort.Canvas.Clear(color.RGBA{})
	systems.BorderSystem(3)
	systems.DrawSystem(win, 3)
	img.Batchers[constants.UIBatch].Draw(data.EditorPane.ViewPort.Canvas)
	img.Batchers[constants.TileBGBatch].Draw(data.EditorPane.ViewPort.Canvas)
	img.Clear()
	data.EditorPane.ViewPort.Canvas.Draw(win, data.EditorPane.ViewPort.Mat)
	// draw block selector
	if data.EditorPane.SelectVis {
		data.EditorPane.BlockSelect.Canvas.Clear(color.RGBA{})
		systems.DrawSystem(win, 4)
		systems.DrawSystem(win, 5)
		img.Batchers[constants.UIBatch].Draw(data.EditorPane.BlockSelect.Canvas)
		img.Batchers[constants.TileBGBatch].Draw(data.EditorPane.BlockSelect.Canvas)
		img.Clear()
		data.EditorPane.BlockSelect.Canvas.Draw(win, data.EditorPane.BlockSelect.Mat)
	}
}

func (s *editorState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}

func (s *editorState) UpdateViews() {
	portPos := pixel.V(viewport.MainCamera.PostCamPos.X+viewport.MainCamera.Rect.W()*0.5, viewport.MainCamera.PostCamPos.Y+viewport.MainCamera.Rect.H()*0.5)
	s.PuzzleView.PortPos = portPos
	s.BorderView.PortPos = portPos
	wRatio := viewport.MainCamera.Rect.W() / s.PuzzleView.Rect.W()
	hRatio := viewport.MainCamera.Rect.H() / s.PuzzleView.Rect.H()
	pickedRatio := wRatio
	if hRatio < wRatio {
		pickedRatio = hRatio
	}
	pickedRatio *= 0.9
	s.PuzzleView.PortSize = pixel.V(pickedRatio, pickedRatio)
	s.BorderView.PortSize = pixel.V(pickedRatio, pickedRatio)
	data.EditorPane.ViewPort.PortSize = pixel.V(pickedRatio, pickedRatio)
	data.EditorPane.ViewPort.PortPos = pixel.V((viewport.MainCamera.Rect.W()-(s.BorderView.Rect.W()-data.EditorPane.ViewPort.Rect.W())*pickedRatio)*0.5, (viewport.MainCamera.Rect.H()+(s.BorderView.Rect.H()-data.EditorPane.ViewPort.Rect.H())*pickedRatio)*0.5)
}
