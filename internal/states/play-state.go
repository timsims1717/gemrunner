package states

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems"
	"gemrunner/internal/ui"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/img"
	"gemrunner/pkg/options"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/sfx"
	"gemrunner/pkg/state"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var (
	PlayState = &playState{}
)

type playState struct {
	*state.AbstractState
}

func (s *playState) Unload(win *pixelgl.Window) {
	ui.ClearDialogsOpen()
	ui.ClearDialogStack()
	systems.DisposeInGameDialogs()
	systems.LevelSessionDispose()
	systems.LevelDispose()
	systems.ClearTemp()
	systems.PuzzleDispose()
	data.CurrPuzzleSet = nil
}

func (s *playState) Load(win *pixelgl.Window) {
	data.EditorDraw = false
	ui.ClearDialogsOpen()
	ui.ClearDialogStack()
	systems.InGameDialogs(win)
	systems.LevelSessionInit()
	systems.LevelInit()
	systems.UpdateViews()
	reanimator.SetFrameRate(constants.FrameRate)
	reanimator.Reset()
	sfx.MusicPlayer.GetStream("game").RepeatTrack(data.CurrLevel.Metadata.MusicTrack)
}

func (s *playState) Update(win *pixelgl.Window) {
	data.P1Input.Update(win, viewport.MainCamera.Mat)
	data.P2Input.Update(win, viewport.MainCamera.Mat)
	data.P3Input.Update(win, viewport.MainCamera.Mat)
	data.P4Input.Update(win, viewport.MainCamera.Mat)
	systems.CursorSystem(true)
	debug.AddText("Play State")
	debug.AddText(fmt.Sprintf("Speed: %d", constants.FrameRate))
	debug.AddText(systems.FormatTimePlayed())
	debug.AddText(fmt.Sprintf("Frame Number: %d", data.CurrLevel.FrameNumber))
	//debug.AddText(fmt.Sprintf("Frame Counter: %d", data.CurrLevel.FrameCounter))
	debug.AddText(fmt.Sprintf("Frame Cycle: %d", data.CurrLevel.FrameCycle))
	//debug.AddTruthText("Frame Change", data.CurrLevel.FrameChange)
	for i, player := range data.CurrLevel.Players {
		if player != nil {
			pos := player.Object.Pos
			debug.AddIntCoords(fmt.Sprintf("Player %d Pos", i+1), int(pos.X), int(pos.Y))
			cx, cy := world.WorldToMap(pos.X, pos.Y)
			debug.AddIntCoords(fmt.Sprintf("Player %d Coords", i+1), cx, cy)
			debug.AddText(fmt.Sprintf("Player %d Score: %d", i+1, data.CurrLevelSess.PlayerStats[i].Score))
			debug.AddText(fmt.Sprintf("Player %d Deaths: %d", i+1, data.CurrLevelSess.PlayerStats[i].Deaths))
			debug.AddText(fmt.Sprintf("Player %d State: %s", i+1, player.State.String()))
			if player.Inventory == nil {
				debug.AddText(fmt.Sprintf("Player %d Inv: Empty", i+1))
			} else {
				item := "unknown"
				cd, ok1 := player.Inventory.GetComponentData(myecs.PickUp)
				if ok1 {
					if pd, ok := cd.(*data.PickUp); ok {
						item = pd.Name
					}
				}
				debug.AddText(fmt.Sprintf("Player %d Inv: %s", i+1, item))
			}
			//debug.AddText(fmt.Sprintf("Player %d # of Tiles: %d", i+1, len(player.StoredBlocks)))
		}
	}

	if reanimator.FRate != constants.FrameRate {
		reanimator.SetFrameRate(constants.FrameRate)
	}
	reanimator.Update()

	// function systems
	systems.PlaySystem()
	systems.FunctionSystem()
	systems.InterpolationSystem()
	ui.DialogStackOpen = len(ui.DialogStack) > 0
	systems.DialogSystem(win)

	if !ui.DialogStackOpen {
		// custom systems
		systems.InGameSystem()
		systems.CharacterActionSystem()
		systems.DynamicSystem()
		systems.CollisionSystem()
		systems.CharacterStateSystem()
		systems.TouchSystem()
		systems.SmashSystem()
		systems.TileSystem()
		systems.MagicSystem()
		systems.TileSpriteSystemPre()
		systems.TileSpriteSystem()
		systems.FloatingTextSystem()
	} else {

	}
	// object systems
	systems.AnimationSystem()
	systems.ParentSystem()
	systems.ObjectSystem()
	systems.ShaderSystem()

	data.BorderView.Update()
	data.PuzzleView.Update()
	data.WorldView.Update()
	data.PuzzleViewNoShader.Update()

	myecs.UpdateManager()
	debug.AddText(fmt.Sprintf("Entity Count: %d", myecs.FullCount))
}

func (s *playState) Draw(win *pixelgl.Window) {
	drawPlayArea(win)
}

func drawPlayArea(win *pixelgl.Window) {
	// draw border
	data.BorderView.Canvas.Clear(constants.ColorBlack)
	systems.DrawBorder(ui.PuzzleBorderObject, ui.PuzzleBorder, data.BorderView.Canvas)
	img.Clear()
	data.BorderView.Draw(win)
	// draw puzzle
	data.WorldView.Canvas.Clear(pixel.RGBA{})
	data.PuzzleView.Canvas.Clear(constants.ColorBlack)
	systems.DrawBatchSystem(data.PuzzleView.Canvas, constants.TileBatch, constants.DrawingLayers)
	img.Clear()
	//data.PuzzleView.Draw(win)
	data.PuzzleView.Draw(data.WorldView.Canvas)
	// draw collapse/regen
	data.PuzzleView.Canvas.Clear(pixel.RGBA{})
	systems.DrawBatchSystem(data.PuzzleView.Canvas, constants.TileBatch, constants.CollapseRegenLayer)
	data.PuzzleView.Canvas.SetComposeMethod(pixel.ComposeRatop)
	systems.DrawBatchSystem(data.PuzzleView.Canvas, constants.TileBatch, constants.CollapseRegenMask)
	data.PuzzleView.Canvas.SetComposeMethod(pixel.ComposeOver)
	//data.PuzzleView.Draw(win)
	data.PuzzleView.Draw(data.WorldView.Canvas)
	// draw effects
	data.PuzzleView.Canvas.Clear(pixel.RGBA{})
	systems.DrawBatchSystem(data.PuzzleView.Canvas, constants.TileBatch, constants.EffectsLayer)
	//data.PuzzleView.Draw(win)
	data.PuzzleView.Draw(data.WorldView.Canvas)
	data.PuzzleViewNoShader.Canvas.Clear(pixel.RGBA{})
	systems.DrawLayerSystem(data.PuzzleViewNoShader.Canvas, 36)
	systems.DrawLayerSystem(data.PuzzleViewNoShader.Canvas, 37)
	// draw debug
	if debug.ShowDebug {
		debug.DrawLines(data.PuzzleViewNoShader.Canvas)
	}
	//data.PuzzleViewNoShader.Draw(data.WorldView.Canvas)
	data.WorldView.Draw(win)
	data.PuzzleViewNoShader.Draw(win)
	// dialog draw system
	systems.DialogDrawSystem(win)
	systems.DrawLayerSystem(win, -10)
	img.Clear()
	systems.TemporarySystem()
	//data.IMDraw.Clear()
	if options.Updated {
		systems.UpdateViews()
	}
}

func (s *playState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}
