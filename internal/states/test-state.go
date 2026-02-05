package states

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems"
	"gemrunner/internal/ui"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/sfx"
	"gemrunner/pkg/state"
	"gemrunner/pkg/viewport"
	"github.com/gopxl/pixel/pixelgl"
)

var (
	TestState = &testState{}
)

type testState struct {
	*state.AbstractState
}

func (s *testState) Unload(win *pixelgl.Window) {
	data.EditorDraw = true
	if data.Editor != nil {
		if data.Editor.PosTop {
			ui.OpenDialog(constants.DialogEditorPanelTop)
			ui.OpenDialog(constants.DialogEditorOptionsBot)
		} else {
			ui.OpenDialog(constants.DialogEditorPanelLeft)
			ui.OpenDialog(constants.DialogEditorOptionsRight)
		}
	}
	ui.CloseDialog(constants.DialogPlayer1Inv)
	ui.CloseDialog(constants.DialogPlayer2Inv)
	ui.CloseDialog(constants.DialogPlayer3Inv)
	ui.CloseDialog(constants.DialogPlayer4Inv)
	ui.CloseDialog(constants.DialogPuzzleTitle)
	ui.CloseDialog(constants.DialogPuzzleTimer)
	systems.DisposeInGameDialogs()
	systems.LevelSessionDispose()
	systems.DisposeCurrLevel()
	systems.ClearTemp()
	systems.InitPuzzle(data.CurrentPlayArea)
	data.CurrPuzzleSet.CurrPuzzle.Update = true
	sfx.MusicPlayer.GetStream("game").Stop()
}

func (s *testState) Load(win *pixelgl.Window) {
	data.EditorDraw = false
	if data.Editor != nil {
		if data.Editor.PosTop {
			ui.CloseDialog(constants.DialogEditorPanelTop)
			ui.CloseDialog(constants.DialogEditorOptionsBot)
		} else {
			ui.CloseDialog(constants.DialogEditorPanelLeft)
			ui.CloseDialog(constants.DialogEditorOptionsRight)
		}
	}
	systems.InGameDialogs(win)
	systems.LevelSessionInit()
	systems.StartLevel(false)
	systems.UpdateViews()
	reanimator.SetFrameRate(constants.Configuration.Gameplay.FrameRate)
	reanimator.Reset()
	sfx.MusicPlayer.GetStream("game").RepeatTrack(data.CurrLevel.Metadata.MusicTrack)
}

func (s *testState) Update(win *pixelgl.Window) {
	data.P1Input.Update(win, viewport.MainCamera.Mat)
	data.P2Input.Update(win, viewport.MainCamera.Mat)
	data.P3Input.Update(win, viewport.MainCamera.Mat)
	data.P4Input.Update(win, viewport.MainCamera.Mat)
	systems.CursorSystem(true)
	debug.AddText("Test State")
	systems.InGameDebugInfo()
	systems.InGameDebugInput()

	if reanimator.FRate != constants.Configuration.Gameplay.FrameRate {
		reanimator.SetFrameRate(constants.Configuration.Gameplay.FrameRate)
	}
	reanimator.Update()

	// function systems
	systems.FunctionSystem()
	systems.InterpolationSystem()
	//systems.AnimationTransitionSystem()
	systems.TestSystem()
	ui.DialogStackOpen = len(ui.DialogStack) > 0
	systems.DialogSystem(win)

	if !ui.DialogStackOpen {
		// custom systems
		systems.InGameSystem()
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

func (s *testState) Draw(win *pixelgl.Window) {
	drawPlayState(win)
}

func (s *testState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}
