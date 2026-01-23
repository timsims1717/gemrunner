package systems

import "gemrunner/internal/data"

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
