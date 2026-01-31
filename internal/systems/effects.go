package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/random"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/util"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

func EffectsSystem() {
	data.ShaderTime += float32(timing.DT)
	// screen shake
	if data.ScreenShake != nil && constants.Configuration.Gameplay.ScreenShake &&
		data.CurrentPlayArea != nil {
		offset, fin := data.ScreenShake.Shake(timing.DT)
		if fin {
			data.CurrentPlayArea.WorldView.Offset = pixel.ZV
		} else {
			data.CurrentPlayArea.WorldView.Offset = offset
		}
	}
	// shadow
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

func ShakeScreen() {
	if constants.Configuration.Gameplay.ScreenShake {
		if data.ScreenShake == nil {
			data.ScreenShake = util.NewShaker(30., 20., 0.5, random.Effects.Int63())
		} else {
			data.ScreenShake.Reset(random.Effects.Int63())
		}
	}
}
