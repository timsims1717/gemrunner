package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/ui"
)

// door options

func OpenDoorOptionsDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			return
		}
		ui.NewDialog(ui.DialogConstructors[constants.DialogDoors])
		doorOptions := ui.Dialogs[constants.DialogDoors]
		for _, ele := range doorOptions.Elements {
			switch ele.Key {
			case "confirm":
				ele.OnClick = ConfirmDoorOptions
			case "door_next":
				ele.OnHold = PuzzleSetViewNextPuzzle(doorOptions)
				ele.OnClick = PuzzleSetViewNextPuzzle(doorOptions)
			case "door_prev":
				ele.OnHold = PuzzleSetViewPrevPuzzle(doorOptions)
				ele.OnClick = PuzzleSetViewPrevPuzzle(doorOptions)
			case "door_end":
				ele.OnHold = PuzzleSetViewGoToEndPuzzle(doorOptions)
				ele.OnClick = PuzzleSetViewGoToEndPuzzle(doorOptions)
			case "door_begin":
				ele.OnHold = PuzzleSetViewGoToBeginPuzzle(doorOptions)
				ele.OnClick = PuzzleSetViewGoToBeginPuzzle(doorOptions)
			case "cancel":
				ele.OnClick = DisposeDialog(constants.DialogDoors)
			}
		}
		UpdateDialogView(doorOptions)
		data.PuzzleSetViewAllowEnd = true
		data.PuzzleSetViewIsMoving = false
		exitIndex := data.CurrPuzzleSet.CurrPuzzle.WrenchTiles[0].Metadata.ExitIndex
		if exitIndex == -1 {
			data.PuzzleSetViewIndex = data.CurrPuzzleSet.PuzzleIndex + 1
		} else {
			data.PuzzleSetViewIndex = exitIndex
		}
		data.PuzzleSetViewPuzzles = make([]int, len(data.CurrPuzzleSet.Puzzles))
		for i := range data.PuzzleSetViewPuzzles {
			data.PuzzleSetViewPuzzles[i] = i
		}
		pzlView := doorOptions.Get("puzzle_set_view")
		CreatePuzzlePreview(pzlView.Get("puzzle_center"), data.PuzzleSetViewIndex)
		CreatePuzzlePreview(pzlView.Get("puzzle_left"), data.PuzzleSetViewIndex-1)
		CreatePuzzlePreview(pzlView.Get("puzzle_right"), data.PuzzleSetViewIndex+1)
		ResetPuzzleSetView(doorOptions)
		PuzzleSetViewNameAndNum(doorOptions, data.PuzzleSetViewIndex)
		ui.OpenDialogInStack(constants.DialogDoors)
	}
}

func ConfirmDoorOptions() {
	if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
		if len(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles) < 1 {
			fmt.Println("WARNING: no tiles selected by wrench")
			ui.CloseDialog(constants.DialogDoors)
			return
		}
		for _, tile := range data.CurrPuzzleSet.CurrPuzzle.WrenchTiles {
			tile.Metadata.ExitIndex = data.PuzzleSetViewIndex
			tile.Metadata.Changed = true
		}
		ui.Dispose(constants.DialogDoors)
		data.CurrPuzzleSet.CurrPuzzle.Update = true
		data.CurrPuzzleSet.CurrPuzzle.Changed = true
	}
}
