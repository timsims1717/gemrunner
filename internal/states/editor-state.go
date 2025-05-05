package states

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/load"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems"
	"gemrunner/internal/ui"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/img"
	"gemrunner/pkg/options"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/state"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var (
	EditorState = &editorState{}
)

type editorState struct {
	*state.AbstractState
}

func (s *editorState) Unload(win *pixelgl.Window) {
	ui.ClearDialogsOpen()
	ui.ClearDialogStack()
	systems.DisposeEditor()
	systems.DisposeEditorDialogs()
	systems.DisposeInGameDialogs()
	systems.PuzzleDispose()
	data.CurrPuzzleSet = nil
}

func (s *editorState) Load(win *pixelgl.Window) {
	data.EditorDraw = true
	ui.ClearDialogsOpen()
	ui.ClearDialogStack()
	systems.EditorDialogs(win)
	systems.InGameDialogs(win)
	if data.CurrPuzzleSet == nil {
		data.CurrPuzzleSet = data.CreatePuzzleSet()
	}
	data.CurrPuzzleSet.SetToFirst()
	systems.EditorInit()
	systems.PuzzleInit()
	systems.UpdateViews()
	reanimator.SetFrameRate(constants.FrameRate)
	reanimator.Reset()
	systems.PushUndoArray(true)
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
	//debug.AddText(fmt.Sprintf("Puzzle Name: %s", data.CurrPuzzleSet.CurrPuzzle.Metadata.Name))
	//debug.AddText(fmt.Sprintf("Puzzle Filename: %s", data.CurrPuzzleSet.CurrPuzzle.Metadata.Filename))
	//debug.AddText(fmt.Sprintf("Puzzle Music Track: %s", data.CurrPuzzleSet.CurrPuzzle.Metadata.MusicTrack))
	debug.AddText(fmt.Sprintf("Undo Stack Size: %d", len(data.CurrPuzzleSet.CurrPuzzle.UndoStack)))
	debug.AddText(fmt.Sprintf("Redo Stack Size: %d", len(data.CurrPuzzleSet.CurrPuzzle.RedoStack)))
	debug.AddTruthText("Puzzle Completed", data.CurrPuzzleSet.CurrPuzzle.Metadata.Completed)
	//t := data.CurrPuzzleSet.CurrPuzzle.Get(x, y)
	//if t != nil {
	//	sprs := systems.GetTileDrawables(t)
	//	if len(sprs) == 1 {
	//		if spr, ok := sprs[0].(*img.Sprite); ok {
	//			debug.AddText(fmt.Sprintf("Tile Sprites: %s", spr.Key))
	//		} else if anim, ok1 := sprs[0].(*reanimator.Tree); ok1 {
	//			debug.AddText(fmt.Sprintf("Tile Sprites: %s", anim.))
	//		}
	//	} else if len(sprs) == 2 {
	//		debug.AddText(fmt.Sprintf("Tile Sprites: %s, %s", sprs[0].Key, sprs[1].Key))
	//	}
	//}

	if data.DebugInput.Get("camUp").JustPressed() || data.DebugInput.Get("camUp").Repeated() {
		data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderY += 0.02
	} else if data.DebugInput.Get("camDown").JustPressed() || data.DebugInput.Get("camDown").Repeated() {
		data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderY -= 0.02
	}
	if data.DebugInput.Get("camRight").JustPressed() || data.DebugInput.Get("camRight").Repeated() {
		data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderX += 0.0001
	} else if data.DebugInput.Get("camLeft").JustPressed() || data.DebugInput.Get("camLeft").Repeated() {
		data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderX -= 0.0001
	}
	if data.DebugInput.Get("debugSP").JustPressed() || data.DebugInput.Get("debugSP").Repeated() {
		//data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderCustom += 0.01
		data.CurrPuzzleSet.CurrPuzzle.SetWidth(data.CurrPuzzleSet.CurrPuzzle.Metadata.Width + 1)
		systems.UpdateViews()
	} else if data.DebugInput.Get("debugSM").JustPressed() || data.DebugInput.Get("debugSM").Repeated() {
		//data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderCustom -= 0.01
		data.CurrPuzzleSet.CurrPuzzle.SetWidth(data.CurrPuzzleSet.CurrPuzzle.Metadata.Width - 1)
		systems.UpdateViews()
	}
	debug.AddText(fmt.Sprintf("Shader Speed: %f, ShaderCustom: %f", data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderSpeed, data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderCustom))
	debug.AddText(fmt.Sprintf("ShaderX: %f, ShaderY: %f", data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderX, data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderY))
	if data.DebugInput.Get("debugTest").JustPressed() {
		dKey := constants.DialogPuzzleSettings
		load.ReloadDialog(dKey)
		systems.CustomizeEditorDialog(dKey)
		systems.UpdateDialogView(ui.Dialogs[dKey])
	}
	if data.DebugInput.Get("switchWorld").JustPressed() {
		data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderMode++
		data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderMode %= constants.ShaderEndOfList
		systems.ChangeWorldShader(data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderMode)
	}
	reanimator.Update()

	// function systems
	systems.FunctionSystem()
	systems.InterpolationSystem()
	if reanimator.FrameSwitch {
		data.Editor.FrameCount++
	}

	ui.DialogStackOpen = len(ui.DialogStack) > 0
	if !ui.DialogStackOpen {
		// custom systems
		systems.TileSpriteSystemPre()
		systems.UpdateEditorModeHotKey()
		systems.PuzzleEditSystem()
		systems.FloatingTextEditorSystem()
	} else {
		// todo: add draw selection here?
	}
	if data.CurrPuzzleSet.CurrPuzzle.Update || data.CurrPuzzleSet.CurrPuzzle.Changed {
		systems.TileSpriteSystem()
		data.CurrPuzzleSet.CurrPuzzle.Update = false
	}
	systems.DialogSystem(win)
	systems.UndoStackSystem()
	// object systems
	systems.ShaderSystem()
	systems.AnimationSystem()
	systems.ParentSystem()
	systems.ObjectSystem()

	//s.UpdateViews()

	data.BorderView.Update()
	data.PuzzleView.Update()
	data.WorldView.Update()
	data.PuzzleViewNoShader.Update()
	data.ScreenView.Update()

	if data.Editor.SelectVis && !ui.Dialogs[constants.DialogEditorBlockSelect].Open {
		ui.OpenDialog(constants.DialogEditorBlockSelect)
	} else if !data.Editor.SelectVis && ui.Dialogs[constants.DialogEditorBlockSelect].Open {
		ui.CloseDialog(constants.DialogEditorBlockSelect)
	}

	myecs.UpdateManager()
	debug.AddText(fmt.Sprintf("Entity Count: %d", myecs.FullCount))
}

func (s *editorState) Draw(win *pixelgl.Window) {
	data.ScreenView.Canvas.Clear(constants.ColorBlack)
	if data.CurrLevel == nil {
		// draw border
		data.BorderView.Canvas.Clear(constants.ColorBlack)
		systems.DrawBorder(ui.PuzzleBorderObject, ui.PuzzleBorder, data.BorderView.Canvas)
		img.Clear()
		data.BorderView.Draw(data.ScreenView.Canvas)
		// draw puzzle
		data.WorldView.Canvas.Clear(pixel.RGBA{})
		data.PuzzleView.Canvas.Clear(constants.ColorBlack)
		systems.DrawLayerSystem(data.PuzzleView.Canvas, 2) // normal tiles
		img.Clear()
		systems.DrawLayerSystem(data.PuzzleView.Canvas, 3) // selected tiles
		img.Clear()
		systems.DrawLayerSystem(data.PuzzleView.Canvas, 4) // ui
		img.Clear()
		//data.PuzzleView.Draw(win)
		data.PuzzleView.Draw(data.WorldView.Canvas)
		data.PuzzleViewNoShader.Canvas.Clear(pixel.RGBA{})
		data.IMDraw.Draw(data.PuzzleViewNoShader.Canvas)
		systems.DrawLayerSystem(data.PuzzleViewNoShader.Canvas, 36)
		systems.DrawLayerSystem(data.PuzzleViewNoShader.Canvas, 37)

		// draw debug
		if debug.ShowDebug {
			debug.DrawLines(data.PuzzleViewNoShader.Canvas)
		}
		//data.PuzzleViewNoShader.Draw(data.WorldView.Canvas)
		data.WorldView.Draw(data.ScreenView.Canvas)
		data.PuzzleViewNoShader.Draw(data.ScreenView.Canvas)
		// dialog draw system
		systems.DialogDrawSystem(data.ScreenView.Canvas)
		systems.DrawLayerSystem(data.ScreenView.Canvas, -10)
		img.Clear()
		systems.TemporarySystem()
		data.IMDraw.Clear()
		data.ScreenView.Draw(win)
		if options.Updated {
			systems.UpdateViews()
		}
	}
}

func (s *editorState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}
