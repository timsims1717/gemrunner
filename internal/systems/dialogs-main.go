package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/ui"
	"github.com/gopxl/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
)

func AddPlayersDialog(win *pixelgl.Window) bool {
	if data.MenuInput.Get("escape").JustPressed() {
		data.MenuInput.Get("escape").Consume()
		if len(data.Players) < 2 {
			return true
		} else {
			var tCntKey, tTxtKey, nTxtKey string
			switch len(data.Players) {
			case 2:
				tCntKey = "selected_p2_cnt"
				tTxtKey = "any_button_p2"
				nTxtKey = "any_button_p3"
			case 3:
				tCntKey = "selected_p3_cnt"
				tTxtKey = "any_button_p3"
				nTxtKey = "any_button_p4"
			case 4:
				tCntKey = "selected_p4_cnt"
				tTxtKey = "any_button_p4"
			}
			data.Players = data.Players[:len(data.Players)-1]
			for _, ele := range ui.Dialogs[constants.DialogAddPlayers].Elements {
				if ele.Key == tCntKey {
					ele.Object.Hidden = true
				} else if ele.Key == nTxtKey {
					ele.Text.Hidden = true
				} else if ele.Key == tTxtKey {
					ele.Text.Hidden = false
				}
			}
		}
	}
	if len(data.Players) < constants.MaxPlayers {
		var playerFound bool
		nextPlayer := data.Player{
			PlayerNum: len(data.Players),
		}
		var nCntKey, nTxtKey, pTxtKey string
		switch len(data.Players) {
		case 1:
			pTxtKey = "any_button_p2"
			nCntKey = "selected_p2_cnt"
			nTxtKey = "any_button_p3"
		case 2:
			pTxtKey = "any_button_p3"
			nCntKey = "selected_p3_cnt"
			nTxtKey = "any_button_p4"
		case 3:
			pTxtKey = "any_button_p4"
			nCntKey = "selected_p4_cnt"
		}
		joysticks := pxginput.GetAllGamepads(win)
		for _, js := range joysticks {
			pressed := pxginput.GetAllJustPressedGamepad(win, js)
			if len(pressed) > 0 {
				playerFound = true
				nextPlayer.Gamepad = js
				break
			}
		}
		if !playerFound && data.MenuInput.Get("space").JustPressed() {
			playerFound = true
			nextPlayer.Keyboard = true
		}
		if playerFound {
			data.Players = append(data.Players, nextPlayer)
			for _, ele := range ui.Dialogs[constants.DialogAddPlayers].Elements {
				if ele.Key == nCntKey {
					ele.Object.Hidden = false
				} else if ele.Key == nTxtKey {
					ele.Text.Hidden = false
				} else if ele.Key == pTxtKey {
					ele.Text.Hidden = true
				}
			}
		}
	}
	return false
}
