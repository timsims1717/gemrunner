package data

import (
	"gemrunner/internal/constants"
	"gemrunner/pkg/object"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

var (
	EditorPanel *editorPane
)

type editorPane struct {
	ViewPort  *viewport.ViewPort
	Entity    *ecs.Entity
	CurrBlock Block
	SelectObj *object.Object
	Offset    pixel.Vec
	Mode      EditorMode
	LastMode  EditorMode
	NoInput   bool

	Hover       bool
	Consume     string
	SelectVis   bool
	SelectTimer *timing.Timer
	SelectQuick bool
	LastCoords  world.Coords

	BlockView   *BlockView
	BlockSelect *viewport.ViewPort
	UndoStack   []*Puzzle
	LastChange  *Puzzle
	RedoStack   []*Puzzle
}

func NewEditorPane() {
	EditorPanel = &editorPane{
		LastCoords: world.Coords{X: -1, Y: -1},
	}
}

type BlockView struct {
	Entity *ecs.Entity
	Object *object.Object
}

func BlockSelectPlacement(b int) pixel.Vec {
	return pixel.V(world.TileSize*float64(b%constants.BlockSelectWidth)+world.TileSize*0.5, -world.TileSize*float64(b/constants.BlockSelectWidth)-world.TileSize*0.5)
}

type EditorMode int

const (
	Brush = iota
	Line
	Square
	Fill
	Erase
	Eyedrop
	Select
	Move
	Cut
	Copy
	Paste
	FlipVertical
	FlipHorizontal
	Pliers
	Undo
	Redo
	Delete
	Save
	Open
	EndModeList
)

func (m EditorMode) String() string {
	switch m {
	case Brush:
		return "Brush"
	case Line:
		return "Line"
	case Square:
		return "Square"
	case Fill:
		return "Fill"
	case Erase:
		return "Erase"
	case Eyedrop:
		return "Eyedrop"
	case Select:
		return "Select"
	case Move:
		return "Move"
	case Cut:
		return "Cut"
	case Copy:
		return "Copy"
	case Paste:
		return "Paste"
	case FlipVertical:
		return "FlipVertical"
	case FlipHorizontal:
		return "FlipHorizontal"
	case Pliers:
		return "Pliers"
	case Undo:
		return "Undo"
	case Redo:
		return "Redo"
	case Delete:
		return "Delete"
	case Save:
		return "Save"
	case Open:
		return "Open"
	}
	return ""
}
