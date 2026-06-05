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
	systems.DisposePuzzle(data.CurrentPlayArea.Puzzle)
	systems.RemoveEditorBoss()
	data.CurrPuzzleSet = nil
}

func (s *editorState) Load(win *pixelgl.Window) {
	data.EditorDraw = true
	ui.ClearDialogsOpen()
	ui.ClearDialogStack()
	systems.EditorDialogs(win)
	if data.CurrPuzzleSet == nil {
		data.CurrPuzzleSet = data.CreatePuzzleSet()
	}
	if data.CurrPuzzleSet.LastEditedPuzzle > 0 && data.CurrPuzzleSet.LastEditedPuzzle < len(data.CurrPuzzleSet.Puzzles) {
		data.CurrPuzzleSet.SetTo(data.CurrPuzzleSet.LastEditedPuzzle)
	} else {
		data.CurrPuzzleSet.SetToFirst()
	}
	if data.CurrentPlayArea == nil {
		data.CurrentPlayArea = systems.CreatePlayArea()
	}
	systems.EditorInit()
	systems.SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	systems.InitPuzzle(data.CurrentPlayArea)
	systems.SetEditorBoss(data.CurrPuzzleSet.CurrPuzzle.Metadata.Boss)
	systems.UpdateViews()
	reanimator.SetFrameRate(constants.Configuration.Gameplay.FrameRate)
	reanimator.Reset()
	systems.PushUndoArray(true)
}

func (s *editorState) Update(win *pixelgl.Window) {
	systems.CursorSystem(false)
	debug.AddText("Editor State")
	debug.AddText(fmt.Sprintf("Editor Mode: %s", data.Editor.Mode.String()))
	debug.AddIntCoords("World", int(data.MenuInput.World.X), int(data.MenuInput.World.Y))
	inPos := data.CurrentPlayArea.PuzzleView.ProjectWorld(data.MenuInput.World)
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
		dKey := constants.DialogBossSettings
		load.ReloadDialog(dKey)
		systems.CustomizeEditorDialog(dKey)
		systems.UpdateDialogView(ui.Dialogs[dKey])
		//constants.Configuration.Graphics.Resolution++
		//constants.Configuration.Graphics.Resolution %= len(options.Resolutions)
		//content.UpdateConfiguration()
	}
	if data.DebugInput.Get("switchWorld").JustPressed() {
		data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderMode++
		data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderMode %= constants.ShaderEndOfList
		systems.ChangeWorldShader(data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderMode)
	}
	reanimator.Update()

	// function systems
	systems.InterpolationSystem()
	systems.FunctionSystem()
	if reanimator.FrameSwitch {
		data.Editor.FrameCount++
	}

	ui.DialogStackOpen = len(ui.DialogStack) > 0
	if !ui.DialogStackOpen {
		// custom systems
		systems.TileSpriteSystemPre()
		systems.UpdateEditorModeHotKey()
		systems.PuzzleEditSystem()
		systems.FloatingTextEditorSystem(data.CurrentPlayArea)
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
	systems.BossEditorSystem()
	systems.AnimationSystem()
	systems.ParentSystem()
	systems.ObjectSystem()
	systems.EffectsSystem()

	//s.UpdateViews()

	data.CurrentPlayArea.BorderView.Update()
	data.CurrentPlayArea.BackgroundView.Update()
	data.CurrentPlayArea.PuzzleView.Update()
	data.CurrentPlayArea.WorldView.Update()
	data.CurrentPlayArea.PuzzleViewNoShader.Update()
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
		data.CurrentPlayArea.WorldView.Canvas.Clear(constants.ColorBlack)
		// draw border
		data.CurrentPlayArea.BorderView.Canvas.Clear(constants.ColorBlack)
		systems.DrawBorder(data.PuzzleBorderObject, data.PuzzleBorder, nil)
		img.Batchers[constants.UIBatch].DrawThenClear(data.CurrentPlayArea.BorderView.Canvas)
		data.CurrentPlayArea.BorderView.Draw(data.ScreenView.Canvas)
		// draw background
		data.CurrentPlayArea.BackgroundView.Canvas.Clear(constants.ColorBlack)
		data.BlackBackground.Draw(data.CurrentPlayArea.BackgroundView.Canvas, data.CurrentPlayArea.Puzzle.Metadata.BackgroundMatrix.Moved(pixel.V(float64(data.CurrentPlayArea.Puzzle.Metadata.Width)*world.TileSize*0.5, float64(data.CurrentPlayArea.Puzzle.Metadata.Height)*world.TileSize*0.5)))
		data.CurrentPlayArea.BackgroundView.Draw(data.CurrentPlayArea.WorldView.Canvas)
		// draw puzzle
		data.CurrentPlayArea.PuzzleView.Canvas.Clear(constants.ColorClear)
		systems.DrawLayerSystem(data.CurrentPlayArea.PuzzleView.Canvas, 8) // background
		img.Clear()
		systems.DrawLayerSystem(data.CurrentPlayArea.PuzzleView.Canvas, 2) // normal tiles
		img.Clear()
		systems.DrawLayerSystem(data.CurrentPlayArea.PuzzleView.Canvas, 3) // selected tiles
		img.Clear()
		systems.DrawLayerSystem(data.CurrentPlayArea.PuzzleView.Canvas, 4) // ui
		img.Clear()
		//data.PuzzleView.Draw(win)
		data.CurrentPlayArea.PuzzleView.Draw(data.CurrentPlayArea.WorldView.Canvas)
		data.CurrentPlayArea.PuzzleViewNoShader.Canvas.Clear(pixel.RGBA{})
		data.CurrentPlayArea.IMDraw.Draw(data.CurrentPlayArea.PuzzleViewNoShader.Canvas)
		systems.DrawLayerSystem(data.CurrentPlayArea.PuzzleViewNoShader.Canvas, 36)
		systems.DrawLayerSystem(data.CurrentPlayArea.PuzzleViewNoShader.Canvas, 37)

		// draw debug
		if debug.ShowDebug {
			debug.DrawLines(data.CurrentPlayArea.PuzzleViewNoShader.Canvas)
		}
		//data.PuzzleViewNoShader.Draw(data.WorldView.Canvas)
		data.CurrentPlayArea.WorldView.Draw(data.ScreenView.Canvas)
		data.CurrentPlayArea.PuzzleViewNoShader.Draw(data.ScreenView.Canvas)
		// dialog draw system
		systems.DialogDrawSystem(data.ScreenView.Canvas)
		systems.DrawLayerSystem(data.ScreenView.Canvas, -10)
		img.Clear()
		systems.TemporarySystem()
		data.CurrentPlayArea.IMDraw.Clear()
	}
}

func (s *editorState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}
