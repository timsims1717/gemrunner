package data

import (
	"github.com/faiface/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
)

var (
	EditorInput = &pxginput.Input{
		Buttons: map[string]*pxginput.ButtonSet{
			"click":          pxginput.NewJoyless(pixelgl.MouseButtonLeft),
			"rightClick":     pxginput.NewJoyless(pixelgl.MouseButtonRight),
			"Brush":          pxginput.NewJoyless(pixelgl.KeyB),
			"Line":           pxginput.NewJoyless(pixelgl.KeyL),
			"Square":         pxginput.NewJoyless(pixelgl.KeyH),
			"Ellipse":        pxginput.NewJoyless(pixelgl.KeyO),
			"Fill":           pxginput.NewJoyless(pixelgl.KeyG),
			"Erase":          pxginput.NewJoyless(pixelgl.KeyE),
			"Eyedrop":        pxginput.NewJoyless(pixelgl.KeyY),
			"Select":         pxginput.NewJoyless(pixelgl.KeyM),
			"ctrlCut":        pxginput.NewJoyless(pixelgl.KeyX),
			"ctrlCopy":       pxginput.NewJoyless(pixelgl.KeyC),
			"ctrlPaste":      pxginput.NewJoyless(pixelgl.KeyV),
			"FlipVertical":   pxginput.NewJoyless(pixelgl.KeyK),
			"FlipHorizontal": pxginput.NewJoyless(pixelgl.KeyU),
			"Pliers":         pxginput.NewJoyless(pixelgl.KeyP),
			"ctrlUndo":       pxginput.NewJoyless(pixelgl.KeyZ),
			"ctrlShiftRedo":  pxginput.NewJoyless(pixelgl.KeyZ),
			"Delete":         pxginput.NewJoyless(pixelgl.KeyDelete),
			//"altEyedrop":     pxginput.NewJoyless(pixelgl.KeyLeftAlt),
			"ctrl":  pxginput.NewJoyless(pixelgl.KeyLeftControl),
			"shift": pxginput.NewJoyless(pixelgl.KeyLeftShift),
		},
		Mode: pxginput.KeyboardMouse,
	}
	DebugInput = &pxginput.Input{
		Buttons: map[string]*pxginput.ButtonSet{
			"debugConsole": pxginput.NewJoyless(pixelgl.KeyGraveAccent),
			"debug":        pxginput.NewJoyless(pixelgl.KeyF3),
			"debugText":    pxginput.NewJoyless(pixelgl.KeyF4),
			"debugMenu":    pxginput.NewJoyless(pixelgl.KeyF7),
			"debugTest":    pxginput.NewJoyless(pixelgl.KeyF8),
			"debugPause":   pxginput.NewJoyless(pixelgl.KeyF9),
			"debugResume":  pxginput.NewJoyless(pixelgl.KeyF10),
			"debugInv":     pxginput.NewJoyless(pixelgl.KeyF11),
			"switchWorld":  pxginput.NewJoyless(pixelgl.KeyTab),
			//"debugSP":      pxginput.NewJoyless(pixelgl.KeyEqual),
			//"debugSM":      pxginput.NewJoyless(pixelgl.KeyMinus),
			//"camUp":        pxginput.NewJoyless(pixelgl.KeyP),
			//"camRight":     pxginput.NewJoyless(pixelgl.KeyApostrophe),
			//"camDown":      pxginput.NewJoyless(pixelgl.KeySemicolon),
			//"camLeft":      pxginput.NewJoyless(pixelgl.KeyL),
		},
		Mode: pxginput.KeyboardMouse,
	}
)
