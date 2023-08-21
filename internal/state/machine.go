package state

import (
	"fmt"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/state"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/viewport"
	"github.com/faiface/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
	"golang.org/x/image/colornames"
)

const (
	EditorStateKey = "editor_state"
)

var (
	EditorState = &editorState{}
	States      = map[string]*state.AbstractState{
		EditorStateKey: state.New(EditorState, false),
	}
)

var (
	switchState = true
	currState   = "unknown"
	nextState   = EditorStateKey
	loading     = false
	loadingDone = false
	done        = make(chan struct{})

	debugInput = &pxginput.Input{
		Buttons: map[string]*pxginput.ButtonSet{
			"debugConsole": pxginput.NewJoyless(pixelgl.KeyGraveAccent),
			"debug":        pxginput.NewJoyless(pixelgl.KeyF3),
			"debugText":    pxginput.NewJoyless(pixelgl.KeyF4),
			"debugMenu":    pxginput.NewJoyless(pixelgl.KeyF7),
			"debugTest":    pxginput.NewJoyless(pixelgl.KeyF8),
			"debugPause":   pxginput.NewJoyless(pixelgl.KeyF9),
			"debugResume":  pxginput.NewJoyless(pixelgl.KeyF10),
			"debugInv":     pxginput.NewJoyless(pixelgl.KeyF11),
			"debugSP":      pxginput.NewJoyless(pixelgl.KeyEqual),
			"debugSM":      pxginput.NewJoyless(pixelgl.KeyMinus),
			"camUp":        pxginput.NewJoyless(pixelgl.KeyP),
			"camRight":     pxginput.NewJoyless(pixelgl.KeyApostrophe),
			"camDown":      pxginput.NewJoyless(pixelgl.KeySemicolon),
			"camLeft":      pxginput.NewJoyless(pixelgl.KeyL),
		},
		Mode: pxginput.KeyboardMouse,
	}
	baseCanvas *pixelgl.Canvas
)

func Update(win *pixelgl.Window) {
	timing.Update()
	updateState()
	if loading {
		select {
		case <-done:
			loading = false
			loadingDone = true
			currState = nextState
		default:
			//LoadingState.Update(win)
		}
	} else {
		debugInput.Update(win, viewport.MainCamera.Mat)
		//if debugInput.Get("debug").JustPressed() {
		//	debug.Debug = !debug.Debug
		//	if debug.Debug {
		//		fmt.Println("DEBUG ON")
		//	} else {
		//		fmt.Println("DEBUG OFF")
		//	}
		//}
		//if debugInput.Get("debugText").JustPressed() {
		//	debug.Text = !debug.Text
		//	if debug.Text {
		//		fmt.Println("DEBUG TEXT ON")
		//	} else {
		//		fmt.Println("DEBUG TEXT OFF")
		//	}
		//}
		//if debugInput.Get("debugInv").JustPressed() {
		//	for _, d := range descent.Descent.Dwarves {
		//		d.Health.Inv = !d.Health.Inv
		//	}
		//}
		//if debugInput.Get("debugMenu").JustPressed() && MenuClosed() {
		//	debugInput.Get("debugMenu").Consume()
		//	OpenMenu(DebugMenu)
		//}
		//if debugInput.Get("debugTest").JustPressed() {
		//	player := descent.Descent.GetPlayers()[0].Player
		//	descent.CreateGnome(descent.Descent.Cave, descent.Descent.Cave.GetTile(player.CamPos.Sub(player.CanvasPos.Sub(debugInput.World))).Transform.Pos)
		//}
		//if debugInput.Get("debugSP").JustPressed() {
		//	if currState == MenuStateKey {
		//		MenuState.splashScale *= 1.2
		//		fmt.Printf("Splash Scale: %f\n", MenuState.splashScale)
		//	} else if descent.Descent.FreeCam {
		//		camera.Cam.ZoomIn(1.)
		//	}
		//}
		//if debugInput.Get("debugSM").JustPressed() {
		//	if currState == MenuStateKey {
		//		MenuState.splashScale /= 1.2
		//		fmt.Printf("Splash Scale: %f\n", MenuState.splashScale)
		//	} else if descent.Descent.FreeCam {
		//		camera.Cam.ZoomIn(-1.)
		//	}
		//}
		//if descent.Descent.FreeCam && len(descent.Descent.Dwarves) > 0 {
		//	if debugInput.Get("freeCamUp").Pressed() {
		//		//camera.Cam.Up()
		//		descent.Descent.Dwarves[0].Player.CamPos.Y += 100. * timing.DT
		//	} else if debugInput.Get("freeCamDown").Pressed() && descent.Descent.FreeCam {
		//		//camera.Cam.Down()
		//		descent.Descent.Dwarves[0].Player.CamPos.Y -= 100. * timing.DT
		//	}
		//	if debugInput.Get("freeCamRight").Pressed() && descent.Descent.FreeCam {
		//		//camera.Cam.Right()
		//		descent.Descent.Dwarves[0].Player.CamPos.X += 100. * timing.DT
		//	} else if debugInput.Get("freeCamLeft").Pressed() && descent.Descent.FreeCam {
		//		//camera.Cam.Left()
		//		descent.Descent.Dwarves[0].Player.CamPos.X -= 100. * timing.DT
		//	}
		//}
		//frame := false
		//if debugInput.Get("debugPause").JustPressed() {
		//	if !debugPause {
		//		fmt.Println("DEBUG PAUSE")
		//		debugPause = true
		//	} else {
		//		frame = true
		//	}
		//} else if debugInput.Get("debugResume").JustPressed() {
		//	fmt.Println("DEBUG RESUME")
		//	debugPause = false
		//}
		//if !debugPause || frame {
		if cState, ok := States[currState]; ok {
			cState.Update(win)
		}
		//}
	}
	viewport.MainCamera.Update()
	myecs.UpdateManager()
}

func Draw(win *pixelgl.Window) {
	img.Clear()
	cState, ok1 := States[currState]
	nState, ok2 := States[nextState]
	if !ok2 {
		panic(fmt.Sprintf("state %s doesn't exist", nextState))
	}
	if loading && nState.ShowLoad || !ok1 {
		//LoadingState.Draw(win)
	} else {
		win.Clear(colornames.Blue)
		cState.Draw(win)
		win.Update()
	}
}

func updateState() {
	if !loading && (currState != nextState || switchState) {
		// uninitialize
		img.FullClear()
		if cState, ok := States[currState]; ok {
			go cState.Unload()
		}
		// initialize
		if nState, ok := States[nextState]; ok {
			go nState.Load(done)
			loading = true
			loadingDone = false
		}
		switchState = false
	}
}

func SwitchState(s string) {
	if !switchState {
		switchState = true
		nextState = s
	}
}
