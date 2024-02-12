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
	Cycle     [constants.MaxPlayers]int
	Priority  int
	NeverFlip bool
}

func NewPickUp(p int, neverFlip bool) *PickUp {
	return &PickUp{
		Priority:  p,
		NeverFlip: neverFlip,
	}
}

type Key struct {
	Object  *object.Object
	Entity  *ecs.Entity
	Sprite  *img.Sprite
	PickUp  *PickUp
	Action  *Interact
	KeyType KeyType
}

type Door struct {
	Object   *object.Object
	Entity   *ecs.Entity
	DoorType DoorType
	Unlock   bool
}

type KeyType int

const (
	PinkKey = iota
	BlueKey
	KeyParity
)

type DoorType int

const (
	PinkOpen = iota
	BlueOpen
	PinkLock
	BlueLock
	PinkUnlock
	BlueUnlock
)
