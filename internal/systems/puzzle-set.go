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
}

func InsertPuzzle() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.Insert(nil, data.CurrPuzzleSet.PuzzleIndex+1)
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
}

func AddPuzzleUp() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	grid := data.CurrPuzzleSet.CurrPuzzle.Grid
	grid.Y++
	data.CurrPuzzleSet.InsertGrid(nil, grid)
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
}

func AddPuzzleDown() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	grid := data.CurrPuzzleSet.CurrPuzzle.Grid
	grid.Y--
	data.CurrPuzzleSet.InsertGrid(nil, grid)
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
}

func AddPuzzleRight() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	grid := data.CurrPuzzleSet.CurrPuzzle.Grid
	grid.X++
	data.CurrPuzzleSet.InsertGrid(nil, grid)
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
}

func AddPuzzleLeft() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	grid := data.CurrPuzzleSet.CurrPuzzle.Grid
	grid.X--
	data.CurrPuzzleSet.InsertGrid(nil, grid)
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
}

func PrevPuzzle() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.Prev()
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
}

func NextPuzzle() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.Next()
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
}

func UpPuzzle() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.Up()
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
}

func DownPuzzle() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.Down()
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
}

func RightPuzzle() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.Right()
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
}

func LeftPuzzle() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.Left()
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
}

func DeletePuzzle() {
	DisposePuzzle(data.CurrentPlayArea.Puzzle)
	data.CurrPuzzleSet.Delete(data.CurrPuzzleSet.PuzzleIndex)
	SetPuzzle(data.CurrentPlayArea, data.CurrPuzzleSet.CurrPuzzle)
	InitPuzzle(data.CurrentPlayArea)
}

func SavePuzzleSet() bool {
	if data.CurrPuzzleSet != nil {
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
