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
	reanimator.SetFrameRate(15)
	reanimator.Reset()
}

func (s *testState) Update(win *pixelgl.Window) {
	data.GameInput.Update(win, viewport.MainCamera.Mat)
	debug.AddText("Test State")
	p1 := data.CurrLevel.Player1
	p1Pos := p1.FauxObj.Pos
	debug.AddIntCoords("Player Pos", int(p1Pos.X), int(p1Pos.Y))
	cx, cy := world.WorldToMap(p1Pos.X, p1Pos.Y)
	debug.AddIntCoords("Player Coords", cx, cy)
	tile := data.CurrLevel.Tiles.Get(cx, cy)
	debug.AddIntCoords("Tile Pos", int(tile.Object.Pos.X), int(tile.Object.Pos.Y))
	debug.AddText(fmt.Sprintf("Player 1 Score: %d", data.CurrLevel.Stats1.Score))

	reanimator.Update()

	// function systems
	systems.FunctionSystem()
	systems.TestSystem()

	systems.DialogSystem()

	if !data.DialogStackOpen {
		// custom systems
		systems.CollisionSystem()
		systems.CharacterSystem()
		systems.CollectSystem()
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
	systems.NewDrawSystem(data.PuzzleView.Canvas, 10) // bg tiles/items/gems
	img.Clear()
	systems.NewDrawSystem(data.PuzzleView.Canvas, 11) // characters
	img.Clear()
	systems.NewDrawSystem(data.PuzzleView.Canvas, 12) // fg tiles
	img.Clear()
	systems.NewDrawSystem(data.PuzzleView.Canvas, 14) // ui
	img.Clear()
	//data.IMDraw.Draw(data.PuzzleView.Canvas)
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