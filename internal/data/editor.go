package data

import (
	"gemrunner/pkg/timing"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"strings"
)

var (
	Editor *editor
)

type editor struct {
	CurrBlock Block
	Offset    pixel.Vec
	Mode      EditorMode
	LastMode  EditorMode
	NoInput   bool
	LastTiles map[Block]*Tile

	Hover        bool
	Consume      string
	SelectVis    bool
	SelectTimer  *timing.Timer
	SelectQuick  bool
	LastCoords   world.Coords
	PosTop       bool
	ModeChanged  bool
	PaletteColor ItemColor

	BlockSelect *viewport.ViewPort
}

func NewEditor() {
	Editor = &editor{
		LastCoords: world.Coords{X: -1, Y: -1},
		LastTiles:  make(map[Block]*Tile),
	}
}

func BlockSelectPlacement(b, w, h int) pixel.Vec {
	wo := float64(w) / 2
	if w%2 == 0 {
		wo -= 0.5
	}
	ho := float64(h) / 2
	if w%2 == 0 {
		ho -= 0.5
	}
	return pixel.V(world.TileSize*(float64(b%w)-wo), world.TileSize*(-float64(b/w)+ho))
}

type EditorMode int

const (
	ModeBrush = iota
	ModeLine
	ModeSquare
	ModeFill
	ModeErase
	ModeEyedrop
	ModeSelect
	ModeMove
	ModeCut
	ModeCopy
	ModePaste
	ModeFlipVertical
	ModeFlipHorizontal
	ModeWrench
	ModeWire
	ModePalette
	ModeText
	ModeUndo
	ModeRedo
	ModeDelete
	ModeSave
	ModeOpen
	EndModeList
)

func (m EditorMode) String() string {
	switch m {
	case ModeBrush:
		return "Brush"
	case ModeLine:
		return "Line"
	case ModeSquare:
		return "Square"
	case ModeFill:
		return "Fill"
	case ModeErase:
		return "Erase"
	case ModeEyedrop:
		return "Eyedrop"
	case ModeSelect:
		return "Select"
	case ModeMove:
		return "Move"
	case ModeCut:
		return "Cut"
	case ModeCopy:
		return "Copy"
	case ModePaste:
		return "Paste"
	case ModeFlipVertical:
		return "FlipVertical"
	case ModeFlipHorizontal:
		return "FlipHorizontal"
	case ModeWrench:
		return "Wrench"
	case ModeWire:
		return "Wire"
	case ModePalette:
		return "Palette"
	case ModeText:
		return "Text"
	case ModeUndo:
		return "Undo"
	case ModeRedo:
		return "Redo"
	case ModeDelete:
		return "Delete"
	case ModeSave:
		return "Save"
	case ModeOpen:
		return "Open"
	}
	return ""
}

func ModeFromSprString(s string) EditorMode {
	ss := strings.Split(s, "_")
	ms := ss[0]
	switch ms {
	case "brush":
		return ModeBrush
	case "line":
		return ModeLine
	case "square":
		return ModeSquare
	case "fill":
		return ModeFill
	case "erase":
		return ModeErase
	case "eyedrop":
		return ModeEyedrop
	case "select":
		return ModeSelect
	case "move":
		return ModeMove
	case "cut":
		return ModeCut
	case "copy":
		return ModeCopy
	case "paste":
		return ModePaste
	case "flipv":
		return ModeFlipVertical
	case "fliph":
		return ModeFlipHorizontal
	case "wrench":
		return ModeWrench
	case "wire":
		return ModeWire
	case "palette":
		return ModePalette
	case "text":
		return ModeText
	case "undo":
		return ModeUndo
	case "redo":
		return ModeRedo
	}
	return EndModeList
}

var BlockList = []Block{
	BlockTurf,
	BlockBedrock,
	BlockFall,
	BlockCracked,
	BlockClose,
	BlockPhase,
	BlockSpike,
	BlockEmpty,

	BlockLadder,
	BlockLadderCracked,
	BlockLadderExit,
	BlockBar,
	BlockHideout,
	BlockTransporter,
	BlockTransporterExit,
	BlockBombLit,

	BlockPlayer1,
	BlockPlayer2,
	BlockPlayer3,
	BlockPlayer4,
	BlockDemon,
	BlockDemonRegen,
	BlockFly,
	BlockFlyRegen,

	BlockGem,
	BlockDoorHidden,
	BlockDoorVisible,
	BlockDoorLocked,
	BlockKey,
	BlockEmpty,
	BlockEmpty,
	BlockEmpty,

	BlockJumpBoots,
	BlockBox,
	BlockBomb,
	BlockJetpack,
	BlockDisguise,
	BlockDrill,
	BlockFlamethrower,
	BlockEmpty,

	BlockReeds,
	BlockFlowers,
	BlockMoss,
	BlockGrass,
	BlockCattail,
	BlockCactus1,
	BlockCactus2,
	BlockVine,

	BlockNest,
	BlockSkull,
	BlockDots,
	BlockBoulder,
	BlockMush1,
	BlockMush2,
	BlockChain,
	BlockGear,
}
