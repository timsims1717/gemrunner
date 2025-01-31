package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/world"
	"strconv"
)

// editor panel mode buttons

func EditorMode(mode data.EditorMode, btn *ui.Element, dialog *ui.Dialog) func() {
	return func() {
		data.Editor.SelectVis = false
		if data.Editor.Mode != mode {
			data.Editor.LastMode = data.Editor.Mode
			data.Editor.ModeChanged = true
			data.Editor.LastCoords = world.Coords{X: -1, Y: -1}
		} else if mode == data.ModePalette {
			data.Editor.LastMode = data.ModeBrush
			data.Editor.ModeChanged = true
			data.Editor.LastCoords = world.Coords{X: -1, Y: -1}
		}
		data.Editor.Mode = mode
		for _, e := range dialog.Elements {
			if e.ElementType == ui.ButtonElement {
				e.Entity.AddComponent(myecs.Drawable, e.Sprite)
			}
		}
		btn.Entity.AddComponent(myecs.Drawable, btn.Sprite2)
		data.CurrPuzzleSet.CurrPuzzle.Update = true
	}
}

// wrench dialogs

func OnOpenCrackTileOptions() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			ui.CloseDialog(constants.DialogCrackedTiles)
			return
		}
		firstTile := data.CurrPuzzleSet.CurrPuzzle.WrenchTiles[0]
		crackDialog := ui.Dialogs[constants.DialogCrackedTiles]
		for _, ele := range crackDialog.Elements {
			switch ele.Key {
			case "cracked_tile_regenerate_check":
				ui.SetChecked(ele, firstTile.Metadata.Regenerate)
			case "cracked_tile_show_check":
				ui.SetChecked(ele, firstTile.Metadata.ShowCrack)
			case "cracked_tile_enemy_check":
				ui.SetChecked(ele, firstTile.Metadata.EnemyCrack)
			case "cracked_tile_title":
				if firstTile.Block == data.BlockCracked {
					ele.Text.SetText("Cracked Turf")
				} else {
					ele.Text.SetText("Cracked Ladder")
				}
			}
		}
	}
}

func ConfirmCrackTileOptions() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			ui.CloseDialog(constants.DialogCrackedTiles)
			return
		}
		crackDialog := ui.Dialogs[constants.DialogCrackedTiles]
		var regen, show, enemy bool
		for _, ele := range crackDialog.Elements {
			switch ele.Key {
			case "cracked_tile_regenerate_check":
				regen = ele.Checked
			case "cracked_tile_show_check":
				show = ele.Checked
			case "cracked_tile_enemy_check":
				enemy = ele.Checked
			}
		}
		for _, tile := range data.CurrPuzzleSet.CurrPuzzle.WrenchTiles {
			tile.Metadata.Regenerate = regen
			tile.Metadata.ShowCrack = show
			tile.Metadata.EnemyCrack = enemy
			tile.Metadata.Changed = true
		}
		ui.CloseDialog(constants.DialogCrackedTiles)
		data.CurrPuzzleSet.CurrPuzzle.Update = true
		data.CurrPuzzleSet.CurrPuzzle.Changed = true
	}
}

// Bomb Options

func OnOpenBombOptions() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			ui.CloseDialog(constants.DialogBomb)
			return
		}
		firstTile := data.CurrPuzzleSet.CurrPuzzle.WrenchTiles[0]
		for _, ele := range ui.Dialogs[constants.DialogBomb].Elements {
			switch ele.Key {
			case "bomb_cross_check":
				ui.SetChecked(ele, firstTile.Metadata.BombCross)
			case "bomb_regenerate_check":
				ui.SetChecked(ele, firstTile.Metadata.Regenerate)
			case "bomb_options_title":
				if firstTile.Block == data.BlockBomb {
					ele.Text.SetText("Bomb Item Options")
				} else {
					ele.Text.SetText("Lit Bomb Options")
				}
			case "bomb_regenerate_delay_input":
				ele.InputType = ui.Numeric
				ui.ChangeText(ele, fmt.Sprintf("%d", firstTile.Metadata.RegenDelay))
			}
		}
	}
}

func ConfirmBombOptions() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			ui.CloseDialog(constants.DialogBomb)
			return
		}
		var regen, cross bool
		var delay int
		for _, ele := range ui.Dialogs[constants.DialogBomb].Elements {
			switch ele.Key {
			case "bomb_cross_check":
				cross = ele.Checked
			case "bomb_regenerate_check":
				regen = ele.Checked
			case "bomb_regenerate_delay_input":
				di, err := strconv.Atoi(ele.Text.Raw)
				if err != nil {
					fmt.Println("WARNING: regen delay not an int:", err)
					di = 0
				}
				delay = di
			}
		}
		for _, tile := range data.CurrPuzzleSet.CurrPuzzle.WrenchTiles {
			tile.Metadata.Regenerate = regen
			tile.Metadata.BombCross = cross
			tile.Metadata.RegenDelay = delay
			tile.Metadata.Changed = true
		}
		ui.CloseDialog(constants.DialogBomb)
		data.CurrPuzzleSet.CurrPuzzle.Update = true
		data.CurrPuzzleSet.CurrPuzzle.Changed = true
	}
}

// Item Options

func customizeItemOptions() {
	itemDlg := ui.Dialogs[constants.DialogItemOptions]
	itemDlg.OnOpen = OnOpenItemOptions
	for _, e := range itemDlg.Elements {
		ele := e
		switch ele.Key {
		case "confirm":
			ele.OnClick = ConfirmItemOptions
		case "item_regenerate_delay_minus":
			ele.OnClick = func() {
				ChangeNumberInput(itemDlg.Get("item_regenerate_delay_input"), -1)
			}
		case "item_regenerate_delay_plus":
			ele.OnClick = func() {
				ChangeNumberInput(itemDlg.Get("item_regenerate_delay_input"), 1)
			}
		case "item_timer_minus":
			ele.OnClick = func() {
				ChangeNumberInput(itemDlg.Get("item_timer_input"), -1)
			}
		case "item_timer_plus":
			ele.OnClick = func() {
				ChangeNumberInput(itemDlg.Get("item_timer_input"), 1)
			}
		case "cancel":
			ele.OnClick = CloseDialog(constants.DialogItemOptions)
		}
	}
}

func OnOpenItemOptions() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			ui.CloseDialog(constants.DialogItemOptions)
			return
		}
		firstTile := data.CurrPuzzleSet.CurrPuzzle.WrenchTiles[0]
		for _, ele := range ui.Dialogs[constants.DialogItemOptions].Elements {
			switch ele.Key {
			case "item_options_title":
				name := "Item"
				switch firstTile.Block {
				case data.BlockJetpack:
					name = "Jetpack"
				case data.BlockFlamethrower:
					name = "Flamethrower"
				case data.BlockDisguise:
					name = "Disguise"
				case data.BlockBomb:
					name = "Bomb"
				}
				ele.Text.SetText(fmt.Sprintf("%s Options", name))
			case "item_timer":
				switch firstTile.Block {
				case data.BlockFlamethrower:
					ele.Text.SetText("Uses")
				default:
					ele.Text.SetText("Timer")
				}
			case "item_regenerate_check":
				ui.SetChecked(ele, firstTile.Metadata.Regenerate)
			case "item_regenerate_delay_input":
				ele.InputType = ui.Numeric
				ui.ChangeText(ele, fmt.Sprintf("%d", firstTile.Metadata.RegenDelay))
			case "item_timer_input":
				ele.InputType = ui.Numeric
				ui.ChangeText(ele, fmt.Sprintf("%d", firstTile.Metadata.Timer))
			}
		}
	}
}

func ConfirmItemOptions() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			ui.CloseDialog(constants.DialogItemOptions)
			return
		}
		var regen bool
		var delay, timer int
		for _, ele := range ui.Dialogs[constants.DialogItemOptions].Elements {
			switch ele.Key {
			case "item_regenerate_check":
				regen = ele.Checked
			case "item_regenerate_delay_input":
				di, err := strconv.Atoi(ele.Text.Raw)
				if err != nil {
					fmt.Println("WARNING: regen delay not an int:", err)
					di = 0
				}
				delay = di
			case "item_timer_input":
				di, err := strconv.Atoi(ele.Text.Raw)
				if err != nil {
					fmt.Println("WARNING: uses not an int:", err)
					di = 0
				}
				timer = di
			}
		}
		for _, tile := range data.CurrPuzzleSet.CurrPuzzle.WrenchTiles {
			tile.Metadata.Regenerate = regen
			tile.Metadata.Timer = timer
			tile.Metadata.RegenDelay = delay
			tile.Metadata.Changed = true
		}
		ui.CloseDialog(constants.DialogItemOptions)
		data.CurrPuzzleSet.CurrPuzzle.Update = true
		data.CurrPuzzleSet.CurrPuzzle.Changed = true
	}
}

// Other

func ChangeNumberInput(in *ui.Element, change int) {
	ChangeNumberInputWithLimits(in, change, 0, 99)
}

func ChangeNumberInputWithLimits(in *ui.Element, change, min, max int) {
	if in == nil {
		fmt.Printf("WARNING: input element is nil\n")
		return
	}
	di, err := strconv.Atoi(in.Text.Raw)
	if err != nil {
		fmt.Printf("WARNING: input %s not an int: %s\n", in.Key, err)
		di = 0
	}
	di += change
	if di < min {
		di = min
	} else if di > max {
		di = max
	}
	ui.ChangeText(in, fmt.Sprintf("%d", di))
}
