package states

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems"
	"gemrunner/internal/ui"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/img"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/sfx"
	"gemrunner/pkg/state"
	"gemrunner/pkg/viewport"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var (
	PlayState = &playState{}
)

type playState struct {
	*state.AbstractState
}

func (s *playState) Unload(win *pixelgl.Window) {
	ui.ClearDialogsOpen()
	ui.ClearDialogStack()
	systems.DisposeInGameDialogs()
	systems.LevelSessionDispose()
	systems.DisposeCurrLevel()
	systems.ClearTemp()
	systems.DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet = nil
}

func (s *playState) Load(win *pixelgl.Window) {
	data.EditorDraw = false
	ui.ClearDialogsOpen()
	ui.ClearDialogStack()
	systems.InGameDialogs(win)
	systems.LevelSessionInit()
	systems.StartLevel(true)
	systems.UpdateViews()
	reanimator.SetFrameRate(constants.Configuration.Gameplay.FrameRate)
	reanimator.Reset()
	sfx.MusicPlayer.GetStream("game").RepeatTrack(data.CurrLevel.Metadata.MusicTrack)
}

func (s *playState) Update(win *pixelgl.Window) {
	data.P1Input.Update(win, viewport.MainCamera.Mat)
	data.P2Input.Update(win, viewport.MainCamera.Mat)
	data.P3Input.Update(win, viewport.MainCamera.Mat)
	data.P4Input.Update(win, viewport.MainCamera.Mat)
	systems.CursorSystem(true)
	debug.AddText("Play State")
	systems.InGameDebugInfo()
	systems.InGameDebugInput()
	systems.PlayDebugInput()

	if reanimator.FRate != constants.Configuration.Gameplay.FrameRate {
		reanimator.SetFrameRate(constants.Configuration.Gameplay.FrameRate)
	}
	reanimator.Update()

	// function systems
	systems.PlayPauseSystem()
	systems.FunctionSystem()
	systems.InterpolationSystem()
	ui.DialogStackOpen = len(ui.DialogStack) > 0
	systems.DialogSystem(win)

	if !ui.DialogStackOpen {
		// custom systems
		systems.InGameSystem()
		systems.PlaySystem()
		systems.CharacterActionSystem()
		systems.DynamicSystem()
		systems.CollisionSystem()
		systems.OutsideMapSystem()
		systems.CharacterStateSystem()
		systems.TouchSystem()
		systems.SmashSystem()
		systems.TileSystem()
		systems.MagicSystem()
		systems.TileSpriteSystemPre()
		systems.TileSpriteSystem()
		systems.FloatingTextSystem()
		systems.AnimationSystem()
	} else {

	}
	// object systems
	systems.ParentSystem()
	systems.ObjectSystem()
	systems.EffectsSystem()

	data.CurrentPlayArea.BorderView.Update()
	data.CurrentPlayArea.PuzzleView.Update()
	data.CurrentPlayArea.WorldView.Update()
	data.CurrentPlayArea.PuzzleViewNoShader.Update()
	data.ScreenView.Update()

	myecs.UpdateManager()
	debug.AddText(fmt.Sprintf("Entity Count: %d", myecs.FullCount))
}

func (s *playState) Draw(win *pixelgl.Window) {
	drawPlayState(win)
}

func drawPlayState(win *pixelgl.Window) {
	data.ScreenView.Canvas.Clear(constants.ColorBlack)
	drawPlayArea(win, data.CurrentPlayArea)
	// dialog draw system
	systems.DialogDrawSystem(data.ScreenView.Canvas)
	systems.DrawLayerSystem(data.ScreenView.Canvas, -10)
	img.Clear()
	systems.TemporarySystem()
}

func drawPlayArea(win *pixelgl.Window, pa *data.PlayArea) {
	// draw border
	pa.BorderView.Canvas.Clear(constants.ColorClear)
	systems.DrawBorder(pa.BorderObject, pa.Border, pa.BorderView.Canvas, pa.Level)
	img.Clear()
	pa.BorderView.Draw(data.ScreenView.Canvas)
	// draw puzzle
	pa.WorldView.Canvas.Clear(constants.ColorClear)
	pa.PuzzleView.Canvas.Clear(constants.ColorBlack)
	systems.DrawBatchSystem(pa.PuzzleView.Canvas, constants.TileBatch, constants.DrawingLayers, pa.LayerOffset)
	img.Clear()
	pa.PuzzleView.Draw(pa.WorldView.Canvas)
	// draw collapse/regen
	pa.PuzzleView.Canvas.Clear(constants.ColorClear)
	systems.DrawBatchSystem(pa.PuzzleView.Canvas, constants.TileBatch, constants.CollapseRegenLayer, pa.LayerOffset)
	pa.PuzzleView.Canvas.SetComposeMethod(pixel.ComposeRatop)
	systems.DrawBatchSystem(pa.PuzzleView.Canvas, constants.TileBatch, constants.CollapseRegenMask, pa.LayerOffset)
	pa.PuzzleView.Canvas.SetComposeMethod(pixel.ComposeOver)
	pa.PuzzleView.Draw(pa.WorldView.Canvas)
	// draw effects
	pa.PuzzleView.Canvas.Clear(constants.ColorClear)
	systems.DrawBatchSystem(pa.PuzzleView.Canvas, constants.TileBatch, constants.EffectsLayer, pa.LayerOffset)
	pa.PuzzleView.Draw(pa.WorldView.Canvas)
	pa.PuzzleViewNoShader.Canvas.Clear(constants.ColorClear)
	systems.DrawLayerSystem(pa.PuzzleViewNoShader.Canvas, 36+pa.LayerOffset)
	systems.DrawLayerSystem(pa.PuzzleViewNoShader.Canvas, 37+pa.LayerOffset)
	// draw debug
	if debug.ShowDebug {
		debug.DrawLines(pa.PuzzleViewNoShader.Canvas)
	}
	pa.WorldView.Draw(data.ScreenView.Canvas)
	pa.PuzzleViewNoShader.Draw(data.ScreenView.Canvas)
}

func (s *playState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}
