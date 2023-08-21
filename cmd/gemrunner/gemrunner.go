package main

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/state"
	"gemrunner/internal/systems"
	"gemrunner/pkg/img"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	world.SetTileSize(constants.TileSize)
	cfg := pixelgl.WindowConfig{
		Title:  "Gem Runner",
		Bounds: pixel.R(0, 0, 1600, 900),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	viewport.MainCamera = viewport.New(win.Canvas())
	viewport.MainCamera.SetRect(pixel.R(0, 0, 1600, 900))

	uiSheet, err := img.LoadSpriteSheet("assets/ui.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.UIBatch, uiSheet, true, true)
	tileSheet, err := img.LoadSpriteSheet("assets/tileset.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.TileBGBatch, tileSheet, true, true)
	img.AddBatcher(constants.TileFGBatch, tileSheet, true, true)

	systems.InitMainBorder()

	timing.Reset()
	for !win.Closed() {
		state.Update(win)
		state.Draw(win)
	}
}

func main() {
	pixelgl.Run(run)
}
