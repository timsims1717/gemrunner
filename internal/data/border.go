package data

import (
	"gemrunner/pkg/object"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
)

var MainBorder *ecs.Entity

var (
	PuzzleBorder       *Border
	PuzzleBorderObject *object.Object
)

type Border struct {
	Width  int
	Height int
	Rect   pixel.Rect
	Empty  bool
	Style  BorderStyle
	Hidden bool

	ExcludeSide Direction
	ExcludeSize int
}

type BorderStyle int

const (
	FancyBorder = iota
	ThinBorder
	ThinBorderReverse
	ThinBorderWhite
	ThinBorderBlue
	ThickBorder
	ThickBorderReverse
	ThickBorderWhite
	ThickBorderBlue
)
