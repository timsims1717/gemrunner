package main

import (
	"gemrunner/embed"
	"gemrunner/internal/constants"
	"gemrunner/internal/content"
	"gemrunner/internal/data"
	"gemrunner/internal/load"
	"gemrunner/internal/states"
	"gemrunner/internal/systems"
	"gemrunner/internal/ui"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/options"
	"gemrunner/pkg/sfx"
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
	constants.WinWidth = 1920
	constants.WinHeight = 1080
	options.RegisterResolution(pixel.V(1920, 1080))
	cfg := pixelgl.WindowConfig{
		Title:  constants.Title,
		Bounds: pixel.R(0, 0, constants.WinWidth, constants.WinHeight),
		VSync:  true,
	}
	options.BilinearFilter = false
	options.VSync = true
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetCursorVisible(false)
	viewport.ILockDefault = true
	viewport.MainCamera = viewport.New(win.Canvas())
	viewport.MainCamera.SetRect(pixel.R(0, 0, constants.WinWidth, constants.WinHeight))
	viewport.MainCamera.CamPos = pixel.V(constants.WinWidth*0.5, constants.WinHeight*0.5)

	data.ScreenView = viewport.New(nil)
	data.ScreenView.SetRect(pixel.R(constants.WinWidth*-0.5, constants.WinHeight*-0.5, constants.WinWidth*0.5, constants.WinHeight*0.5))
	data.ScreenView.PortPos = viewport.MainCamera.CamPos
	data.ScreenView.CamPos = pixel.V(constants.WinWidth*0.5, constants.WinHeight*0.5)

	state.Register(constants.MainMenuKey, state.New(states.MainMenuState))
	state.Register(constants.EditorStateKey, state.New(states.EditorState))
	state.Register(constants.TestStateKey, state.New(states.TestState))
	state.Register(constants.PlayStateKey, state.New(states.PlayState))
	//state.PushState(states.TestStateKey)
	//filename := fmt.Sprintf("%s/%s", constants.PuzzlesDir, "Get Those Gems.puzzle")
	//err = systems.OpenPuzzleFile(filename)
	//if err != nil {
	//	panic(err)
	//}

	//mainFont, err := typeface.LoadTTF("embed/Jive_Talking.ttf", 128.)
	mainFont, err := typeface.LoadBytes(embed.JiveTalking, 128.)
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

	//sh, err := shaders.LoadFileToString("embed/puzzle-shader.frag.glsl")
	//if err != nil {
	//	panic(err)
	//}
	//data.ColorShader = sh
	data.ColorShader = embed.ColorShader
	data.PuzzleShader = embed.PuzzleShader
	data.WorldShader = embed.WorldShader
	data.ScreenShader = embed.ScreenShader

	debug.Initialize(&viewport.MainCamera.PostCamPos)
	debug.Release = constants.Release
	debug.Version = constants.Version
	debug.Build = constants.Build
	debug.ShowText = false
	debug.ShowDebug = false
	debug.Verbose = true

	object.ILock = true

	content.LoadFavoritesFile()

	load.Music()
	load.SoundEffects()

	ui.ScrollSpeed = constants.ScrollSpeed
	load.DialogConstructors()
	load.InitEditorConstructors()
	load.InitEditorPanels()
	load.InitMainMenuConstructors()
	systems.InitMainBorder()
	systems.CursorInit()

	//data.ScreenView.Canvas.SetFragmentShader(data.ScreenShader)
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
		if data.DebugInput.Get("debug").JustPressed() {
			debug.ShowDebug = !debug.ShowDebug
			debug.ShowText = !debug.ShowText
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
