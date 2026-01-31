package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/content"
	"gemrunner/internal/data"
	"gemrunner/internal/random"
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
		data.CurrLevelSess.LevelMap = make(map[int]data.LevelCompletion)
		data.CurrLevelSess.PuzzleFile = data.CurrPuzzleSet.Metadata.Name
		data.CurrLevelSess.Filename = fmt.Sprintf("%s%s", data.CurrPuzzleSet.Metadata.Name, constants.SaveExt)
	} else {
		data.CurrPuzzleSet.SetTo(data.CurrLevelSess.PuzzleIndex)
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
	if constants.Configuration.Gameplay.ShowTimer {
		ui.OpenDialog(constants.DialogPuzzleTimer)
	}
}

func StartLevel(record bool) {
	if data.CurrentPlayArea == nil {
		data.CurrentPlayArea = CreatePlayArea()
	}
	data.CurrLevelSess.PuzzleIndex = data.CurrPuzzleSet.PuzzleIndex
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitLevel(data.CurrentPlayArea)
	//InitPlayArea(data.CurrentPlayArea)

	levelSeed := random.RandomSeed()
	random.SetLevelSeed(levelSeed)

	// set up recording stuff
	data.CurrentPlayArea.Level.Recording = data.CurrReplay == nil && record
	data.CurrentPlayArea.Level.SaveRecord = constants.Configuration.Gameplay.AlwaysRecord
	if data.CurrentPlayArea.Level.Recording {
		data.CurrentPlayArea.Level.LevelReplay = &data.LevelReplay{
			PuzzleSet:   data.CurrPuzzleSet.Metadata.Name,
			Filename:    data.CurrPuzzleSet.Metadata.Filename,
			ReplayFile:  content.ReplayFile(data.CurrPuzzleSet.Metadata.Name, data.CurrPuzzleSet.PuzzleIndex),
			PuzzleNum:   data.CurrPuzzleSet.PuzzleIndex,
			StartCoords: data.CurrLevelSess.StartCoords,
			Seed:        levelSeed,
		}
		data.CurrentPlayArea.Level.ReplayFrame = data.ReplayFrame{}
	} else if data.CurrReplay != nil {
		data.CurrReplay.FrameIndex = 0
		data.CurrentPlayArea.Level.StartCoords = data.CurrReplay.StartCoords
	}
}

func LevelSessionDispose() {
	data.CurrLevelSess = nil
}

func LevelTransitionSystem() {

}
