package data

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/world"
	"github.com/beefsack/go-astar"
	"github.com/bytearena/ecs"
)

type Tile struct {
	Block        Block             `json:"tile"`
	LastBlock    Block             `json:"-"`
	Metadata     TileMetadata      `json:"metadata"`
	Flags        TileFlags         `json:"-"`
	Coords       world.Coords      `json:"-"`
	Object       *object.Object    `json:"-"`
	Update       bool              `json:"-"`
	Entity       *ecs.Entity       `json:"-"`
	Mask         *ecs.Entity       `json:"-"`
	Counter      int               `json:"-"`
	Live         bool              `json:"-"`
	AltBlock     int               `json:"alt"`
	FloatingText *FloatingText     `json:"-"`
	TextData     *FloatingTextData `json:"text,omitempty"`
}

func (t *Tile) SpriteString() string {
	colSuffixKey := constants.SprColorYellow
	colSuffixTools := ""
	colSuffixTrans := ""
	switch t.Metadata.Color {
	case NonPlayerYellow:
		colSuffixKey = constants.SprColorYellow
	case NonPlayerBrown:
		colSuffixKey = constants.SprColorBrown
	case NonPlayerGray:
		colSuffixKey = constants.SprColorGray
	case NonPlayerCyan:
		colSuffixKey = constants.SprColorCyan
	case NonPlayerRed:
		colSuffixKey = constants.SprColorRed
		colSuffixTrans = constants.SprColorRed
	case PlayerBlue:
		colSuffixKey = constants.SprColorBlue
		colSuffixTools = constants.SprColorBlue
		colSuffixTrans = constants.SprColorBlue
	case PlayerGreen:
		colSuffixKey = constants.SprColorGreen
		colSuffixTools = constants.SprColorGreen
		colSuffixTrans = constants.SprColorGreen
	case PlayerPurple:
		colSuffixKey = constants.SprColorPurple
		colSuffixTools = constants.SprColorPurple
		colSuffixTrans = constants.SprColorPurple
	case PlayerOrange:
		colSuffixKey = constants.SprColorOrange
		colSuffixTools = constants.SprColorOrange
		colSuffixTrans = constants.SprColorOrange
	}
	switch t.Block {
	case BlockTurf, BlockFall, BlockCracked, BlockClose, BlockHideout,
		BlockLadderTurf, BlockLadderCrackedTurf, BlockLadderExitTurf:
		if CurrPuzzleSet != nil && CurrPuzzleSet.CurrPuzzle.Metadata.WorldSprite != "" {
			return CurrPuzzleSet.CurrPuzzle.Metadata.WorldSprite
		}
		return constants.WorldSprites[constants.WorldMoss]
	case BlockBedrock, BlockPhase:
		if CurrPuzzleSet != nil && CurrPuzzleSet.CurrPuzzle.Metadata.WorldSprite != "" {
			return fmt.Sprintf("%s_%s", CurrPuzzleSet.CurrPuzzle.Metadata.WorldSprite, constants.TileBedrock)
		}
		return fmt.Sprintf("%s_%s", constants.WorldSprites[constants.WorldMoss], constants.TileBedrock)
	case BlockSpike:
		if CurrPuzzleSet != nil && CurrPuzzleSet.CurrPuzzle.Metadata.WorldSprite != "" {
			return fmt.Sprintf("%s_%s", CurrPuzzleSet.CurrPuzzle.Metadata.WorldSprite, constants.TileSpike)
		}
		return fmt.Sprintf("%s_%s", constants.WorldSprites[constants.WorldMoss], constants.TileSpike)
	case BlockLiquid:
		if CurrPuzzleSet != nil && CurrPuzzleSet.CurrPuzzle.Metadata.WorldLiquid != "" {
			return fmt.Sprintf("%s_%s", constants.TileLiquid, CurrPuzzleSet.CurrPuzzle.Metadata.WorldLiquid)
		}
		return constants.TileLiquid
	case BlockLadder:
		return constants.TileLadderMiddle
	case BlockLadderCracked:
		return constants.TileLadderCrackM
	case BlockLadderExit:
		return constants.TileExitLadderM
	case BlockBar:
		return constants.TileBar
	case BlockTransporterExit:
		return constants.TileTransporterExit
	case BlockTransporter:
		return constants.TileTransporter + colSuffixTrans
	case BlockJumpBoots:
		return constants.ItemJumpBoots + colSuffixTools
	case BlockBox:
		return constants.ItemBox + colSuffixTools
	case BlockJetpack:
		return constants.ItemJetpack + colSuffixTools
	case BlockDisguise:
		return constants.ItemDisguise + colSuffixTools
	case BlockDrill:
		return constants.ItemDrill + colSuffixTools
	case BlockFlamethrower:
		return constants.ItemFlamethrower + colSuffixTools
	case BlockBomb:
		return constants.ItemBomb + colSuffixTools
	case BlockBombLit:
		return constants.ItemBombLit
	case BlockPlayer1:
		return constants.CharPlayer1
	case BlockPlayer2:
		return constants.CharPlayer2
	case BlockPlayer3:
		return constants.CharPlayer3
	case BlockPlayer4:
		return constants.CharPlayer4
	case BlockGem:
		return constants.ItemGem + colSuffixKey
	case BlockKey:
		return constants.ItemKey + colSuffixKey
	case BlockDoorHidden:
		return constants.TileDoorHidden + colSuffixKey
	case BlockDoorVisible:
		return constants.TileDoorVisible + colSuffixKey
	case BlockDoorLocked:
		return constants.TileDoorLocked + colSuffixKey
	case BlockDemon:
		return constants.CharDemon
	case BlockFly:
		return constants.CharFly
	case BlockDemonRegen:
		return constants.TileDemonRegen
	case BlockFlyRegen:
		return constants.TileFlyRegen
	case BlockReeds:
		return constants.DoodadReeds
	case BlockFlowers:
		return constants.DoodadFlowers
	case BlockMoss:
		return constants.DoodadMoss
	case BlockGrass:
		return constants.DoodadGrass
	case BlockCattail:
		return constants.DoodadCattail
	case BlockCactus1:
		return constants.DoodadCactus1
	case BlockCactus2:
		return constants.DoodadCactus2
	case BlockVine:
		return constants.DoodadVine
	case BlockNest:
		return constants.DoodadNest
	case BlockSkull:
		return constants.DoodadSkull
	case BlockDots:
		return constants.DoodadDots
	case BlockBoulder:
		return constants.DoodadBoulder
	case BlockMush1:
		return constants.DoodadMush1
	case BlockMush2:
		return constants.DoodadMush2
	case BlockChain:
		return constants.DoodadChain
	case BlockGear:
		return constants.DoodadGear
	}
	return constants.TileEmpty
}

func (t *Tile) Copy() *Tile {
	return &Tile{
		Block:    t.Block,
		AltBlock: t.AltBlock,
		Coords:   t.Coords,
		Metadata: CopyMetadata(t.Metadata),
		TextData: t.TextData.Copy(),
	}
}

func (t *Tile) CopyInto(c *Tile) {
	c.Block = t.Block
	c.AltBlock = t.AltBlock
	c.Object.Flip = t.Metadata.Flipped
	c.Metadata = CopyMetadata(t.Metadata)
	c.TextData = t.TextData.Copy()
	CreateFloatingText(c, c.TextData)
}

func (t *Tile) ToEmpty() {
	t.Block = BlockEmpty
	t.Metadata = DefaultMetadata()
	t.Flags = DefaultFlags()
	t.TextData = nil
	if t.FloatingText != nil {
		myecs.Manager.DisposeEntity(t.FloatingText.Entity)
		myecs.Manager.DisposeEntity(t.FloatingText.ShEntity)
		t.FloatingText = nil
	}
}

func (t *Tile) IsEmpty() bool {
	return t != nil && !(t.IsLadder() ||
		t.Block == BlockBar ||
		t.Block == BlockTurf ||
		t.Block == BlockBedrock ||
		t.Block == BlockFall ||
		t.Block == BlockPhase ||
		t.Block == BlockCracked ||
		t.Block == BlockSpike)
}

func (t *Tile) IsSolid() bool {
	return t == nil || (!t.Flags.Collapse &&
		!t.IsLadder() &&
		(t.Block == BlockTurf ||
			t.Block == BlockBedrock ||
			t.Block == BlockFall ||
			t.Block == BlockPhase ||
			t.Block == BlockCracked ||
			t.Block == BlockClose ||
			t.Block == BlockSpike ||
			(t.Block == BlockLadderExitTurf && !CurrLevel.DoorsOpen)))
}

func (t *Tile) IsNilOrSolid() bool {
	return t == nil || (!t.Flags.Collapse &&
		!t.IsLadder() &&
		(t.Block == BlockTurf ||
			t.Block == BlockBedrock ||
			t.Block == BlockFall ||
			t.Block == BlockPhase ||
			t.Block == BlockCracked ||
			t.Block == BlockClose ||
			t.Block == BlockSpike ||
			(t.Block == BlockLadderExitTurf && !CurrLevel.DoorsOpen)))
}

func (t *Tile) IsRunnable() bool {
	return t == nil || (!t.Flags.Collapse &&
		(t.Block == BlockTurf ||
			t.Block == BlockBedrock ||
			t.Block == BlockLadderTurf ||
			t.Block == BlockLadderCrackedTurf ||
			t.Block == BlockLadderExitTurf ||
			t.Block == BlockFall ||
			t.Block == BlockPhase ||
			t.Block == BlockCracked ||
			t.Block == BlockClose ||
			t.Block == BlockSpike))
}
 
// IsBlock used for tile connectivity
func (t *Tile) IsBlock() bool {
	return t == nil ||
		((t.Block == BlockTurf ||
			t.Block == BlockBedrock ||
			t.Block == BlockLadderTurf ||
			t.Block == BlockLadderCrackedTurf ||
			t.Block == BlockLadderExitTurf ||
			t.Block == BlockFall ||
			t.Block == BlockPhase ||
			t.Block == BlockCracked ||
			t.Block == BlockClose ||
			t.Block == BlockSpike) &&
			(t.Flags.Regen ||
				!(t.Flags.Collapse && t.Counter > constants.CollapseCounter && !t.Flags.BareFangs)))
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

var PlayerAbove bool

func (t *Tile) PathNeighbors() []astar.Pather {
	if CurrLevel == nil {
		return []astar.Pather{}
	}
	var neighbors []astar.Pather
	// Down
	d := CurrLevel.Get(t.Coords.X, t.Coords.Y-1)
	if !PlayerAbove && d != nil && (d.IsEmpty() || d.IsLadder()) {
		neighbors = append(neighbors, d)
	}
	notFalling := d.IsSolid() || d.IsLadder() || t.IsLadder()
	// Left
	l := CurrLevel.Get(t.Coords.X-1, t.Coords.Y)
	//lb := CurrLevel.Tiles.Get(t.Coords.X-1, t.Coords.Y-1)
	//if notFalling && l != nil && !l.IsSolid() &&
	//	(l.Ladder || lb == nil || lb.IsSolid()) {
	//	neighbors = append(neighbors, l)
	//}
	if notFalling && !l.IsSolid() {
		neighbors = append(neighbors, l)
	}
	// Right
	r := CurrLevel.Get(t.Coords.X+1, t.Coords.Y)
	//rb := CurrLevel.Tiles.Get(t.Coords.X+1, t.Coords.Y-1)
	//if notFalling && r != nil && !r.IsSolid() &&
	//	(r.Ladder || rb == nil || rb.IsSolid()) {
	//	neighbors = append(neighbors, r)
	//}
	if notFalling && !r.IsSolid() {
		neighbors = append(neighbors, r)
	}
	// Up
	u := CurrLevel.Get(t.Coords.X, t.Coords.Y+1)
	if PlayerAbove && notFalling && !u.IsSolid() && t.IsLadder() {
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
	Cracked     bool `json:"-"`
	Collapse    bool `json:"-"`
	BareFangs   bool `json:"-"`
	Regen       bool `json:"-"`
	PhaseChange bool `json:"-"`
	LCracked    bool `json:"-"`
	LCollapse   bool `json:"-"`
	Using       bool `json:"-"`
}

func DefaultFlags() TileFlags {
	return TileFlags{}
}

type TileMetadata struct {
	Flipped     bool           `json:"flipped,omitempty"`
	EnemyCrack  bool           `json:"enemyCrack,omitempty"`
	Regenerate  bool           `json:"regenerate,omitempty"`
	RegenDelay  int            `json:"regenDelay,omitempty"`
	Timer       int            `json:"timer,omitempty"`
	BombCross   bool           `json:"bombCross,omitempty"`
	LinkedTiles []world.Coords `json:"linkedTiles,regenTiles,omitempty"`
	Phase       int            `json:"phase,omitempty"`
	ShowCrack   bool           `json:"showCrack,omitempty"`
	Changed     bool           `json:"-"`
	ExitIndex   int            `json:"exitIndex,omitempty"`
	Color       ItemColor      `json:"itemColor"`
}

func DefaultMetadata() TileMetadata {
	return TileMetadata{
		Regenerate: true,
		Timer:      5,
		ExitIndex:  -1,
	}
}

func CopyMetadata(m TileMetadata) TileMetadata {
	cm := TileMetadata{
		Flipped:     m.Flipped,
		EnemyCrack:  m.EnemyCrack,
		Regenerate:  m.Regenerate,
		RegenDelay:  m.RegenDelay,
		Timer:       m.Timer,
		BombCross:   m.BombCross,
		LinkedTiles: nil,
		Phase:       m.Phase,
		ShowCrack:   m.ShowCrack,
		ExitIndex:   m.ExitIndex,
		Color:       m.Color,
	}
	for _, ln := range m.LinkedTiles {
		cm.LinkedTiles = append(cm.LinkedTiles, ln)
	}
	return cm
}
