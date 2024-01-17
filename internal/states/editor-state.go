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
	"gemrunner/pkg/timing"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel/pixelgl"
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
	debug.AddText(fmt.Sprintf("Editor Mode: %s", data.Editor.Mode.String()))
	debug.AddIntCoords("World", int(data.EditorInput.World.X), int(data.EditorInput.World.Y))
	inPos := data.PuzzleView.ProjectWorld(data.EditorInput.World)
	debug.AddIntCoords("Puzzle View In", int(inPos.X), int(inPos.Y))
	debug.AddIntCoords("BlockSelect Pos", int(data.Editor.BlockSelect.PortPos.X), int(data.Editor.BlockSelect.PortPos.Y))

	x, y := world.WorldToMap(inPos.X, inPos.Y)
	debug.AddIntCoords("Puzzle Coords", x, y)
	debug.AddIntCoords("Last Coords", data.Editor.LastCoords.X, data.Editor.LastCoords.Y)
	debug.AddText(fmt.Sprintf("NoInput: %t", data.Editor.NoInput))
	debug.AddText(fmt.Sprintf("SelectVis: %t", data.Editor.SelectVis))

	if data.DebugInput.Get("camUp").Pressed() {
		data.Editor.BlockSelect.PortPos.Y += 100. * timing.DT
	} else if data.DebugInput.Get("camDown").Pressed() {
		data.Editor.BlockSelect.PortPos.Y -= 100. * timing.DT
	}
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
	if data.DebugInput.Get("debugTest").JustPressed() {

	}

	systems.DialogSystem()

	if data.DebugInput.Get("switchWorld").JustPressed() {
		if data.CurrPuzzle != nil {
			data.CurrPuzzle.WorldNumber++
			if data.CurrPuzzle.WorldNumber >= constants.WorldCustom {
				data.CurrPuzzle.WorldNumber %= constants.WorldCustom
			}
			systems.ChangeWorld(data.CurrPuzzle.WorldNumber)
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
	} else {

	}
	// object systems
	systems.ParentSystem()
	systems.ObjectSystem()
	systems.AnimationSystem()

	//s.UpdateViews()

	data.BorderView.Update()
	data.PuzzleView.Update()

	if data.Editor.SelectVis && !data.Dialogs["block_select"].Open {
		systems.OpenDialog("block_select")
	} else if !data.Editor.SelectVis && data.Dialogs["block_select"].Open {
		data.CloseDialog("block_select")
	}

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
	// dialog draw system
	systems.DialogDrawSystem(win)
	systems.TemporarySystem()
	data.IMDraw.Clear()
	if options.Updated {
		systems.UpdateViews()
	}
}

func (s *editorState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}
