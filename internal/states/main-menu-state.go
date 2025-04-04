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
	"gemrunner/pkg/options"
	"gemrunner/pkg/sfx"
	"gemrunner/pkg/state"
	"github.com/gopxl/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
)

var (
	MainMenuState = &mainMenuState{}
)

type mainMenuState struct {
	*state.AbstractState
}

func (s *mainMenuState) Unload(win *pixelgl.Window) {
	ui.ClearDialogsOpen()
	ui.ClearDialogStack()
	systems.DisposeMainDialogs()
	sfx.MusicPlayer.GetStream("game").Stop()
}

func (s *mainMenuState) Load(win *pixelgl.Window) {
	ui.ClearDialogsOpen()
	ui.ClearDialogStack()
	systems.MainDialogs(win)
	systems.PreLoadCustomPuzzleList()
	ui.OpenDialog(constants.DialogMainMenu)
	systems.UpdateViews()
	sfx.MusicPlayer.GetStream("game").RepeatTrack(constants.TrackMenu)
}

func (s *mainMenuState) Update(win *pixelgl.Window) {
	systems.CursorSystem(false)
	debug.AddText("Main Menu State")
	debug.AddIntCoords("World", int(data.MenuInput.World.X), int(data.MenuInput.World.Y))
	switch data.MenuInputUsed {
	case pxginput.Any:
		debug.AddText("Input Type: Any")
	case pxginput.KeyboardMouse:
		debug.AddText("Input Type: Keyboard/Mouse")
	case pxginput.Gamepad:
		debug.AddText(fmt.Sprintf("Input Type: Gamepad %d", data.MainJoystick))
	}
	debug.AddText(fmt.Sprintf("Number of Players: %d", len(data.Players)))
	if data.DebugInput.Get("debugTest").JustPressed() {
		err := content.LoadReplay("APlayTest_0_2025.02.19.14.09.01.replay")
		if err != nil {
			panic(err)
		}
		err = systems.OpenPuzzleSet(data.CurrReplay.PuzzleSet + constants.PuzzleExt)
		if err != nil {
			panic(err)
		}
		data.CurrPuzzleSet.SetTo(data.CurrReplay.PuzzleNum)
		if data.CurrPuzzleSet.Metadata.NumPlayers < 1 {
			data.CurrPuzzleSet.Metadata.NumPlayers = data.CurrPuzzleSet.CurrPuzzle.NumPlayers()
		}
		systems.PuzzleInit()
		state.SwitchState(constants.PlayStateKey)
	}

	// function systems
	systems.FunctionSystem()
	systems.InterpolationSystem()

	ui.DialogStackOpen = len(ui.DialogStack) > 0
	if !ui.DialogStackOpen {
	} else {
	}
	systems.UpdateInputType(win)
	systems.DialogSystem(win)
	// object systems
	systems.AnimationSystem()
	systems.ParentSystem()
	systems.ObjectSystem()

	data.ScreenView.Update()

	myecs.UpdateManager()
	debug.AddText(fmt.Sprintf("Entity Count: %d", myecs.FullCount))
}

func (s *mainMenuState) Draw(win *pixelgl.Window) {
	data.ScreenView.Canvas.Clear(constants.ColorBlack)
	// draw border
	//data.BorderView.Canvas.Clear(constants.ColorBlack)
	//systems.BorderSystem(1)
	//img.Batchers[constants.UIBatch].Draw(data.BorderView.Canvas)
	//img.Clear()
	//data.BorderView.Draw(win)
	// dialog draw system
	systems.DialogDrawSystem(data.ScreenView.Canvas)
	systems.DrawLayerSystem(data.ScreenView.Canvas, -10)
	img.Clear()
	data.ScreenView.Draw(win)
	systems.TemporarySystem()
	//data.IMDraw.Clear()
	if options.Updated {
		systems.UpdateViews()
	}
}

func (s *mainMenuState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}
