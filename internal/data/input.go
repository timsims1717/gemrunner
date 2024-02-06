package data

import (
	"github.com/gopxl/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
)

var (
	MenuInput = &pxginput.Input{
		Buttons: map[string]*pxginput.ButtonSet{
			"click":          pxginput.NewJoyless(pixelgl.MouseButtonLeft),
			"rightClick":     pxginput.NewJoyless(pixelgl.MouseButtonRight),
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
			"debugResume":  pxginput.NewJoyless(pixelgl.KeyF10),
			"debugInv":     pxginput.NewJoyless(pixelgl.KeyF11),
			"switchWorld":  pxginput.NewJoyless(pixelgl.KeyTab),
			//"debugSP":      pxginput.NewJoyless(pixelgl.KeyEqual),
			//"debugSM":      pxginput.NewJoyless(pixelgl.KeyMinus),
			"camUp":    pxginput.NewJoyless(pixelgl.KeyKP8),
			"camRight": pxginput.NewJoyless(pixelgl.KeyKP6),
			"camDown":  pxginput.NewJoyless(pixelgl.KeyKP5),
			"camLeft":  pxginput.NewJoyless(pixelgl.KeyKP4),
		},
		Mode: pxginput.KeyboardMouse,
	}
)
