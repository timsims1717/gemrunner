package systems

import (
	"gemrunner/internal/data"
	"gemrunner/pkg/state"
	"time"
)

func TestSystem() {
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
		data.CurrLevel.Failed = true
	}
	if data.MenuInput.Get("escape").JustPressed() ||
		data.CurrLevel.Complete || data.CurrLevel.Failed {
		if data.CurrLevel.Complete {
			if data.CurrPuzzleSet.CurrPuzzle != nil {
				data.CurrPuzzleSet.CurrPuzzle.Metadata.Completed = true
			}
		}
		state.PopState()
	}
}
