package data

import (
	"gemrunner/pkg/object"
	"gemrunner/pkg/timing"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
)

var (
	P1Input = &pxginput.Input{
		Key: "p1_input",
		Buttons: map[string]*pxginput.ButtonSet{
			"left":      pxginput.NewWithButtons(pixelgl.KeyKP4, pixelgl.ButtonDpadLeft).AddAxis(pixelgl.AxisLeftX, false),
			"right":     pxginput.NewWithButtons(pixelgl.KeyKP6, pixelgl.ButtonDpadRight).AddAxis(pixelgl.AxisLeftX, true),
			"up":        pxginput.NewWithButtons(pixelgl.KeyKP8, pixelgl.ButtonDpadUp).AddAxis(pixelgl.AxisLeftY, false),
			"down":      pxginput.NewWithButtons(pixelgl.KeyKP5, pixelgl.ButtonDpadDown).AddAxis(pixelgl.AxisLeftY, true),
			"jump":      pxginput.NewWithButtons(pixelgl.KeyKP0, pixelgl.ButtonA),
			"pickUp":    pxginput.NewWithButtons(pixelgl.KeyKP1, pixelgl.ButtonX),
			"action":    pxginput.NewWithButtons(pixelgl.KeyKP2, pixelgl.ButtonB),
			"digLeft":   pxginput.New().AddKey(pixelgl.KeyKP7).AddAxis(pixelgl.AxisLeftTrigger, true),
			"digRight":  pxginput.New().AddKey(pixelgl.KeyKP9).AddAxis(pixelgl.AxisRightTrigger, true),
			"kill":      pxginput.NewWithButtons(pixelgl.KeyKPSubtract, pixelgl.ButtonBack),
			"pause":     pxginput.New().AddKey(pixelgl.KeyEscape).AddButton(pixelgl.ButtonStart),
			"speedUp":   pxginput.NewWithButtons(pixelgl.KeyEqual, pixelgl.ButtonRightBumper),
			"speedDown": pxginput.NewWithButtons(pixelgl.KeyMinus, pixelgl.ButtonLeftBumper),
		},
		Mode:     pxginput.Any,
		Deadzone: 0.1,
	}
	P2Input = &pxginput.Input{
		Key: "p2_input",
		Buttons: map[string]*pxginput.ButtonSet{
			"left":     pxginput.NewWithButtons(pixelgl.KeyA, pixelgl.ButtonDpadLeft).AddAxis(pixelgl.AxisLeftX, false),
			"right":    pxginput.NewWithButtons(pixelgl.KeyD, pixelgl.ButtonDpadRight).AddAxis(pixelgl.AxisLeftX, true),
			"up":       pxginput.NewWithButtons(pixelgl.KeyW, pixelgl.ButtonDpadUp).AddAxis(pixelgl.AxisLeftY, false),
			"down":     pxginput.NewWithButtons(pixelgl.KeyS, pixelgl.ButtonDpadDown).AddAxis(pixelgl.AxisLeftY, true),
			"jump":     pxginput.NewWithButtons(pixelgl.KeySpace, pixelgl.ButtonA),
			"pickUp":   pxginput.NewWithButtons(pixelgl.KeyLeftShift, pixelgl.ButtonX),
			"action":   pxginput.NewWithButtons(pixelgl.KeyLeftControl, pixelgl.ButtonB),
			"digLeft":  pxginput.New().AddKey(pixelgl.KeyQ).AddAxis(pixelgl.AxisLeftTrigger, true),
			"digRight": pxginput.New().AddKey(pixelgl.KeyE).AddAxis(pixelgl.AxisRightTrigger, true),
			"kill":     pxginput.NewWithButtons(pixelgl.KeyGraveAccent, pixelgl.ButtonBack),
			"pause":    pxginput.New().AddButton(pixelgl.ButtonStart),
		},
		Mode:     pxginput.Any,
		Deadzone: 0.1,
	}
	P3Input = &pxginput.Input{
		Key: "p3_input",
		Buttons: map[string]*pxginput.ButtonSet{
			"left":     pxginput.NewWithButtons(pixelgl.KeyG, pixelgl.ButtonDpadLeft).AddAxis(pixelgl.AxisLeftX, false),
			"right":    pxginput.NewWithButtons(pixelgl.KeyJ, pixelgl.ButtonDpadRight).AddAxis(pixelgl.AxisLeftX, true),
			"up":       pxginput.NewWithButtons(pixelgl.KeyY, pixelgl.ButtonDpadUp).AddAxis(pixelgl.AxisLeftY, false),
			"down":     pxginput.NewWithButtons(pixelgl.KeyH, pixelgl.ButtonDpadDown).AddAxis(pixelgl.AxisLeftY, true),
			"jump":     pxginput.NewWithButtons(pixelgl.KeyV, pixelgl.ButtonA),
			"pickUp":   pxginput.NewWithButtons(pixelgl.KeyB, pixelgl.ButtonX),
			"action":   pxginput.NewWithButtons(pixelgl.KeyN, pixelgl.ButtonB),
			"digLeft":  pxginput.New().AddKey(pixelgl.KeyT).AddAxis(pixelgl.AxisLeftTrigger, true),
			"digRight": pxginput.New().AddKey(pixelgl.KeyU).AddAxis(pixelgl.AxisRightTrigger, true),
			"kill":     pxginput.NewWithButtons(pixelgl.Key8, pixelgl.ButtonBack),
			"pause":    pxginput.New().AddButton(pixelgl.ButtonStart),
		},
		Mode:     pxginput.Any,
		Deadzone: 0.1,
	}
	P4Input = &pxginput.Input{
		Key: "p4_input",
		Buttons: map[string]*pxginput.ButtonSet{
			"left":     pxginput.NewWithButtons(pixelgl.KeyL, pixelgl.ButtonDpadLeft).AddAxis(pixelgl.AxisLeftX, false),
			"right":    pxginput.NewWithButtons(pixelgl.KeyApostrophe, pixelgl.ButtonDpadRight).AddAxis(pixelgl.AxisLeftX, true),
			"up":       pxginput.NewWithButtons(pixelgl.KeyP, pixelgl.ButtonDpadUp).AddAxis(pixelgl.AxisLeftY, false),
			"down":     pxginput.NewWithButtons(pixelgl.KeySemicolon, pixelgl.ButtonDpadDown).AddAxis(pixelgl.AxisLeftY, true),
			"jump":     pxginput.NewWithButtons(pixelgl.KeyRightAlt, pixelgl.ButtonA),
			"pickUp":   pxginput.NewWithButtons(pixelgl.KeyComma, pixelgl.ButtonX),
			"action":   pxginput.NewWithButtons(pixelgl.KeyPeriod, pixelgl.ButtonB),
			"digLeft":  pxginput.New().AddKey(pixelgl.KeyO).AddAxis(pixelgl.AxisLeftTrigger, true),
			"digRight": pxginput.New().AddKey(pixelgl.KeyLeftBracket).AddAxis(pixelgl.AxisRightTrigger, true),
			"kill":     pxginput.NewWithButtons(pixelgl.KeyBackspace, pixelgl.ButtonBack),
			"pause":    pxginput.New().AddButton(pixelgl.ButtonStart),
		},
		Mode:     pxginput.Any,
		Deadzone: 0.1,
	}
	MenuInput = &pxginput.Input{
		Buttons: map[string]*pxginput.ButtonSet{
			"click":          pxginput.NewJoyless(pixelgl.MouseButtonLeft),
			"rightClick":     pxginput.NewJoyless(pixelgl.MouseButtonRight),
			"space":          pxginput.NewJoyless(pixelgl.KeySpace),
			"escape":         pxginput.NewJoyless(pixelgl.KeyEscape),
			"left":           pxginput.NewJoyless(pixelgl.KeyLeft).AddButton(pixelgl.ButtonDpadLeft),
			"right":          pxginput.NewJoyless(pixelgl.KeyRight).AddButton(pixelgl.ButtonDpadRight),
			"up":             pxginput.NewJoyless(pixelgl.KeyUp).AddButton(pixelgl.ButtonDpadUp),
			"down":           pxginput.NewJoyless(pixelgl.KeyDown).AddButton(pixelgl.ButtonDpadDown),
			"backspace":      pxginput.NewJoyless(pixelgl.KeyBackspace),
			"delete":         pxginput.NewJoyless(pixelgl.KeyDelete),
			"enter":          pxginput.NewJoyless(pixelgl.KeyEnter),
			"home":           pxginput.NewJoyless(pixelgl.KeyHome),
			"end":            pxginput.NewJoyless(pixelgl.KeyEnd),
			"Brush":          pxginput.NewJoyless(pixelgl.KeyB),
			"Line":           pxginput.NewJoyless(pixelgl.KeyL),
			"Square":         pxginput.NewJoyless(pixelgl.KeyH),
			"Fill":           pxginput.NewJoyless(pixelgl.KeyG),
			"Erase":          pxginput.NewJoyless(pixelgl.KeyE),
			"Eyedrop":        pxginput.NewJoyless(pixelgl.KeyY),
			"Select":         pxginput.NewJoyless(pixelgl.KeyM),
			"ctrlCut":        pxginput.NewJoyless(pixelgl.KeyX),
			"ctrlCopy":       pxginput.NewJoyless(pixelgl.KeyC),
			"ctrlPaste":      pxginput.NewJoyless(pixelgl.KeyV),
			"FlipVertical":   pxginput.NewJoyless(pixelgl.KeyK),
			"FlipHorizontal": pxginput.NewJoyless(pixelgl.KeyU),
			"Wrench":         pxginput.NewJoyless(pixelgl.KeyP),
			"Wire":           pxginput.NewJoyless(pixelgl.KeyI),
			"Text":           pxginput.NewJoyless(pixelgl.KeyT),
			"ctrlUndo":       pxginput.NewJoyless(pixelgl.KeyZ),
			"ctrlShiftRedo":  pxginput.NewJoyless(pixelgl.KeyZ),
			"ctrlSave":       pxginput.NewJoyless(pixelgl.KeyS),
			"ctrlOpen":       pxginput.NewJoyless(pixelgl.KeyO),
			"Delete":         pxginput.NewJoyless(pixelgl.KeyDelete),
			//"altEyedrop":     pxginput.NewJoyless(pixelgl.KeyLeftAlt),
			"ctrl":   pxginput.NewJoyless(pixelgl.KeyLeftControl),
			"rCtrl":  pxginput.NewJoyless(pixelgl.KeyRightControl),
			"shift":  pxginput.NewJoyless(pixelgl.KeyLeftShift),
			"rShift": pxginput.NewJoyless(pixelgl.KeyRightShift),
		},
		Mode: pxginput.KeyboardMouse,
	}
	DebugInput = &pxginput.Input{
		Buttons: map[string]*pxginput.ButtonSet{
			"debugConsole": pxginput.NewJoyless(pixelgl.KeyGraveAccent),
			"debug":        pxginput.NewJoyless(pixelgl.KeyF3),
			"debugText":    pxginput.NewJoyless(pixelgl.KeyF4),
			"fullscreen":   pxginput.NewJoyless(pixelgl.KeyF5),
			"fuzzy":        pxginput.NewJoyless(pixelgl.KeyF6),
			"debugMenu":    pxginput.NewJoyless(pixelgl.KeyF7),
			"debugTest":    pxginput.NewJoyless(pixelgl.KeyF8),
			"debugPause":   pxginput.NewJoyless(pixelgl.KeyF9),
			"debugFrame":   pxginput.NewJoyless(pixelgl.KeyF10),
			"debugInv":     pxginput.NewJoyless(pixelgl.KeyF11),
			"switchWorld":  pxginput.NewJoyless(pixelgl.KeyTab),
			"beBadGuyP1":   pxginput.NewJoyless(pixelgl.Key1),
			"beBadGuyP2":   pxginput.NewJoyless(pixelgl.Key2),
			"beBadGuyP3":   pxginput.NewJoyless(pixelgl.Key3),
			"beBadGuyP4":   pxginput.NewJoyless(pixelgl.Key4),
			"debugSP":      pxginput.NewJoyless(pixelgl.KeyKPAdd),
			"debugSM":      pxginput.NewJoyless(pixelgl.KeyKPSubtract),
			"camUp":        pxginput.NewJoyless(pixelgl.KeyUp),
			"camRight":     pxginput.NewJoyless(pixelgl.KeyRight),
			"camDown":      pxginput.NewJoyless(pixelgl.KeyDown),
			"camLeft":      pxginput.NewJoyless(pixelgl.KeyLeft),
		},
		Mode: pxginput.KeyboardMouse,
	}
	HideCursorTimer *timing.Timer
	CursorObj       *object.Object
	CursorEntity    *ecs.Entity
)
