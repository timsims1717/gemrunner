package data

import (
	"bytes"
	"encoding/json"
	"gemrunner/internal/constants"
	"gemrunner/pkg/object"
	"gemrunner/pkg/world"
	"github.com/beefsack/go-astar"
	"github.com/bytearena/ecs"
)

type Block int

const (
	BlockTurf = iota
	BlockFall
	BlockLadder
	BlockDoorPink
	BlockLockPink
	BlockDoorBlue
	BlockLockBlue
	BlockPlayer1
	BlockDemon
	BlockBox
	BlockKeyPink
	BlockKeyBlue
	BlockGem
	BlockChain
	BlockReeds
	BlockFlowers
	BlockEmpty
)

func (b Block) String() string {
	switch b {
	case BlockTurf, BlockFall:
		if CurrPuzzle != nil && CurrPuzzle.Metadata.WorldSprite != "" {
			return CurrPuzzle.Metadata.WorldSprite
		}
		return constants.WorldSprites[constants.WorldRock]
	case BlockLadder:
		return constants.TileLadderMiddle
	case BlockDoorPink:
		return constants.TileDoorPink
	case BlockLockPink:
		return constants.TileLockPink
	case BlockDoorBlue:
		return constants.TileDoorBlue
	case BlockLockBlue:
		return constants.TileLockBlue
	case BlockPlayer1:
		return constants.CharPlayer1
	case BlockDemon:
		return constants.CharDemon
	case BlockBox:
		return constants.ItemBox
	case BlockKeyPink:
		return constants.ItemKeyPink
	case BlockKeyBlue:
		return constants.ItemKeyBlue
	case BlockGem:
		return constants.ItemGem
	case BlockChain:
		return constants.DoodadChain
	case BlockReeds:
		return constants.DoodadReeds
	case BlockFlowers:
		return constants.DoodadFlowers
	}
	return "empty"
}

var toID = map[string]Block{
	constants.TileTurf:         BlockTurf,
	constants.TileFall:         BlockFall,
	constants.TileLadderMiddle: BlockLadder,
	constants.TileDoorPink:     BlockDoorPink,
	constants.TileLockPink:     BlockLockPink,
	constants.TileDoorBlue:     BlockDoorBlue,
	constants.TileLockBlue:     BlockLockBlue,
	constants.CharPlayer1:      BlockPlayer1,
	constants.CharDemon:        BlockDemon,
	constants.ItemBox:          BlockBox,
	constants.ItemKeyPink:      BlockKeyPink,
	constants.ItemKeyBlue:      BlockKeyBlue,
	constants.ItemGem:          BlockGem,
	constants.DoodadChain:      BlockChain,
	constants.DoodadReeds:      BlockReeds,
	constants.DoodadFlowers:    BlockFlowers,
	constants.TileEmpty:        BlockEmpty,
}

func (b Block) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	switch b {
	case BlockTurf:
		buffer.WriteString(constants.TileTurf)
	case BlockFall:
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
	Entity *ecs.Entity    `json:"-"`
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
	t.Block = BlockEmpty
	t.Ladder = false
}

func (t *Tile) Solid() bool {
	return !t.Ladder && (t.Block == BlockTurf || t.Block == BlockFall)
}

// a* implementation

func (t *Tile) PathNeighbors() []astar.Pather {
	if CurrLevel == nil {
		return []astar.Pather{}
	}
	var neighbors []astar.Pather
	// Down
	d := CurrLevel.Tiles.Get(t.Coords.X, t.Coords.Y-1)
	if d != nil && !d.Solid() {
		neighbors = append(neighbors, d)
	}
	notFalling := d == nil || d.Solid() || d.Ladder || t.Ladder
	// Left
	l := CurrLevel.Tiles.Get(t.Coords.X-1, t.Coords.Y)
	//lb := CurrLevel.Tiles.Get(t.Coords.X-1, t.Coords.Y-1)
	//if notFalling && l != nil && !l.Solid() &&
	//	(l.Ladder || lb == nil || lb.Solid()) {
	//	neighbors = append(neighbors, l)
	//}
	if notFalling && l != nil && !l.Solid() {
		neighbors = append(neighbors, l)
	}
	// Right
	r := CurrLevel.Tiles.Get(t.Coords.X+1, t.Coords.Y)
	//rb := CurrLevel.Tiles.Get(t.Coords.X+1, t.Coords.Y-1)
	//if notFalling && r != nil && !r.Solid() &&
	//	(r.Ladder || rb == nil || rb.Solid()) {
	//	neighbors = append(neighbors, r)
	//}
	if notFalling && r != nil && !r.Solid() {
		neighbors = append(neighbors, r)
	}
	// Up
	u := CurrLevel.Tiles.Get(t.Coords.X, t.Coords.Y+1)
	if notFalling && u != nil && !u.Solid() && t.Ladder {
		neighbors = append(neighbors, u)
	}
	return neighbors
}

func (t *Tile) PathNeighborCost(to astar.Pather) float64 {
	return 1.
}

func (t *Tile) PathEstimatedCost(to astar.Pather) float64 {
	return 1.
}
