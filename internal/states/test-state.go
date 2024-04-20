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
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/sfx"
	"gemrunner/pkg/state"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var (
	TestState = &testState{}
)

type testState struct {
	*state.AbstractState
}

func (s *testState) Unload() {
	if data.Editor != nil {
		if data.Editor.PosTop {
			data.OpenDialog("editor_panel_top")
			data.OpenDialog("editor_options_bot")
		} else {
			data.OpenDialog("editor_panel_left")
			data.OpenDialog("editor_options_right")
		}
	}
	systems.LevelDispose()
	systems.ClearTemp()
	data.EditorDraw = true
	data.CurrPuzzle.Update = true
	sfx.MusicPlayer.GetStream("game").Stop()
}

func (s *testState) Load() {
	if data.Editor != nil {
		if data.Editor.PosTop {
			data.CloseDialog("editor_panel_top")
			data.CloseDialog("editor_options_bot")
		} else {
			data.CloseDialog("editor_panel_left")
			data.CloseDialog("editor_options_right")
		}
	}
	systems.LevelInit()
	systems.UpdateViews()
	data.EditorDraw = false
	reanimator.SetFrameRate(constants.FrameRate)
	reanimator.Reset()
	sfx.MusicPlayer.GetStream("game").RepeatTrack(data.CurrLevel.Metadata.MusicTrack)
}

func (s *testState) Update(win *pixelgl.Window) {
	data.P1Input.Update(win, viewport.MainCamera.Mat)
	data.P2Input.Update(win, viewport.MainCamera.Mat)
	data.P3Input.Update(win, viewport.MainCamera.Mat)
	data.P4Input.Update(win, viewport.MainCamera.Mat)
	debug.AddText("Test State")
	debug.AddText(fmt.Sprintf("Speed: %d", constants.FrameRate))
	for i, player := range data.CurrLevel.Players {
		if player != nil {
			pos := player.Object.Pos
			debug.AddIntCoords(fmt.Sprintf("Player %d Pos", i+1), int(pos.X), int(pos.Y))
			cx, cy := world.WorldToMap(pos.X, pos.Y)
			debug.AddIntCoords(fmt.Sprintf("Player %d Coords", i+1), cx, cy)
			debug.AddText(fmt.Sprintf("Player %d Score: %d", i+1, data.CurrLevel.Stats[i].Score))
			debug.AddText(fmt.Sprintf("Player %d State: %s", i+1, player.State.String()))
		}
	}

	if reanimator.FRate != constants.FrameRate {
		reanimator.SetFrameRate(constants.FrameRate)
	}
	reanimator.Update()

	// function systems
	systems.FunctionSystem()
	systems.TestSystem()
	data.DialogStackOpen = len(data.DialogStack) > 0
	systems.DialogSystem()

	if !data.DialogStackOpen {
		// custom systems
		systems.InGameSystem()
		systems.CharacterActionSystem()
		systems.DynamicSystem()
		systems.CollisionSystem()
		systems.CharacterStateSystem()
		systems.TouchSystem()
		systems.SmashSystem()
		systems.TileSystem()
		systems.TileSpriteSystemPre()
		systems.TileSpriteSystem()
	} else {

	}
	// object systems
	systems.ParentSystem()
	systems.ObjectSystem()
	systems.AnimationSystem()

	data.BorderView.Update()
	data.PuzzleView.Update()

	myecs.UpdateManager()
	debug.AddText(fmt.Sprintf("Entity Count: %d", myecs.FullCount))
}

func (s *testState) Draw(win *pixelgl.Window) {
	// draw border
	data.BorderView.Canvas.Clear(constants.ColorBlack)
	systems.BorderSystem(1)
	img.Batchers[constants.UIBatch].Draw(data.BorderView.Canvas)
	img.Clear()
	data.BorderView.Draw(win)
	// draw puzzle
	data.PuzzleView.Canvas.Clear(constants.ColorBlack)
	systems.DrawBatchSystem(data.PuzzleView.Canvas, constants.TileBatch, constants.DrawingLayers)
	img.Clear()
	data.PuzzleView.Draw(win)
	// draw collapse/regen
	data.PuzzleView.Canvas.Clear(pixel.RGBA{})
	systems.DrawBatchSystem(data.PuzzleView.Canvas, constants.TileBatch, constants.CollapseRegenLayer)
	data.PuzzleView.Canvas.SetComposeMethod(pixel.ComposeRatop)
	systems.DrawBatchSystem(data.PuzzleView.Canvas, constants.TileBatch, constants.CollapseRegenMask)
	data.PuzzleView.Canvas.SetComposeMethod(pixel.ComposeOver)
	data.PuzzleView.Draw(win)
	// draw debug
	if debug.ShowDebug {
		data.PuzzleView.Canvas.Clear(pixel.RGBA{})
		debug.DrawLines(data.PuzzleView.Canvas)
		data.PuzzleView.Draw(win)
	}
	// dialog draw system
	systems.DialogDrawSystem(win)
	systems.TemporarySystem()
	//data.IMDraw.Clear()
	if options.Updated {
		systems.UpdateViews()
	}
}

func (s *testState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}
