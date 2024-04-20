package data

import (
	"gemrunner/pkg/timing"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"strings"
)

var (
	Editor *editor
)

type editor struct {
	CurrBlock Block
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
	PosTop      bool

	BlockSelect *viewport.ViewPort
}

func NewEditor() {
	Editor = &editor{
		LastCoords: world.Coords{X: -1, Y: -1},
	}
}

func BlockSelectPlacement(b, w, h int) pixel.Vec {
	wo := float64(w) / 2
	if w%2 == 0 {
		wo -= 0.5
	}
	ho := float64(h) / 2
	if w%2 == 0 {
		ho -= 0.5
	}
	return pixel.V(world.TileSize*(float64(b%w)-wo), world.TileSize*(-float64(b/w)+ho))
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
	Wrench
	Wire
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
	case Wrench:
		return "Wrench"
	case Wire:
		return "Wire"
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

func ModeFromSprString(s string) EditorMode {
	ss := strings.Split(s, "_")
	ms := ss[0]
	switch ms {
	case "brush":
		return Brush
	case "line":
		return Line
	case "square":
		return Square
	case "fill":
		return Fill
	case "erase":
		return Erase
	case "eyedrop":
		return Eyedrop
	case "select":
		return Select
	case "move":
		return Move
	case "cut":
		return Cut
	case "copy":
		return Copy
	case "paste":
		return Paste
	case "flipv":
		return FlipVertical
	case "fliph":
		return FlipHorizontal
	case "wrench":
		return Wrench
	case "wire":
		return Wire
	case "undo":
		return Undo
	case "redo":
		return Redo
	}
	return EndModeList
}
