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

func PopulateCustomPuzzleList(pzlList *ui.Element) {
	localPlay := ui.Dialogs[constants.DialogPlayLocal]
	pzlList.Object.Hidden = false
	totalElements := len(data.PuzzleSetSortedList) * 2
	totalPuzzles := len(data.PuzzleSetSortedList)
	includeTop := 0.
	includeBot := 0.
	for i := 0; i < totalPuzzles; i++ {
		pzIndex := i
		ctIndex := pzIndex * 2
		favIndex := ctIndex + 1
		y := float64(pzIndex)*-34 + 15
		if len(pzlList.Elements) <= ctIndex {
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
				Key:         fmt.Sprintf(load.FavoriteItem.Key, pzIndex),
				SprKey:      load.FavoriteItem.SprKey,
				SprKey2:     load.FavoriteItem.SprKey2,
				Batch:       load.FavoriteItem.Batch,
				Position:    pixel.V(load.FavoriteItem.Position.X, y),
				ElementType: load.FavoriteItem.ElementType,
			}
			entry.Left = fsi.Key
			if pzIndex == 0 {
				entry.Up = "play_custom_tab"
				fsi.Up = "play_custom_tab"
			} else {
				entry.Up = fmt.Sprintf(load.PuzzleListEntry.Key, pzIndex-1)
				fsi.Up = fmt.Sprintf(load.FavoriteItem.Key, pzIndex-1)
			}
			if pzIndex == len(data.PuzzleSetSortedList)-1 {
				entry.Down = "cancel"
				fsi.Down = "cancel"
			} else {
				entry.Down = fmt.Sprintf(load.PuzzleListEntry.Key, pzIndex+1)
				fsi.Down = fmt.Sprintf(load.FavoriteItem.Key, pzIndex+1)
			}
			fsi.Right = entry.Key
			fav := ui.CreateCheckboxElement(fsi, localPlay, pzlList.ViewPort)
			ct := ui.CreateContainer(entry, localPlay, pzlList.ViewPort)
			pzlList.Elements = append(pzlList.Elements, ct)
			pzlList.Elements = append(pzlList.Elements, fav)
		}
		ct := pzlList.Elements[ctIndex]
		fav := pzlList.Elements[favIndex]
		pz := data.PuzzleSetSortedList[pzIndex]
		ui.SetChecked(fav, pz.Favorite)
		fav.OnClick = func() {
			ui.SetChecked(fav, !fav.Checked)
			if fav.Checked {
				data.FavoritesList = append(data.FavoritesList, pz.UUID.String())
			} else {
				data.FavoritesList = util.RemoveStrUO(pz.UUID.String(), data.FavoritesList)
			}
			go content.SaveFavoritesFile()
		}
		fav.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, pzlList.ViewPort, func(hvc *data.HoverClick) {
			if localPlay.Open && localPlay.Active {
				click := hvc.Input.Get("click")
				if hvc.Hover && click.JustPressed() {
					fav.OnClick()
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
		ct.OnClick = func() {
			// todo: if there is a game to continue, switch to the continue button instead
			localPlay.SetFocus(localPlay.Get("custom_new_game"), true)
		}
		ct.OnFocus = func(focused bool) {
			if focused {
				if data.SelectedPuzzleIndex != pzIndex {
					data.SelectedPuzzleIndex = pzIndex
					for i2, ie := range pzlList.Elements {
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
				}
				ui.MoveScrollToInclude(pzlList, y+ct.Object.Rect.H()*0.5, y-ct.Object.Rect.H()*0.5)
			}
		}
		ct.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, pzlList.ViewPort, func(hvc *data.HoverClick) {
			if localPlay.Open && localPlay.Active {
				click := hvc.Input.Get("click")
				if hvc.Hover && click.JustPressed() {
					ct.OnFocus(true)
					click.Consume()
				}
			}
		}))
		if pzIndex == data.SelectedPuzzleIndex {
			includeTop = y + ct.Object.Rect.H()*0.5
			includeBot = y - ct.Object.Rect.H()*0.5
		}
	}
	if len(pzlList.Elements) > totalElements {
		for i := len(pzlList.Elements) - 1; i >= totalElements; i-- {
			e2 := pzlList.Elements[i]
			if e2.ElementType == ui.ContainerElement {
				for _, ctE := range e2.Elements {
					myecs.Manager.DisposeEntity(ctE.Entity)
				}
			}
			myecs.Manager.DisposeEntity(e2.Entity)
		}
		if totalElements > 0 {
			pzlList.Elements = pzlList.Elements[:totalElements]
		} else {
			pzlList.Elements = []*ui.Element{}
		}
	}

	ui.UpdateScrollBounds(pzlList)
	ui.MoveScrollToInclude(pzlList, includeTop, includeBot)
}
