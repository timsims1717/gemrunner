package data

import (
	"gemrunner/internal/constants"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"github.com/bytearena/ecs"
)

type Interact struct {
	Fn func(*Level, int, *Dynamic, *ecs.Entity)
}

func NewInteract(fn func(*Level, int, *Dynamic, *ecs.Entity)) *Interact {
	return &Interact{Fn: fn}
}

type PickUp struct {
	Name        string
	Cycle       [constants.MaxPlayers]int
	Priority    int
	NeverFlip   bool
	Inventory   int
	NoInventory bool
}

func NewPickUp(name string, p int, neverFlip bool) *PickUp {
	return &PickUp{
		Name:      name,
		Priority:  p,
		NeverFlip: neverFlip,
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
	Locked
	Unlocked
)
