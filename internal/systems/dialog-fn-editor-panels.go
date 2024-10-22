package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"strconv"
	"strings"
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
		theTile := data.CurrPuzzleSet.CurrPuzzle.WrenchTiles[0]
		ftDialog := ui.Dialogs[constants.DialogFloatingText]
		if theTile.FText == nil {
			data.SelectedTextColor = pixel.ToRGBA(constants.ColorWhite)
			data.SelectedShadowColor = pixel.ToRGBA(constants.ColorBlue)
			ui.ChangeText(ftDialog.Get("floating_text_value"), "")
			for _, ele := range ftDialog.Elements {
				if strings.Contains(ele.Key, "_check_color") {
					ui.SetChecked(ele, false)
				} else if strings.Contains(ele.Key, "_check_shadow") {
					ui.SetChecked(ele, false)
				}
			}
			ui.SetChecked(ftDialog.Get("white_check_color"), true)
			ui.SetChecked(ftDialog.Get("blue_check_shadow"), true)
			ui.SetChecked(ftDialog.Get("floating_text_shadow_check"), true)
			ui.SetChecked(ftDialog.Get("floating_text_show_check"), false)
			ui.SetChecked(ftDialog.Get("floating_text_bob_check"), true)
			ui.ChangeText(ftDialog.Get("floating_text_time_input"), "0")
			// add alignment here
		} else {
			data.SelectedTextColor = theTile.FText.Color
			data.SelectedShadowColor = theTile.FText.ShadowCol
			ui.ChangeText(ftDialog.Get("floating_text_value"), theTile.FText.Raw)
			for _, ele := range ftDialog.Elements {
				if strings.Contains(ele.Key, "_check_color") {
					updateColorCheckbox(ele, theTile.FText.Color)
				} else if strings.Contains(ele.Key, "_check_shadow") {
					updateColorCheckbox(ele, theTile.FText.ShadowCol)
				}
			}
			ui.SetChecked(ftDialog.Get("floating_text_shadow_check"), theTile.FText.HasShadow)
			ui.SetChecked(ftDialog.Get("floating_text_show_check"), theTile.FText.Prox)
			ui.SetChecked(ftDialog.Get("floating_text_bob_check"), theTile.FText.Bob)
			timerInput := ftDialog.Get("floating_text_time_input")
			timerInput.InputType = ui.Numeric
			ui.ChangeText(timerInput, strconv.Itoa(theTile.FText.Timer))
			// add alignment here
		}
	}
}

func ConfirmFloatingText() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			ui.CloseDialog(constants.DialogFloatingText)
			return
		}
		theTile := data.CurrPuzzleSet.CurrPuzzle.WrenchTiles[0]
		ftDialog := ui.Dialogs[constants.DialogFloatingText]

		timer, err := strconv.Atoi(ftDialog.Get("floating_text_time_input").Value)
		if err != nil {
			fmt.Println("WARNING: time input not an int:", err)
			timer = 0
		}

		if theTile.FText == nil {
			theTile.FText = data.NewFloatingText().
				WithTile(theTile).
				WithPos(theTile.Object.Pos)
		}

		theTile.FText.WithText(ftDialog.Get("floating_text_value").Value).
			WithColor(data.SelectedTextColor).
			WithShadow(data.SelectedShadowColor).
			WithTimer(timer)
		if !ftDialog.Get("floating_text_shadow_check").Checked {
			theTile.FText.RemoveShadow()
		}
		theTile.FText.Prox = ftDialog.Get("floating_text_show_check").Checked
		theTile.FText.Bob = ftDialog.Get("floating_text_bob_check").Checked

		ui.CloseDialog(constants.DialogFloatingText)
		data.CurrPuzzleSet.CurrPuzzle.Update = true
		data.CurrPuzzleSet.CurrPuzzle.Changed = true
		PushUndoArray(true)
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
					fmt.Println("WARNING: timer not an int:", err)
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

// Disguise Options

func OnOpenDisguiseOptions() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			ui.CloseDialog(constants.DialogDisguise)
			return
		}
		firstTile := data.CurrPuzzleSet.CurrPuzzle.WrenchTiles[0]
		for _, ele := range ui.Dialogs[constants.DialogDisguise].Elements {
			switch ele.Key {
			case "disguise_regenerate_check":
				ui.SetChecked(ele, firstTile.Metadata.Regenerate)
			case "disguise_regenerate_delay_input":
				ele.InputType = ui.Numeric
				ui.ChangeText(ele, fmt.Sprintf("%d", firstTile.Metadata.RegenDelay))
			case "disguise_timer_input":
				ele.InputType = ui.Numeric
				ui.ChangeText(ele, fmt.Sprintf("%d", firstTile.Metadata.Timer))
			}
		}
	}
}

func ConfirmDisguiseOptions() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			ui.CloseDialog(constants.DialogDisguise)
			return
		}
		var regen bool
		var delay, timer int
		for _, ele := range ui.Dialogs[constants.DialogDisguise].Elements {
			switch ele.Key {
			case "disguise_regenerate_check":
				regen = ele.Checked
			case "disguise_regenerate_delay_input":
				di, err := strconv.Atoi(ele.Text.Raw)
				if err != nil {
					fmt.Println("WARNING: regen delay not an int:", err)
					di = 0
				}
				delay = di
			case "disguise_timer_input":
				di, err := strconv.Atoi(ele.Text.Raw)
				if err != nil {
					fmt.Println("WARNING: timer not an int:", err)
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
		ui.CloseDialog(constants.DialogDisguise)
		data.CurrPuzzleSet.CurrPuzzle.Update = true
		data.CurrPuzzleSet.CurrPuzzle.Changed = true
		PushUndoArray(true)
	}
}

// Other

func ChangeNumberInput(in *ui.Element, change int) {
	if in == nil {
		fmt.Printf("WARNING: input element %s is nil\n", in.Key)
		return
	}
	di, err := strconv.Atoi(in.Text.Raw)
	if err != nil {
		fmt.Printf("WARNING: input %s not an int: %s\n", in.Key, err)
		di = 0
	}
	di += change
	if di < 0 {
		di = 0
	} else if di > 99 {
		di = 99
	}
	ui.ChangeText(in, fmt.Sprintf("%d", di))
}

func updateColorCheckbox(x *ui.Element, col pixel.RGBA) {
	key := x.Key[:strings.LastIndex(x.Key, "_")]
	switch key {
	case "red_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorRed))
	case "orange_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorOrange))
	case "green_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorGreen))
	case "cyan_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorCyan))
	case "blue_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorBlue))
	case "purple_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorPurple))
	case "pink_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorPink))
	case "yellow_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorYellow))
	case "gold_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorGold))
	case "brown_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorBrown))
	case "tan_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorTan))
	case "light_gray_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorLightGray))
	case "gray_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorGray))
	case "burnt_check":
		ui.SetChecked(x, col == pixel.ToRGBA(constants.ColorBurnt))
	}
}
