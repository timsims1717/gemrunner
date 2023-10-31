package data

import (
	"bytes"
	"encoding/json"
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

var toID = map[string]Block{
	constants.TileTurf:         Turf,
	constants.TileFall:         Fall,
	constants.TileLadderMiddle: Ladder,
	constants.TileDoorPink:     DoorPink,
	constants.TileLockPink:     LockPink,
	constants.TileDoorBlue:     DoorBlue,
	constants.TileLockBlue:     LockBlue,
	constants.CharPlayer1:      Player1,
	constants.CharDevil:        Devil,
	constants.ItemBox:          Box,
	constants.ItemKeyPink:      KeyPink,
	constants.ItemKeyBlue:      KeyBlue,
	constants.ItemGem:          Gem,
	constants.DoodadChain:      Chain,
	constants.DoodadReeds:      Reeds,
	constants.DoodadFlowers:    Flowers,
	constants.TileEmpty:        Empty,
}

func (b Block) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	switch b {
	case Turf:
		buffer.WriteString(constants.TileTurf)
	case Fall:
		buffer.WriteString(constants.TileFall)
	default:
		buffer.WriteString(b.String())
	}
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (b *Block) UnmarshalJSON(bts []byte) error {
	var j string
	err := json.Unmarshal(bts, &j)
	if err != nil {
		return err
	}
	*b = toID[j]
	return nil
}

type Tile struct {
	Block  Block          `json:"tile"`
	Ladder bool           `json:"ladder"`
	Coords world.Coords   `json:"-"`
	Object *object.Object `json:"-"`
	Update bool           `json:"-"`
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
