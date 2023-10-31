package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/pkg/img"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
)

func DialogSystem() {
	if len(data.DialogStack) > 0 {
		data.DialogOpen = true
	} else {
		data.DialogOpen = false
	}
	closeKey := ""
	for i, dialog := range data.DialogStack {
		dialog.BorderVP.Update()
		dialog.ViewPort.Update()
		dialog.BorderObject.Layer = 100 + i
		for _, btn := range dialog.Buttons {
			btn.Object.Layer = 100 + i
		}
		if i == len(data.DialogStack)-1 {
			// update this dialog box with input
			if data.MainInput.Get("escape").JustPressed() {
				data.MainInput.Get("escape").Consume()
				closeKey = dialog.Key
			}
		}
	}
	if closeKey != "" {
		CloseDialogbox(closeKey)
	}
}

func OpenDialogbox(key string) {
	dialog, ok := data.Dialogs[key]
	if !ok {
		fmt.Printf("Warning: OpenDialog: %s not registered\n", key)
		return
	}
	dialog.Open = true
	data.DialogStack = append(data.DialogStack, dialog)
}

func CloseDialogbox(key string) {
	dialog, ok := data.Dialogs[key]
	if !ok {
		fmt.Printf("Warning: CloseDialog: %s not registered\n", key)
		return
	}
	dialog.Open = false
	index := -1
	for i, d := range data.DialogStack {
		if d.Key == key {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Printf("Warning: CloseDialog: %s not open\n", key)
		return
	} else if len(data.DialogStack) == 1 {
		data.ClearDialogStack()
	} else {
		data.DialogStack = append(data.DialogStack[:index], data.DialogStack[index+1:]...)
	}
}

func DialogDrawSystem(win *pixelgl.Window) {
	for i, dialog := range data.DialogStack {
		dialog.BorderVP.Canvas.Clear(color.RGBA{})
		BorderSystem(100 + i)
		img.Batchers[constants.UIBatch].Draw(dialog.BorderVP.Canvas)
		dialog.BorderVP.Draw(win)
		img.Clear()
		dialog.ViewPort.Canvas.Clear(constants.ColorBlack)
		DrawSystem(win, 100+i)
		img.Batchers[constants.UIBatch].Draw(dialog.ViewPort.Canvas)
		dialog.ViewPort.Draw(win)
		img.Clear()
	}
}
