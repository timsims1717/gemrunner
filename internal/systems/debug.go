package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/world"
)

func GoToLevelLeft() {
	if data.CurrLevel != nil {
		g := data.CurrLevel.Puzzle.Grid
		g.X--
		if lvl := data.CurrLevelSess.PuzzleSet.GetGridPuzzle(g); lvl != nil {
			hi := data.CurrLevelSess.PuzzleSet.GetGrid(g)
			data.CurrLevel.Complete = true
			data.CurrLevel.ExitIndex = hi
			data.CurrLevel.StartCoords = nil
		}
	}
}

func GoToLevelRight() {
	if data.CurrLevel != nil {
		g := data.CurrLevel.Puzzle.Grid
		g.X++
		if lvl := data.CurrLevelSess.PuzzleSet.GetGridPuzzle(g); lvl != nil {
			hi := data.CurrLevelSess.PuzzleSet.GetGrid(g)
			data.CurrLevel.Complete = true
			data.CurrLevel.ExitIndex = hi
			data.CurrLevel.StartCoords = nil
		}
	}
}

func GoToLevelUp() {
	if data.CurrLevel != nil {
		g := data.CurrLevel.Puzzle.Grid
		g.Y++
		if lvl := data.CurrLevelSess.PuzzleSet.GetGridPuzzle(g); lvl != nil {
			hi := data.CurrLevelSess.PuzzleSet.GetGrid(g)
			data.CurrLevel.Complete = true
			data.CurrLevel.ExitIndex = hi
			data.CurrLevel.StartCoords = nil
		}
	}
}

func GoToLevelDown() {
	if data.CurrLevel != nil {
		g := data.CurrLevel.Puzzle.Grid
		g.Y--
		if lvl := data.CurrLevelSess.PuzzleSet.GetGridPuzzle(g); lvl != nil {
			hi := data.CurrLevelSess.PuzzleSet.GetGrid(g)
			data.CurrLevel.Complete = true
			data.CurrLevel.ExitIndex = hi
			data.CurrLevel.StartCoords = nil
		}
	}
}

func OpenDoors() {
	data.CurrLevel.DoorsOpen = true
	data.CurrLevel.Continuity = data.CurrLevelSess.PuzzleSet.Metadata.Continuity != data.NoContinuity
}

func InGameDebugInfo() {
	debug.AddText(fmt.Sprintf("Speed: %d", constants.Configuration.Gameplay.FrameRate))
	debug.AddText(FormatTimePlayed())
	debug.AddText(fmt.Sprintf("Frame Number: %d", data.CurrLevel.FrameNumber))
	debug.AddText(fmt.Sprintf("Frame Counter: %d", data.CurrLevel.FrameCounter))
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
				item := player.Inventory.Name
				debug.AddText(fmt.Sprintf("Player %d Inv: %s", i+1, item))
			}
			//debug.AddText(fmt.Sprintf("Player %d # of Tiles: %d", i+1, len(player.StoredBlocks)))
		}
	}
}

func InGameDebugInput() {
	if data.DebugInput.Get("debugTest").JustPressed() {
		for _, player := range data.CurrLevel.Players {
			if player != nil {
				player.SmallBombs++
			}
		}
		data.DebugInput.Get("debugTest").Consume()
	}
	if data.DebugInput.Get("debugInv").JustPressed() {
		OpenDoors()
		data.DebugInput.Get("debugInv").Consume()
	}
}

func PlayDebugInput() {
	if data.DebugInput.Get("ctrl").Pressed() {
		if data.DebugInput.Get("debugLevelUp").JustPressed() {
			GoToLevelUp()
			data.DebugInput.Get("debugLevelUp").Consume()
		} else if data.DebugInput.Get("debugLevelDown").JustPressed() {
			GoToLevelDown()
			data.DebugInput.Get("debugLevelDown").Consume()
		} else if data.DebugInput.Get("debugLevelLeft").JustPressed() {
			GoToLevelLeft()
			data.DebugInput.Get("debugLevelLeft").Consume()
		} else if data.DebugInput.Get("debugLevelRight").JustPressed() {
			GoToLevelRight()
			data.DebugInput.Get("debugLevelRight").Consume()
		}
	}
}

func TransitionDebugInput() {
	if data.DebugInput.Get("debugLevelUp").Pressed() {
		data.OtherPlayArea.BorderView.PortPos.Y += 200. * timing.DT
	} else if data.DebugInput.Get("debugLevelDown").Pressed() {
		data.OtherPlayArea.BorderView.PortPos.Y -= 200. * timing.DT
	} else if data.DebugInput.Get("debugLevelLeft").Pressed() {
		data.OtherPlayArea.BorderView.PortPos.X -= 200. * timing.DT
	} else if data.DebugInput.Get("debugLevelRight").Pressed() {
		data.OtherPlayArea.BorderView.PortPos.X += 200. * timing.DT
	}
}
