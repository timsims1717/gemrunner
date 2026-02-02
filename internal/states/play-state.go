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
		if data.LevelTrans {
			systems.LevelTransitionSystem()
		}
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
	drawPlayArea(win)
}

func drawPlayArea(win *pixelgl.Window) {
	data.ScreenView.Canvas.Clear(constants.ColorBlack)
	// draw border
	data.CurrentPlayArea.BorderView.Canvas.Clear(constants.ColorBlack)
	systems.DrawBorder(data.PuzzleBorderObject, data.PuzzleBorder, data.CurrentPlayArea.BorderView.Canvas, true)
	img.Clear()
	data.CurrentPlayArea.BorderView.Draw(data.ScreenView.Canvas)
	// draw puzzle
	data.CurrentPlayArea.WorldView.Canvas.Clear(pixel.RGBA{})
	data.CurrentPlayArea.PuzzleView.Canvas.Clear(constants.ColorBlack)
	systems.DrawBatchSystem(data.CurrentPlayArea.PuzzleView.Canvas, constants.TileBatch, constants.DrawingLayers)
	img.Clear()
	data.CurrentPlayArea.PuzzleView.Draw(data.CurrentPlayArea.WorldView.Canvas)
	// draw collapse/regen
	data.CurrentPlayArea.PuzzleView.Canvas.Clear(pixel.RGBA{})
	systems.DrawBatchSystem(data.CurrentPlayArea.PuzzleView.Canvas, constants.TileBatch, constants.CollapseRegenLayer)
	data.CurrentPlayArea.PuzzleView.Canvas.SetComposeMethod(pixel.ComposeRatop)
	systems.DrawBatchSystem(data.CurrentPlayArea.PuzzleView.Canvas, constants.TileBatch, constants.CollapseRegenMask)
	data.CurrentPlayArea.PuzzleView.Canvas.SetComposeMethod(pixel.ComposeOver)
	data.CurrentPlayArea.PuzzleView.Draw(data.CurrentPlayArea.WorldView.Canvas)
	// draw effects
	data.CurrentPlayArea.PuzzleView.Canvas.Clear(pixel.RGBA{})
	systems.DrawBatchSystem(data.CurrentPlayArea.PuzzleView.Canvas, constants.TileBatch, constants.EffectsLayer)
	data.CurrentPlayArea.PuzzleView.Draw(data.CurrentPlayArea.WorldView.Canvas)
	data.CurrentPlayArea.PuzzleViewNoShader.Canvas.Clear(pixel.RGBA{})
	systems.DrawLayerSystem(data.CurrentPlayArea.PuzzleViewNoShader.Canvas, 36)
	systems.DrawLayerSystem(data.CurrentPlayArea.PuzzleViewNoShader.Canvas, 37)
	// draw debug
	if debug.ShowDebug {
		debug.DrawLines(data.CurrentPlayArea.PuzzleViewNoShader.Canvas)
	}
	data.CurrentPlayArea.WorldView.Draw(data.ScreenView.Canvas)
	data.CurrentPlayArea.PuzzleViewNoShader.Draw(data.ScreenView.Canvas)
	// dialog draw system
	systems.DialogDrawSystem(data.ScreenView.Canvas)
	systems.DrawLayerSystem(data.ScreenView.Canvas, -10)
	img.Clear()
	systems.TemporarySystem()
}

func (s *playState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}
