package data

import (
	"github.com/gopxl/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
)

var (
	P1Input = &pxginput.Input{
		Key: "p1_input",
		Buttons: map[string]*pxginput.ButtonSet{
			"left":      pxginput.New(pixelgl.KeyKP4, pixelgl.ButtonDpadLeft),
			"right":     pxginput.New(pixelgl.KeyKP6, pixelgl.ButtonDpadRight),
			"up":        pxginput.New(pixelgl.KeyKP8, pixelgl.ButtonDpadUp),
			"down":      pxginput.New(pixelgl.KeyKP5, pixelgl.ButtonDpadDown),
			"jump":      pxginput.New(pixelgl.KeyKP0, pixelgl.ButtonA),
			"pickUp":    pxginput.New(pixelgl.KeyKP1, pixelgl.ButtonX),
			"stow":      pxginput.New(pixelgl.KeyKP3, pixelgl.ButtonY),
			"action":    pxginput.New(pixelgl.KeyKP2, pixelgl.ButtonB),
			"kill":      pxginput.New(pixelgl.KeyK, pixelgl.ButtonBack),
			"speedUp":   pxginput.New(pixelgl.KeyEqual, pixelgl.ButtonRightBumper),
			"speedDown": pxginput.New(pixelgl.KeyMinus, pixelgl.ButtonLeftBumper),
		},
		Mode: pxginput.Any,
	}
	P2Input = &pxginput.Input{
		Key: "p2_input",
		Buttons: map[string]*pxginput.ButtonSet{
			"left":   pxginput.New(pixelgl.KeyA, pixelgl.ButtonDpadLeft),
			"right":  pxginput.New(pixelgl.KeyD, pixelgl.ButtonDpadRight),
			"up":     pxginput.New(pixelgl.KeyW, pixelgl.ButtonDpadUp),
			"down":   pxginput.New(pixelgl.KeyS, pixelgl.ButtonDpadDown),
			"jump":   pxginput.New(pixelgl.KeySpace, pixelgl.ButtonA),
			"pickUp": pxginput.New(pixelgl.KeyLeftShift, pixelgl.ButtonX),
			"stow":   pxginput.New(pixelgl.KeyC, pixelgl.ButtonY),
			"action": pxginput.New(pixelgl.KeyLeftControl, pixelgl.ButtonB),
			"kill":   pxginput.New(pixelgl.KeyX, pixelgl.ButtonBack),
		},
		Mode: pxginput.Any,
	}
	P3Input = &pxginput.Input{
		Key: "p3_input",
		Buttons: map[string]*pxginput.ButtonSet{
			"left":   pxginput.New(pixelgl.KeyG, pixelgl.ButtonDpadLeft),
			"right":  pxginput.New(pixelgl.KeyJ, pixelgl.ButtonDpadRight),
			"up":     pxginput.New(pixelgl.KeyY, pixelgl.ButtonDpadUp),
			"down":   pxginput.New(pixelgl.KeyH, pixelgl.ButtonDpadDown),
			"jump":   pxginput.New(pixelgl.KeyV, pixelgl.ButtonA),
			"pickUp": pxginput.New(pixelgl.KeyB, pixelgl.ButtonX),
			"stow":   pxginput.New(pixelgl.KeyM, pixelgl.ButtonY),
			"action": pxginput.New(pixelgl.KeyN, pixelgl.ButtonB),
			"kill":   pxginput.New(pixelgl.KeyM, pixelgl.ButtonBack),
		},
		Mode: pxginput.Any,
	}
	P4Input = &pxginput.Input{
		Key: "p4_input",
		Buttons: map[string]*pxginput.ButtonSet{
			"left":   pxginput.New(pixelgl.KeyL, pixelgl.ButtonDpadLeft),
			"right":  pxginput.New(pixelgl.KeyApostrophe, pixelgl.ButtonDpadRight),
			"up":     pxginput.New(pixelgl.KeyP, pixelgl.ButtonDpadUp),
			"down":   pxginput.New(pixelgl.KeySemicolon, pixelgl.ButtonDpadDown),
			"jump":   pxginput.New(pixelgl.KeyRightAlt, pixelgl.ButtonA),
			"pickUp": pxginput.New(pixelgl.KeyComma, pixelgl.ButtonX),
			"stow":   pxginput.New(pixelgl.KeySlash, pixelgl.ButtonY),
			"action": pxginput.New(pixelgl.KeyPeriod, pixelgl.ButtonB),
			"kill":   pxginput.New(pixelgl.KeySlash, pixelgl.ButtonBack),
		},
		Mode: pxginput.Any,
	}
	MenuInput = &pxginput.Input{
		Buttons: map[string]*pxginput.ButtonSet{
			"click":          pxginput.NewJoyless(pixelgl.MouseButtonLeft),
			"rightClick":     pxginput.NewJoyless(pixelgl.MouseButtonRight),
			"escape":         pxginput.NewJoyless(pixelgl.KeyEscape),
			"left":           pxginput.NewJoyless(pixelgl.KeyLeft),
			"right":          pxginput.NewJoyless(pixelgl.KeyRight),
			"up":             pxginput.NewJoyless(pixelgl.KeyUp),
			"down":           pxginput.NewJoyless(pixelgl.KeyDown),
			"backspace":      pxginput.NewJoyless(pixelgl.KeyBackspace),
			"delete":         pxginput.NewJoyless(pixelgl.KeyDelete),
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
			//"debugSP":      pxginput.NewJoyless(pixelgl.KeyEqual),
			//"debugSM":      pxginput.NewJoyless(pixelgl.KeyMinus),
			//"camUp":    pxginput.NewJoyless(pixelgl.KeyKP8),
			//"camRight": pxginput.NewJoyless(pixelgl.KeyKP6),
			//"camDown":  pxginput.NewJoyless(pixelgl.KeyKP5),
			//"camLeft":  pxginput.NewJoyless(pixelgl.KeyKP4),
		},
		Mode: pxginput.KeyboardMouse,
	}
)
