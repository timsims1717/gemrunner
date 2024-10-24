package data

import (
	"gemrunner/internal/constants"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
)

type Interact struct {
	Fn func(int, *Dynamic, *ecs.Entity)
}

func NewInteract(fn func(int, *Dynamic, *ecs.Entity)) *Interact {
	return &Interact{Fn: fn}
}

type PickUp struct {
	Name      string
	Cycle     [constants.MaxPlayers]int
	Priority  int
	Inventory int
}

func NewPickUp(name string, p int) *PickUp {
	return &PickUp{
		Name:      name,
		Priority:  p,
		Inventory: -1,
	}
}

type Key struct {
	Object *object.Object
	Entity *ecs.Entity
	Sprite *img.Sprite
	PickUp *PickUp
	Action *Interact
	Color  string
}

type Door struct {
	Object   *object.Object
	Entity   *ecs.Entity
	Color    string
	DoorType DoorType
	Unlock   bool
}

type DoorType int

const (
	Opened = iota
	Closed
	Locked
	Unlocked
)

type Bomb struct {
	Object   *object.Object
	Entity   *ecs.Entity
	Draws    []interface{}
	Anim     *reanimator.Tree
	SymSpr   *img.Sprite
	PickUp   *PickUp
	Action   *Interact
	Name     string
	Metadata TileMetadata
	Origin   world.Coords
	LitKey   string
	Regen    bool
	Waiting  bool
	Color    string
}

type Jetpack struct {
	Object   *object.Object
	Entity   *ecs.Entity
	Anim     *reanimator.Tree
	PickUp   *PickUp
	Action   *Interact
	Name     string
	Metadata TileMetadata
	Origin   world.Coords
	Color    string
	Counter  int
	Using    bool
	Regen    bool
	Waiting  bool
}

type Disguise struct {
	Object   *object.Object
	Entity   *ecs.Entity
	Anim     *reanimator.Tree
	PickUp   *PickUp
	Action   *Interact
	Name     string
	Metadata TileMetadata
	Origin   world.Coords
	Color    string
	Counter  int
	Using    bool
	Doff     bool
	Regen    bool
	Waiting  bool
}
