package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gemrunner/internal/constants"
)

type Block int

const (
	BlockTurf = iota
	BlockBedrock
	BlockFall
	BlockCracked
	BlockClose
	BlockPhase
	BlockSpike

	BlockLadder
	BlockLadderCracked
	BlockLadderExit
	BlockBar
	BlockHideout

	BlockPlayer1
	BlockPlayer2
	BlockPlayer3
	BlockPlayer4
	BlockDemon
	BlockDemonRegen
	BlockFly
	BlockFlyRegen

	BlockGem
	BlockDoorHidden
	BlockDoorVisible
	BlockDoorLocked
	BlockKey

	BlockJumpBoots
	BlockBox
	BlockJetpack
	BlockBomb
	BlockBombLit
	BlockDisguise
	BlockDrill
	BlockFlamethrower

	BlockReeds
	BlockFlowers
	BlockMoss
	BlockGrass
	BlockCattail
	BlockCactus1
	BlockCactus2
	BlockVine

	BlockNest
	BlockSkull
	BlockDots
	BlockBoulder
	BlockMush1
	BlockMush2
	BlockChain
	BlockGear

	BlockEmpty

	BlockLadderTurf
	BlockLadderCrackedTurf
	BlockLadderExitTurf
)

func (b Block) String() string {
	switch b {
	case BlockTurf:
		return constants.TileTurf
	case BlockFall:
		return constants.TileFall
	case BlockCracked:
		return constants.TileCracked
	case BlockClose:
		return constants.TileClose
	case BlockHideout:
		return constants.TileHideout
	case BlockLadderTurf:
		return constants.TileLadderTurf
	case BlockLadderCrackedTurf:
		return constants.TileLadderCrackedTurf
	case BlockLadderExitTurf:
		return constants.TileLadderExitTurf
	case BlockBedrock:
		return constants.TileBedrock
	case BlockPhase:
		return constants.TilePhase
	case BlockSpike:
		return constants.TileSpike
	case BlockLadder:
		return constants.TileLadderMiddle
	case BlockLadderCracked:
		return constants.TileLadderCrackM
	case BlockLadderExit:
		return constants.TileExitLadderM
	case BlockBar:
		return constants.TileBar
	case BlockJumpBoots:
		return constants.ItemJumpBoots
	case BlockBox:
		return constants.ItemBox
	case BlockJetpack:
		return constants.ItemJetpack
	case BlockDisguise:
		return constants.ItemDisguise
	case BlockDrill:
		return constants.ItemDrill
	case BlockFlamethrower:
		return constants.ItemFlamethrower
	case BlockBomb:
		return constants.ItemBomb
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
		return constants.ItemGem
	case BlockKey:
		return constants.ItemKey
	case BlockDoorHidden:
		return constants.TileDoorHidden
	case BlockDoorVisible:
		return constants.TileDoorVisible
	case BlockDoorLocked:
		return constants.TileDoorLocked
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

func (b Block) SpriteString() string {
	colSuffixKey := constants.SprColorYellow
	colSuffixTools := ""
	if Editor != nil {
		switch Editor.PaletteColor {
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
		case PlayerBlue:
			colSuffixKey = constants.SprColorBlue
			colSuffixTools = constants.SprColorBlue
		case PlayerGreen:
			colSuffixKey = constants.SprColorGreen
			colSuffixTools = constants.SprColorGreen
		case PlayerPurple:
			colSuffixKey = constants.SprColorPurple
			colSuffixTools = constants.SprColorPurple
		case PlayerOrange:
			colSuffixKey = constants.SprColorOrange
			colSuffixTools = constants.SprColorOrange
		}
	}
	switch b {
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
	default:
		return b.String()
	}
}

var toID = map[string]Block{
	constants.TileTurf:              BlockTurf,
	constants.TileBedrock:           BlockBedrock,
	constants.TileFall:              BlockFall,
	constants.TileCracked:           BlockCracked,
	constants.TileClose:             BlockClose,
	constants.TilePhase:             BlockPhase,
	constants.TileSpike:             BlockSpike,
	constants.TileLadderMiddle:      BlockLadder,
	constants.TileLadderCrackM:      BlockLadderCracked,
	constants.TileExitLadderM:       BlockLadderExit,
	constants.TileLadderTurf:        BlockLadderTurf,
	constants.TileLadderCrackedTurf: BlockLadderCrackedTurf,
	constants.TileLadderExitTurf:    BlockLadderExitTurf,
	constants.TileBar:               BlockBar,
	constants.TileHideout:           BlockHideout,
	constants.ItemJumpBoots:         BlockJumpBoots,
	constants.ItemBox:               BlockBox,
	constants.ItemJetpack:           BlockJetpack,
	constants.ItemDisguise:          BlockDisguise,
	constants.ItemDrill:             BlockDrill,
	constants.ItemFlamethrower:      BlockFlamethrower,
	constants.ItemBomb:              BlockBomb,
	constants.ItemBombLit:           BlockBombLit,
	constants.CharPlayer1:           BlockPlayer1,
	constants.CharPlayer2:           BlockPlayer2,
	constants.CharPlayer3:           BlockPlayer3,
	constants.CharPlayer4:           BlockPlayer4,
	constants.TileDoorHidden:        BlockDoorHidden,
	constants.TileDoorVisible:       BlockDoorVisible,
	constants.TileDoorLocked:        BlockDoorLocked,
	constants.ItemGem:               BlockGem,
	constants.ItemKey:               BlockKey,
	constants.CharDemon:             BlockDemon,
	constants.CharFly:               BlockFly,
	constants.TileDemonRegen:        BlockDemonRegen,
	constants.TileFlyRegen:          BlockFlyRegen,
	constants.DoodadReeds:           BlockReeds,
	constants.DoodadFlowers:         BlockFlowers,
	constants.DoodadMoss:            BlockMoss,
	constants.DoodadGrass:           BlockGrass,
	constants.DoodadCattail:         BlockCattail,
	constants.DoodadCactus1:         BlockCactus1,
	constants.DoodadCactus2:         BlockCactus2,
	constants.DoodadVine:            BlockVine,
	constants.DoodadNest:            BlockNest,
	constants.DoodadSkull:           BlockSkull,
	constants.DoodadDots:            BlockDots,
	constants.DoodadBoulder:         BlockBoulder,
	constants.DoodadMush1:           BlockMush1,
	constants.DoodadMush2:           BlockMush2,
	constants.DoodadChain:           BlockChain,
	constants.DoodadGear:            BlockGear,
	constants.TileEmpty:             BlockEmpty,
}

func (b Block) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(b.String())
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
