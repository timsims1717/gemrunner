package systems

import (
	"fmt"
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
		dialog.BorderVP.Update()
		dialog.ViewPort.Update()
		dialog.BorderObject.Layer = layer
		for _, btn := range dialog.Buttons {
			btn.Object.Layer = layer
		}
		// update this dialog box with input
		if !data.DialogStackOpen {

		}
		layer++
	}
	closeKey := ""
	for i, dialog := range data.DialogStack {
		dialog.BorderVP.Update()
		dialog.ViewPort.Update()
		dialog.BorderObject.Layer = layer
		for _, btn := range dialog.Buttons {
			btn.Object.Layer = layer
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
		CloseDialog(closeKey)
	}
}

func OpenDialog(key string) {
	dialog, ok := data.Dialogs[key]
	if !ok {
		fmt.Printf("Warning: OpenDialog: %s not registered\n", key)
		return
	}
	dialog.Open = true
	data.DialogsOpen = append(data.DialogsOpen, dialog)
}

func OpenDialogInStack(key string) {
	dialog, ok := data.Dialogs[key]
	if !ok {
		fmt.Printf("Warning: OpenDialog: %s not registered\n", key)
		return
	}
	dialog.Open = true
	data.DialogStack = append(data.DialogStack, dialog)
}

func CloseDialog(key string) {
	dialog, ok := data.Dialogs[key]
	if !ok {
		fmt.Printf("Warning: CloseDialog: %s not registered\n", key)
		return
	}
	dialog.Open = false
	index := -1
	stack := false
	for i, d := range data.DialogsOpen {
		if d.Key == key {
			index = i
			break
		}
	}
	for i, d := range data.DialogStack {
		if d.Key == key {
			index = i
			stack = true
			break
		}
	}
	if index == -1 {
		fmt.Printf("Warning: CloseDialog: %s not open\n", key)
		return
	} else {
		if stack {
			if len(data.DialogStack) == 1 {
				data.ClearDialogStack()
			} else {
				data.DialogStack = append(data.DialogStack[:index], data.DialogStack[index+1:]...)
			}
		} else {
			if len(data.DialogsOpen) == 1 {
				data.ClearDialogsOpen()
			} else {
				data.DialogsOpen = append(data.DialogsOpen[:index], data.DialogsOpen[index+1:]...)
			}
		}
	}
}

func DialogDrawSystem(win *pixelgl.Window) {
	layer := 100
	for _, dialog := range data.DialogsOpen {
		dialog.BorderVP.Canvas.Clear(color.RGBA{})
		BorderSystem(layer)
		img.Batchers[constants.UIBatch].Draw(dialog.BorderVP.Canvas)
		dialog.BorderVP.Draw(win)
		img.Clear()
		dialog.ViewPort.Canvas.Clear(constants.ColorBlack)
		DrawSystem(win, layer)
		img.Batchers[constants.UIBatch].Draw(dialog.ViewPort.Canvas)
		dialog.ViewPort.Draw(win)
		img.Clear()
		layer++
	}
	for _, dialog := range data.DialogStack {
		dialog.BorderVP.Canvas.Clear(color.RGBA{})
		BorderSystem(layer)
		img.Batchers[constants.UIBatch].Draw(dialog.BorderVP.Canvas)
		dialog.BorderVP.Draw(win)
		img.Clear()
		dialog.ViewPort.Canvas.Clear(constants.ColorBlack)
		DrawSystem(win, layer)
		img.Batchers[constants.UIBatch].Draw(dialog.ViewPort.Canvas)
		dialog.ViewPort.Draw(win)
		img.Clear()
		layer++
	}
}
