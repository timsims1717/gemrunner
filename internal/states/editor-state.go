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
	if data.CurrPuzzleSet == nil {
		data.CurrPuzzleSet = data.CreatePuzzleSet()
	}
	data.CurrPuzzleSet.SetToFirst()
	systems.EditorInit()
	systems.PuzzleInit()
	systems.UpdateViews()
	systems.PushUndoArray(true)
	data.EditorDraw = true
}

func (s *editorState) Update(win *pixelgl.Window) {
	systems.CursorSystem(false)
	debug.AddText("Editor State")
	debug.AddText(fmt.Sprintf("Editor Mode: %s", data.Editor.Mode.String()))
	debug.AddIntCoords("World", int(data.MenuInput.World.X), int(data.MenuInput.World.Y))
	inPos := data.PuzzleView.ProjectWorld(data.MenuInput.World)
	debug.AddIntCoords("Puzzle View In", int(inPos.X), int(inPos.Y))

	x, y := world.WorldToMap(inPos.X, inPos.Y)
	debug.AddIntCoords("Puzzle Coords", x, y)
	debug.AddIntCoords("Last Coords", data.Editor.LastCoords.X, data.Editor.LastCoords.Y)
	debug.AddText(fmt.Sprintf("Puzzle Name: %s", data.CurrPuzzleSet.CurrPuzzle.Metadata.Name))
	debug.AddText(fmt.Sprintf("Puzzle Filename: %s", data.CurrPuzzleSet.CurrPuzzle.Metadata.Filename))
	debug.AddText(fmt.Sprintf("Puzzle Music Track: %s", data.CurrPuzzleSet.CurrPuzzle.Metadata.MusicTrack))
	debug.AddTruthText("Puzzle Completed", data.CurrPuzzleSet.CurrPuzzle.Metadata.Completed)
	t := data.CurrPuzzleSet.CurrPuzzle.Tiles.Get(x, y)
	if t != nil {
		sprs := systems.GetTileSprites(t)
		if len(sprs) == 1 {
			debug.AddText(fmt.Sprintf("Tile Sprites: %s", sprs[0].Key))
		} else if len(sprs) == 2 {
			debug.AddText(fmt.Sprintf("Tile Sprites: %s, %s", sprs[0].Key, sprs[1].Key))
		}
	}

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
		data.CurrPuzzleSet.CurrPuzzle.Metadata.Completed = true
	}

	if data.DebugInput.Get("switchWorld").JustPressed() {
		systems.ChangeWorldToNext()
	}

	// function systems
	systems.FunctionSystem()

	data.DialogStackOpen = len(data.DialogStack) > 0
	if !data.DialogStackOpen {
		// custom systems
		systems.TileSpriteSystemPre()
		systems.UpdateEditorModeHotKey()
		systems.PuzzleEditSystem()
	} else {
		// todo: add draw selection here?
	}
	if data.CurrPuzzleSet.CurrPuzzle.Update {
		systems.TileSpriteSystem()
		data.CurrPuzzleSet.CurrPuzzle.Update = false
	}
	systems.DialogSystem()
	// object systems
	systems.ParentSystem()
	systems.ObjectSystem()
	systems.AnimationSystem()

	//s.UpdateViews()

	data.BorderView.Update()
	data.PuzzleView.Update()
	data.PuzzleViewNoShader.Update()

	if data.Editor.SelectVis && !data.Dialogs[constants.DialogEditorBlockSelect].Open {
		data.OpenDialog(constants.DialogEditorBlockSelect)
	} else if !data.Editor.SelectVis && data.Dialogs[constants.DialogEditorBlockSelect].Open {
		data.CloseDialog(constants.DialogEditorBlockSelect)
	}

	myecs.UpdateManager()
	debug.AddText(fmt.Sprintf("Entity Count: %d", myecs.FullCount))
}

func (s *editorState) Draw(win *pixelgl.Window) {
	if data.CurrLevel == nil {
		// draw border
		data.BorderView.Canvas.Clear(constants.ColorBlack)
		systems.BorderSystem(1)
		img.Batchers[constants.UIBatch].Draw(data.BorderView.Canvas)
		img.Clear()
		data.BorderView.Draw(win)
		// draw puzzle
		data.PuzzleView.Canvas.Clear(constants.ColorBlack)
		systems.DrawLayerSystem(data.PuzzleView.Canvas, 2) // normal tiles
		img.Clear()
		systems.DrawLayerSystem(data.PuzzleView.Canvas, 3) // selected tiles
		img.Clear()
		systems.DrawLayerSystem(data.PuzzleView.Canvas, 4) // ui
		img.Clear()
		data.PuzzleView.Draw(win)
		data.PuzzleViewNoShader.Canvas.Clear(color.RGBA{})
		data.IMDraw.Draw(data.PuzzleViewNoShader.Canvas)
		data.PuzzleViewNoShader.Draw(win)
		// dialog draw system
		systems.DialogDrawSystem(win)
		systems.DrawLayerSystem(win, -10)
		img.Clear()
		systems.TemporarySystem()
		data.IMDraw.Clear()
		if options.Updated {
			systems.UpdateViews()
		}
	}
}

func (s *editorState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}
