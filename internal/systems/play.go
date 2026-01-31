package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/content"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/sfx"
	"gemrunner/pkg/state"
	"gemrunner/pkg/util"
	"gemrunner/pkg/world"
	"time"
)

func PlayPauseSystem() {
	if !ui.DialogStackOpen && (data.MenuInput.Get("escape").JustPressed() ||
		data.P1Input.Get("pause").JustPressed() ||
		data.P2Input.Get("pause").JustPressed() ||
		data.P3Input.Get("pause").JustPressed() ||
		data.P4Input.Get("pause").JustPressed()) {
		data.MenuInput.Get("escape").Consume()
		data.P1Input.Get("pause").Consume()
		data.P2Input.Get("pause").Consume()
		data.P3Input.Get("pause").Consume()
		data.P4Input.Get("pause").Consume()
		ui.OpenDialogInStack(constants.DialogPauseMenu)
	}
}

func PlaySystem() {
	if data.MenuInput.Get("record").JustPressed() {
		data.CurrLevel.SaveRecord = true
	}
	data.CurrLevelSess.TimePlayed = time.Since(data.CurrLevelSess.LevelStart)
	if constants.Configuration.Gameplay.ShowTimer {
		UpdatePuzzleTimer()
	}
	allDead := true
	for _, p := range data.CurrLevel.Players {
		if p == nil {
			continue
		}
		if p.State != data.Dead && p.State != data.Waiting {
			allDead = false
		}
	}
	if allDead {
		DropScore()
		Restart()
	}
	if data.CurrLevel.Complete {
		if completion, ok := data.CurrLevelSess.LevelMap[data.CurrLevelSess.PuzzleIndex]; ok {
			data.CurrLevelSess.GemsCollected = append(data.CurrLevelSess.GemsCollected, completion.GemsCollected...)
		}
		data.CurrLevelSess.LevelMap[data.CurrLevelSess.PuzzleIndex] = data.LevelCompletion{
			Index:         data.CurrLevelSess.PuzzleIndex,
			GemsCollected: data.CurrLevelSess.GemsCollected,
			Completed:     data.CurrLevel.DoorsOpen,
			Continuity:    data.CurrLevelSess.PuzzleSet.Metadata.Adventure && data.CurrLevelSess.PuzzleSet.Metadata.Continuity != data.NoContinuity,
		}
		if data.CurrLevel.Recording && data.CurrLevel.SaveRecord {
			go content.SaveReplay(data.CurrLevel.LevelReplay)
		}
		UpdateScoreAndInv()
		data.CurrLevelSess.LastPuzzle = data.CurrLevelSess.PuzzleIndex
		data.CurrLevelSess.StartCoords = data.CurrLevel.StartCoords
		exitIndex := data.CurrLevel.ExitIndex
		if exitIndex == -1 {
			exitIndex = data.CurrLevelSess.PuzzleIndex + 1
		}

		if data.CurrLevelSess.PuzzleSet.Metadata.Adventure {
			next := data.CurrLevelSess.PuzzleSet.Puzzles[exitIndex]
			if next.Grid == data.CurrLevel.Puzzle.Grid {
				// go to the same puzzle
				GoToLevel(exitIndex)
			} else if util.Abs(next.Grid.X-data.CurrLevel.Puzzle.Grid.X)+
				util.Abs(next.Grid.Y-data.CurrLevel.Puzzle.Grid.Y) == 1 {
				// next puzzle is next to this one
				GoToLevel(exitIndex)
			} else {
				// next puzzle is "far away"
				ui.OpenDialogInStack(constants.DialogAdventureTrans)
			}
		} else {
			GoToLevel(exitIndex)
		}
	}
}

func Restart() {
	record := data.CurrLevel.Recording
	DisposeCurrLevel()
	ClearTemp()
	StartLevel(record)
	UpdateViews()
	data.EditorDraw = false
	reanimator.SetFrameRate(constants.Configuration.Gameplay.FrameRate)
	reanimator.Reset()
	sfx.MusicPlayer.GetStream("game").RepeatTrack(data.CurrLevel.Metadata.MusicTrack)
	go content.SaveSaveGame()
}

func GoToLevel(i int) {
	record := data.CurrLevel.Recording
	DisposeCurrLevel()
	ClearTemp()
	for {
		if i >= len(data.CurrPuzzleSet.Puzzles) {
			SetComplete()
			return
		}
		data.CurrPuzzleSet.SetTo(i)
		if data.CurrPuzzleSet.CurrPuzzle.NumPlayers() > 0 {
			break
		}
		i++
	}
	data.CurrLevelSess.PuzzleIndex = data.CurrPuzzleSet.PuzzleIndex
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	StartLevel(record)
	UpdateViews()
	data.EditorDraw = false
	reanimator.SetFrameRate(constants.Configuration.Gameplay.FrameRate)
	reanimator.Reset()
	sfx.MusicPlayer.GetStream("game").RepeatTrack(data.CurrLevel.Metadata.MusicTrack)
	go content.SaveSaveGame()
}

func SetComplete() {
	state.SwitchState(constants.MainMenuKey)
}

func DropScore() {
	data.CurrLevelSess.TotalTime = data.CurrLevelSess.TotalTime + data.CurrLevelSess.TimePlayed
	data.CurrLevelSess.LevelStart = time.Now()
	for _, stats := range data.CurrLevelSess.PlayerStats {
		if stats != nil {
			stats.LScore = 0
			stats.LGems = 0
		}
	}
	data.CurrLevelSess.GemsCollected = []world.Coords{}
}

func UpdateScoreAndInv() {
	data.CurrLevelSess.TotalTime = data.CurrLevelSess.TotalTime + data.CurrLevelSess.TimePlayed
	data.CurrLevelSess.LevelStart = time.Now()
	for i, stats := range data.CurrLevelSess.PlayerStats {
		if stats != nil {
			stats.Score += stats.LScore
			stats.Gems += stats.LGems
			stats.LScore = 0
			stats.LGems = 0
		}
		if p := data.CurrLevel.Players[i]; p != nil {
			stats.CurrBombs = p.SmallBombs
			if p.Inventory != nil {
				stats.Inventory = p.Inventory
				stats.Inventory.Entity.RemoveComponent(myecs.Temp)
			} else {
				if stats.Inventory != nil && stats.Inventory.Entity != nil {
					stats.Inventory.Entity.AddComponent(myecs.Temp, myecs.ClearFlag(false))
				}
				stats.Inventory = nil
			}
		}
	}
	data.CurrLevelSess.GemsCollected = []world.Coords{}
}
