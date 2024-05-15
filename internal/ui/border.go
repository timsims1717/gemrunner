package ui

import (
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
)

var MainBorder *ecs.Entity

type Border struct {
	Width  int
	Height int
	Rect   pixel.Rect
	Empty  bool
	Style  BorderStyle
}

type BorderStyle int

const (
	FancyBorder = iota
	ThinBorder
	ThinBorderWhite
	ThinBorderBlue
)
