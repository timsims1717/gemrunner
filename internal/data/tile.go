package data

import (
	"gemrunner/internal/constants"
	"gemrunner/pkg/object"
	"gemrunner/pkg/world"
)

type Block int

const (
	Turf = iota
	Fall
	Ladder
	DoorPink
	LockPink
	DoorBlue
	LockBlue
	Player1
	Devil
	Box
	KeyPink
	KeyBlue
	Gem
	Chain
	Reeds
	Flowers
	Empty
)

func (b Block) String() string {
	switch b {
	case Turf, Fall:
		if CurrPuzzle != nil && CurrPuzzle.World != "" {
			return CurrPuzzle.World
		}
		return constants.WorldRock
	case Ladder:
		return constants.TileLadderMiddle
	case DoorPink:
		return constants.TileDoorPink
	case LockPink:
		return constants.TileLockPink
	case DoorBlue:
		return constants.TileDoorBlue
	case LockBlue:
		return constants.TileLockBlue
	case Player1:
		return constants.CharPlayer1
	case Devil:
		return constants.CharDevil
	case Box:
		return constants.ItemBox
	case KeyPink:
		return constants.ItemKeyPink
	case KeyBlue:
		return constants.ItemKeyBlue
	case Gem:
		return constants.ItemGem
	case Chain:
		return constants.DoodadChain
	case Reeds:
		return constants.DoodadReeds
	case Flowers:
		return constants.DoodadFlowers
	}
	return "empty"
}

type Tile struct {
	Block  Block
	Ladder bool
	Coords world.Coords
	Object *object.Object
	Update bool
}

func (t *Tile) Copy() *Tile {
	return &Tile{
		Block:  t.Block,
		Ladder: t.Ladder,
		Coords: t.Coords,
	}
}

func (t *Tile) CopyInto(c *Tile) {
	c.Block = t.Block
	c.Ladder = t.Ladder
}

func (t *Tile) Empty() {
	t.Block = Empty
	t.Ladder = false
}
