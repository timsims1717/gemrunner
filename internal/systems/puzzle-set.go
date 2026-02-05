package systems

import (
	"fmt"
	"gemrunner/internal/content"
	"gemrunner/internal/data"
	"github.com/pkg/errors"
)

func NewPuzzleSet() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet = data.CreatePuzzleSet()
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
}

func AddPuzzle() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.AppendNew()
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
	data.CurrPuzzleSet.NeedToSave = true
}

func InsertPuzzle() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.Insert(nil, data.CurrPuzzleSet.PuzzleIndex+1)
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
	data.CurrPuzzleSet.NeedToSave = true
}

func AddPuzzleUp() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	grid := data.CurrPuzzleSet.CurrPuzzle.Grid
	grid.Y++
	data.CurrPuzzleSet.InsertGrid(nil, grid)
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
	data.CurrPuzzleSet.NeedToSave = true
}

func AddPuzzleDown() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	grid := data.CurrPuzzleSet.CurrPuzzle.Grid
	grid.Y--
	data.CurrPuzzleSet.InsertGrid(nil, grid)
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
	data.CurrPuzzleSet.NeedToSave = true
}

func AddPuzzleRight() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	grid := data.CurrPuzzleSet.CurrPuzzle.Grid
	grid.X++
	data.CurrPuzzleSet.InsertGrid(nil, grid)
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
	data.CurrPuzzleSet.NeedToSave = true
}

func AddPuzzleLeft() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	grid := data.CurrPuzzleSet.CurrPuzzle.Grid
	grid.X--
	data.CurrPuzzleSet.InsertGrid(nil, grid)
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
	data.CurrPuzzleSet.NeedToSave = true
}

func PrevPuzzle() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.Prev()
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
	data.CurrPuzzleSet.NeedToSave = true
}

func NextPuzzle() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.Next()
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
	data.CurrPuzzleSet.NeedToSave = true
}

func UpPuzzle() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.Up()
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
	data.CurrPuzzleSet.NeedToSave = true
}

func DownPuzzle() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.Down()
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
	data.CurrPuzzleSet.NeedToSave = true
}

func RightPuzzle() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.Right()
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
	data.CurrPuzzleSet.NeedToSave = true
}

func LeftPuzzle() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.Left()
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
	data.CurrPuzzleSet.NeedToSave = true
}

func DeletePuzzle() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.Delete(data.CurrPuzzleSet.PuzzleIndex)
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
	data.CurrPuzzleSet.NeedToSave = true
}

func SavePuzzleSet() bool {
	if data.CurrPuzzleSet != nil {
		data.CurrPuzzleSet.LastEditedPuzzle = data.CurrPuzzleSet.PuzzleIndex
		if data.CurrPuzzleSet.Metadata.Name == "" {
			fmt.Println("ERROR: puzzle set has no name")
			return false
		}
		if data.CurrPuzzleSet.Metadata.Filename == "" {
			data.CurrPuzzleSet.Metadata.Filename = fmt.Sprintf("%s.puzzle", data.CurrPuzzleSet.Metadata.Name)
		}
		err := content.SavePuzzleSetToFile()
		if err != nil {
			fmt.Println("ERROR:", err)
			return false
		}
		data.CurrPuzzleSet.NeedToSave = false
		for _, pzl := range data.CurrPuzzleSet.Puzzles {
			pzl.Changed = false
		}
		return true
	} else {
		fmt.Println("ERROR: no puzzle set to save")
		return false
	}
}

func OpenPuzzleSet(filename string) error {
	err := content.OpenPuzzleSetFile(filename)
	if err != nil {
		return err
	}
	data.CurrPuzzleSet.SetToFirst()
	if data.CurrPuzzleSet.Metadata.NumPlayers < 1 {
		data.CurrPuzzleSet.Metadata.NumPlayers = data.CurrPuzzleSet.CurrPuzzle.NumPlayers()
	}
	return nil
}

func CombinePuzzleSet(filename string) error {
	pzlSet, err := content.OpenPuzzleSetFileRt(filename)
	if err != nil {
		return err
	} else if pzlSet == nil {
		return errors.New("no puzzle set to combine")
	}
	oIndex := data.CurrPuzzleSet.PuzzleIndex + 1
	for _, pzl := range pzlSet.Puzzles {
		data.CurrPuzzleSet.Insert(pzl, data.CurrPuzzleSet.PuzzleIndex+1)
	}
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.SetTo(oIndex)
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
	return nil
}
