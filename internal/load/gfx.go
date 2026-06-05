package load

import (
	"gemrunner/embed"
	"gemrunner/internal/data"
	"gemrunner/pkg/img"
	"github.com/gopxl/pixel"
	"github.com/pkg/errors"
)

func Shaders() {
	data.ColorShader = embed.ColorShader
	data.PuzzleShader = embed.PuzzleShader
	data.WorldShader = embed.WorldShader
	data.ScreenShader = embed.ScreenShader
	data.BossBlobShader = embed.BossBlobShader
}

func Backgrounds() error {
	errMsg := "loading backgrounds"
	blackBg, err := img.LoadImage("assets/backgrounds/black_background.png")
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	data.BlackBackground = pixel.NewSprite(blackBg, blackBg.Bounds())
	return nil
}
