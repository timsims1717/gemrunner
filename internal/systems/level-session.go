package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/ui"
	"time"
)

func LevelSessionInit() {
	if data.CurrPuzzleSet == nil {
		panic("no puzzle set loaded to start")
	}
	numPlayers := data.CurrPuzzleSet.CurrPuzzle.NumPlayers()
	if data.CurrPuzzleSet.Metadata.NumPlayers < numPlayers {
		data.CurrPuzzleSet.Metadata.NumPlayers = numPlayers
	}
	if data.CurrLevelSess == nil {
		data.CurrLevelSess = &data.LevelSession{}
		for p := 0; p < data.CurrPuzzleSet.Metadata.NumPlayers; p++ {
			data.CurrLevelSess.PlayerStats[p] = data.NewStats()
		}
		data.CurrLevelSess.PuzzleFile = data.CurrPuzzleSet.Metadata.Name
		data.CurrLevelSess.Filename = fmt.Sprintf("%s%s", data.CurrPuzzleSet.Metadata.Name, constants.SaveExt)
	} else {
		data.CurrPuzzleSet.SetTo(data.CurrLevelSess.PuzzleIndex)
		PuzzleInit()
	}
	data.CurrLevelSess.Metadata = data.CurrPuzzleSet.Metadata
	data.CurrLevelSess.PuzzleSet = data.CurrPuzzleSet
	data.CurrLevelSess.LevelStart = time.Now()
	for p := 0; p < data.CurrPuzzleSet.Metadata.NumPlayers; p++ {
		switch p {
		case 0:
			ui.OpenDialog(constants.DialogPlayer1Inv)
		case 1:
			ui.OpenDialog(constants.DialogPlayer2Inv)
		case 2:
			ui.OpenDialog(constants.DialogPlayer3Inv)
		case 3:
			ui.OpenDialog(constants.DialogPlayer4Inv)
		}
	}
	if data.CurrPuzzleSet.CurrPuzzle.Metadata.Name != "" {
		ui.OpenDialog(constants.DialogPuzzleTitle)
	}
}

func LevelSessionDispose() {
	data.CurrLevelSess = nil
}
