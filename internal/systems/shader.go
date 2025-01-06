package systems

import (
	"gemrunner/internal/data"
	"gemrunner/pkg/timing"
)

func ShaderSystem() {
	if data.CurrPuzzleSet != nil {
		data.CurrPuzzleSet.Elapsed += float32(timing.DT)
	}
}
