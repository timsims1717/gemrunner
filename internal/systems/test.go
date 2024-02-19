package systems

import (
	"gemrunner/internal/data"
	"gemrunner/pkg/state"
)

func TestSystem() {
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
			if data.CurrPuzzle != nil {
				data.CurrPuzzle.Metadata.Completed = true
			}
		}
		state.PopState()
	}
}
