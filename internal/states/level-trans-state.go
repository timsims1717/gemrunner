package states

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/content"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems"
	"gemrunner/internal/ui"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/img"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/sfx"
	"gemrunner/pkg/state"
	"github.com/gopxl/pixel/pixelgl"
)

var (
	LevelTransState = &levelTransState{}
)

type levelTransState struct {
	*state.AbstractState
}

func (s *levelTransState) Unload(win *pixelgl.Window) {
	data.LevelTrans = false
	ui.ClearDialogsOpen()
	ui.ClearDialogStack()
	systems.DisposeLevel(data.CurrentPlayArea.Level)
	data.CurrentPlayArea.Level = nil
	//systems.DisposePlayArea(data.CurrentPlayArea)
	data.CurrentPlayArea = data.OtherPlayArea
	data.CurrentPlayArea.Border.ExcludeSide = data.NoDirection
	data.CurrentPlayArea.Border.ExcludeSize = 0
	data.OtherPlayArea = nil

	systems.RandAndRecord(data.CurrLevel.Recording)
	systems.UpdateViews()
	reanimator.SetFrameRate(constants.Configuration.Gameplay.FrameRate)
	reanimator.Reset()
	sfx.MusicPlayer.GetStream("game").RepeatTrack(data.CurrLevel.Metadata.MusicTrack)
	go content.SaveSaveGame()
}

func (s *levelTransState) Load(win *pixelgl.Window) {
	data.LevelTrans = true
	data.EditorDraw = false
	ui.ClearDialogsOpen()
	ui.ClearDialogStack()
	systems.ClearTemp()
	data.CurrentPlayArea.LayerOffset = 1000
	systems.UpdateLevelLayer(data.CurrentPlayArea.Level, 1000)
	data.CurrLevelSess.PuzzleIndex = data.CurrPuzzleSet.PuzzleIndex
	if data.OtherPlayArea == nil {
		data.OtherPlayArea = systems.CreatePlayArea()
	}
	systems.SetPuzzle(data.OtherPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	systems.InitLevel(data.OtherPlayArea)
	systems.UpdateViews()
	systems.SetTransitionBorders()
	systems.StartInterpolation()
}

func (s *levelTransState) Update(win *pixelgl.Window) {
	systems.CursorSystem(true)
	debug.AddText("Level Transition State")
	systems.InGameDebugInfo()
	systems.TransitionDebugInput()

	systems.InterpolationSystem()
	systems.LevelTransitionSystem()

	systems.TileSpriteSystemPre()
	systems.TileSpriteSystem()

	// object systems
	systems.ParentSystem()
	systems.ObjectSystem()
	systems.EffectsSystem()

	data.CurrentPlayArea.BorderView.Update()
	data.CurrentPlayArea.PuzzleView.Update()
	data.CurrentPlayArea.WorldView.Update()
	data.CurrentPlayArea.PuzzleViewNoShader.Update()
	data.OtherPlayArea.BorderView.Update()
	data.OtherPlayArea.PuzzleView.Update()
	data.OtherPlayArea.WorldView.Update()
	data.OtherPlayArea.PuzzleViewNoShader.Update()
	data.ScreenView.Update()

	myecs.UpdateManager()
	debug.AddText(fmt.Sprintf("Entity Count: %d", myecs.FullCount))
}

func (s *levelTransState) Draw(win *pixelgl.Window) {
	data.ScreenView.Canvas.Clear(constants.ColorBlack)
	drawPlayArea(win, data.CurrentPlayArea)
	drawPlayArea(win, data.OtherPlayArea)
	// dialog draw system
	systems.DialogDrawSystem(data.ScreenView.Canvas)
	systems.DrawLayerSystem(data.ScreenView.Canvas, -10)
	img.Clear()
	systems.TemporarySystem()
}

func (s *levelTransState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}
