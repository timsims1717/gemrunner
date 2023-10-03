package main

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/states"
	"gemrunner/internal/systems"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/img"
	"gemrunner/pkg/state"
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
	viewport.MainCamera.CamPos = pixel.V(1600*0.5, 900*0.5)

	state.Register(states.EditorStateKey, state.New(states.EditorState))

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

	debug.Initialize(&viewport.MainCamera.PostCamPos)
	debug.Text = true

	systems.InitMainBorder()

	win.Show()
	timing.Reset()
	for !win.Closed() {
		timing.Update()
		debug.Clear()
		data.DebugInput.Update(win, viewport.MainCamera.Mat)
		//options.WindowUpdate(win)
		//if options.Updated {
		//	viewport.MainCamera.CamPos = pixel.V(viewport.MainCamera.Rect.W()*0.5, viewport.MainCamera.Rect.H()*0.5)
		//}

		state.Update(win)
		viewport.MainCamera.Update()
		state.Draw(win)

		//win.SetSmooth(false)
		debug.Draw(win)
		//win.SetSmooth(true)
		//win.SetSmooth(options.BilinearFilter)

		//sfx.MusicPlayer.Update()
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
