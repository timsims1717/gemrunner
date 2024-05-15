package systems

import (
	"gemrunner/internal/data"
	"github.com/gopxl/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
)

func UpdateInputType(win *pixelgl.Window) {
	mouseMoved := data.MenuInput.MouseMoved
	if mouseMoved || data.MenuInput.Get("click").JustPressed() {
		data.MenuInputUsed = pxginput.KeyboardMouse
		return
	}
	if data.MainJoystick == -1 {
		joysticks := pxginput.GetAllGamepads(win)
		for _, js := range joysticks {
			pressed := pxginput.GetAllJustPressedGamepad(win, js)
			if len(pressed) > 0 {
				data.MainJoystick = int(js)
				data.MenuInputUsed = pxginput.Gamepad
				return
			}
		}
	} else {
		pressed := pxginput.GetAllJustPressedGamepad(win, pixelgl.Joystick(data.MainJoystick))
		if len(pressed) > 0 {
			data.MenuInputUsed = pxginput.Gamepad
			return
		}
	}
}
