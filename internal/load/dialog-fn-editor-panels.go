package load

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems"
	"gemrunner/pkg/world"
	"strconv"
)

// editor panel mode buttons

func EditorMode(mode data.EditorMode, btn *data.Button, dialog *data.Dialog) func() {
	return func() {
		data.Editor.SelectVis = false
		if data.Editor.Mode != mode {
			data.Editor.LastMode = data.Editor.Mode
			data.Editor.LastCoords = world.Coords{X: -1, Y: -1}
		}
		data.Editor.Mode = mode
		for _, e := range dialog.Elements {
			if b, ok := e.(*data.Button); ok {
				b.Entity.AddComponent(myecs.Drawable, b.Sprite)
			}
		}
		btn.Entity.AddComponent(myecs.Drawable, btn.ClickSpr)
		data.CurrPuzzleSet.CurrPuzzle.Update = true
	}
}

// wrench dialogs

func OnOpenCrackTileOptions() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			data.CloseDialog(constants.DialogCrackedTiles)
			return
		}
		firstTile := data.CurrPuzzleSet.CurrPuzzle.WrenchTiles[0]
		crackDialog := data.Dialogs[constants.DialogCrackedTiles]
		for _, ele := range crackDialog.Elements {
			if x, ok := ele.(*data.Checkbox); ok {
				switch x.Key {
				case "cracked_tile_regenerate_check":
					data.SetChecked(x, firstTile.Metadata.Regenerate)
				case "cracked_tile_show_check":
					data.SetChecked(x, firstTile.Metadata.ShowCrack)
				case "cracked_tile_enemy_check":
					data.SetChecked(x, firstTile.Metadata.EnemyCrack)
				}
			} else if t, okT := ele.(*data.Text); okT {
				if t.Key == "cracked_tile_title" {
					if firstTile.Block == data.BlockCracked {
						t.Text.SetText("Cracked Turf")
					} else {
						t.Text.SetText("Cracked Ladder")
					}
				}
			}
		}
	}
}

func ConfirmCrackTileOptions() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			data.CloseDialog(constants.DialogCrackedTiles)
			return
		}
		crackDialog := data.Dialogs[constants.DialogCrackedTiles]
		var regen, show, enemy bool
		for _, ele := range crackDialog.Elements {
			if x, ok := ele.(*data.Checkbox); ok {
				switch x.Key {
				case "cracked_tile_regenerate_check":
					regen = x.Checked
				case "cracked_tile_show_check":
					show = x.Checked
				case "cracked_tile_enemy_check":
					enemy = x.Checked
				}
			}
		}
		for _, tile := range data.CurrPuzzleSet.CurrPuzzle.WrenchTiles {
			tile.Metadata.Regenerate = regen
			tile.Metadata.ShowCrack = show
			tile.Metadata.EnemyCrack = enemy
			tile.Metadata.Changed = true
		}
		data.CloseDialog(constants.DialogCrackedTiles)
		data.CurrPuzzleSet.CurrPuzzle.Update = true
		data.CurrPuzzleSet.CurrPuzzle.Changed = true
		systems.PushUndoArray(true)
	}
}

// Bomb Options

func OnOpenBombOptions() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			data.CloseDialog(constants.DialogBomb)
			return
		}
		firstTile := data.CurrPuzzleSet.CurrPuzzle.WrenchTiles[0]
		bombDialog := data.Dialogs[constants.DialogBomb]
		for _, ele := range bombDialog.Elements {
			if x, ok := ele.(*data.Checkbox); ok {
				switch x.Key {
				case "bomb_cross_check":
					data.SetChecked(x, firstTile.Metadata.BombCross)
				case "bomb_regenerate_check":
					data.SetChecked(x, firstTile.Metadata.Regenerate)
				}
			} else if t, okT := ele.(*data.Text); okT {
				if t.Key == "bomb_options_title" {
					if firstTile.Block == data.BlockBomb {
						t.Text.SetText("Bomb Item Options")
					} else {
						t.Text.SetText("Lit Bomb Options")
					}
				}
			} else if i, okI := ele.(*data.Input); okI {
				if i.Key == "bomb_regenerate_delay_input" {
					i.NumbersOnly = true
					data.ChangeText(i, fmt.Sprintf("%d", firstTile.Metadata.RegenDelay))
				}
			}
		}
	}
}

func ConfirmBombOptions() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			data.CloseDialog(constants.DialogBomb)
			return
		}
		bombDialog := data.Dialogs[constants.DialogBomb]
		var regen, cross bool
		var delay int
		for _, ele := range bombDialog.Elements {
			if x, ok := ele.(*data.Checkbox); ok {
				switch x.Key {
				case "bomb_cross_check":
					cross = x.Checked
				case "bomb_regenerate_check":
					regen = x.Checked
				}
			} else if i, okI := ele.(*data.Input); okI {
				if i.Key == "bomb_regenerate_delay_input" {
					di, err := strconv.Atoi(i.Text.Raw)
					if err != nil {
						fmt.Println("WARNING: regen delay not an int:", err)
						di = 0
					}
					delay = di
				}
			}
		}
		for _, tile := range data.CurrPuzzleSet.CurrPuzzle.WrenchTiles {
			tile.Metadata.Regenerate = regen
			tile.Metadata.BombCross = cross
			tile.Metadata.RegenDelay = delay
			tile.Metadata.Changed = true
		}
		data.CloseDialog(constants.DialogBomb)
		data.CurrPuzzleSet.CurrPuzzle.Update = true
		data.CurrPuzzleSet.CurrPuzzle.Changed = true
		systems.PushUndoArray(true)
	}
}

func IncOrDecBombInput(inc bool) {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		bombDialog := data.Dialogs[constants.DialogBomb]
		for _, ele := range bombDialog.Elements {
			if i, okI := ele.(*data.Input); okI {
				if i.Key == "bomb_regenerate_delay_input" {
					di, err := strconv.Atoi(i.Text.Raw)
					if err != nil {
						fmt.Println("WARNING: regen delay not an int:", err)
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
					data.ChangeText(i, fmt.Sprintf("%d", di))
				}
			}
		}
	}
}
