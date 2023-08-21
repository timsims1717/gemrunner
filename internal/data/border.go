package data

import "github.com/bytearena/ecs"

var MainBorder *ecs.Entity

type Border struct {
	Width  int
	Height int
	Empty  bool
}
