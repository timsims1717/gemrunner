package states

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/img"
	"gemrunner/pkg/options"
	"gemrunner/pkg/state"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel/pixelgl"
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
	data.CurrPuzzle = data.CreateBlankPuzzle()
	systems.PuzzleInit()
	systems.EditorInit()
	systems.UpdateViews()
}

func (s *editorState) Update(win *pixelgl.Window) {
	data.EditorInput.Update(win, viewport.MainCamera.Mat)
	debug.AddText("Editor State")
	debug.AddText(fmt.Sprintf("Editor Mode: %s", data.EditorPanel.Mode.String()))
	debug.AddIntCoords("World", int(data.EditorInput.World.X), int(data.EditorInput.World.Y))
	inPos := data.PuzzleView.ProjectWorld(data.EditorInput.World)
	debug.AddIntCoords("Puzzle View In", int(inPos.X), int(inPos.Y))
	//inPos = data.EditorPanel.ViewPort.ProjectWorld(data.EditorInput.World)
	//debug.AddIntCoords("Editor Pane In", int(inPos.X), int(inPos.Y))
	//inPos = data.EditorPanel.BlockSelect.ProjectWorld(data.EditorInput.World)
	//debug.AddIntCoords("Block Select In", int(inPos.X), int(inPos.Y))

	x, y := world.WorldToMap(inPos.X, inPos.Y)
	debug.AddIntCoords("Puzzle Coords", x, y)
	debug.AddIntCoords("Last Coords", data.EditorPanel.LastCoords.X, data.EditorPanel.LastCoords.Y)
	debug.AddText(fmt.Sprintf("NoInput: %t", data.EditorPanel.NoInput))

	//if data.DebugInput.Get("camUp").Pressed() {
	//	data.PuzzleView.PortSize.Y += 1. * timing.DT
	//} else if data.DebugInput.Get("camDown").Pressed() {
	//	data.PuzzleView.PortSize.Y -= 1. * timing.DT
	//}
	//if data.DebugInput.Get("camRight").Pressed() {
	//	data.PuzzleView.PortSize.X += 1. * timing.DT
	//} else if data.DebugInput.Get("camLeft").Pressed() {
	//	data.PuzzleView.PortSize.X -= 1. * timing.DT
	//}
	//if data.DebugInput.Get("debugSP").JustPressed() {
	//	data.PuzzleView.ZoomIn(1.)
	//} else if data.DebugInput.Get("debugSM").JustPressed() {
	//	data.PuzzleView.ZoomIn(-1.)
	//}

	systems.DialogSystem()

	if data.DebugInput.Get("switchWorld").JustPressed() {
		if data.CurrPuzzle != nil {
			switch data.CurrPuzzle.World {
			case constants.WorldRock:
				systems.ChangeWorld(constants.WorldSlate, constants.ColorOrange, constants.ColorRed)
			case constants.WorldSlate:
				systems.ChangeWorld(constants.WorldBrick, constants.ColorRed, constants.ColorBlue)
			case constants.WorldBrick:
				systems.ChangeWorld(constants.WorldGravel, constants.ColorBlue, constants.ColorGray)
			default:
				systems.ChangeWorld(constants.WorldRock, constants.ColorGray, constants.ColorGreen)
			}
			data.CurrPuzzle.Update = true
		}
	}

	// function systems
	systems.FunctionSystem()

	if !data.DialogStackOpen {
		// custom systems
		systems.TileSpriteSystemPre()
		systems.UpdateEditorModeHotKey()
		systems.PuzzleEditSystem()
		systems.TileSpriteSystem()
	}
	// object systems
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
	data.BorderView.Canvas.Clear(constants.ColorBlack)
	systems.BorderSystem(1)
	img.Batchers[constants.UIBatch].Draw(data.BorderView.Canvas)
	img.Clear()
	data.BorderView.Draw(win)
	// draw puzzle
	data.PuzzleView.Canvas.Clear(constants.ColorBlack)
	systems.DrawSystem(win, 2) // normal tiles
	img.Batchers[constants.BGBatch].Draw(data.PuzzleView.Canvas)
	img.Batchers[constants.FGBatch].Draw(data.PuzzleView.Canvas)
	img.Clear()
	systems.DrawSystem(win, 3) // selected tiles
	img.Batchers[constants.BGBatch].Draw(data.PuzzleView.Canvas)
	img.Batchers[constants.FGBatch].Draw(data.PuzzleView.Canvas)
	img.Clear()
	systems.DrawSystem(win, 4) // ui
	img.Batchers[constants.UIBatch].Draw(data.PuzzleView.Canvas)
	img.Clear()
	data.IMDraw.Draw(data.PuzzleView.Canvas)
	data.PuzzleView.Draw(win)
	// draw editor panel
	data.EditorPanel.ViewPort.Canvas.Clear(color.RGBA{})
	systems.BorderSystem(10)
	systems.DrawSystem(win, 10)
	img.Batchers[constants.UIBatch].Draw(data.EditorPanel.ViewPort.Canvas)
	img.Batchers[constants.BGBatch].Draw(data.EditorPanel.ViewPort.Canvas)
	img.Clear()
	data.EditorPanel.ViewPort.Draw(win)
	// draw block selector
	if data.EditorPanel.SelectVis {
		data.EditorPanel.BlockSelect.Canvas.Clear(color.RGBA{})
		systems.DrawSystem(win, 11)
		img.Batchers[constants.UIBatch].Draw(data.EditorPanel.BlockSelect.Canvas)
		img.Batchers[constants.BGBatch].Draw(data.EditorPanel.BlockSelect.Canvas)
		img.Clear()
		systems.DrawSystem(win, 12)
		img.Batchers[constants.UIBatch].Draw(data.EditorPanel.BlockSelect.Canvas)
		img.Clear()
		data.EditorPanel.BlockSelect.Draw(win)
	}
	// dialog draw system
	systems.DialogDrawSystem(win)
	data.IMDraw.Draw(win)
	systems.TemporarySystem()
	data.IMDraw.Clear()
	if options.Updated {
		systems.UpdateViews()
	}
}

func (s *editorState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}
