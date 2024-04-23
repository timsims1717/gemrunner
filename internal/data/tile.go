package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/pkg/object"
	"gemrunner/pkg/world"
	"github.com/beefsack/go-astar"
	"github.com/bytearena/ecs"
)

type Block int

const (
	BlockTurf = iota
	BlockBedrock
	BlockFall
	BlockCracked
	BlockPhase
	BlockSpike
	BlockLadder
	BlockLadderCracked
	BlockLadderExit

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
	BlockBox
	BlockChain
	BlockReeds
	BlockFlowers

	BlockDemon
	BlockDemonRegen
	BlockFly
	BlockFlyRegen
	BlockEmpty

	BlockLadderTurf
	BlockLadderCrackedTurf
	BlockLadderExitTurf
)

func (b Block) String() string {
	switch b {
	case BlockTurf, BlockFall, BlockCracked, BlockPhase,
		BlockLadderTurf, BlockLadderCrackedTurf, BlockLadderExitTurf:
		if CurrPuzzle != nil && CurrPuzzle.Metadata.WorldSprite != "" {
			return CurrPuzzle.Metadata.WorldSprite
		}
		return constants.WorldSprites[constants.WorldMoss]
	case BlockBedrock:
		if CurrPuzzle != nil && CurrPuzzle.Metadata.WorldSprite != "" {
			return fmt.Sprintf("%s_%s", CurrPuzzle.Metadata.WorldSprite, constants.TileBedrock)
		}
		return fmt.Sprintf("%s_%s", constants.WorldSprites[constants.WorldMoss], constants.TileBedrock)
	case BlockSpike:
		if CurrPuzzle != nil && CurrPuzzle.Metadata.WorldSprite != "" {
			return fmt.Sprintf("%s_%s", CurrPuzzle.Metadata.WorldSprite, constants.TileSpike)
		}
		return fmt.Sprintf("%s_%s", constants.WorldSprites[constants.WorldMoss], constants.TileSpike)
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
	case BlockDemonRegen:
		return constants.TileDemonRegen
	case BlockFlyRegen:
		return constants.TileFlyRegen
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
	constants.TileTurf:              BlockTurf,
	constants.TileBedrock:           BlockBedrock,
	constants.TileFall:              BlockFall,
	constants.TileCracked:           BlockCracked,
	constants.TilePhase:             BlockPhase,
	constants.TileSpike:             BlockSpike,
	constants.TileLadderMiddle:      BlockLadder,
	constants.TileLadderCrackM:      BlockLadderCracked,
	constants.TileExitLadderM:       BlockLadderExit,
	constants.TileLadderTurf:        BlockLadderTurf,
	constants.TileLadderCrackedTurf: BlockLadderCrackedTurf,
	constants.TileLadderExitTurf:    BlockLadderExitTurf,
	constants.ItemBox:               BlockBox,
	constants.CharPlayer1:           BlockPlayer1,
	constants.CharPlayer2:           BlockPlayer2,
	constants.CharPlayer3:           BlockPlayer3,
	constants.CharPlayer4:           BlockPlayer4,
	constants.ItemKeyYellow:         BlockKeyYellow,
	constants.ItemKeyOrange:         BlockKeyOrange,
	constants.ItemKeyGray:           BlockKeyGray,
	constants.ItemKeyCyan:           BlockKeyCyan,
	constants.ItemKeyBlue:           BlockKeyBlue,
	constants.ItemKeyGreen:          BlockKeyGreen,
	constants.ItemKeyPurple:         BlockKeyPurple,
	constants.ItemKeyBrown:          BlockKeyBrown,
	constants.TileDoorYellow:        BlockDoorYellow,
	constants.TileDoorOrange:        BlockDoorOrange,
	constants.TileDoorGray:          BlockDoorGray,
	constants.TileDoorCyan:          BlockDoorCyan,
	constants.TileDoorBlue:          BlockDoorBlue,
	constants.TileDoorGreen:         BlockDoorGreen,
	constants.TileDoorPurple:        BlockDoorPurple,
	constants.TileDoorBrown:         BlockDoorBrown,
	constants.TileLockYellow:        BlockLockYellow,
	constants.TileLockOrange:        BlockLockOrange,
	constants.TileLockGray:          BlockLockGray,
	constants.TileLockCyan:          BlockLockCyan,
	constants.TileLockBlue:          BlockLockBlue,
	constants.TileLockGreen:         BlockLockGreen,
	constants.TileLockPurple:        BlockLockPurple,
	constants.TileLockBrown:         BlockLockBrown,
	constants.ItemGemYellow:         BlockGemYellow,
	constants.ItemGemOrange:         BlockGemOrange,
	constants.ItemGemGray:           BlockGemGray,
	constants.ItemGemCyan:           BlockGemCyan,
	constants.ItemGemBlue:           BlockGemBlue,
	constants.ItemGemGreen:          BlockGemGreen,
	constants.ItemGemPurple:         BlockGemPurple,
	constants.ItemGemBrown:          BlockGemBrown,
	constants.CharDemon:             BlockDemon,
	constants.CharFly:               BlockFly,
	constants.TileDemonRegen:        BlockDemonRegen,
	constants.TileFlyRegen:          BlockFlyRegen,
	constants.DoodadChain:           BlockChain,
	constants.DoodadReeds:           BlockReeds,
	constants.DoodadFlowers:         BlockFlowers,
	constants.TileEmpty:             BlockEmpty,
}

func (b Block) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	switch b {
	case BlockTurf:
		buffer.WriteString(constants.TileTurf)
	case BlockBedrock:
		buffer.WriteString(constants.TileBedrock)
	case BlockFall:
		buffer.WriteString(constants.TileFall)
	case BlockPhase:
		buffer.WriteString(constants.TilePhase)
	case BlockCracked:
		buffer.WriteString(constants.TileCracked)
	case BlockLadderTurf:
		buffer.WriteString(constants.TileLadderTurf)
	case BlockLadderCrackedTurf:
		buffer.WriteString(constants.TileLadderCrackedTurf)
	case BlockLadderExitTurf:
		buffer.WriteString(constants.TileLadderExitTurf)
	case BlockSpike:
		buffer.WriteString(constants.TileSpike)
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
		var ji int
		err = json.Unmarshal(bts, &ji)
		if err != nil {
			var jb bool
			err = json.Unmarshal(bts, &jb)
			if err != nil {
				return err
			}
			*b = toID[constants.TileLadderMiddle]
		}
		*b = Block(ji)
	}
	*b = toID[j]
	return nil
}

//func (b *Block) UnmarshalJSON(bts []byte) error {
//	var j string
//	err := json.Unmarshal(bts, &j)
//	if err != nil {
//		return err
//	}
//	*b = toID[j]
//	return nil
//}

type Tile struct {
	Block    Block          `json:"tile"`
	Metadata TileMetadata   `json:"metadata"`
	Flags    TileFlags      `json:"-"`
	Coords   world.Coords   `json:"-"`
	Object   *object.Object `json:"-"`
	Update   bool           `json:"-"`
	Entity   *ecs.Entity    `json:"-"`
	Mask     *ecs.Entity    `json:"-"`
	Counter  int            `json:"-"`
	Live     bool           `json:"-"`
}

func (t *Tile) Copy() *Tile {
	return &Tile{
		Block:    t.Block,
		Coords:   t.Coords,
		Metadata: t.Metadata,
	}
}

func (t *Tile) CopyInto(c *Tile) {
	c.Block = t.Block
	c.Object.Flip = t.Metadata.Flipped
	c.Metadata = CopyMetadata(t.Metadata)
}

func (t *Tile) ToEmpty() {
	t.Block = BlockEmpty
	t.Metadata = DefaultMetadata()
	t.Flags = DefaultFlags()
}

func (t *Tile) IsEmpty() bool {
	return !(t.IsLadder() ||
		t.Block == BlockTurf ||
		t.Block == BlockBedrock ||
		t.Block == BlockFall ||
		t.Block == BlockCracked ||
		t.Block == BlockSpike ||
		t.Block == BlockPhase)
}

func (t *Tile) IsSolid() bool {
	return !t.Flags.Collapse &&
		!t.IsLadder() &&
		(t.Block == BlockTurf ||
			t.Block == BlockBedrock ||
			t.Block == BlockFall ||
			t.Block == BlockCracked ||
			t.Block == BlockSpike)
}

func (t *Tile) IsNilOrSolid() bool {
	return t == nil || (!t.Flags.Collapse &&
		!t.IsLadder() &&
		(t.Block == BlockTurf ||
			t.Block == BlockBedrock ||
			t.Block == BlockFall ||
			t.Block == BlockCracked ||
			t.Block == BlockSpike))
}

func (t *Tile) IsBlock() bool {
	return t == nil ||
		((t.Block == BlockTurf ||
			t.Block == BlockBedrock ||
			t.Block == BlockLadderTurf ||
			t.Block == BlockLadderCrackedTurf ||
			t.Block == BlockLadderExitTurf ||
			t.Block == BlockFall ||
			t.Block == BlockCracked ||
			t.Block == BlockSpike) &&
			(t.Flags.Regen ||
				!(t.Flags.Collapse && t.Counter > constants.CollapseCounter)))
}

func (t *Tile) IsLadder() bool {
	if t.Live {
		return !t.Flags.LCollapse && (t.Block == BlockLadder ||
			t.Block == BlockLadderTurf ||
			(t.Block == BlockLadderCracked && !t.Flags.LCollapse) ||
			(t.Block == BlockLadderCrackedTurf && !t.Flags.LCollapse) ||
			t.Block == BlockLadderCrackedTurf ||
			(t.Block == BlockLadderExitTurf && CurrLevel.DoorsOpen) ||
			(t.Block == BlockLadderExit && CurrLevel.DoorsOpen))
	} else {
		return t.Block == BlockLadder ||
			t.Block == BlockLadderTurf ||
			t.Block == BlockLadderExit ||
			t.Block == BlockLadderExitTurf ||
			t.Block == BlockLadderCracked ||
			t.Block == BlockLadderCrackedTurf
	}
}

func (t *Tile) CanDig() bool {
	return !t.Flags.Collapse && t.Block == BlockTurf
}

// a* implementation

func (t *Tile) PathNeighbors() []astar.Pather {
	if CurrLevel == nil {
		return []astar.Pather{}
	}
	var neighbors []astar.Pather
	// Down
	d := CurrLevel.Tiles.Get(t.Coords.X, t.Coords.Y-1)
	if d != nil && !d.IsSolid() {
		neighbors = append(neighbors, d)
	}
	notFalling := d == nil || d.IsSolid() || d.IsLadder() || t.IsLadder()
	// Left
	l := CurrLevel.Tiles.Get(t.Coords.X-1, t.Coords.Y)
	//lb := CurrLevel.Tiles.Get(t.Coords.X-1, t.Coords.Y-1)
	//if notFalling && l != nil && !l.IsSolid() &&
	//	(l.Ladder || lb == nil || lb.IsSolid()) {
	//	neighbors = append(neighbors, l)
	//}
	if notFalling && l != nil && !l.IsSolid() {
		neighbors = append(neighbors, l)
	}
	// Right
	r := CurrLevel.Tiles.Get(t.Coords.X+1, t.Coords.Y)
	//rb := CurrLevel.Tiles.Get(t.Coords.X+1, t.Coords.Y-1)
	//if notFalling && r != nil && !r.IsSolid() &&
	//	(r.Ladder || rb == nil || rb.IsSolid()) {
	//	neighbors = append(neighbors, r)
	//}
	if notFalling && r != nil && !r.IsSolid() {
		neighbors = append(neighbors, r)
	}
	// Up
	u := CurrLevel.Tiles.Get(t.Coords.X, t.Coords.Y+1)
	if notFalling && u != nil && !u.IsSolid() && t.IsLadder() {
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

type TileFlags struct {
	Cracked   bool `json:"-"`
	Collapse  bool `json:"-"`
	Regen     bool `json:"-"`
	LCracked  bool `json:"-"`
	LCollapse bool `json:"-"`
}

func DefaultFlags() TileFlags {
	return TileFlags{}
}

type TileMetadata struct {
	Flipped     bool           `json:"flipped"`
	EnemyCrack  bool           `json:"enemyCrack"`
	Regenerate  bool           `json:"regenerate"`
	LinkedTiles []world.Coords `json:"linkedTiles,regenTiles"`
	Phase       int            `json:"phase"`
	ShowCrack   bool           `json:"showCrack"`
	Changed     bool           `json:"changed"`
}

func DefaultMetadata() TileMetadata {
	return TileMetadata{
		Regenerate: true,
	}
}

func CopyMetadata(m TileMetadata) TileMetadata {
	cm := TileMetadata{
		Flipped:     m.Flipped,
		EnemyCrack:  m.EnemyCrack,
		Regenerate:  m.Regenerate,
		LinkedTiles: nil,
		Phase:       m.Phase,
		ShowCrack:   m.ShowCrack,
	}
	for _, ln := range m.LinkedTiles {
		cm.LinkedTiles = append(cm.LinkedTiles, ln)
	}
	return cm
}
