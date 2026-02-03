package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/gween64/ease"
	"gemrunner/pkg/object"
	"math"
)

func StartInterpolation() {
	oldGrid := data.CurrentPlayArea.Level.Puzzle.Grid
	newGrid := data.OtherPlayArea.Level.Puzzle.Grid
	xMove := math.Abs((data.CurrentPlayArea.WorldView.Canvas.Bounds().W() + data.OtherPlayArea.WorldView.Canvas.Bounds().W()) * 0.5)
	yMove := math.Abs((data.CurrentPlayArea.WorldView.Canvas.Bounds().H() + data.OtherPlayArea.WorldView.Canvas.Bounds().H()) * 0.5)
	var interA, interB, interAB, interBB *object.Interpolation
	if oldGrid.X < newGrid.X { // going right, moving play areas left
		data.OtherPlayArea.WorldView.PortPos.X = data.CurrentPlayArea.WorldView.PortPos.X + xMove
		data.OtherPlayArea.BorderView.PortPos.X = data.CurrentPlayArea.BorderView.PortPos.X + xMove
		interA = object.NewInterpolation(object.InterpolateCustom).
			SetValue(&data.CurrentPlayArea.WorldView.PortPos.X).
			SetGween(data.CurrentPlayArea.WorldView.PortPos.X, data.CurrentPlayArea.WorldView.PortPos.X-xMove, constants.LevelTransSpeed, ease.InOutQuad)
		interAB = object.NewInterpolation(object.InterpolateCustom).
			SetValue(&data.CurrentPlayArea.BorderView.PortPos.X).
			SetGween(data.CurrentPlayArea.BorderView.PortPos.X, data.CurrentPlayArea.BorderView.PortPos.X-xMove, constants.LevelTransSpeed, ease.InOutQuad)
		interB = object.NewInterpolation(object.InterpolateCustom).
			SetValue(&data.OtherPlayArea.WorldView.PortPos.X).
			SetGween(data.OtherPlayArea.WorldView.PortPos.X, data.OtherPlayArea.WorldView.PortPos.X-xMove, constants.LevelTransSpeed, ease.InOutQuad)
		interBB = object.NewInterpolation(object.InterpolateCustom).
			SetValue(&data.OtherPlayArea.BorderView.PortPos.X).
			SetGween(data.OtherPlayArea.BorderView.PortPos.X, data.OtherPlayArea.BorderView.PortPos.X-xMove, constants.LevelTransSpeed, ease.InOutQuad)
	} else if oldGrid.X > newGrid.X { // going left, moving play areas right
		data.OtherPlayArea.WorldView.PortPos.X = data.CurrentPlayArea.WorldView.PortPos.X - xMove
		data.OtherPlayArea.BorderView.PortPos.X = data.CurrentPlayArea.BorderView.PortPos.X - xMove
		interA = object.NewInterpolation(object.InterpolateCustom).
			SetValue(&data.CurrentPlayArea.WorldView.PortPos.X).
			SetGween(data.CurrentPlayArea.WorldView.PortPos.X, data.CurrentPlayArea.WorldView.PortPos.X+xMove, constants.LevelTransSpeed, ease.InOutQuad)
		interAB = object.NewInterpolation(object.InterpolateCustom).
			SetValue(&data.CurrentPlayArea.BorderView.PortPos.X).
			SetGween(data.CurrentPlayArea.BorderView.PortPos.X, data.CurrentPlayArea.BorderView.PortPos.X+xMove, constants.LevelTransSpeed, ease.InOutQuad)
		interB = object.NewInterpolation(object.InterpolateCustom).
			SetValue(&data.OtherPlayArea.WorldView.PortPos.X).
			SetGween(data.OtherPlayArea.WorldView.PortPos.X, data.OtherPlayArea.WorldView.PortPos.X+xMove, constants.LevelTransSpeed, ease.InOutQuad)
		interBB = object.NewInterpolation(object.InterpolateCustom).
			SetValue(&data.OtherPlayArea.BorderView.PortPos.X).
			SetGween(data.OtherPlayArea.BorderView.PortPos.X, data.OtherPlayArea.BorderView.PortPos.X+xMove, constants.LevelTransSpeed, ease.InOutQuad)
	} else if oldGrid.Y < newGrid.Y { // going up, moving play areas down
		data.OtherPlayArea.WorldView.PortPos.Y = data.CurrentPlayArea.WorldView.PortPos.Y + yMove
		data.OtherPlayArea.BorderView.PortPos.Y = data.CurrentPlayArea.BorderView.PortPos.Y + yMove
		interA = object.NewInterpolation(object.InterpolateCustom).
			SetValue(&data.CurrentPlayArea.WorldView.PortPos.Y).
			SetGween(data.CurrentPlayArea.WorldView.PortPos.Y, data.CurrentPlayArea.WorldView.PortPos.Y-yMove, constants.LevelTransSpeed, ease.InOutQuad)
		interAB = object.NewInterpolation(object.InterpolateCustom).
			SetValue(&data.CurrentPlayArea.BorderView.PortPos.Y).
			SetGween(data.CurrentPlayArea.BorderView.PortPos.Y, data.CurrentPlayArea.BorderView.PortPos.Y-yMove, constants.LevelTransSpeed, ease.InOutQuad)
		interB = object.NewInterpolation(object.InterpolateCustom).
			SetValue(&data.OtherPlayArea.WorldView.PortPos.Y).
			SetGween(data.OtherPlayArea.WorldView.PortPos.Y, data.OtherPlayArea.WorldView.PortPos.Y-yMove, constants.LevelTransSpeed, ease.InOutQuad)
		interBB = object.NewInterpolation(object.InterpolateCustom).
			SetValue(&data.OtherPlayArea.BorderView.PortPos.Y).
			SetGween(data.OtherPlayArea.BorderView.PortPos.Y, data.OtherPlayArea.BorderView.PortPos.Y-yMove, constants.LevelTransSpeed, ease.InOutQuad)
	} else if oldGrid.Y > newGrid.Y { // going down, moving play areas up
		data.OtherPlayArea.WorldView.PortPos.Y = data.CurrentPlayArea.WorldView.PortPos.Y - yMove
		data.OtherPlayArea.BorderView.PortPos.Y = data.CurrentPlayArea.BorderView.PortPos.Y - yMove
		interA = object.NewInterpolation(object.InterpolateCustom).
			SetValue(&data.CurrentPlayArea.WorldView.PortPos.Y).
			SetGween(data.CurrentPlayArea.WorldView.PortPos.Y, data.CurrentPlayArea.WorldView.PortPos.Y+yMove, constants.LevelTransSpeed, ease.InOutQuad)
		interAB = object.NewInterpolation(object.InterpolateCustom).
			SetValue(&data.CurrentPlayArea.BorderView.PortPos.Y).
			SetGween(data.CurrentPlayArea.BorderView.PortPos.Y, data.CurrentPlayArea.BorderView.PortPos.Y+yMove, constants.LevelTransSpeed, ease.InOutQuad)
		interB = object.NewInterpolation(object.InterpolateCustom).
			SetValue(&data.OtherPlayArea.WorldView.PortPos.Y).
			SetGween(data.OtherPlayArea.WorldView.PortPos.Y, data.OtherPlayArea.WorldView.PortPos.Y+yMove, constants.LevelTransSpeed, ease.InOutQuad)
		interBB = object.NewInterpolation(object.InterpolateCustom).
			SetValue(&data.OtherPlayArea.BorderView.PortPos.Y).
			SetGween(data.OtherPlayArea.BorderView.PortPos.Y, data.OtherPlayArea.BorderView.PortPos.Y+yMove, constants.LevelTransSpeed, ease.InOutQuad)
	}
	data.CurrentPlayArea.BorderEntity.
		AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interAB})
	data.OtherPlayArea.BorderEntity.
		AddComponent(myecs.Interpolation, []*object.Interpolation{interB, interBB})
}

func SetTransitionBorders() {
	oldGrid := data.CurrentPlayArea.Level.Puzzle.Grid
	newGrid := data.OtherPlayArea.Level.Puzzle.Grid
	if oldGrid.X < newGrid.X { // going right, moving play areas left
		data.CurrentPlayArea.Border.ExcludeSide = data.Right
		data.CurrentPlayArea.Border.ExcludeSize = data.OtherPlayArea.Puzzle.Metadata.Height + 2
		data.OtherPlayArea.Border.ExcludeSide = data.Left
		data.OtherPlayArea.Border.ExcludeSize = data.CurrentPlayArea.Puzzle.Metadata.Height + 2
	} else if oldGrid.X > newGrid.X { // going left, moving play areas right
		data.CurrentPlayArea.Border.ExcludeSide = data.Left
		data.CurrentPlayArea.Border.ExcludeSize = data.OtherPlayArea.Puzzle.Metadata.Height + 2
		data.OtherPlayArea.Border.ExcludeSide = data.Right
		data.OtherPlayArea.Border.ExcludeSize = data.CurrentPlayArea.Puzzle.Metadata.Height + 2
	} else if oldGrid.Y < newGrid.Y { // going up, moving play areas down
		data.CurrentPlayArea.Border.ExcludeSide = data.Up
		data.CurrentPlayArea.Border.ExcludeSize = data.OtherPlayArea.Puzzle.Metadata.Width + 2
		data.OtherPlayArea.Border.ExcludeSide = data.Down
		data.OtherPlayArea.Border.ExcludeSize = data.CurrentPlayArea.Puzzle.Metadata.Width + 2
	} else if oldGrid.Y > newGrid.Y { // going down, moving play areas up
		data.CurrentPlayArea.Border.ExcludeSide = data.Down
		data.CurrentPlayArea.Border.ExcludeSize = data.OtherPlayArea.Puzzle.Metadata.Width + 2
		data.OtherPlayArea.Border.ExcludeSide = data.Up
		data.OtherPlayArea.Border.ExcludeSize = data.CurrentPlayArea.Puzzle.Metadata.Width + 2
	}
}
