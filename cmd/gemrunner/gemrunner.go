package main

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/load"
	"gemrunner/internal/states"
	"gemrunner/internal/systems"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/options"
	"gemrunner/pkg/sfx"
	"gemrunner/pkg/shaders"
	"gemrunner/pkg/state"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/typeface"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/gopxl/pixel/text"
)

func run() {
	world.SetTileSize(constants.TileSize)
	cfg := pixelgl.WindowConfig{
		Title:  "Gem Runner",
		Bounds: pixel.R(0, 0, 1600, 900),
		VSync:  true,
	}
	options.BilinearFilter = false
	options.VSync = true
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	viewport.ILockDefault = true
	viewport.MainCamera = viewport.New(win.Canvas())
	viewport.MainCamera.SetRect(pixel.R(0, 0, 1600, 900))
	viewport.MainCamera.CamPos = pixel.V(1600*0.5, 900*0.5)

	state.Register(states.EditorStateKey, state.New(states.EditorState))
	state.Register(states.TestStateKey, state.New(states.TestState))
	//state.PushState(states.TestStateKey)
	//filename := fmt.Sprintf("%s/%s", constants.PuzzlesDir, "Get Those Gems.puzzle")
	//err = systems.OpenPuzzle(filename)
	//if err != nil {
	//	panic(err)
	//}

	mainFont, err := typeface.LoadTTF("assets/Jive_Talking.ttf", 128.)
	typeface.Atlases["main"] = text.NewAtlas(mainFont, text.ASCII)

	uiSheet, err := img.LoadSpriteSheet("assets/ui.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.UIBatch, uiSheet, true, true)
	tileSheet, err := img.LoadSpriteSheet("assets/tileset.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.TileBatch, tileSheet, true, true)

	sh, err := shaders.LoadFileToString("assets/shaders/puzzle-shader.frag.glsl")
	if err != nil {
		panic(err)
	}
	data.PuzzleShader = sh

	debug.Initialize(&viewport.MainCamera.PostCamPos)
	debug.ShowText = true
	debug.ShowDebug = true

	object.ILock = true

	load.Music()

	load.InitConstructors()
	load.Dialogs(win)
	systems.InitMainBorder()

	win.Show()
	timing.Reset()
	for !win.Closed() {
		timing.Update()
		debug.Clear()
		options.WindowUpdate(win)
		if options.Updated {
			viewport.MainCamera.CamPos = pixel.V(viewport.MainCamera.Rect.W()*0.5, viewport.MainCamera.Rect.H()*0.5)
		}

		data.MenuInput.Update(win, viewport.MainCamera.Mat)
		data.DebugInput.Update(win, viewport.MainCamera.Mat)
		if data.DebugInput.Get("debugPause").JustPressed() {
			state.ToggleDebugPause()
		}
		if data.DebugInput.Get("debugFrame").JustPressed() || data.DebugInput.Get("debugFrame").Repeated() {
			state.DebugFrameAdvance()
		}
		if data.DebugInput.Get("fullscreen").JustPressed() {
			options.FullScreen = !options.FullScreen
		}
		if data.DebugInput.Get("fuzzy").JustPressed() {
			options.BilinearFilter = !options.BilinearFilter
		}
		if data.DebugInput.Get("debugText").JustPressed() {
			debug.ShowText = !debug.ShowText
		}
		if data.DebugInput.Get("debug").JustPressed() {
			debug.ShowDebug = !debug.ShowDebug
		}

		state.Update(win)
		win.Clear(constants.ColorBlack)
		viewport.MainCamera.Update()
		state.Draw(win)

		win.SetSmooth(false)
		debug.DrawText(win)
		debug.DrawFPS(win)
		win.SetSmooth(options.BilinearFilter)

		sfx.MusicPlayer.Update()
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
