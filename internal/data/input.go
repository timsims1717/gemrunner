package data

import (
	"github.com/faiface/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
)

var (
	EditorInput = &pxginput.Input{
		Buttons: map[string]*pxginput.ButtonSet{
			"click":      pxginput.NewJoyless(pixelgl.MouseButtonLeft),
			"rightClick": pxginput.NewJoyless(pixelgl.MouseButtonRight),
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
			"debugSP":      pxginput.NewJoyless(pixelgl.KeyEqual),
			"debugSM":      pxginput.NewJoyless(pixelgl.KeyMinus),
			"camUp":        pxginput.NewJoyless(pixelgl.KeyP),
			"camRight":     pxginput.NewJoyless(pixelgl.KeyApostrophe),
			"camDown":      pxginput.NewJoyless(pixelgl.KeySemicolon),
			"camLeft":      pxginput.NewJoyless(pixelgl.KeyL),
		},
		Mode: pxginput.KeyboardMouse,
	}
)
