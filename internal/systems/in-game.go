package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
)

func InGameSystem() {
	if data.P1Input.Get("speedUp").JustPressed() {
		constants.FrameRate += constants.FrameRateInt
		if constants.FrameRate > constants.FrameRateMax {
			constants.FrameRate = constants.FrameRateMax
		}
	} else if data.P1Input.Get("speedDown").JustPressed() {
		constants.FrameRate -= constants.FrameRateInt
		if constants.FrameRate < constants.FrameRateMin {
			constants.FrameRate = constants.FrameRateMin
		}
	}
	if data.P1Input.Get("p1_kill").JustPressed() {
		data.CurrLevel.Players[0].Flags.Hit = true
		data.CurrLevel.Players[0].State = data.Hit
	}
	if data.P2Input.Get("p2_kill").JustPressed() {
		data.CurrLevel.Players[1].Flags.Hit = true
		data.CurrLevel.Players[1].State = data.Hit
	}
}
