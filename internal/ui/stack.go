package ui

import "fmt"

func ClearDialogStack() {
	for _, dialog := range DialogStack {
		dialog.Open = false
		dialog.Active = false
	}
	DialogStack = []*Dialog{}
}

func ClearDialogsOpen() {
	for _, dialog := range DialogsOpen {
		dialog.Open = false
		dialog.Active = false
	}
	DialogsOpen = []*Dialog{}
}

func OpenDialog(key string) {
	dialog, ok := Dialogs[key]
	if !ok {
		fmt.Printf("Warning: OpenDialog: %s not registered\n", key)
		return
	}
	dialog.Open = true
	DialogsOpen = append(DialogsOpen, dialog)
	if dialog.OnOpen != nil {
		dialog.OnOpen()
	}
}

func OpenDialogInStack(key string) {
	dialog, ok := Dialogs[key]
	if !ok {
		fmt.Printf("Warning: OpenDialog: %s not registered\n", key)
		return
	}
	dialog.Open = true
	DialogStack = append(DialogStack, dialog)
	if dialog.OnOpen != nil {
		dialog.OnOpen()
	}
}

func SetCloseSpcFn(key string, fn func()) {
	dialog, ok := Dialogs[key]
	if !ok {
		fmt.Printf("Warning: SetCloseSpcFn: %s not registered\n", key)
		return
	}
	dialog.OnCloseSpc = fn
}

func SetOnClick(dlgKey, btnKey string, fn func()) {
	dialog, ok := Dialogs[dlgKey]
	if !ok {
		fmt.Printf("Warning: SetOnClick: Dialog %s not registered\n", dlgKey)
		return
	}
	for _, ele := range dialog.Elements {
		if ele.ElementType == ButtonElement && ele.Key == btnKey {
			ele.OnClick = fn
			return
		}
	}
	fmt.Printf("Warning: SetOnClick: Button %s not registered in Dialog %s\n", btnKey, dlgKey)
}

func SetTempOnClick(dlgKey, btnKey string, fn func()) {
	dialog, ok := Dialogs[dlgKey]
	if !ok {
		fmt.Printf("Warning: SetTempOnClick: Dialog %s not registered\n", dlgKey)
		return
	}
	for _, ele := range dialog.Elements {
		if ele.ElementType == ButtonElement && ele.Key == btnKey {
			oldFn := ele.OnClick
			ele.OnClick = func() {
				fn()
				ele.OnClick = oldFn
				oldFn()
			}
			return
		}
	}
	fmt.Printf("Warning: SetTempOnClick: Button %s not registered in Dialog %s\n", btnKey, dlgKey)
}

func CloseDialog(key string) {
	dialog, ok := Dialogs[key]
	if !ok {
		fmt.Printf("Warning: CloseDialog: %s not registered\n", key)
		return
	}
	dialog.Open = false
	dialog.Active = false
	index := -1
	stack := false
	for i, d := range DialogsOpen {
		if d.Key == key {
			index = i
			break
		}
	}
	for i, d := range DialogStack {
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
			if len(DialogStack) == 1 {
				ClearDialogStack()
			} else {
				DialogStack = append(DialogStack[:index], DialogStack[index+1:]...)
			}
		} else {
			if len(DialogsOpen) == 1 {
				ClearDialogsOpen()
			} else {
				DialogsOpen = append(DialogsOpen[:index], DialogsOpen[index+1:]...)
			}
		}
		if dialog.OnClose != nil {
			dialog.OnClose()
		}
		if dialog.OnCloseSpc != nil {
			dialog.OnCloseSpc()
			dialog.OnCloseSpc = nil
		}
	}
}
