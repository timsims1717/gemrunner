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

var EditorPanel *editorPane

type editorPane struct {
	ViewPort  *viewport.ViewPort
	Entity    *ecs.Entity
	CurrBlock Block
	SelectObj *object.Object
	Offset    pixel.Vec

	Hover       bool
	Consume     string
	SelectVis   bool
	SelectTimer *timing.Timer
	SelectQuick bool

	BlockView   *BlockView
	BlockSelect *viewport.ViewPort
}

func NewEditorPane() {
	EditorPanel = &editorPane{}
}

type BlockView struct {
	Entity *ecs.Entity
	Object *object.Object
}

func BlockSelectPlacement(b int) pixel.Vec {
	return pixel.V(world.TileSize*float64(b%constants.BlockSelectWidth)+world.TileSize*0.5, -world.TileSize*float64(b/constants.BlockSelectWidth)-world.TileSize*0.5)
}
