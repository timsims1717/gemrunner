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
	Cycle     [constants.MaxPlayers]int
	Priority  int
	Inventory int
	Color     ItemColor
}

func NewPickUp(p int, color ItemColor) *PickUp {
	return &PickUp{
		Priority:  p,
		Inventory: -1,
		Color:     color,
	}
}

type BasicItem struct {
	Name     string
	Object   *object.Object
	Entity   *ecs.Entity
	Sprite   *img.Sprite
	Anim     *reanimator.Tree
	PickUp   *PickUp
	Action   *Interact
	Color    ItemColor
	Metadata TileMetadata
	Origin   world.Coords
	Using    bool
	Regen    bool
	Waiting  bool
	Counter  int
	Uses     int
}

type Door struct {
	Item     *BasicItem
	DoorType DoorType
	Unlock   bool
}

type Transporter struct {
	Item  *BasicItem
	Dest  *Tile
	BarE  *ecs.Entity
	BarO  *object.Object
	BarT  *reanimator.Tree
	BarUp bool
}

type ItemColor int

const (
	ColorDefault = iota
	NonPlayerYellow
	NonPlayerBrown
	NonPlayerGray
	NonPlayerCyan
	NonPlayerRed
	PlayerBlue
	PlayerGreen
	PlayerPurple
	PlayerOrange
)

func (ic ItemColor) String() string {
	switch ic {
	case ColorDefault:
		return "default"
	case NonPlayerYellow:
		return "yellow"
	case NonPlayerBrown:
		return "brown"
	case NonPlayerGray:
		return "gray"
	case NonPlayerCyan:
		return "cyan"
	case PlayerBlue:
		return "blue"
	case PlayerGreen:
		return "green"
	case PlayerPurple:
		return "purple"
	case PlayerOrange:
		return "orange"
	case NonPlayerRed:
		return "red"
	}
	return ""
}

func (ic ItemColor) SpriteString() string {
	switch ic {
	case NonPlayerYellow:
		return "_yellow"
	case NonPlayerBrown:
		return "_brown"
	case NonPlayerGray:
		return "_gray"
	case NonPlayerCyan:
		return "_cyan"
	case PlayerBlue:
		return "_blue"
	case PlayerGreen:
		return "_green"
	case PlayerPurple:
		return "_purple"
	case PlayerOrange:
		return "_orange"
	case NonPlayerRed:
		return "_red"
	default:
		return "_yellow"
	}
}

type DoorType int

const (
	Hidden = iota
	Visible
	Locked
	Unlocked
)

type Bomb struct {
	Item   *BasicItem
	Draws  []interface{}
	SymSpr *img.Sprite
	LitKey string
}

type Disguise struct {
	Item *BasicItem
	Doff bool
}
