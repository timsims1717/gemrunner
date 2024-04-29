package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/pkg/img"
	"gemrunner/pkg/util"
	"gemrunner/pkg/viewport"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"image/color"
)

func DialogSystem() {
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
	layer += 100
	for key, dialog := range data.Dialogs {
		if !util.ContainsStr(key, updated) {
			layer = UpdateDialogLayers(dialog, layer)
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
	nextLayer := UpdateSubElements(dialog.Elements, layer)
	return nextLayer
}

func UpdateSubElements(elements []interface{}, layer int) int {
	nextLayer := layer + 1
	for _, e := range elements {
		if spr, okS := e.(*data.SprElement); okS {
			spr.Object.Layer = layer
		} else if btn, okB := e.(*data.Button); okB {
			btn.Object.Layer = layer
		} else if xbx, okX := e.(*data.Checkbox); okX {
			xbx.Object.Layer = layer
		} else if in, okI := e.(*data.Input); okI {
			in.Layer = nextLayer
			in.BorderVP.Update()
			in.BorderObject.Layer = in.Layer
			in.ViewPort.Update()
			in.Text.Obj.Layer = in.Layer
			in.CaretObj.Layer = in.Layer
			nextLayer++
		} else if txt, okT := e.(*data.Text); okT {
			txt.Text.Obj.Layer = layer
		} else if scr, okScr := e.(*data.Scroll); okScr {
			scr.BorderVP.Update()
			scr.BorderObject.Layer = nextLayer
			scr.Object.Layer = nextLayer
			scr.ViewPort.Update()
			scr.Layer = nextLayer
			nextLayer = UpdateSubElements(scr.Elements, scr.Layer)
		} else if ct, okCt := e.(*data.Container); okCt {
			ct.BorderVP.Update()
			ct.BorderObject.Layer = nextLayer
			ct.Object.Layer = nextLayer
			ct.ViewPort.Update()
			ct.Layer = nextLayer
			nextLayer = UpdateSubElements(ct.Elements, ct.Layer)
		}
	}
	return nextLayer
}

func UpdateDialogLayers(dialog *data.Dialog, layer int) int {
	dialog.Layer = layer
	if !dialog.NoBorder {
		//dialog.BorderVP.Update()
		dialog.BorderObject.Layer = layer
	}
	//dialog.ViewPort.Update()
	layer = UpdateSubElementLayers(dialog.Elements, layer)
	layer++
	return layer
}

func UpdateSubElementLayers(elements []interface{}, layer int) int {
	nextLayer := layer + 1
	for _, e := range elements {
		if spr, okS := e.(*data.SprElement); okS {
			spr.Object.Layer = layer
		} else if btn, okB := e.(*data.Button); okB {
			btn.Object.Layer = layer
		} else if xbx, okX := e.(*data.Checkbox); okX {
			xbx.Object.Layer = layer
		} else if in, okI := e.(*data.Input); okI {
			in.Layer = nextLayer
			//in.BorderVP.Update()
			in.BorderObject.Layer = in.Layer
			//in.ViewPort.Update()
			in.Text.Obj.Layer = in.Layer
			in.CaretObj.Layer = in.Layer
			nextLayer++
		} else if txt, okT := e.(*data.Text); okT {
			txt.Text.Obj.Layer = layer
		} else if scr, okScr := e.(*data.Scroll); okScr {
			//scr.BorderVP.Update()
			scr.BorderObject.Layer = nextLayer
			scr.Object.Layer = nextLayer
			//scr.ViewPort.Update()
			scr.Layer = nextLayer
			nextLayer = UpdateSubElementLayers(scr.Elements, scr.Layer)
		} else if ct, okCt := e.(*data.Container); okCt {
			//ct.BorderVP.Update()
			ct.BorderObject.Layer = nextLayer
			ct.Object.Layer = nextLayer
			//ct.ViewPort.Update()
			ct.Layer = nextLayer
			nextLayer = UpdateSubElementLayers(ct.Elements, ct.Layer)
		}
	}
	return nextLayer
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
	// draw elements w/no sub elements
	dialog.ViewPort.Canvas.Clear(constants.ColorBlack)
	DrawLayerSystem(dialog.ViewPort.Canvas, dialog.Layer)
	img.Clear()
	// draw elements w/sub elements
	for _, e := range dialog.Elements {
		DrawSubElements(e, dialog.ViewPort)
	}
	dialog.ViewPort.Draw(win)
	img.Clear()
}

// DrawSubElements draws the border and sub elements of ui elements
// with sub elements.
func DrawSubElements(element interface{}, vp *viewport.ViewPort) {
	if in, okI := element.(*data.Input); okI {
		// draw border
		in.BorderVP.Canvas.Clear(color.RGBA{})
		BorderSystem(in.Layer)
		img.Batchers[constants.UIBatch].Draw(in.BorderVP.Canvas)
		in.BorderVP.Draw(vp.Canvas)
		img.Clear()
		// draw input
		in.ViewPort.Canvas.Clear(pixel.RGBA{})
		DrawLayerSystem(in.ViewPort.Canvas, in.Layer)
		in.ViewPort.Draw(vp.Canvas)
		img.Clear()
	} else if scr, okScr := element.(*data.Scroll); okScr {
		// draw border
		scr.BorderVP.Canvas.Clear(color.RGBA{})
		BorderSystem(scr.Layer)
		img.Batchers[constants.UIBatch].Draw(scr.BorderVP.Canvas)
		scr.BorderVP.Draw(vp.Canvas)
		img.Clear()
		// draw scroll elements
		scr.ViewPort.Canvas.Clear(pixel.RGBA{})
		DrawLayerSystem(scr.ViewPort.Canvas, scr.Layer)
		img.Clear()
		for _, e := range scr.Elements {
			DrawSubElements(e, scr.ViewPort)
		}
		scr.ViewPort.Draw(vp.Canvas)
		img.Clear()
	} else if ct, okCt := element.(*data.Container); okCt {
		// draw border
		ct.BorderVP.Canvas.Clear(color.RGBA{})
		BorderSystem(ct.Layer)
		img.Batchers[constants.UIBatch].Draw(ct.BorderVP.Canvas)
		ct.BorderVP.Draw(vp.Canvas)
		img.Clear()
		// draw container elements
		ct.ViewPort.Canvas.Clear(pixel.RGBA{})
		DrawLayerSystem(ct.ViewPort.Canvas, ct.Layer)
		img.Clear()
		for _, e := range ct.Elements {
			DrawSubElements(e, ct.ViewPort)
		}
		ct.ViewPort.Draw(vp.Canvas)
		img.Clear()
	}
}
