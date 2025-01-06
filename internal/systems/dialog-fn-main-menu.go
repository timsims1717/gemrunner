package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/content"
	"gemrunner/internal/data"
	"gemrunner/internal/load"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/state"
	"gemrunner/pkg/util"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
	"strings"
)

func QuitGame(win *pixelgl.Window) func() {
	return func() {
		win.SetClosed(true)
	}
}

func StartEditor() {
	state.SwitchState(constants.EditorStateKey)
}

func StartCustomContinue() {
	if data.CustomPuzzleListLoaded && data.SelectedPuzzleIndex > -1 &&
		data.SelectedPuzzleIndex < len(data.PuzzleSetFileList) {
		err := OpenPuzzleSet(data.PuzzleSetFileList[data.SelectedPuzzleIndex].Filename)
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}
		svg := data.PuzzleSetFileList[data.SelectedPuzzleIndex].Filename
		svg = strings.Replace(svg, constants.PuzzleExt, constants.SaveExt, -1)
		err = content.LoadSaveGame(svg)
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}
		ui.ClearDialogStack()
		state.SwitchState(constants.PlayStateKey)
	}
}

func StartCustomNew() {
	if data.CustomPuzzleListLoaded && data.SelectedPuzzleIndex > -1 &&
		data.SelectedPuzzleIndex < len(data.PuzzleSetFileList) {
		err := OpenPuzzleSet(data.PuzzleSetFileList[data.SelectedPuzzleIndex].Filename)
		if err != nil {
			fmt.Println("ERROR:", err)
			return
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
		case "play_custom_tab":
			ele.Get("custom_tab_text_shadow").Text.Show()
			ele.Border.Style = ui.ThinBorderWhite
		case "play_main_tab":
			ele.Get("main_tab_text_shadow").Text.Hide()
			ele.Border.Style = ui.ThinBorderBlue
		case "main_tab_display":
		case "custom_tab_display":
			for _, e1 := range ele.Elements {
				ele1 := e1
				switch ele1.Key {
				case "custom_puzzle_list":
					ele1.Object.Hidden = true
					data.CustomPuzzleListLoaded = false
					go func() {
						content.LoadLocalPuzzleList()
						content.OrganizeLocalPuzzles(content.DefaultFilters)
						PopulateCustomPuzzleList(ele1)
						data.CustomPuzzleListLoaded = true
						for _, e2 := range ele.Elements {
							switch e2.Key {
							case "custom_puzzle_loading":
								e2.Object.Hidden = true
							case "custom_puzzle_list",
								"custom_puzzle_list_scroll_bar",
								"custom_puzzle_list_scroll_up",
								"custom_puzzle_list_scroll_down":
								e2.Object.Hidden = false
							}
						}
					}()
				case "custom_puzzle_loading":
					ele1.Object.Hidden = false
				case "custom_puzzle_list_scroll_bar",
					"custom_puzzle_list_scroll_up",
					"custom_puzzle_list_scroll_down":
					ele1.Object.Hidden = true
				}
			}
		}
	}
	ui.OpenDialogInStack(constants.DialogPlayLocal)
}

func PreLoadCustomPuzzleList() {
	customTab := ui.Dialogs[constants.DialogPlayLocal].Get("custom_tab_display")
	pzlList := customTab.Get("custom_puzzle_list")
	go func() {
		content.LoadLocalPuzzleList()
		content.OrganizeLocalPuzzles(content.DefaultFilters)
		PopulateCustomPuzzleList(pzlList)
	}()
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
			ele.Text.Obj.Hidden = ele.Key != "any_button_p2"
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

func PopulateCustomPuzzleList(ele *ui.Element) {
	localPlay := ui.Dialogs[constants.DialogPlayLocal]
	ele.Object.Hidden = false
	totalElements := len(data.PuzzleSetSortedList) * 2
	totalPuzzles := len(data.PuzzleSetSortedList)
	for i := 0; i < totalPuzzles; i++ {
		pzIndex := i
		ctIndex := pzIndex * 2
		favIndex := ctIndex + 1
		y := float64(pzIndex)*-34 + 15
		if len(ele.Elements) <= ctIndex {
			entry := ui.ElementConstructor{
				Key:         fmt.Sprintf(load.PuzzleListEntry.Key, pzIndex),
				Width:       load.PuzzleListEntry.Width,
				Height:      load.PuzzleListEntry.Height,
				HelpText:    load.PuzzleListEntry.HelpText,
				Position:    pixel.V(0, y),
				ElementType: load.PuzzleListEntry.ElementType,
			}
			pts := ui.ElementConstructor{
				Key:         load.PuzzleTitleShadowItem.Key,
				Color:       load.PuzzleTitleShadowItem.Color,
				Position:    load.PuzzleTitleShadowItem.Position,
				ElementType: load.PuzzleTitleShadowItem.ElementType,
				Anchor:      load.PuzzleTitleShadowItem.Anchor,
			}
			entry.SubElements = append(entry.SubElements, pts)
			pti := ui.ElementConstructor{
				Key:         load.PuzzleTitleItem.Key,
				Color:       load.PuzzleTitleItem.Color,
				Position:    load.PuzzleTitleItem.Position,
				ElementType: load.PuzzleTitleItem.ElementType,
				Anchor:      load.PuzzleTitleItem.Anchor,
			}
			entry.SubElements = append(entry.SubElements, pti)
			pn := ui.ElementConstructor{
				Key:         load.PuzzleAuthorItem.Key,
				Color:       load.PuzzleAuthorItem.Color,
				Position:    load.PuzzleAuthorItem.Position,
				ElementType: load.PuzzleAuthorItem.ElementType,
				Anchor:      load.PuzzleAuthorItem.Anchor,
			}
			entry.SubElements = append(entry.SubElements, pn)
			pzns := ui.ElementConstructor{
				Key:         load.PuzzleNumberSpr.Key,
				SprKey:      load.PuzzleNumberSpr.SprKey,
				Batch:       load.PuzzleNumberSpr.Batch,
				Position:    load.PuzzleNumberSpr.Position,
				ElementType: load.PuzzleNumberSpr.ElementType,
			}
			entry.SubElements = append(entry.SubElements, pzns)
			plns := ui.ElementConstructor{
				Key:         load.PlayerNumberSpr.Key,
				SprKey:      load.PlayerNumberSpr.SprKey,
				Batch:       load.PlayerNumberSpr.Batch,
				Position:    load.PlayerNumberSpr.Position,
				ElementType: load.PlayerNumberSpr.ElementType,
			}
			entry.SubElements = append(entry.SubElements, plns)
			ni := ui.ElementConstructor{
				Key:         load.NumberItems.Key,
				Color:       load.NumberItems.Color,
				Position:    load.NumberItems.Position,
				ElementType: load.NumberItems.ElementType,
				Anchor:      load.NumberItems.Anchor,
			}
			entry.SubElements = append(entry.SubElements, ni)
			fsi := ui.ElementConstructor{
				Key:         load.FavoriteItem.Key,
				SprKey:      load.FavoriteItem.SprKey,
				SprKey2:     load.FavoriteItem.SprKey2,
				Batch:       load.FavoriteItem.Batch,
				Position:    pixel.V(load.FavoriteItem.Position.X, y),
				ElementType: load.FavoriteItem.ElementType,
			}
			fav := ui.CreateCheckboxElement(fsi, localPlay, ele.ViewPort)
			ct := ui.CreateContainer(entry, localPlay, ele.ViewPort)
			ele.Elements = append(ele.Elements, ct)
			ele.Elements = append(ele.Elements, fav)
		}
		ct := ele.Elements[ctIndex]
		fav := ele.Elements[favIndex]
		pz := data.PuzzleSetSortedList[pzIndex]
		ui.SetChecked(fav, pz.Favorite)
		fav.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, ele.ViewPort, func(hvc *data.HoverClick) {
			if localPlay.Open && localPlay.Active {
				click := hvc.Input.Get("click")
				if hvc.Hover && click.JustPressed() {
					ui.SetChecked(fav, !fav.Checked)
					if fav.Checked {
						data.FavoritesList = append(data.FavoritesList, pz.UUID.String())
					} else {
						data.FavoritesList = util.RemoveStrUO(pz.UUID.String(), data.FavoritesList)
					}
					go content.SaveFavoritesFile()
					click.Consume()
				}
			}
		}))
		for _, cEle := range ct.Elements {
			switch cEle.Key {
			case "puzzle_title":
				cEle.Text.SetText(pz.Name)
			case "puzzle_author":
				cEle.Text.SetText(pz.Author)
			case "puzzle_title_shadow":
				cEle.Text.SetText(pz.Name)
				cEle.Text.Obj.Hidden = pzIndex != data.SelectedPuzzleIndex
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
										ce2.Text.Obj.Hidden = ih != data.SelectedPuzzleIndex
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
	ui.MoveToScrollTop(ele)
}
