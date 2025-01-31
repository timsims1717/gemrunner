package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/world"
)

func ShaderSystem() {
	if data.CurrPuzzleSet != nil {
		data.CurrPuzzleSet.Elapsed += float32(timing.DT)
	}
	if data.CurrLevel != nil {
		for p := 0; p < constants.MaxPlayers; p++ {
			if p < data.CurrPuzzleSet.CurrPuzzle.NumPlayers() {
				data.CurrLevel.PLoc[p][0] = float32(data.CurrLevel.Players[p].Object.Pos.X / (float64(data.CurrLevel.Metadata.Width) * world.TileSize))
				data.CurrLevel.PLoc[p][1] = float32(data.CurrLevel.Players[p].Object.Pos.Y / (float64(data.CurrLevel.Metadata.Height) * world.TileSize))
			} else {
				data.CurrLevel.PLoc[p][0] = -1
				data.CurrLevel.PLoc[p][1] = -1
			}
		}
	}
}
