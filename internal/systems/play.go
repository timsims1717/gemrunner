package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/sfx"
	"gemrunner/pkg/state"
)

func PlaySystem() {
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
		Restart()
	}
	if data.CurrLevel.Complete {
		if data.CurrPuzzleSet.PuzzleIndex == len(data.CurrPuzzleSet.Puzzles)-1 {
			SetComplete()
		} else {
			NextLevel()
		}
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
}

func NextLevel() {
	LevelDispose()
	ClearTemp()
	data.CurrPuzzleSet.Next()
	PuzzleInit()
	UpdatePuzzleShaders()
	LevelInit()
	UpdateViews()
	data.EditorDraw = false
	reanimator.SetFrameRate(constants.FrameRate)
	reanimator.Reset()
	sfx.MusicPlayer.GetStream("game").RepeatTrack(data.CurrLevel.Metadata.MusicTrack)
}

func SetComplete() {
	LevelDispose()
	ClearTemp()
	PuzzleDispose()
	data.CurrPuzzleSet = nil
	state.SwitchState(constants.MainMenuKey)
}
