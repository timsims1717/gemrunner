package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/pkg/img"
	"github.com/gopxl/pixel/pixelgl"
	"image/color"
)

func DialogSystem() {
	if len(data.DialogStack) > 0 {
		data.DialogStackOpen = true
	} else {
		data.DialogStackOpen = false
	}
	layer := 100
	for _, dialog := range data.DialogsOpen {
		if !dialog.NoBorder {
			dialog.BorderVP.Update()
			dialog.BorderObject.Layer = layer
		}
		dialog.ViewPort.Update()
		for _, e := range dialog.Elements {
			if spr, okS := e.(*data.SprElement); okS {
				spr.Object.Layer = layer
			} else if btn, okB := e.(*data.Button); okB {
				btn.Object.Layer = layer
			} else if txt, okT := e.(*data.Text); okT {
				txt.Text.Obj.Layer = layer
			}
		}
		// update this dialog box with input
		if !data.DialogStackOpen {

		}
		layer++
	}
	closeKey := ""
	for i, dialog := range data.DialogStack {
		if !dialog.NoBorder {
			dialog.BorderVP.Update()
			dialog.BorderObject.Layer = layer
		}
		dialog.ViewPort.Update()
		for _, e := range dialog.Elements {
			if spr, okS := e.(*data.SprElement); okS {
				spr.Object.Layer = layer
			} else if btn, okB := e.(*data.Button); okB {
				btn.Object.Layer = layer
			} else if txt, okT := e.(*data.Text); okT {
				txt.Text.Obj.Layer = layer
			}
		}
		if i == len(data.DialogStack)-1 {
			// update this dialog box with input
			if data.MainInput.Get("escape").JustPressed() {
				data.MainInput.Get("escape").Consume()
				closeKey = dialog.Key
			}
		}
		layer++
	}
	if closeKey != "" {
		data.CloseDialog(closeKey)
	}
}

func DialogDrawSystem(win *pixelgl.Window) {
	layer := 100
	for _, dialog := range data.DialogsOpen {
		if !dialog.NoBorder {
			dialog.BorderVP.Canvas.Clear(color.RGBA{})
			BorderSystem(layer)
			img.Batchers[constants.UIBatch].Draw(dialog.BorderVP.Canvas)
			dialog.BorderVP.Draw(win)
			img.Clear()
		}
		dialog.ViewPort.Canvas.Clear(constants.ColorBlack)
		NewDrawSystem(dialog.ViewPort.Canvas, layer)
		dialog.ViewPort.Draw(win)
		img.Clear()
		layer++
	}
	for _, dialog := range data.DialogStack {
		if !dialog.NoBorder {
			dialog.BorderVP.Canvas.Clear(color.RGBA{})
			BorderSystem(layer)
			img.Batchers[constants.UIBatch].Draw(dialog.BorderVP.Canvas)
			dialog.BorderVP.Draw(win)
			img.Clear()
		}
		dialog.ViewPort.Canvas.Clear(constants.ColorBlack)
		NewDrawSystem(dialog.ViewPort.Canvas, layer)
		dialog.ViewPort.Draw(win)
		img.Clear()
		layer++
	}
}
