package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/pkg/img"
	"gemrunner/pkg/util"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"image/color"
)

func DialogSystem() {
	if len(data.DialogStack) > 0 {
		data.DialogStackOpen = true
	} else {
		data.DialogStackOpen = false
	}
	var updated []string
	layer := 100
	for _, dialog := range data.DialogsOpen {
		dialog.Active = !data.DialogStackOpen
		updated = append(updated, dialog.Key)
		layer = UpdateDialog(dialog, layer)
	}
	closeKey := ""
	for i, dialog := range data.DialogStack {
		dialog.Active = i == len(data.DialogStack)-1
		if dialog.Active {
			// update this dialog box with input
			if data.MenuInput.Get("escape").JustPressed() {
				data.MenuInput.Get("escape").Consume()
				closeKey = dialog.Key
				continue
			}
		}
		updated = append(updated, dialog.Key)
		layer = UpdateDialog(dialog, layer)
	}
	if closeKey != "" {
		data.CloseDialog(closeKey)
	}
	for key, dialog := range data.Dialogs {
		if !util.ContainsStr(key, updated) {
			UpdateDialog(dialog, layer+100)
		}
	}
}

func UpdateDialog(dialog *data.Dialog, layer int) int {
	dialog.Layer = layer
	if !dialog.NoBorder {
		dialog.BorderVP.Update()
		dialog.BorderObject.Layer = layer
	}
	dialog.ViewPort.Update()
	layer = UpdateSubElementLayers(dialog.Elements, layer)
	layer++
	return layer
}

func UpdateSubElementLayers(elements []interface{}, layer int) int {
	for _, e := range elements {
		if spr, okS := e.(*data.SprElement); okS {
			spr.Object.Layer = layer
		} else if btn, okB := e.(*data.Button); okB {
			btn.Object.Layer = layer
		} else if in, okI := e.(*data.Input); okI {
			in.Layer = layer + 1
			in.BorderVP.Update()
			in.BorderObject.Layer = in.Layer
			in.ViewPort.Update()
			in.Text.Obj.Layer = in.Layer
			in.CaretObj.Layer = in.Layer
		} else if txt, okT := e.(*data.Text); okT {
			txt.Text.Obj.Layer = layer
		} else if scr, okScr := e.(*data.Scroll); okScr {
			//scr.BarObject.Layer = layer
			//scr.UpObject.Layer = layer
			//scr.DwnObject.Layer = layer
			scr.BorderVP.Update()
			scr.BorderObject.Layer = layer + 1
			scr.ViewPort.Update()
			scr.Layer = layer + 1
			return UpdateSubElementLayers(scr.Elements, scr.Layer)
		}
	}
	return layer
}

func DialogDrawSystem(win *pixelgl.Window) {
	for _, dialog := range data.DialogsOpen {
		DrawDialog(dialog, win)
	}
	for _, dialog := range data.DialogStack {
		DrawDialog(dialog, win)
	}
}

func DrawDialog(dialog *data.Dialog, win *pixelgl.Window) {
	if !dialog.NoBorder {
		dialog.BorderVP.Canvas.Clear(color.RGBA{})
		BorderSystem(dialog.Layer)
		img.Batchers[constants.UIBatch].Draw(dialog.BorderVP.Canvas)
		dialog.BorderVP.Draw(win)
		img.Clear()
	}
	dialog.ViewPort.Canvas.Clear(constants.ColorBlack)
	DrawLayerSystem(dialog.ViewPort.Canvas, dialog.Layer)
	img.Clear()
	for _, e := range dialog.Elements {
		if in, okI := e.(*data.Input); okI {
			in.BorderVP.Canvas.Clear(color.RGBA{})
			BorderSystem(in.Layer)
			img.Batchers[constants.UIBatch].Draw(in.BorderVP.Canvas)
			in.BorderVP.Draw(dialog.ViewPort.Canvas)
			img.Clear()
			DrawSubElement(e)
			in.ViewPort.Draw(dialog.ViewPort.Canvas)
		} else if scr, okScr := e.(*data.Scroll); okScr {
			scr.BorderVP.Canvas.Clear(color.RGBA{})
			BorderSystem(scr.Layer)
			img.Batchers[constants.UIBatch].Draw(scr.BorderVP.Canvas)
			scr.BorderVP.Draw(dialog.ViewPort.Canvas)
			img.Clear()
			DrawSubElement(e)
			scr.ViewPort.Draw(dialog.ViewPort.Canvas)
		}
	}
	dialog.ViewPort.Draw(win)
	img.Clear()
}

func DrawSubElement(element interface{}) {
	if in, okI := element.(*data.Input); okI {
		in.ViewPort.Canvas.Clear(pixel.RGBA{})
		DrawLayerSystem(in.ViewPort.Canvas, in.Layer)
		img.Clear()
	} else if scr, okScr := element.(*data.Scroll); okScr {
		scr.ViewPort.Canvas.Clear(pixel.RGBA{})
		DrawLayerSystem(scr.ViewPort.Canvas, scr.Layer)
		img.Clear()
		for _, e := range scr.Elements {
			if _, okScr2 := e.(*data.Scroll); okScr2 {
				DrawSubElement(e)
			}
		}
	}
}
