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
	"gemrunner/pkg/world"
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
	debug.AddText(fmt.Sprintf("Speed: %d", constants.Configuration.Gameplay.FrameRate))
	debug.AddText(systems.FormatTimePlayed())
	debug.AddText(fmt.Sprintf("Frame Number: %d", data.CurrLevel.FrameNumber))
	//debug.AddText(fmt.Sprintf("Frame Counter: %d", data.CurrLevel.FrameCounter))
	debug.AddText(fmt.Sprintf("Frame Cycle: %d", data.CurrLevel.FrameCycle))
	//debug.AddTruthText("Frame Change", data.CurrLevel.FrameChange)
	for i, player := range data.CurrLevel.Players {
		if player != nil {
			pos := player.Object.Pos
			debug.AddIntCoords(fmt.Sprintf("Player %d Pos", i+1), int(pos.X), int(pos.Y))
			cx, cy := world.WorldToMap(pos.X, pos.Y)
			debug.AddIntCoords(fmt.Sprintf("Player %d Coords", i+1), cx, cy)
			debug.AddText(fmt.Sprintf("Player %d Score: %d", i+1, data.CurrLevelSess.PlayerStats[i].Score))
			debug.AddText(fmt.Sprintf("Player %d Deaths: %d", i+1, data.CurrLevelSess.PlayerStats[i].Deaths))
			debug.AddText(fmt.Sprintf("Player %d State: %s", i+1, player.State.String()))
			if player.Inventory == nil {
				debug.AddText(fmt.Sprintf("Player %d Inv: Empty", i+1))
			} else {
				item := player.Inventory.Name
				debug.AddText(fmt.Sprintf("Player %d Inv: %s", i+1, item))
			}
			//debug.AddText(fmt.Sprintf("Player %d # of Tiles: %d", i+1, len(player.StoredBlocks)))
		}
	}

	if data.DebugInput.Get("debugTest").JustPressed() {
		for _, player := range data.CurrLevel.Players {
			if player != nil {
				player.SmallBombs++
			}
		}
		data.DebugInput.Get("debugTest").Consume()
	}
	if data.DebugInput.Get("debugInv").JustPressed() {
		systems.OpenDoors()
		data.DebugInput.Get("debugInv").Consume()
	}
	if data.DebugInput.Get("ctrl").Pressed() {
		if data.DebugInput.Get("debugLevelUp").JustPressed() {
			systems.GoToLevelUp()
			data.DebugInput.Get("debugLevelUp").Consume()
		} else if data.DebugInput.Get("debugLevelDown").JustPressed() {
			systems.GoToLevelDown()
			data.DebugInput.Get("debugLevelDown").Consume()
		}
		if data.DebugInput.Get("debugLevelLeft").JustPressed() {
			systems.GoToLevelLeft()
			data.DebugInput.Get("debugLevelLeft").Consume()
		}
		if data.DebugInput.Get("debugLevelRight").JustPressed() {
			systems.GoToLevelRight()
			data.DebugInput.Get("debugLevelRight").Consume()
		}
	}

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
