package states

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/img"
	"gemrunner/pkg/state"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
)

var (
	EditorState = &editorState{}
)

type editorState struct {
	*state.AbstractState
}

func (s *editorState) Unload() {

}

func (s *editorState) Load() {
	data.CurrPuzzle = data.CreateTestPuzzle()
	systems.PuzzleInit()
	systems.EditorInit()
	s.UpdateViews()
}

func (s *editorState) Update(win *pixelgl.Window) {
	debug.AddText("Editor State")
	debug.AddIntCoords("World", int(data.EditorInput.World.X), int(data.EditorInput.World.Y))
	inPos := data.PuzzleView.ProjectWorld(data.EditorInput.World)
	debug.AddIntCoords("Puzzle View In", int(inPos.X), int(inPos.Y))
	inPos = data.EditorPanel.ViewPort.ProjectWorld(data.EditorInput.World)
	debug.AddIntCoords("Editor Pane In", int(inPos.X), int(inPos.Y))
	inPos = data.EditorPanel.BlockSelect.ProjectWorld(data.EditorInput.World)
	debug.AddIntCoords("Block Select In", int(inPos.X), int(inPos.Y))

	if data.DebugInput.Get("camUp").Pressed() {
		data.PuzzleView.PortSize.Y += 1. * timing.DT
	} else if data.DebugInput.Get("camDown").Pressed() {
		data.PuzzleView.PortSize.Y -= 1. * timing.DT
	}
	if data.DebugInput.Get("camRight").Pressed() {
		data.PuzzleView.PortSize.X += 1. * timing.DT
	} else if data.DebugInput.Get("camLeft").Pressed() {
		data.PuzzleView.PortSize.X -= 1. * timing.DT
	}
	if data.DebugInput.Get("debugSP").JustPressed() {
		data.PuzzleView.ZoomIn(1.)
	} else if data.DebugInput.Get("debugSM").JustPressed() {
		data.PuzzleView.ZoomIn(-1.)
	}

	data.EditorInput.Update(win, viewport.MainCamera.Mat)

	//data.EditorPanel.Hover = data.EditorPanel.ViewPort.PointInside(data.EditorInput.World) || (data.EditorPanel.SelectVis && data.EditorPanel.BlockSelect.PointInside(data.EditorInput.World))
	//if data.EditorPanel.Hover {
	//	pos := data.EditorPanel.ViewPort.ProjectedOut(data.EditorInput.World)
	//	bsObj := data.EditorPanel.BlockView.Object
	//	if bsObj.Rect.Moved(bsObj.Pos).Contains(pos) {
	//		if data.EditorInput.Get("click").Pressed() {
	//			data.EditorPanel.SelectVis = true
	//		}
	//	}
	//} else {
	//	data.EditorPanel.SelectVis = false
	//	if data.EditorInput.Get("rightClick").Pressed() {
	//		if s.PuzzleView.PointInside(data.EditorInput.World) {
	//			projPos := util.ProjectedPoint(data.EditorInput.World, s.PuzzleView.Rect, s.PuzzleView.Mat)
	//			x, y := world.WorldToMap(projPos.X, projPos.Y)
	//			coords := world.Coords{X: x, Y: y}
	//			systems.DeleteBlock(coords)
	//		}
	//	} else if data.EditorInput.Get("click").Pressed() {
	//		if s.PuzzleView.PointInside(data.EditorInput.World) {
	//			projPos := util.ProjectedPoint(data.EditorInput.World, s.PuzzleView.Rect, s.PuzzleView.Mat)
	//			x, y := world.WorldToMap(projPos.X, projPos.Y)
	//			coords := world.Coords{X: x, Y: y}
	//			systems.ChangeBlock(coords, data.RedRock)
	//		}
	//	}
	//}

	systems.TemporarySystem()
	systems.FunctionSystem()
	systems.PuzzleEditSystem()
	systems.TileSpriteSystem()
	systems.ParentSystem()
	systems.ObjectSystem()
	systems.AnimationSystem()

	//s.UpdateViews()

	data.BorderView.Update()
	data.PuzzleView.Update()
	data.EditorPanel.ViewPort.Update()
	bSelect := data.EditorPanel.BlockSelect
	bSelect.PortSize = data.EditorPanel.ViewPort.PortSize
	bSelect.PortPos = data.EditorPanel.ViewPort.PortPos
	bSelect.PortPos.X += (bSelect.Canvas.Bounds().W() + world.TileSize + 8.) * bSelect.PortSize.X * 0.5
	bSelect.PortPos.Y -= (data.EditorPanel.ViewPort.Canvas.Bounds().H() - world.TileSize + 6.) * bSelect.PortSize.Y * 0.5
	data.EditorPanel.BlockSelect.Update()
	myecs.UpdateManager()
	debug.AddText(fmt.Sprintf("Entity Count: %d", myecs.FullCount))
}

func (s *editorState) Draw(win *pixelgl.Window) {
	// draw border
	data.BorderView.Canvas.Clear(constants.BlackColor)
	systems.BorderSystem(1)
	img.Batchers[constants.UIBatch].Draw(data.BorderView.Canvas)
	img.Clear()
	data.BorderView.Canvas.Draw(win, data.BorderView.Mat)
	// draw puzzle
	data.PuzzleView.Canvas.Clear(constants.BlackColor)
	systems.DrawSystem(win, 2)
	img.Batchers[constants.TileBGBatch].Draw(data.PuzzleView.Canvas)
	img.Batchers[constants.TileFGBatch].Draw(data.PuzzleView.Canvas)
	img.Clear()
	data.PuzzleView.Canvas.Draw(win, data.PuzzleView.Mat)
	// draw editor panel
	data.EditorPanel.ViewPort.Canvas.Clear(color.RGBA{})
	systems.BorderSystem(3)
	systems.DrawSystem(win, 3)
	img.Batchers[constants.UIBatch].Draw(data.EditorPanel.ViewPort.Canvas)
	img.Batchers[constants.TileBGBatch].Draw(data.EditorPanel.ViewPort.Canvas)
	img.Clear()
	data.EditorPanel.ViewPort.Canvas.Draw(win, data.EditorPanel.ViewPort.Mat)
	// draw block selector
	if data.EditorPanel.SelectVis {
		data.EditorPanel.BlockSelect.Canvas.Clear(color.RGBA{})
		systems.DrawSystem(win, 4)
		systems.DrawSystem(win, 5)
		img.Batchers[constants.UIBatch].Draw(data.EditorPanel.BlockSelect.Canvas)
		img.Batchers[constants.TileBGBatch].Draw(data.EditorPanel.BlockSelect.Canvas)
		img.Clear()
		data.EditorPanel.BlockSelect.Canvas.Draw(win, data.EditorPanel.BlockSelect.Mat)
	}
}

func (s *editorState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}

func (s *editorState) UpdateViews() {
	data.PuzzleView.PortPos = viewport.MainCamera.PostCamPos
	data.BorderView.PortPos = viewport.MainCamera.PostCamPos
	wRatio := viewport.MainCamera.Rect.W() / data.PuzzleView.Rect.W()
	hRatio := viewport.MainCamera.Rect.H() / data.PuzzleView.Rect.H()
	pickedRatio := wRatio
	if hRatio < wRatio {
		pickedRatio = hRatio
	}
	pickedRatio *= 0.9
	data.PuzzleView.PortSize = pixel.V(pickedRatio, pickedRatio)
	data.BorderView.PortSize = pixel.V(pickedRatio, pickedRatio)
	data.EditorPanel.ViewPort.PortSize = pixel.V(pickedRatio, pickedRatio)
	data.EditorPanel.ViewPort.PortPos = pixel.V((viewport.MainCamera.Rect.W()-(data.BorderView.Rect.W()-data.EditorPanel.ViewPort.Rect.W())*pickedRatio)*0.5, (viewport.MainCamera.Rect.H()+(data.BorderView.Rect.H()-data.EditorPanel.ViewPort.Rect.H())*pickedRatio)*0.5)
}
