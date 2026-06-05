package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/random"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/util"
	"gemrunner/pkg/world"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

func EffectsSystem() {
	data.ShaderTime += float32(timing.DT)
	data.ShaderTime64 += timing.DT
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

func setCanvasShaderColorsFromMD(canvas *pixelgl.Canvas, metadata data.PuzzleMetadata) {
	canvas.SetUniform("uPrimary", util.RGBAToVec3(metadata.PrimaryColor))
	canvas.SetUniform("uSecondary", util.RGBAToVec3(metadata.SecondaryColor))
	canvas.SetUniform("uDoodad", util.RGBAToVec3(metadata.DoodadColor))
	canvas.SetUniform("uGoop", util.RGBAToVec3(metadata.GoopColor))
	canvas.SetUniform("uLiquidPrimary", util.RGBAToVec3(metadata.LiquidPrimaryColor))
	canvas.SetUniform("uLiquidSecondary", util.RGBAToVec3(metadata.LiquidSecondaryColor))
}

func setCanvasShaderColors(canvas *pixelgl.Canvas, pc, sc, dc, gc, lpc, lsc pixel.RGBA) {
	canvas.SetUniform("uPrimary", util.RGBAToVec3(pc))
	canvas.SetUniform("uSecondary", util.RGBAToVec3(sc))
	canvas.SetUniform("uDoodad", util.RGBAToVec3(dc))
	canvas.SetUniform("uGoop", util.RGBAToVec3(gc))
	canvas.SetUniform("uLiquidPrimary", util.RGBAToVec3(lpc))
	canvas.SetUniform("uLiquidSecondary", util.RGBAToVec3(lsc))
}

func setCanvasShaderColorsDefault(canvas *pixelgl.Canvas) {
	canvas.SetUniform("uPrimary", mgl32.Vec3{})
	canvas.SetUniform("uSecondary", mgl32.Vec3{})
	canvas.SetUniform("uDoodad", mgl32.Vec3{})
	canvas.SetUniform("uGoop", mgl32.Vec3{})
	canvas.SetUniform("uLiquidPrimary", mgl32.Vec3{})
	canvas.SetUniform("uLiquidSecondary", mgl32.Vec3{})
}
