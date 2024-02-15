package systems

import (
	"gemrunner/internal/data"
	"gemrunner/pkg/state"
)

func TestSystem() {
	if data.CurrLevel.Players[0].State == data.Dead {
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
