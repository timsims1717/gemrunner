package load

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/content"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems"
	"gemrunner/internal/ui"
	"gemrunner/pkg/state"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
)

func QuitGame(win *pixelgl.Window) func() {
	return func() {
		win.SetClosed(true)
	}
}

func StartEditor() {
	state.SwitchState(constants.EditorStateKey)
}

func StartPlayNew() {
	if data.SelectedPuzzleIndex > -1 &&
		data.SelectedPuzzleIndex < len(data.PuzzleSetFileList) {
		filename := fmt.Sprintf("%s/%s", constants.PuzzlesDir, data.PuzzleSetFileList[data.SelectedPuzzleIndex].Filename)
		err := systems.OpenPuzzle(filename)
		if err != nil {
			err := systems.OpenPuzzleSet(filename)
			if err != nil {
				fmt.Println("ERROR:", err)
				return
			}
		}
		ui.ClearDialogStack()
		state.SwitchState(constants.PlayStateKey)
	}
}

func ConfirmAddPlayers() {
	localPlay := ui.Dialogs[constants.DialogPlayLocal]
	data.SelectedPuzzleIndex = 0
	for _, e := range localPlay.Elements {
		ele := e
		switch ele.Key {
		//case "play_custom_tab":
		//	for _, ele2 := range ele.Elements {
		//		if ele2.Key == "custom_tab_text_shadow" {
		//			ele2.Text.Hidden = true
		//		}
		//	}
		case "play_main_tab":
			for _, ele2 := range ele.Elements {
				if ele2.Key == "main_tab_text_shadow" {
					ele2.Text.Hidden = true
				}
			}
		case "custom_puzzle_list":
			err := content.LoadPuzzleContent()
			if err != nil {
				fmt.Println("ERROR:", err)
				ele.Elements = []*ui.Element{}
				ui.UpdateScrollBounds(ele)
				return
			}
			totalElements := len(data.PuzzleSetFileList) * 2
			totalPuzzles := len(data.PuzzleSetFileList)
			//xPos := ele.ViewPort.CamPos.X - ele.ViewPort.Rect.W()*0.5 + 4
			for i := 0; i < totalPuzzles; i++ {
				pzIndex := i
				ctIndex := pzIndex * 2
				favIndex := ctIndex + 1
				y := float64(pzIndex)*-34 + 15
				if len(ele.Elements) <= ctIndex {
					entry := ui.ElementConstructor{
						Key:         fmt.Sprintf(puzzleListEntry.Key, pzIndex),
						Width:       puzzleListEntry.Width,
						Height:      puzzleListEntry.Height,
						HelpText:    puzzleListEntry.HelpText,
						Position:    pixel.V(0, y),
						ElementType: puzzleListEntry.ElementType,
					}
					pts := ui.ElementConstructor{
						Key:         puzzleTitleShadowItem.Key,
						Color:       puzzleTitleShadowItem.Color,
						Position:    puzzleTitleShadowItem.Position,
						ElementType: puzzleTitleShadowItem.ElementType,
					}
					entry.SubElements = append(entry.SubElements, pts)
					pti := ui.ElementConstructor{
						Key:         puzzleTitleItem.Key,
						Color:       puzzleTitleItem.Color,
						Position:    puzzleTitleItem.Position,
						ElementType: puzzleTitleItem.ElementType,
					}
					entry.SubElements = append(entry.SubElements, pti)
					pzns := ui.ElementConstructor{
						Key:         puzzleNumberSpr.Key,
						SprKey:      puzzleNumberSpr.SprKey,
						Batch:       puzzleNumberSpr.Batch,
						Position:    puzzleNumberSpr.Position,
						ElementType: puzzleNumberSpr.ElementType,
					}
					entry.SubElements = append(entry.SubElements, pzns)
					plns := ui.ElementConstructor{
						Key:         playerNumberSpr.Key,
						SprKey:      playerNumberSpr.SprKey,
						Batch:       playerNumberSpr.Batch,
						Position:    playerNumberSpr.Position,
						ElementType: playerNumberSpr.ElementType,
					}
					entry.SubElements = append(entry.SubElements, plns)
					ni := ui.ElementConstructor{
						Key:         numberItems.Key,
						Color:       numberItems.Color,
						Position:    numberItems.Position,
						ElementType: numberItems.ElementType,
					}
					entry.SubElements = append(entry.SubElements, ni)
					fsi := ui.ElementConstructor{
						Key:         favoriteItem.Key,
						SprKey:      favoriteItem.SprKey,
						SprKey2:     favoriteItem.SprKey2,
						Batch:       favoriteItem.Batch,
						Position:    pixel.V(favoriteItem.Position.X, y),
						ElementType: favoriteItem.ElementType,
					}
					fav := ui.CreateCheckboxElement(fsi, localPlay, ele.ViewPort)
					ct := ui.CreateContainer(entry, localPlay, ele.ViewPort)
					ele.Elements = append(ele.Elements, ct)
					ele.Elements = append(ele.Elements, fav)
				}
				ct := ele.Elements[ctIndex]
				fav := ele.Elements[favIndex]
				pz := data.PuzzleSetFileList[pzIndex]
				ui.SetChecked(fav, pz.Favorite)
				fav.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, ele.ViewPort, func(hvc *data.HoverClick) {
					if localPlay.Open && localPlay.Active {
						click := hvc.Input.Get("click")
						if hvc.Hover && click.JustPressed() {
							ui.SetChecked(fav, !fav.Checked)
							click.Consume()
						}
					}
				}))
				for _, cEle := range ct.Elements {
					switch cEle.Key {
					case "puzzle_title":
						cEle.Text.SetText(fmt.Sprintf("%s\n %s", pz.Name, pz.Author))
					case "puzzle_title_shadow":
						cEle.Text.SetText(pz.Name)
						cEle.Text.Hidden = pzIndex != data.SelectedPuzzleIndex
					case "number_text":
						cEle.Text.SetText(fmt.Sprintf("%d %d", pz.NumPlayers, pz.NumPuzzles))
					}
				}
				ct.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, ele.ViewPort, func(hvc *data.HoverClick) {
					if localPlay.Open && localPlay.Active {
						click := hvc.Input.Get("click")
						if hvc.Hover && click.JustPressed() {
							if data.SelectedPuzzleIndex != pzIndex {
								data.SelectedPuzzleIndex = pzIndex
								for i2, ie := range ele.Elements {
									ih := i2 / 2
									if ie.ElementType == ui.ContainerElement {
										for _, ce2 := range ie.Elements {
											switch ce2.Key {
											case "puzzle_title_shadow":
												ce2.Text.Hidden = ih != data.SelectedPuzzleIndex
											}
										}
									}
								}
								click.Consume()
							}
						}
					}
				}))
			}
			if len(ele.Elements) > totalElements {
				for i := len(ele.Elements) - 1; i >= totalElements; i-- {
					e2 := ele.Elements[i]
					if e2.ElementType == ui.ContainerElement {
						for _, ctE := range e2.Elements {
							myecs.Manager.DisposeEntity(ctE.Entity)
							myecs.Manager.DisposeEntity(ctE.BorderEntity)
						}
					}
					myecs.Manager.DisposeEntity(e2.Entity)
					myecs.Manager.DisposeEntity(e2.BorderEntity)
				}
				if totalElements > 0 {
					ele.Elements = ele.Elements[:totalElements]
				} else {
					ele.Elements = []*ui.Element{}
				}
			}
			ui.UpdateScrollBounds(ele)
		}
	}
	ui.OpenDialogInStack(constants.DialogPlayLocal)
}

func OpenAddPlayers() {
	for _, ele := range ui.Dialogs[constants.DialogAddPlayers].Elements {
		if ele.ElementType == ui.ContainerElement {
			switch ele.Key {
			case "selected_p1_cnt":
				ele.Object.Hidden = false
			default:
				ele.Object.Hidden = true
			}
		} else if ele.ElementType == ui.TextElement {
			if ele.Key == "any_button_p2" {
				ele.Text.Hidden = false
			} else {
				ele.Text.Hidden = true
			}
		}
	}
	p1 := data.Player{}
	if data.MenuInputUsed == pxginput.Gamepad && data.MainJoystick > 0 {
		p1.Gamepad = pixelgl.Joystick(data.MainJoystick)
	} else {
		p1.Keyboard = true
	}
	data.Players = []data.Player{p1}
	ui.OpenDialogInStack(constants.DialogAddPlayers)
}
