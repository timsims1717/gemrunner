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
	BlockCracked
	BlockPhase
	BlockLadder
	BlockLadderCracked
	BlockLadderExit
	BlockBox

	BlockPlayer1
	BlockKeyBlue
	BlockPlayer2
	BlockKeyGreen
	BlockPlayer3
	BlockKeyPurple
	BlockPlayer4
	BlockKeyBrown

	BlockDoorBlue
	BlockLockBlue
	BlockDoorGreen
	BlockLockGreen
	BlockDoorPurple
	BlockLockPurple
	BlockDoorBrown
	BlockLockBrown

	BlockDoorYellow
	BlockLockYellow
	BlockDoorOrange
	BlockLockOrange
	BlockDoorGray
	BlockLockGray
	BlockDoorCyan
	BlockLockCyan

	BlockGemYellow
	BlockKeyYellow
	BlockGemOrange
	BlockKeyOrange
	BlockGemGray
	BlockKeyGray
	BlockGemCyan
	BlockKeyCyan

	BlockGemBlue
	BlockGemGreen
	BlockGemPurple
	BlockGemBrown

	BlockDemon
	BlockFly
	BlockChain
	BlockReeds
	BlockFlowers
	BlockEmpty
)

func (b Block) String() string {
	switch b {
	case BlockTurf, BlockFall, BlockCracked, BlockPhase:
		if CurrPuzzle != nil && CurrPuzzle.Metadata.WorldSprite != "" {
			return CurrPuzzle.Metadata.WorldSprite
		}
		return constants.WorldSprites[constants.WorldRock]
	case BlockLadder:
		return constants.TileLadderMiddle
	case BlockLadderCracked:
		return constants.TileLadderCrackM
	case BlockLadderExit:
		return constants.TileExitLadderM
	case BlockBox:
		return constants.ItemBox
	case BlockPlayer1:
		return constants.CharPlayer1
	case BlockPlayer2:
		return constants.CharPlayer2
	case BlockPlayer3:
		return constants.CharPlayer3
	case BlockPlayer4:
		return constants.CharPlayer4
	case BlockKeyBlue:
		return constants.ItemKeyBlue
	case BlockKeyGreen:
		return constants.ItemKeyGreen
	case BlockKeyPurple:
		return constants.ItemKeyPurple
	case BlockKeyBrown:
		return constants.ItemKeyBrown
	case BlockKeyYellow:
		return constants.ItemKeyYellow
	case BlockKeyOrange:
		return constants.ItemKeyOrange
	case BlockKeyGray:
		return constants.ItemKeyGray
	case BlockKeyCyan:
		return constants.ItemKeyCyan
	case BlockDoorBlue:
		return constants.TileDoorBlue
	case BlockDoorGreen:
		return constants.TileDoorGreen
	case BlockDoorPurple:
		return constants.TileDoorPurple
	case BlockDoorBrown:
		return constants.TileDoorBrown
	case BlockDoorYellow:
		return constants.TileDoorYellow
	case BlockDoorOrange:
		return constants.TileDoorOrange
	case BlockDoorGray:
		return constants.TileDoorGray
	case BlockDoorCyan:
		return constants.TileDoorCyan
	case BlockLockBlue:
		return constants.TileLockBlue
	case BlockLockGreen:
		return constants.TileLockGreen
	case BlockLockPurple:
		return constants.TileLockPurple
	case BlockLockBrown:
		return constants.TileLockBrown
	case BlockLockYellow:
		return constants.TileLockYellow
	case BlockLockOrange:
		return constants.TileLockOrange
	case BlockLockGray:
		return constants.TileLockGray
	case BlockLockCyan:
		return constants.TileLockCyan
	case BlockGemBlue:
		return constants.ItemGemBlue
	case BlockGemGreen:
		return constants.ItemGemGreen
	case BlockGemPurple:
		return constants.ItemGemPurple
	case BlockGemBrown:
		return constants.ItemGemBrown
	case BlockGemYellow:
		return constants.ItemGemYellow
	case BlockGemOrange:
		return constants.ItemGemOrange
	case BlockGemGray:
		return constants.ItemGemGray
	case BlockGemCyan:
		return constants.ItemGemCyan
	case BlockDemon:
		return constants.CharDemon
	case BlockFly:
		return constants.CharFly
	case BlockChain:
		return constants.DoodadChain
	case BlockReeds:
		return constants.DoodadReeds
	case BlockFlowers:
		return constants.DoodadFlowers
	}
	return constants.TileEmpty
}

var toID = map[string]Block{
	constants.TileTurf:         BlockTurf,
	constants.TileFall:         BlockFall,
	constants.TileCracked:      BlockCracked,
	constants.TilePhase:        BlockPhase,
	constants.TileLadderMiddle: BlockLadder,
	constants.TileLadderCrackM: BlockLadderCracked,
	constants.TileExitLadderM:  BlockLadderExit,
	constants.ItemBox:          BlockBox,
	constants.CharPlayer1:      BlockPlayer1,
	constants.CharPlayer2:      BlockPlayer2,
	constants.CharPlayer3:      BlockPlayer3,
	constants.CharPlayer4:      BlockPlayer4,
	constants.ItemKeyYellow:    BlockKeyYellow,
	constants.ItemKeyOrange:    BlockKeyOrange,
	constants.ItemKeyGray:      BlockKeyGray,
	constants.ItemKeyCyan:      BlockKeyCyan,
	constants.ItemKeyBlue:      BlockKeyBlue,
	constants.ItemKeyGreen:     BlockKeyGreen,
	constants.ItemKeyPurple:    BlockKeyPurple,
	constants.ItemKeyBrown:     BlockKeyBrown,
	constants.TileDoorYellow:   BlockDoorYellow,
	constants.TileDoorOrange:   BlockDoorOrange,
	constants.TileDoorGray:     BlockDoorGray,
	constants.TileDoorCyan:     BlockDoorCyan,
	constants.TileDoorBlue:     BlockDoorBlue,
	constants.TileDoorGreen:    BlockDoorGreen,
	constants.TileDoorPurple:   BlockDoorPurple,
	constants.TileDoorBrown:    BlockDoorBrown,
	constants.TileLockYellow:   BlockLockYellow,
	constants.TileLockOrange:   BlockLockOrange,
	constants.TileLockGray:     BlockLockGray,
	constants.TileLockCyan:     BlockLockCyan,
	constants.TileLockBlue:     BlockLockBlue,
	constants.TileLockGreen:    BlockLockGreen,
	constants.TileLockPurple:   BlockLockPurple,
	constants.TileLockBrown:    BlockLockBrown,
	constants.ItemGemYellow:    BlockGemYellow,
	constants.ItemGemOrange:    BlockGemOrange,
	constants.ItemGemGray:      BlockGemGray,
	constants.ItemGemCyan:      BlockGemCyan,
	constants.ItemGemBlue:      BlockGemBlue,
	constants.ItemGemGreen:     BlockGemGreen,
	constants.ItemGemPurple:    BlockGemPurple,
	constants.ItemGemBrown:     BlockGemBrown,
	constants.CharDemon:        BlockDemon,
	constants.CharFly:          BlockFly,
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
	case BlockPhase:
		buffer.WriteString(constants.TilePhase)
	case BlockCracked:
		buffer.WriteString(constants.TileCracked)
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
	Block    Block          `json:"tile"`
	Ladder   bool           `json:"ladder"`
	Metadata TileMetadata   `json:"metadata"`
	Coords   world.Coords   `json:"-"`
	Object   *object.Object `json:"-"`
	Update   bool           `json:"-"`
	Entity   *ecs.Entity    `json:"-"`
}

func (t *Tile) Copy() *Tile {
	return &Tile{
		Block:    t.Block,
		Ladder:   t.Ladder,
		Coords:   t.Coords,
		Metadata: t.Metadata,
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

type TileMetadata struct {
	Flipped bool `json:"flipped"`
}
