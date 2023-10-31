package load

import (
	"gemrunner/internal/data"
	"github.com/faiface/pixel"
)

func Dialogs() {
	data.OpenPuzzleConstructor = &data.DialogConstructor{
		Key:    "open_puzzle",
		Width:  7,
		Height: 7,
		Buttons: []data.ButtonConstructor{
			{
				SprKey:      "cancel_btn",
				ClickSprKey: "cancel_btn_click",
				HelpText:    "Cancel",
				Position:    pixel.V(32, -32),
			},
		},
	}

	data.NewDialog(data.OpenPuzzleConstructor)

	for key, dialog := range data.Dialogs {
		for _, btn := range dialog.Buttons {
			switch btn.Key {
			case "cancel_btn":
				btn.OnClick = CancelDialog(key)
			}
		}
	}
}
