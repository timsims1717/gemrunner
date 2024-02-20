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
	"gemrunner/pkg/state"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
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
	reanimator.SetFrameRate(15)
	reanimator.Reset()
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

	systems.DialogSystem()

	if !data.DialogStackOpen {
		// custom systems
		systems.InGameSystem()
		systems.CharacterActionSystem()
		systems.DynamicSystem()
		systems.CollisionSystem()
		systems.CharacterStateSystem()
		systems.CollectSystem()
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
	//systems.DrawLayerSystem(data.PuzzleView.Canvas, 10) // bg tiles, ladders
	//img.Clear()
	//systems.DrawLayerSystem(data.PuzzleView.Canvas, 11) // gems
	//img.Clear()
	//systems.DrawLayerSystem(data.PuzzleView.Canvas, 18) // other reanimator
	//img.Clear()
	//systems.DrawLayerSystem(data.PuzzleView.Canvas, 19) // fg tiles, liquid, crushers
	//img.Clear()
	//systems.DrawLayerSystem(data.PuzzleView.Canvas, 20) // ui
	//img.Clear()
	//data.IMDraw.Draw(data.PuzzleView.Canvas)
	if debug.ShowDebug {
		debug.DrawLines(data.PuzzleView.Canvas)
	}
	data.PuzzleView.Draw(win)
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
