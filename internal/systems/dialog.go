package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/ui"
	"gemrunner/pkg/img"
	"gemrunner/pkg/util"
	"gemrunner/pkg/viewport"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"image/color"
)

func DialogSystem(win *pixelgl.Window) {
	var updated []string
	layer := 100
	for _, dialog := range ui.DialogsOpen {
		dialog.Active = !ui.DialogStackOpen
		updated = append(updated, dialog.Key)
		layer = UpdateDialog(dialog, layer)
	}
	closeKey := ""
	for i, dialog := range ui.DialogStack {
		dialog.Active = i == len(ui.DialogStack)-1
		if dialog.Active {
			closeDlg := false
			// update this dialog box with input
			switch dialog.Key {
			case constants.DialogAddPlayers:
				closeDlg = AddPlayersDialog(win)
			default:
				if data.MenuInput.Get("escape").JustPressed() {
					data.MenuInput.Get("escape").Consume()
					closeDlg = true
				}
			}
			if closeDlg {
				closeKey = dialog.Key
			}
		}
		updated = append(updated, dialog.Key)
		layer = UpdateDialog(dialog, layer)
	}
	if closeKey != "" {
		ui.CloseDialog(closeKey)
	}
	layer += 100
	for key, dialog := range ui.Dialogs {
		if !util.ContainsStr(key, updated) {
			//UpdateDialogLayer99(dialog)
			layer = UpdateDialogLayers(dialog, layer)
		}
	}
}

func UpdateDialog(dialog *ui.Dialog, layer int) int {
	dialog.Loaded = true
	dialog.Layer = layer
	if !dialog.NoBorder {
		dialog.BorderVP.Update()
		dialog.BorderObject.Layer = layer
	}
	dialog.ViewPort.Update()
	nextLayer := UpdateSubElements(dialog.Elements, dialog.ViewPort, layer)
	return nextLayer
}

func UpdateSubElements(elements []*ui.Element, vp *viewport.ViewPort, layer int) int {
	nextLayer := layer + 1
	for _, e := range elements {
		e.Object.Unloaded = !vp.RectInside(e.Object.Rect.Moved(e.Object.Pos))
		switch e.ElementType {
		case ui.SpriteElement, ui.ButtonElement, ui.CheckboxElement:
			e.Object.Layer = layer
		case ui.TextElement:
			e.Text.Obj.Layer = layer
		case ui.InputElement:
			e.Layer = nextLayer
			e.BorderObject.Layer = e.Layer
			e.Text.Obj.Layer = e.Layer
			e.CaretObj.Layer = e.Layer
			if !e.Object.Hidden && !e.Object.Unloaded {
				e.BorderVP.Update()
				e.ViewPort.Update()
			}
			nextLayer++
		case ui.ScrollElement, ui.ContainerElement:
			e.BorderObject.Layer = nextLayer
			e.Object.Layer = nextLayer
			e.Layer = nextLayer
			if !e.Object.Hidden && !e.Object.Unloaded {
				e.BorderVP.Update()
				e.ViewPort.Update()
				nextLayer = UpdateSubElements(e.Elements, e.ViewPort, e.Layer)
			} else {
				nextLayer = UpdateSubElementLayers(e.Elements, e.Layer)
			}
		}
	}
	return nextLayer
}

func UpdateDialogLayers(dialog *ui.Dialog, layer int) int {
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

func UpdateSubElementLayers(elements []*ui.Element, layer int) int {
	nextLayer := layer + 1
	for _, e := range elements {
		switch e.ElementType {
		case ui.SpriteElement, ui.ButtonElement, ui.CheckboxElement:
			e.Object.Layer = layer
		case ui.TextElement:
			e.Text.Obj.Layer = layer
		case ui.InputElement:
			e.Layer = nextLayer
			//e.BorderVP.Update()
			e.BorderObject.Layer = e.Layer
			//e.ViewPort.Update()
			e.Text.Obj.Layer = e.Layer
			e.CaretObj.Layer = e.Layer
			nextLayer++
		case ui.ScrollElement, ui.ContainerElement:
			//e.BorderVP.Update()
			e.BorderObject.Layer = nextLayer
			e.Object.Layer = nextLayer
			//e.ViewPort.Update()
			e.Layer = nextLayer
			nextLayer = UpdateSubElementLayers(e.Elements, e.Layer)
		}
	}
	return nextLayer
}

func UpdateDialogLayer99(dialog *ui.Dialog) {
	dialog.Layer = 99
	if !dialog.NoBorder {
		//dialog.BorderVP.Update()
		dialog.BorderObject.Layer = 99
	}
	//dialog.ViewPort.Update()
	UpdateSubElementLayer99(dialog.Elements)
}

func UpdateSubElementLayer99(elements []*ui.Element) {
	for _, e := range elements {
		switch e.ElementType {
		case ui.SpriteElement, ui.ButtonElement, ui.CheckboxElement:
			e.Object.Layer = 99
		case ui.TextElement:
			e.Text.Obj.Layer = 99
		case ui.InputElement:
			e.Layer = 99
			//e.BorderVP.Update()
			e.BorderObject.Layer = e.Layer
			//e.ViewPort.Update()
			e.Text.Obj.Layer = e.Layer
			e.CaretObj.Layer = e.Layer
		case ui.ScrollElement, ui.ContainerElement:
			//e.BorderVP.Update()
			e.BorderObject.Layer = 99
			e.Object.Layer = 99
			//e.ViewPort.Update()
			e.Layer = 99
			UpdateSubElementLayer99(e.Elements)
		}
	}
}

func DialogDrawSystem(win *pixelgl.Window) {
	for _, dialog := range ui.DialogsOpen {
		DrawDialog(dialog, win)
	}
	for _, dialog := range ui.DialogStack {
		DrawDialog(dialog, win)
	}
}

func DrawDialog(dialog *ui.Dialog, win *pixelgl.Window) {
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
func DrawSubElements(element *ui.Element, vp *viewport.ViewPort) {
	if element == nil {
		return
	}
	if element.Object.Hidden {
		return
	}
	if element.Object.Unloaded {
		return
	}
	switch element.ElementType {
	case ui.InputElement:
		// draw border
		element.BorderVP.Canvas.Clear(color.RGBA{})
		BorderSystem(element.Layer)
		img.Batchers[constants.UIBatch].Draw(element.BorderVP.Canvas)
		element.BorderVP.Draw(vp.Canvas)
		img.Clear()
		// draw input
		element.ViewPort.Canvas.Clear(pixel.RGBA{})
		DrawLayerSystem(element.ViewPort.Canvas, element.Layer)
		element.ViewPort.Draw(vp.Canvas)
		img.Clear()
	case ui.ScrollElement, ui.ContainerElement:
		// draw border
		element.BorderVP.Canvas.Clear(color.RGBA{})
		BorderSystem(element.Layer)
		img.Batchers[constants.UIBatch].Draw(element.BorderVP.Canvas)
		element.BorderVP.Draw(vp.Canvas)
		img.Clear()
		// draw container elements
		element.ViewPort.Canvas.Clear(pixel.RGBA{})
		DrawLayerSystem(element.ViewPort.Canvas, element.Layer)
		img.Clear()
		for _, e := range element.Elements {
			DrawSubElements(e, element.ViewPort)
		}
		element.ViewPort.Draw(vp.Canvas)
		img.Clear()
	}
}
