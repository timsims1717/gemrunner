package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/content"
	"gemrunner/internal/data"
	"gemrunner/internal/ui"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/sfx"
	"gemrunner/pkg/state"
	"time"
)

func PlaySystem() {
	data.CurrLevelSess.TimePlayed = time.Since(data.CurrLevelSess.LevelStart)
	allDead := true
	for _, p := range data.CurrLevel.Players {
		if p == nil {
			continue
		}
		if p.State != data.Dead {
			allDead = false
		}
	}
	if allDead {
		DropScore()
		Restart()
	}
	if data.CurrLevel.Complete {
		UpdateScore()
		if data.CurrPuzzleSet.PuzzleIndex == len(data.CurrPuzzleSet.Puzzles)-1 {
			SetComplete()
		} else {
			NextLevel()
		}
	}
	if data.MenuInput.Get("escape").JustPressed() ||
		data.P1Input.Get("pause").JustPressed() ||
		data.P2Input.Get("pause").JustPressed() ||
		data.P3Input.Get("pause").JustPressed() ||
		data.P4Input.Get("pause").JustPressed() {
		data.MenuInput.Get("escape").Consume()
		data.P1Input.Get("pause").Consume()
		data.P2Input.Get("pause").Consume()
		data.P3Input.Get("pause").Consume()
		data.P4Input.Get("pause").Consume()
		ui.OpenDialogInStack(constants.DialogPauseMenu)
	}
}

func Restart() {
	LevelDispose()
	ClearTemp()
	LevelInit()
	UpdateViews()
	data.EditorDraw = false
	reanimator.SetFrameRate(constants.FrameRate)
	reanimator.Reset()
	sfx.MusicPlayer.GetStream("game").RepeatTrack(data.CurrLevel.Metadata.MusicTrack)
	go content.SaveSaveGame()
}

func NextLevel() {
	LevelDispose()
	ClearTemp()
	data.CurrPuzzleSet.Next()
	data.CurrLevelSess.PuzzleIndex = data.CurrPuzzleSet.PuzzleIndex
	PuzzleInit()
	LevelInit()
	UpdateViews()
	data.EditorDraw = false
	reanimator.SetFrameRate(constants.FrameRate)
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
}

func UpdateScore() {
	data.CurrLevelSess.TotalTime = data.CurrLevelSess.TotalTime + data.CurrLevelSess.TimePlayed
	data.CurrLevelSess.LevelStart = time.Now()
	for _, stats := range data.CurrLevelSess.PlayerStats {
		if stats != nil {
			stats.Score += stats.LScore
			stats.Gems += stats.LGems
			stats.LScore = 0
			stats.LGems = 0
		}
	}
}
