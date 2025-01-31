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
	"gemrunner/pkg/world"
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
	systems.LevelSessionDispose()
	systems.LevelDispose()
	systems.ClearTemp()
	systems.PuzzleViewInit()
	systems.UpdateEditorShaders()
	systems.UpdatePuzzleShaders()
	systems.ChangeWorldShader(data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderMode)
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
	systems.LevelSessionInit()
	systems.LevelInit()
	systems.UpdateViews()
	reanimator.SetFrameRate(constants.FrameRate)
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
	debug.AddText(fmt.Sprintf("Speed: %d", constants.FrameRate))
	debug.AddText(fmt.Sprintf("Frame Number: %d", data.CurrLevel.FrameNumber))
	debug.AddText(fmt.Sprintf("Frame Counter: %d", data.CurrLevel.FrameCounter))
	debug.AddText(fmt.Sprintf("Frame Cycle: %d", data.CurrLevel.FrameCycle))
	debug.AddTruthText("Frame Change", data.CurrLevel.FrameChange)
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
				item := "unknown"
				cd, ok1 := player.Inventory.GetComponentData(myecs.PickUp)
				if ok1 {
					if pd, ok := cd.(*data.PickUp); ok {
						item = pd.Name
					}
				}
				debug.AddText(fmt.Sprintf("Player %d Inv: %s", i+1, item))
			}
			debug.AddText(fmt.Sprintf("Player %d # of Tiles: %d", i+1, len(player.StoredBlocks)))
		}
	}

	if reanimator.FRate != constants.FrameRate {
		reanimator.SetFrameRate(constants.FrameRate)
	}
	reanimator.Update()

	// function systems
	systems.FunctionSystem()
	systems.InterpolationSystem()
	systems.TestSystem()
	ui.DialogStackOpen = len(ui.DialogStack) > 0
	systems.DialogSystem(win)

	if !ui.DialogStackOpen {
		// custom systems
		systems.InGameSystem()
		systems.CharacterActionSystem()
		systems.DynamicSystem()
		systems.CollisionSystem()
		systems.CharacterStateSystem()
		systems.TouchSystem()
		systems.SmashSystem()
		systems.TileSystem()
		systems.MagicSystem()
		systems.TileSpriteSystemPre()
		systems.TileSpriteSystem()
		systems.FloatingTextSystem()
	} else {

	}
	// object systems
	systems.AnimationSystem()
	systems.ParentSystem()
	systems.ObjectSystem()
	systems.ShaderSystem()

	data.BorderView.Update()
	data.PuzzleView.Update()
	data.WorldView.Update()
	data.PuzzleViewNoShader.Update()

	myecs.UpdateManager()
	debug.AddText(fmt.Sprintf("Entity Count: %d", myecs.FullCount))
}

func (s *testState) Draw(win *pixelgl.Window) {
	drawPlayArea(win)
}

func (s *testState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}
