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

// floating text

func OnOpenFloatingText() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			ui.CloseDialog(constants.DialogFloatingText)
			return
		}
		//firstTile := data.CurrPuzzleSet.CurrPuzzle.WrenchTiles[0]
		//ftDialog := ui.Dialogs[constants.DialogFloatingText]
		//for _, ele := range ftDialog.Elements {
		//	switch ele.Key {
		//
		//	}
		//}
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
		PushUndoArray(true)
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
		PushUndoArray(true)
	}
}

func IncOrDecBombRegen(inc bool) {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		bombDialog := ui.Dialogs[constants.DialogBomb]
		for _, ele := range bombDialog.Elements {
			if ele.Key == "bomb_regenerate_delay_input" {
				IncOrDecNumberInput(ele, inc)
			}
		}
	}
}

// Jetpack Options

func OnOpenJetpackOptions() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			ui.CloseDialog(constants.DialogJetpack)
			return
		}
		firstTile := data.CurrPuzzleSet.CurrPuzzle.WrenchTiles[0]
		for _, ele := range ui.Dialogs[constants.DialogJetpack].Elements {
			switch ele.Key {
			case "jetpack_regenerate_check":
				ui.SetChecked(ele, firstTile.Metadata.Regenerate)
			case "jetpack_regenerate_delay_input":
				ele.InputType = ui.Numeric
				ui.ChangeText(ele, fmt.Sprintf("%d", firstTile.Metadata.RegenDelay))
			case "jetpack_timer_input":
				ele.InputType = ui.Numeric
				ui.ChangeText(ele, fmt.Sprintf("%d", firstTile.Metadata.Timer))
			}
		}
	}
}

func ConfirmJetpackOptions() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			ui.CloseDialog(constants.DialogJetpack)
			return
		}
		var regen bool
		var delay, timer int
		for _, ele := range ui.Dialogs[constants.DialogJetpack].Elements {
			switch ele.Key {
			case "jetpack_regenerate_check":
				regen = ele.Checked
			case "jetpack_regenerate_delay_input":
				di, err := strconv.Atoi(ele.Text.Raw)
				if err != nil {
					fmt.Println("WARNING: regen delay not an int:", err)
					di = 0
				}
				delay = di
			case "jetpack_timer_input":
				di, err := strconv.Atoi(ele.Text.Raw)
				if err != nil {
					fmt.Println("WARNING: regen delay not an int:", err)
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
		ui.CloseDialog(constants.DialogJetpack)
		data.CurrPuzzleSet.CurrPuzzle.Update = true
		data.CurrPuzzleSet.CurrPuzzle.Changed = true
		PushUndoArray(true)
	}
}

func IncOrDecJetpackRegen(inc bool) {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		bombDialog := ui.Dialogs[constants.DialogJetpack]
		for _, ele := range bombDialog.Elements {
			if ele.Key == "jetpack_regenerate_delay_input" {
				IncOrDecNumberInput(ele, inc)
			}
		}
	}
}

func IncOrDecJetpackTimer(inc bool) {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		bombDialog := ui.Dialogs[constants.DialogJetpack]
		for _, ele := range bombDialog.Elements {
			if ele.Key == "jetpack_timer_input" {
				IncOrDecNumberInput(ele, inc)
			}
		}
	}
}

func IncOrDecNumberInput(in *ui.Element, inc bool) {
	di, err := strconv.Atoi(in.Text.Raw)
	if err != nil {
		fmt.Printf("WARNING: input %s not an int: %s\n", in.Key, err)
		di = 0
	}
	if inc {
		di++
	} else {
		di--
	}
	if di < 0 {
		di = 0
	} else if di > 99 {
		di = 99
	}
	ui.ChangeText(in, fmt.Sprintf("%d", di))
}
